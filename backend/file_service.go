package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileSystemService 文件系统服务实现
type FileSystemService struct {
	ctx context.Context
}

// NewFileSystemService 创建文件系统服务
func NewFileSystemService() *FileSystemService {
	return &FileSystemService{}
}

// SetContext 设置上下文（由适配器调用）
func (f *FileSystemService) SetContext(ctx context.Context) {
	f.ctx = ctx
}

// OpenFolderDialog 打开文件夹选择对话框（使用 Wails 原生 API）
func (f *FileSystemService) OpenFolderDialog() (string, error) {
	if f.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	// 使用 Wails 原生的目录选择对话框
	path, err := runtime.OpenDirectoryDialog(f.ctx, runtime.OpenDialogOptions{
		Title: "选择项目文件夹",
	})

	if err != nil {
		return "", err
	}

	// 用户取消时返回空字符串
	if path == "" {
		return "", fmt.Errorf("user cancelled")
	}

	return path, nil
}

// OpenFileDialog 打开文件选择对话框（使用 Wails 原生 API）
func (f *FileSystemService) OpenFileDialog() (string, error) {
	if f.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	path, err := runtime.OpenFileDialog(f.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})

	if err != nil {
		return "", err
	}

	if path == "" {
		return "", fmt.Errorf("user cancelled")
	}

	return path, nil
}

// SaveFileDialog 保存文件对话框（使用 Wails 原生 API）
func (f *FileSystemService) SaveFileDialog() (string, error) {
	if f.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	path, err := runtime.SaveFileDialog(f.ctx, runtime.SaveDialogOptions{
		Title: "保存文件",
	})

	if err != nil {
		return "", err
	}

	if path == "" {
		return "", fmt.Errorf("user cancelled")
	}

	return path, nil
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
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDir 列出目录内容（支持懒加载）
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

// SetProjectRoot 设置项目根目录
func (f *FileSystemService) SetProjectRoot(path string) error {
	// 验证路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	
	// 切换到该目录
	return os.Chdir(path)
}

// CreateFile 创建新文件
func (f *FileSystemService) CreateFile(path string) error {
	// 确保父目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 检查文件是否已存在
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists")
	}

	// 创建空文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	return nil
}

// CreateDirectory 创建新目录
func (f *FileSystemService) CreateDirectory(path string) error {
	// 检查目录是否已存在
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("directory already exists")
	}

	// 递归创建目录
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	return nil
}

// DeleteFileOrDirectory 删除文件或目录
func (f *FileSystemService) DeleteFileOrDirectory(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist")
	}

	// 删除文件或目录（递归）
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return nil
}

// RenameFileOrDirectory 重命名文件或目录
func (f *FileSystemService) RenameFileOrDirectory(oldPath, newPath string) error {
	// 检查旧路径是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist")
	}

	// 检查新路径是否已存在
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("target path already exists")
	}

	// 确保父目录存在
	newDir := filepath.Dir(newPath)
	if err := os.MkdirAll(newDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// 重命名
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename: %v", err)
	}

	return nil
}

// MoveFileOrDirectory 移动文件或目录
func (f *FileSystemService) MoveFileOrDirectory(sourcePath, targetPath string) error {
	// 复用重命名逻辑
	return f.RenameFileOrDirectory(sourcePath, targetPath)
}

// GetFileStats 获取文件统计信息
func (f *FileSystemService) GetFileStats(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %v", err)
	}

	return &FileInfo{
		Name:    info.Name(),
		Path:    path,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime().Unix(),
	}, nil
}

// SearchFiles 在目录中搜索文件（支持关键词过滤）
func (f *FileSystemService) SearchFiles(rootPath, keyword string, maxResults int) ([]FileInfo, error) {
	var results []FileInfo
	
	err := filepath.WalkDir(rootPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil // 跳过错误
		}

		// 跳过隐藏文件和 node_modules
		if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 匹配关键词
		if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(keyword)) {
			info := FileInfo{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
				Path:  path,
			}

			if !entry.IsDir() {
				fileInfo, _ := entry.Info()
				info.Size = fileInfo.Size()
				info.ModTime = fileInfo.ModTime().Unix()
			}

			results = append(results, info)
			
			// 限制结果数量
			if len(results) >= maxResults {
				return filepath.SkipDir
			}
		}

		return nil
	})

	return results, err
}

// ReadFileWithEncoding 读取文件（支持检测编码，预留 UTF-8 转码功能）
func (f *FileSystemService) ReadFileWithEncoding(path string) (string, error) {
	// 当前实现假设文件为 UTF-8
	// 未来可以添加编码检测和转换逻辑
	content, err := f.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	// 检查是否包含 BOM 并移除
	if len(content) >= 3 && content[0:3] == "\xef\xbb\xbf" {
		content = content[3:]
	}
	
	return content, nil
}

// GetFileExtension 获取文件扩展名
func (f *FileSystemService) GetFileExtension(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

// IsTextFile 判断是否为文本文件
func (f *FileSystemService) IsTextFile(path string) bool {
	ext := strings.ToLower(f.GetFileExtension(path))
	
	// 常见的文本文件扩展名
	textExtensions := map[string]bool{
		"txt": true, "md": true, "json": true, "xml": true, "html": true, "htm": true,
		"css": true, "js": true, "ts": true, "vue": true, "jsx": true, "tsx": true,
		"go": true, "py": true, "java": true, "c": true, "cpp": true, "h": true,
		"rb": true, "php": true, "sh": true, "yaml": true, "yml": true, "toml": true,
		"sql": true, "graphql": true, "rs": true, "swift": true, "kt": true,
		"log": true, "cfg": true, "ini": true, "conf": true, "env": true,
		"gitignore": true, "dockerfile": true, "makefile": true,
	}
	
	return textExtensions[ext]
}

// BackupFile 备份文件（添加 .bak 后缀）
func (f *FileSystemService) BackupFile(path string) error {
	backupPath := path + ".bak"
	
	// 如果备份文件已存在，先删除
	os.Remove(backupPath)
	
	// 复制文件
	content, err := f.ReadFile(path)
	if err != nil {
		return err
	}
	
	return f.WriteFile(backupPath, content)
}

// CopyFileOrDirectory 复制文件或目录
func (f *FileSystemService) CopyFileOrDirectory(sourcePath, targetPath string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source path does not exist: %v", err)
	}

	if sourceInfo.IsDir() {
		return f.copyDirectory(sourcePath, targetPath)
	}
	
	return f.copyFile(sourcePath, targetPath)
}

// copyFile 复制单个文件
func (f *FileSystemService) copyFile(source, target string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	// 确保目标目录存在
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(target, content, 0644)
}

// copyDirectory 递归复制目录
func (f *FileSystemService) copyDirectory(source, target string) error {
	// 创建目标目录
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		targetPath := filepath.Join(target, entry.Name())

		if entry.IsDir() {
			if err := f.copyDirectory(sourcePath, targetPath); err != nil {
				return err
			}
		} else {
			if err := f.copyFile(sourcePath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDirectoryTree 获取完整的目录树（用于初始化加载）
func (f *FileSystemService) GetDirectoryTree(path string, depth int) ([]FileInfo, error) {
	if depth <= 0 {
		return nil, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
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

// TouchFile 更新文件修改时间或创建空文件
func (f *FileSystemService) TouchFile(path string) error {
	// 如果文件不存在，创建空文件
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return f.CreateFile(path)
	}

	// 更新修改时间
	now := time.Now()
	if err := os.Chtimes(path, now, now); err != nil {
		return fmt.Errorf("failed to update file time: %v", err)
	}

	return nil
}
