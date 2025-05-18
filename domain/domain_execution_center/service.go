package domainexecutioncenter

import (
	httpTypes "HttpScheduleBE/api/types"
)

type Service struct {
	r Repository
}

func NewExecutionCenterService(r Repository) *Service {
	r.Migration()
	return &Service{r: r}
}

// func (s *Service) StartExecution(taskID uint) error {
// 	execution := &TaskExecution{
// 		TaskID:    taskID,
// 		Status:    "running",
// 		StartTime: time.Now().Format("2006-01-02 15:04:05"),
// 	}
// 	return s.r.CreateExecution(execution)
// }

// func (s *Service) StopExecution(executionID uint) error {
// 	execution, err := s.r.GetExecutionByID(executionID)
// 	if err != nil {
// 		return err
// 	}
// 	execution.Status = "stopped"
// 	execution.EndTime = time.Now().Format("2006-01-02 15:04:05")
// 	return s.r.UpdateExecution(execution)
// }

func (s *Service) GetAllExecution() (*[]httpTypes.ResponseExecutionCenter, error) {
	executions,err := s.r.GetExecutions()
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

