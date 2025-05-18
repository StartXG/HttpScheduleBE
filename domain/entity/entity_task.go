package entity

import (

	"gorm.io/gorm"
)

type TaskCenter struct {
	*gorm.Model
	TaskName       string `json:"task_name" gorm:"column:task_name;type:varchar(255);not null;comment:任务名称"`
	TaskCron       string `json:"task_cron" gorm:"column:task_cron;type:varchar(255);not null;comment:任务cron表达式"`
	TaskUrl        string `json:"task_url" gorm:"column:task_url;type:varchar(255);not null;comment:任务请求地址"`
	TaskMethod     string `json:"task_method" gorm:"column:task_method;type:varchar(255);not null;comment:任务请求方法"`
	TaskHeader     string `json:"task_header" gorm:"column:task_header;type:varchar(255);not null;comment:任务请求头"`
	TaskBody       string `json:"task_body" gorm:"column:task_body;type:varchar(255);not null;comment:任务请求体"`
	// TaskStatus     bool `json:"task_status" gorm:"column:task_status;type:varchar(255);not null;comment:任务状态"`
	IsTaskEnabled bool `json:"is_task_enabled" gorm:"column:is_task_enabled;type:varchar(255);not null;comment:任务是否启用"`
	TaskRemark     string `json:"task_remark" gorm:"column:task_remark;type:varchar(255);not null;comment:任务备注"`
}

func (TaskCenter) TableName() string {
	return "task_center"
}