package database

import (
	"HttpScheduleBE/config"
	domainExecutionCenter "HttpScheduleBE/domain/domain_execution_center"
	domainTaskCenter "HttpScheduleBE/domain/domain_task_center"
	"fmt"
)

type Databases struct {
	TaskCenterRepository      *domainTaskCenter.Repository
	ExecutionCenterRepository *domainExecutionCenter.Repository
}

var AppConfig = &config.Config{}

func CreateDBs(cfg *config.Config) *Databases {
	AppConfig = cfg
	DSN := cfg.DatabaseUser + ":" + cfg.DatabasePassword + "@tcp(" + cfg.DatabaseHost + ":" + fmt.Sprint(cfg.DatabasePort) + ")/" + cfg.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := NewMySQLDb(DSN)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	return &Databases{
		TaskCenterRepository:      domainTaskCenter.NewTaskCenterRepository(db),
		ExecutionCenterRepository: domainExecutionCenter.NewExecutionCenterRepository(db),
	}
}
