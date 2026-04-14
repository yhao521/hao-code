package main

import (
	"os"
	"path/filepath"
	"strings"
)

// FileSystemService 文件系统服务实现
type FileSystemService struct{}

// NewFileSystemService 创建文件系统服务
func NewFileSystemService() *FileSystemService {
	return &FileSystemService{}
}

// ReadFile 读取文件内容
func (f *FileSystemService) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile 写入文件内容
func (f *FileSystemService) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDir 列出目录内容
func (f *FileSystemService) ListDir(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		// 跳过隐藏文件和 node_modules
		if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules" {
			continue
		}

		info := FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Path:  filepath.Join(path, entry.Name()),
		}

		if !entry.IsDir() {
			fileInfo, _ := entry.Info()
			info.Size = fileInfo.Size()
			info.ModTime = fileInfo.ModTime().Unix()
		}

		files = append(files, info)
	}

	return files, nil
}

// GetProjectRoot 获取项目根目录
func (f *FileSystemService) GetProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}
