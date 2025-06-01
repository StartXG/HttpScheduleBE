package internal

import (
	"HttpScheduleBE/entity"
	"HttpScheduleBE/pkgs/executor"
	"HttpScheduleBE/services/execution/repo"
	"fmt"
)


func LogRecord(data <-chan executor.ExecuteResultForRecord, repository *repo.Repository) {
	for {
		select {
		case logData := <-data:
			var execution entity.ExecutionCenter
			// 将 ExecuteResultForRecord 转换为 ExecutionCenter
			execution.TaskID = logData.TaskID
			execution.Status = "completed" // 假设状态为 completed
			execution.StartTime = logData.StartTime.Format("2006-01-02 15:04:05")
			execution.EndTime = logData.EndTime.Format("2006-01-02 15:04:05")
			execution.ErrorLog = logData.ErrMsg
			// if logData.Result != "" {
			// 	execution.ErrorLog = logData.Result
			// } 
			if err := repository.CreateExecution(&execution); err != nil {
				fmt.Printf("Failed to create execution record: %v\n", err)
			} else {
				fmt.Printf("Execution record created successfully: %+v\n", execution)
			}
		default:
			// 如果没有数据，继续等待
			continue
		}
	}
}