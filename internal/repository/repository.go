package repository

import (
	"github.com/AronditFire/todo-app/entity"
	"gorm.io/gorm"
)

type TaskList interface {
	CreateTask(userID int, task entity.Task) error
	GetAllTask(userID int) ([]entity.Task, error)
	GetTaskByID(userID, id int) (entity.Task, error)
	UpdateTask(userID, taskId int, desc string) error
	DeleteTask(userID, taskID int) error
}

type Authorization interface {
	CreateUser(userReg entity.UserRegisterRequest) error
	GetUser(username string) (entity.User, error)
	GetUserByID(id int) (entity.User, error)
}

type Repository struct {
	TaskList
	Authorization
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		TaskList:      NewTaskRepo(db),
		Authorization: NewAuthRepo(db),
	}
}
