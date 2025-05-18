package entity

import "gorm.io/gorm"

type ExecutionCenter struct {
	*gorm.Model
	TaskID    uint   `json:"task_id" gorm:"column:task_id;not null;comment:关联的任务ID"`
	Status    string `json:"status" gorm:"column:status;type:varchar(255);not null;comment:任务执行状态"`
	StartTime string `json:"start_time" gorm:"column:start_time;type:varchar(255);not null;comment:任务开始时间"`
	EndTime   string `json:"end_time" gorm:"column:end_time;type:varchar(255);comment:任务结束时间"`
	ErrorLog  string `json:"error_log" gorm:"column:error_log;type:text;comment:错误日志"`
}

func (ExecutionCenter) TableName() string {
	return "execution_center"
}
