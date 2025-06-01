package database

import (
	"HttpScheduleBE/config"
	ExecutionRepo "HttpScheduleBE/services/execution/repo"
	TaskRepo "HttpScheduleBE/services/task/repo"
	"fmt"
)

type Databases struct {
	TaskCenterRepository      *TaskRepo.Repository
	ExecutionCenterRepository *ExecutionRepo.Repository
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
		TaskCenterRepository:      TaskRepo.NewTaskCenterRepository(db),
		ExecutionCenterRepository: ExecutionRepo.NewExecutionCenterRepository(db),
	}
}
