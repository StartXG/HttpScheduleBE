package service

import (
	"HttpScheduleBE/services/execution/repo"
	httpTypes "HttpScheduleBE/services/execution/types"
)

type ExecutionCenter interface {
	GetAllExecution() (*[]httpTypes.ResponseExecutionCenter, error)
}

type Service struct {
	r repo.Repository
}

func NewExecutionCenterService(r repo.Repository) *Service {
	if err := r.Migration(); err != nil {
		panic("Failed to migrate ExecutionCenter repository: " + err.Error())
	}
	return &Service{r: r}
}

func (s *Service) GetAllExecution() (*[]httpTypes.ResponseExecutionCenter, error) {
	executions, err := s.r.GetExecutions()
	if err != nil {
		return nil, err
	}
	var httpExecutions []httpTypes.ResponseExecutionCenter
	for _, execution := range *executions {
		httpExecutions = append(httpExecutions, httpTypes.ResponseExecutionCenter{
			TaskID:    execution.TaskID,
			Status:    execution.Status,
			StartTime: execution.StartTime,
			EndTime:   execution.EndTime,
			ErrorLog:  execution.ErrorLog,
		})
	}
	return &httpExecutions, nil
}
