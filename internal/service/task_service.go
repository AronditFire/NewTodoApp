package service

import (
	"errors"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
)

type TaskService struct {
	repo repository.TaskList
}

func NewTaskService(repo repository.TaskList) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(userID int, task entity.Task) error {
	if (len(task.Description) > 0) || (len(task.Description) < 1000) {
		return s.repo.CreateTask(userID, task)
	} else {
		return errors.New("Invalid description length!")
	}
}

func (s *TaskService) GetAllTask(userID int) ([]entity.Task, error) {
	return s.repo.GetAllTask(userID)
}

func (s *TaskService) GetTaskByID(userID, id int) (entity.Task, error) {
	if id > 0 {
		return s.repo.GetTaskByID(userID, id)
	} else {
		return entity.Task{}, errors.New("Invalid id while trying to get task by ID")
	}
}

func (s *TaskService) UpdateTask(userID, taskId int, desc string) error {
	if (len(desc) > 0) || (len(desc) < 1000) {
		return s.repo.UpdateTask(userID, taskId, desc)
	} else {
		return errors.New("Invalid description length to update!")
	}

}

func (s *TaskService) DeleteTask(userID, taskID int) error {
	if taskID > 0 {
		return s.repo.DeleteTask(userID, taskID)
	} else {
		return errors.New("Invalid id while trying to delete task")
	}
}
