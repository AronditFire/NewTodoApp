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

func (s *TaskService) CreateTask(task entity.Task) error {
	if (len(task.Description) > 0) || (len(task.Description) < 1000) {
		return s.repo.CreateTask(task)
	} else {
		return errors.New("Invalid description length!")
	}
}

func (s *TaskService) GetAllTask() ([]entity.Task, error) {
	return s.repo.GetAllTask()
}

func (s *TaskService) GetTaskByID(id int) (entity.Task, error) {
	if id > 0 {
		return s.repo.GetTaskByID(id)
	} else {
		return entity.Task{}, errors.New("Invalid id while trying to get task by ID")
	}
}

func (s *TaskService) UpdateTask(id int, desc string) error {
	if (len(desc) > 0) || (len(desc) < 1000) {
		return s.repo.UpdateTask(id, desc)
	} else {
		return errors.New("Invalid description length to update!")
	}

}

func (s *TaskService) DeleteTask(id int) error {
	if id > 0 {
		return s.repo.DeleteTask(id)
	} else {
		return errors.New("Invalid id while trying to delete task")
	}
}
