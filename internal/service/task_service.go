package service

import (
	"errors"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/cache"
)

type TaskService struct {
	crepo cache.TaskList
}

func NewTaskService(crepo cache.TaskList) *TaskService {
	return &TaskService{crepo: crepo}
}

func (s *TaskService) CreateTask(userID int, task entity.Task) (int, error) {
	if len(task.Description) > 0 && len(task.Description) < 1000 {
		return s.crepo.CreateTask(userID, task)
	} else {
		return 0, errors.New("Invalid description length!")
	}
}

func (s *TaskService) GetAllTask(userID int) ([]entity.Task, error) {
	return s.crepo.GetAllTask(userID)
}

func (s *TaskService) GetTaskByID(userID, id int) (entity.Task, error) {
	if id > 0 {
		return s.crepo.GetTaskByID(userID, id)
	} else {
		return entity.Task{}, errors.New("Invalid id while trying to get task by ID")
	}
}

func (s *TaskService) UpdateTask(userID, taskId int, desc string) error {
	if (len(desc) > 0) && (len(desc) < 1000) {
		return s.crepo.UpdateTask(userID, taskId, desc)
	} else {
		return errors.New("Invalid description length to update!")
	}

}

func (s *TaskService) DeleteTask(userID, taskID int) error {
	if taskID > 0 {
		return s.crepo.DeleteTask(userID, taskID)
	} else {
		return errors.New("Invalid id while trying to delete task")
	}
}
