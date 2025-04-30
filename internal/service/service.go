package service

import (
	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
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
	LoginUser(userLogin entity.UserAuthRequest) (string, string, error)
	ParseAccessToken(accessTokenStr string) (*tokenClaims, error)
	ParseRefreshToken(refreshTokenStr string) (int, error)
	RenewTokens(id int) (string, string, error)
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
