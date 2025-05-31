package cache

import (
	"context"
	"log"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func InitRedis() (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if err := rds.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rds, nil
}

func CloseRedisConnection(rds *redis.Client) error {
	if err := rds.Close(); err != nil {
		return err
	}

	return nil
}

type TaskList interface {
	CreateTask(userID int, task entity.Task) (int, error)
	GetAllTask(userID int) ([]entity.Task, error)
	GetTaskByID(userID, id int) (entity.Task, error)
	UpdateTask(userID, taskId int, desc string) error
	DeleteTask(userID, taskID int) error
	SaveTasksToCache(ctx context.Context, userID int, tasks []entity.Task) error
	SaveTaskToCache(ctx context.Context, userID int, task entity.Task) error
}

type RedisRepository struct {
	TaskList
}

func NewRedisRepository(rdb *redis.Client, repo *repository.Repository) *RedisRepository {
	return &RedisRepository{
		TaskList: NewTaskCache(rdb, repo),
	}
}
