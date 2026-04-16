package backend

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// TaskDefinition 定义一个可执行的任务
type TaskDefinition struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Type        string `json:"type"` // npm, make, custom
	Description string `json:"description,omitempty"`
}

// DetectTasks 自动检测项目中的任务
func DetectTasks(rootPath string) ([]TaskDefinition, error) {
	var tasks []TaskDefinition

	// 1. 检测 package.json
	pkgPath := filepath.Join(rootPath, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			for name, cmd := range pkg.Scripts {
				tasks = append(tasks, TaskDefinition{
					Name:        name,
					Command:     "npm run " + name,
					Type:        "npm",
					Description: cmd,
				})
			}
		}
	}

	// 2. 检测 Makefile
	mkPath := filepath.Join(rootPath, "Makefile")
	if data, err := os.ReadFile(mkPath); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			// 简单的 Makefile target 识别：以字母开头，后面跟冒号
			if strings.Contains(line, ":") && !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "#") {
				parts := strings.Split(line, ":")
				target := strings.TrimSpace(parts[0])
				if target != "" && !strings.Contains(target, "=") {
					tasks = append(tasks, TaskDefinition{
						Name:    target,
						Command: "make " + target,
						Type:    "make",
					})
				}
			}
		}
	}

	return tasks, nil
}
