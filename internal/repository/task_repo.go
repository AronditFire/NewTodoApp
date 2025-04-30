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

func (r *TaskRepo) CreateTask(userID int, task entity.Task) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	task.UserID = userID

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (r *TaskRepo) GetAllTask(userID int) ([]entity.Task, error) {
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

	if err := tx.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return tasks, tx.Commit().Error
}

func (r *TaskRepo) GetTaskByID(userID, id int) (entity.Task, error) {
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

	if err := tx.Where("user_id = ? AND id = ?", userID, id).First(&task).Error; err != nil {
		tx.Rollback()
		return entity.Task{}, err
	}

	return task, tx.Commit().Error
}

func (r *TaskRepo) UpdateTask(userID, taskId int, desc string) error {
	var task entity.Task

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("user_id = ? AND id = ?", userID, taskId).First(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	task.Description = desc

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *TaskRepo) DeleteTask(userID, taskID int) error {
	var task entity.Task
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("user_id = ? AND id = ?", userID, taskID).First(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
