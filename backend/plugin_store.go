package backend

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type PluginStore struct {
	rootDir string
}

func NewPluginStore(rootDir string) *PluginStore {
	return &PluginStore{rootDir: rootDir}
}

// InstallFromURL 从 URL 下载并安装插件（支持 .zip 格式）
func (s *PluginStore) InstallFromURL(url string) (*PluginManifest, error) {
	// 1. 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	// 2. 创建临时文件
	tmpFile, err := os.CreateTemp("", "plugin-*.zip")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, err
	}
	tmpFile.Close()

	// 3. 解压并安装
	return s.installFromZip(tmpFile.Name())
}

// InstallFromFile 从本地 zip 文件安装插件
func (s *PluginStore) InstallFromFile(zipPath string) (*PluginManifest, error) {
	return s.installFromZip(zipPath)
}

// installFromZip 核心解压与安装逻辑
func (s *PluginStore) installFromZip(zipPath string) (*PluginManifest, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var manifest *PluginManifest
	pluginName := ""

	for _, f := range r.File {
		// 简单处理：假设根目录即为插件名，或者从 package.json 读取
		if f.Name == "package.json" || filepath.Base(f.Name) == "package.json" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			var m PluginManifest
			if err := json.NewDecoder(rc).Decode(&m); err == nil {
				manifest = &m
				pluginName = m.Name
			}
			rc.Close()
		}
	}

	if manifest == nil {
		return nil, fmt.Errorf("invalid plugin: package.json not found")
	}

	// 确定安装路径
	installPath := filepath.Join(s.rootDir, pluginName)
	os.MkdirAll(installPath, 0755)

	// 再次遍历解压所有文件
	for _, f := range r.File {
		path := filepath.Join(installPath, f.Name)

		// 安全检查：防止 Zip Slip 漏洞
		if !filepath.HasPrefix(path, filepath.Clean(installPath)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return nil, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return nil, err
		}
	}

	return manifest, nil
}

// Uninstall 卸载插件
func (s *PluginStore) Uninstall(pluginName string) error {
	pluginPath := filepath.Join(s.rootDir, pluginName)
	return os.RemoveAll(pluginPath)
}

// ListInstalled 列出所有已安装的插件名称
func (s *PluginStore) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(s.rootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
