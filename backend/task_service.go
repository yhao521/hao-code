package backend

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

type TaskService struct {
	ctx context.Context
}

type TaskItem struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Source  string `json:"source"` // e.g., "package.json", "Taskfile.yml"
}

func NewTaskService(ctx context.Context) *TaskService {
	return &TaskService{ctx: ctx}
}

// GetTasks 扫描工作区根目录下的任务配置文件并返回任务列表
func (s *TaskService) GetTasks(rootPath string) ([]TaskItem, error) {
	var tasks []TaskItem

	// 1. 解析 package.json
	pkgPath := filepath.Join(rootPath, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			for name, cmd := range pkg.Scripts {
				tasks = append(tasks, TaskItem{
					Name:    name,
					Command: cmd,
					Source:  "package.json",
				})
			}
		}
	}

	// 2. 解析 Taskfile.yml (简单实现：仅识别顶层 key，实际需引入 yaml 库)
	taskfilePath := filepath.Join(rootPath, "Taskfile.yml")
	if _, err := os.Stat(taskfilePath); err == nil {
		// 这里为了简化，暂时不引入 gopkg.in/yaml.v3，仅作为占位符
		// 实际生产中应使用 yaml.Unmarshal 解析 tasks 字段
		tasks = append(tasks, TaskItem{
			Name:    "taskfile-detect",
			Command: "echo 'Taskfile.yml detected (parsing not fully implemented yet)'",
			Source:  "Taskfile.yml",
		})
	}

	return tasks, nil
}

// RunTask 在指定目录下执行命令
func (s *TaskService) RunTask(rootPath string, command string) error {
	// 使用 sh -c 或 cmd /c 来执行复杂的脚本字符串
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.CommandContext(s.ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(s.ctx, "sh", "-c", command)
	}
	cmd.Dir = rootPath

	// 注意：在 Wails 中，通常我们会将输出重定向到 WebSocket 或 PTY
	// 这里为了演示，我们直接运行。如果需要实时输出，需要结合 TerminalService
	_, err := cmd.CombinedOutput()
	return err
}

func isWindows() bool {
	return filepath.Separator == '\\'
}
