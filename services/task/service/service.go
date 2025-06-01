package service

import (
	"HttpScheduleBE/entity"
	"HttpScheduleBE/services/task/repo"
	httpTypes "HttpScheduleBE/services/task/types"
	"fmt"
)

type TaskCenter interface {
	CreateTask(req *httpTypes.RequestTaskCenter) error
	UpdateTask(taskId string, req *httpTypes.RequestTaskCenter) error
	DeleteTask(id string) error
	GetAllTasks() ([]httpTypes.ResponseTaskCenter, error)
}

type Service struct {
	r repo.Repository
}

func NewTaskCenterService(r repo.Repository) *Service {
	if err := r.Migration(); err != nil {
		panic("Failed to migrate TaskCenter repository: " + err.Error())
	}
	return &Service{
		r: r,
	}
}

func (s *Service) CreateTask(req *httpTypes.RequestTaskCenter) error {
	if err := s.r.CreateTask(&entity.TaskCenter{
		TaskName:      req.TaskName,
		TaskUrl:       req.TaskUrl,
		TaskMethod:    req.TaskMethod,
		TaskHeader:    req.TaskHeader,
		TaskBody:      req.TaskBody,
		TaskCron:      req.TaskCron,
		TaskRemark:    req.TaskRemark,
		IsTaskEnabled: req.IsTaskEnabled,
	}); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (s *Service) UpdateTask(taskId string, req *httpTypes.RequestTaskCenter) error {
	task, err := s.r.GetTaskByID(taskId)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	fmt.Println(task)
	task.TaskName = req.TaskName
	task.TaskUrl = req.TaskUrl
	task.TaskMethod = req.TaskMethod
	task.TaskHeader = req.TaskHeader
	task.TaskBody = req.TaskBody
	task.TaskCron = req.TaskCron
	task.TaskRemark = req.TaskRemark
	task.IsTaskEnabled = req.IsTaskEnabled
	fmt.Println(task)

	return s.r.UpdateTask(task)
}

func (s *Service) DeleteTask(id string) error {
	_, err := s.r.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	return s.r.DeleteTask(id)
}

func (s *Service) GetAllTasks() ([]httpTypes.ResponseTaskCenter, error) {
	tasks, err := s.r.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	var responseTasks []httpTypes.ResponseTaskCenter
	for _, task := range tasks {
		responseTasks = append(responseTasks, httpTypes.ResponseTaskCenter{
			TaskId:        task.ID,
			TaskName:      task.TaskName,
			TaskUrl:       task.TaskUrl,
			TaskMethod:    task.TaskMethod,
			TaskHeader:    task.TaskHeader,
			TaskBody:      task.TaskBody,
			TaskCron:      task.TaskCron,
			TaskRemark:    task.TaskRemark,
			IsTaskEnabled: task.IsTaskEnabled,
		})
	}
	// Return the response tasks
	return responseTasks, nil
}
