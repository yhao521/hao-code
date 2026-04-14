package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to Hao-Code Editor!", name)
}

// ReadFile reads content from a file
func (a *App) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func (a *App) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDir lists files in a directory
func (a *App) ListDir(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info := FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Path:    filepath.Join(path, entry.Name()),
		}
		
		if !entry.IsDir() {
			fileInfo, _ := entry.Info()
			info.Size = fileInfo.Size()
		}
		
		files = append(files, info)
	}

	return files, nil
}

// GetProjectRoot returns the project root directory
func (a *App) GetProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}

// FileInfo represents information about a file
type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size,omitempty"`
	IsDir bool   `json:"isDir"`
}
