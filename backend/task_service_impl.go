package backend

import (
	"os/exec"
)

// TaskService 任务服务，用于管理和执行项目任务
type TaskService struct{}

// NewTaskService 创建任务服务实例
func NewTaskService() *TaskService {
	return &TaskService{}
}

// GetTasks 获取项目中的任务列表
func (s *TaskService) GetTasks(rootPath string) ([]TaskItem, error) {
	taskDefs, err := DetectTasks(rootPath)
	if err != nil {
		return nil, err
	}

	// 将 TaskDefinition 转换为 TaskItem
	var taskItems []TaskItem
	for _, def := range taskDefs {
		taskItems = append(taskItems, TaskItem{
			Name:        def.Name,
			Command:     def.Command,
			Type:        def.Type,
			Description: def.Description,
		})
	}

	return taskItems, nil
}

// RunTask 运行指定任务
func (s *TaskService) RunTask(rootPath string, command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = rootPath
	return cmd.Run()
}
