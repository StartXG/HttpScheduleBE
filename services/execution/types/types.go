package types

type ResponseExecutionCenter struct {
	TaskID    uint   `json:"task_id"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	ErrorLog  string `json:"error_log"`
}

type ResponseExecutingTask struct {
	TaskID uint   `json:"task_id"`
	Status string `json:"status"`
	Name   string `json:"name"`
}
