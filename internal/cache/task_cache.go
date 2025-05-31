package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
	"github.com/redis/go-redis/v9"
)

const TTL = 60 // время жизни кэша в секундах

type TaskCache struct {
	rdb  *redis.Client
	repo repository.TaskList
}

func NewTaskCache(rdb *redis.Client, repo repository.TaskList) *TaskCache {
	return &TaskCache{
		rdb:  rdb,
		repo: repo,
	}
}

func (r *TaskCache) SaveTasksToCache(ctx context.Context, userID int, tasks []entity.Task) error {
	key := fmt.Sprintf("user:%d:tasks", userID)

	hm := make(map[string]interface{}, len(tasks))
	for _, t := range tasks {
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		// поле―строка, значение―JSON
		hm[fmt.Sprint(t.ID)] = b
	}
	// Записываем сразу всё:
	if err := r.rdb.HSet(ctx, key, hm).Err(); err != nil {
		return err
	}
	// Ставим TTL в 60 секунд

	if err := r.rdb.Expire(ctx, key, time.Second*time.Duration(TTL)).Err(); err != nil {
		return err
	}
	return nil
}

func (r *TaskCache) SaveTaskToCache(ctx context.Context, userID int, task entity.Task) error {
	key := fmt.Sprintf("user:%d:tasks", userID)

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	pipe := r.rdb.TxPipeline()
	pipe.HSet(ctx, key, fmt.Sprint(task.ID), taskJSON)
	pipe.Expire(ctx, key, time.Second*TTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline while saving one task in redis: %w", err)
	}

	return nil
}

func (r *TaskCache) CreateTask(userID int, task entity.Task) (int, error) {
	key := fmt.Sprintf("user:%d:tasks", userID)

	id, err := r.repo.CreateTask(userID, task)
	if err != nil {
		return 0, err
	}

	task.ID = id
	task.UserID = userID

	raw, err := json.Marshal(task)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal task: %w", err)
	}

	pipe := r.rdb.TxPipeline()
	pipe.HSet(ctx, key, fmt.Sprint(id), raw)
	pipe.Expire(ctx, key, TTL*time.Second)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return task.ID, nil
}

func (r *TaskCache) GetAllTask(userID int) ([]entity.Task, error) {
	key := fmt.Sprintf("user:%d:tasks", userID)

	data, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from cache: %w", err)
	}
	var tasks []entity.Task
	if len(data) > 0 {
		// Кеш есть, десериализуем
		for _, raw := range data {
			var t entity.Task
			if err := json.Unmarshal([]byte(raw), &t); err != nil {
				return nil, err
			}
			tasks = append(tasks, t)
		}
		if err = r.rdb.Expire(ctx, key, time.Second*TTL).Err(); err != nil {
			return nil, fmt.Errorf("failed to set expiration for cache key: %w", err)
		}
		return tasks, nil
	}

	// Кеша нет, получаем из репозитория
	tasks, err = r.repo.GetAllTask(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from repository: %w", err)
	}

	if err := r.SaveTasksToCache(ctx, userID, tasks); err != nil {
		return nil, fmt.Errorf("failed to save tasks to cache: %w", err)
	}

	return tasks, nil
}

func (r *TaskCache) GetTaskByID(userID, id int) (entity.Task, error) {
	key := fmt.Sprintf("user:%d:tasks", userID)

	data, err := r.rdb.HGet(ctx, key, fmt.Sprint(id)).Result()
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get task from cache: %w", err)
	}
	var task entity.Task
	if len(data) > 0 {

		if err := json.Unmarshal([]byte(data), &task); err != nil {
			return entity.Task{}, err
		}

		if err = r.rdb.Expire(ctx, key, time.Second*TTL).Err(); err != nil {
			return entity.Task{}, fmt.Errorf("failed to set expiration for cache key: %w", err)
		}
		return task, nil
	}

	// Кеша нет, получаем из репозитория
	task, err = r.repo.GetTaskByID(userID, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get tasks from repository: %w", err)
	}

	if err := r.SaveTaskToCache(ctx, userID, task); err != nil {
		return entity.Task{}, fmt.Errorf("failed to save tasks to cache: %w", err)
	}

	return task, nil
}

func (r *TaskCache) UpdateTask(userID, taskId int, desc string) error {
	if err := r.repo.UpdateTask(userID, taskId, desc); err != nil {
		return fmt.Errorf("failed to update task in repository: %w", err)
	}

	task, err := r.repo.GetTaskByID(userID, taskId)
	if err != nil {
		return fmt.Errorf("failed to get updated task: %w", err)
	}

	if err := r.SaveTaskToCache(ctx, userID, task); err != nil {
		return fmt.Errorf("failed to save updated task to cache: %w", err)
	}

	return nil
}

func (r *TaskCache) DeleteTask(userID, taskID int) error {
	if err := r.repo.DeleteTask(userID, taskID); err != nil {
		return fmt.Errorf("failed to delete task in repository: %w", err)
	}

	key := fmt.Sprintf("user:%d:tasks", userID)

	if err := r.rdb.HDel(ctx, key, fmt.Sprint(taskID)).Err(); err != nil {
		return fmt.Errorf("failed to delete task from cache: %w", err)
	}

	return nil
}
