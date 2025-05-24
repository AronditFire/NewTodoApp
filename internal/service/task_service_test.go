package service

import (
	"testing"

	"github.com/AronditFire/todo-app/entity"
	mock_repository "github.com/AronditFire/todo-app/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockTaskList, userID int, task entity.Task)

	tests := []struct {
		name          string
		userID        int
		task          entity.Task
		mockBehavior  mockBehavior
		expectedError string
	}{
		{
			name:   "Success",
			userID: 1,
			task:   entity.Task{Description: "Test"},
			mockBehavior: func(r *mock_repository.MockTaskList, userID int, task entity.Task) {
				r.EXPECT().CreateTask(userID, gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Invalid Description Length",
			userID:        1,
			task:          entity.Task{Description: ""}, // Empty description
			mockBehavior:  func(r *mock_repository.MockTaskList, userID int, task entity.Task) {},
			expectedError: "Invalid description length!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockTaskList(ctrl)
			tt.mockBehavior(mockRepo, tt.userID, tt.task)

			service := NewTaskService(mockRepo)

			err := service.CreateTask(tt.userID, tt.task)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllTask(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockTaskList, userID int)

	tests := []struct {
		name          string
		userID        int
		mockBehavior  mockBehavior
		expectedTasks []entity.Task
		expectedError string
	}{
		{
			name:   "Success",
			userID: 1,
			mockBehavior: func(r *mock_repository.MockTaskList, userID int) {
				r.EXPECT().GetAllTask(userID).Return([]entity.Task{
					{ID: 1, Description: "Test Task"},
				}, nil)
			},
			expectedTasks: []entity.Task{
				{ID: 1, Description: "Test Task"},
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockTaskList(ctrl)
			tt.mockBehavior(mockRepo, tt.userID)

			service := NewTaskService(mockRepo)

			gotTasks, err := service.GetAllTask(tt.userID)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTasks, gotTasks)
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockTaskList, userID int, id int)

	tests := []struct {
		name          string
		userID        int
		taskID        int
		mockBehavior  mockBehavior
		expectedTask  entity.Task
		expectedError string
	}{
		{
			name:   "Success",
			userID: 1,
			taskID: 1,
			mockBehavior: func(r *mock_repository.MockTaskList, userID int, id int) {
				r.EXPECT().GetTaskByID(userID, id).Return(
					entity.Task{ID: 1, Description: "Test"}, nil,
				)
			},
			expectedTask:  entity.Task{ID: 1, Description: "Test"},
			expectedError: "",
		},
		{
			name:          "Bad task ID",
			userID:        1,
			taskID:        0,
			mockBehavior:  func(r *mock_repository.MockTaskList, userID int, id int) {},
			expectedTask:  entity.Task{},
			expectedError: "Invalid id while trying to get task by ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockTaskList(ctrl)
			tt.mockBehavior(mockRepo, tt.userID, tt.taskID)

			service := NewTaskService(mockRepo)

			gotTask, err := service.GetTaskByID(tt.userID, tt.taskID)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask, gotTask)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockTaskList, userID, taskID int, desc string)

	tests := []struct {
		name          string
		userID        int
		taskID        int
		desc          string
		mockBehavior  mockBehavior
		expectedError string
	}{
		{
			name:   "Success",
			userID: 1,
			taskID: 1,
			desc:   "Updated Task Description",
			mockBehavior: func(r *mock_repository.MockTaskList, userID, taskID int, desc string) {
				r.EXPECT().UpdateTask(userID, taskID, desc).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Invalid Description Length",
			userID:        1,
			taskID:        1,
			desc:          "", // Empty description
			mockBehavior:  func(r *mock_repository.MockTaskList, userID, taskID int, desc string) {},
			expectedError: "Invalid description length to update!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockTaskList(ctrl)
			tt.mockBehavior(mockRepo, tt.userID, tt.taskID, tt.desc)

			service := NewTaskService(mockRepo)

			err := service.UpdateTask(tt.userID, tt.taskID, tt.desc)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestDeleteTask(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockTaskList, userID, taskID int)

	tests := []struct {
		name          string
		userID        int
		taskID        int
		mockBehavior  mockBehavior
		expectedError string
	}{
		{
			name:   "Success",
			userID: 1,
			taskID: 1,
			mockBehavior: func(r *mock_repository.MockTaskList, userID, taskID int) {
				r.EXPECT().DeleteTask(userID, taskID).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Bad task ID",
			userID:        1,
			taskID:        0,
			mockBehavior:  func(r *mock_repository.MockTaskList, userID int, id int) {},
			expectedError: "Invalid id while trying to delete task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockTaskList(ctrl)
			tt.mockBehavior(mockRepo, tt.userID, tt.taskID)

			service := NewTaskService(mockRepo)

			err := service.DeleteTask(tt.userID, tt.taskID)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
