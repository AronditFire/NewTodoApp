package service

import (
	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
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
	LoginUser(userLogin entity.UserAuthRequest) (entity.User, error)
}

type Service struct {
	TaskList
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		TaskList:      NewTaskService(repo.TaskList),
		Authorization: NewAuthService(repo.Authorization),
	}
}
