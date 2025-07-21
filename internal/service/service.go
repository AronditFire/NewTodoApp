package service

import (
	"net/http"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/cache"
	"github.com/AronditFire/todo-app/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TaskList interface {
	CreateTask(userID int, task entity.Task) (int, error)
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
	ParseAccessToken(accessTokenStr string) (*TokenClaims, error)
	ParseRefreshToken(refreshTokenStr string) (int, error)
	RenewTokens(id int) (string, string, error)
	GoogleLogin() string
	GetClientGoogle(code string) (*http.Client, error)
}

type ParsingJSON interface {
	ParseJSON(bindfile entity.BindFile) error
	GetJsonTable() ([]map[string]any, error)
}

type Service struct {
	TaskList
	Authorization
	ParsingJSON
}

func NewService(crepo *cache.RedisRepository, repo *repository.Repository, id, secret, rURL string) *Service {
	return &Service{
		TaskList:      NewTaskService(crepo.TaskList),
		Authorization: NewAuthService(repo.Authorization, id, secret, rURL),
		ParsingJSON:   NewParseService(repo.ParsingJSON),
	}
}
