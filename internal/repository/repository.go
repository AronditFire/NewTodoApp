package repository

import (
	"github.com/AronditFire/todo-app/entity"
	"gorm.io/gorm"
)

type TaskList interface {
	CreateTask(task entity.Task) error
	GetAllTask() ([]entity.Task, error)
	GetTaskByID(id int) (entity.Task, error)
	UpdateTask(id int, desc string) error
	DeleteTask(id int) error
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
