package repository

import (
	"github.com/AronditFire/todo-app/entity"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(task entity.Task) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (r *TaskRepo) GetAllTask() ([]entity.Task, error) {
	var tasks []entity.Task

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Find(&tasks).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return tasks, tx.Commit().Error
}

func (r *TaskRepo) GetTaskByID(id int) (entity.Task, error) {
	var task entity.Task

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entity.Task{}, err
	}

	if err := tx.First(&task, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return entity.Task{}, err
	}

	return task, tx.Commit().Error
}

func (r *TaskRepo) UpdateTask(id int, desc string) error {
	var task entity.Task
	task.ID = id
	task.Description = desc

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *TaskRepo) DeleteTask(id int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Delete(&entity.Task{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
