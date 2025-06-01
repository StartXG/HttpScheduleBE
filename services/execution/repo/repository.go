package repo

import (
	"HttpScheduleBE/entity"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewExecutionCenterRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migration() error {
	return r.db.AutoMigrate(&entity.ExecutionCenter{})
}

func (r *Repository) GetExecutions() (*[]entity.ExecutionCenter, error) {
	var execution []entity.ExecutionCenter
	err := r.db.Order("created_at desc").Find(&execution).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *Repository) CreateExecution(execution *entity.ExecutionCenter) error {
	return r.db.Create(execution).Error
}

// func (r *Repository) CreateExecution(execution *TaskExecution) error {
// 	return r.db.Create(execution).Error
// }

// func (r *Repository) UpdateExecution(execution *TaskExecution) error {
// 	return r.db.Save(execution).Error
// }

// func (r *Repository) GetExecutionByID(id uint) (*TaskExecution, error) {
// 	var execution TaskExecution
// 	err := r.db.First(&execution, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &execution, nil
// }

// func (r *Repository) GetExecutionsByTaskID(taskID uint) ([]TaskExecution, error) {
// 	var executions []TaskExecution
// 	err := r.db.Where("task_id = ?", taskID).Find(&executions).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return executions, nil
// }
