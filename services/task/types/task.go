package types

type RequestTaskCenter struct {
	TaskName      string `json:"task_name" binding:"required"`
	TaskCron      string `json:"task_cron" binding:"required"`
	TaskUrl       string `json:"task_url" binding:"required"`
	TaskMethod    string `json:"task_method" default:"GET"`
	TaskHeader    string `json:"task_header"`
	TaskBody      string `json:"task_body"`
	IsTaskEnabled bool   `json:"is_task_enabled" default:"false"`
	TaskRemark    string `json:"task_remark"`
}

type ResponseTaskCenter struct {
	TaskId        uint   `json:"task_id"`
	TaskName      string `json:"task_name"`
	TaskCron      string `json:"task_cron"`
	TaskUrl       string `json:"task_url"`
	TaskMethod    string `json:"task_method"`
	TaskHeader    string `json:"task_header"`
	TaskBody      string `json:"task_body"`
	IsTaskEnabled bool   `json:"is_task_enabled" default:"false"`
	TaskRemark    string `json:"task_remark"`
}
