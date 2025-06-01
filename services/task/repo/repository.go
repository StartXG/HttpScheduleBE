package repo

import (
	"HttpScheduleBE/entity"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewTaskCenterRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migration() error {
	return r.db.AutoMigrate(&entity.TaskCenter{})
}

func (r *Repository) CreateTask(task *entity.TaskCenter) error {
	return r.db.Create(task).Error
}

func (r *Repository) GetTaskByID(id string) (*entity.TaskCenter, error) {
	var task entity.TaskCenter
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) UpdateTask(task *entity.TaskCenter) error {
	return r.db.Save(task).Error
}

func (r *Repository) DeleteTask(id string) error {
	task := &entity.TaskCenter{}
	if err := r.db.First(task, id).Error; err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	return r.db.Delete(task).Error
}

func (r *Repository) GetAllTasks() ([]entity.TaskCenter, error) {
	var tasks []entity.TaskCenter
	err := r.db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
