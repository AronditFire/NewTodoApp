package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/service"
	mock_service "github.com/AronditFire/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_getAllTask(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockTaskList, userID int)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		userID               int
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name: "Success",
			mockBehavior: func(s *mock_service.MockTaskList, userID int) {
				s.EXPECT().GetAllTask(userID).Return([]entity.Task{
					{
						ID:          1,
						Description: "test",
						UserID:      1,
					},
				}, nil)
			},
			userID:               1,
			expectedStatusCode:   200,
			expectedBodyResponse: `{"data":[{"id":1,"description":"test","UserID":1}]}`,
		},
		{
			name: "Failed to Get Tasks",
			mockBehavior: func(s *mock_service.MockTaskList, userID int) {
				s.EXPECT().GetAllTask(userID).Return(nil, errors.New(`{"error":"Could not get tasks for this user"}`))
			},
			userID:               1,
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not get tasks for this user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTaskList(c)
			tt.mockBehavior(repo, tt.userID)

			service := &service.Service{TaskList: repo}
			handler := Handler{service}

			r := gin.New()
			r.GET("/api/", func(c *gin.Context) {
				c.Set(userCtx, tt.userID)
				handler.getAllTasks(c)
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/", nil)

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}

func TestHandler_getTaskByID(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockTaskList, userID, id int)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		userID               int
		requestPath          string
		taskID               int
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name: "Success",
			mockBehavior: func(s *mock_service.MockTaskList, userID, id int) {
				s.EXPECT().GetTaskByID(userID, id).Return(entity.Task{
					ID:          1,
					Description: "test",
					UserID:      1,
				}, nil)
			},
			userID:               1,
			requestPath:          "/api/1",
			taskID:               1,
			expectedStatusCode:   200,
			expectedBodyResponse: `{"id":1,"description":"test","UserID":1}`,
		},
		{
			name:                 "Bad Param ID",
			mockBehavior:         func(s *mock_service.MockTaskList, userID, id int) {},
			userID:               1,
			requestPath:          "/api/abc",
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid task id"}`,
		},
		{
			name: "Server Error - Task not found",
			mockBehavior: func(s *mock_service.MockTaskList, userID, id int) {
				s.EXPECT().GetTaskByID(userID, id).Return(entity.Task{}, errors.New(`"error":"Could not get task by ID"`))
			},
			userID:               1,
			requestPath:          "/api/1",
			taskID:               1,
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not get task by ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTaskList(c)
			tt.mockBehavior(repo, tt.userID, tt.taskID)

			service := &service.Service{TaskList: repo}
			handler := Handler{service}

			r := gin.New()
			r.GET("/api/:id", func(c *gin.Context) {
				c.Set(userCtx, tt.userID)
				handler.getTaskByID(c)
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}

func TestHandler_createTask(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockTaskList, userID int, task entity.Task)

	tests := []struct {
		name                 string
		inputBody            string
		inputTask            entity.Task
		mockBehavior         mockBehavior
		userID               int
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:      "Success",
			inputBody: `{"id":1,"description":"test","UserID":1}`,
			inputTask: entity.Task{
				ID:          1,
				Description: "test",
				UserID:      1,
			},
			mockBehavior: func(s *mock_service.MockTaskList, userID int, task entity.Task) {
				s.EXPECT().CreateTask(userID, task).Return(nil)
			},
			userID:               1,
			expectedStatusCode:   201,
			expectedBodyResponse: `{"message":"created"}`,
		},
		{
			name:                 "Bad Request",
			inputBody:            `{"descripti":""}`, // changed to invalid JSON to trigger binding error
			inputTask:            entity.Task{},
			mockBehavior:         func(s *mock_service.MockTaskList, userID int, task entity.Task) {},
			userID:               1,
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"Could not unbind request while creating task"}`,
		},
		{
			name:      "Server Error - Create Error",
			inputBody: `{"id":1,"description":"test","UserID":1}`,
			inputTask: entity.Task{
				ID:          1,
				Description: "test",
				UserID:      1,
			},
			mockBehavior: func(s *mock_service.MockTaskList, userID int, task entity.Task) {
				s.EXPECT().CreateTask(userID, task).Return(errors.New(`"error":"Could not create task in database"`))
			},
			userID:               1,
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not create task in database"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTaskList(c)
			tt.mockBehavior(repo, tt.userID, tt.inputTask)

			service := &service.Service{TaskList: repo}
			handler := &Handler{service}

			r := gin.New()
			r.POST("/api", func(c *gin.Context) {
				c.Set(userCtx, tt.userID)
				handler.createTask(c)
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBufferString(tt.inputBody))

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}

func TestHandler_updateTask(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockTaskList, userID, taskID int, desc string)

	tests := []struct {
		name                 string
		inputBody            string
		inputTask            entity.TaskRequest
		mockBehavior         mockBehavior
		taskID               int
		requestPath          string
		userID               int
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:      "Success",
			inputBody: `{"description":"newTest"}`,
			inputTask: entity.TaskRequest{
				Description: "newTest",
			},
			mockBehavior: func(s *mock_service.MockTaskList, userID, taskID int, desc string) {
				s.EXPECT().UpdateTask(userID, taskID, desc).Return(nil)
			},
			userID:               1,
			taskID:               1,
			requestPath:          "/api/1",
			expectedStatusCode:   200,
			expectedBodyResponse: `{"message":"updated"}`,
		},
		{
			name:      "Bad Param ID",
			inputBody: `{"description":"newTest"}`,
			inputTask: entity.TaskRequest{
				Description: "newTest",
			},
			mockBehavior:         func(s *mock_service.MockTaskList, userID, taskID int, desc string) {},
			userID:               1,
			requestPath:          "/api/abc",
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid task id"}`,
		},
		{
			name:                 "Invalid Body Request",
			inputBody:            `{"descriptio":"newTest"}`,
			inputTask:            entity.TaskRequest{},
			mockBehavior:         func(s *mock_service.MockTaskList, userID, taskID int, desc string) {},
			userID:               1,
			taskID:               1,
			requestPath:          "/api/1",
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"Could not unbind request while updating task"}`,
		},
		{
			name:      "Server Error",
			inputBody: `{"description":"newTest"}`,
			inputTask: entity.TaskRequest{
				Description: "newTest",
			},
			mockBehavior: func(s *mock_service.MockTaskList, userID, taskID int, desc string) {
				s.EXPECT().UpdateTask(userID, taskID, desc).Return(errors.New(`"error":"Could not update task in database"`))
			},
			userID:               1,
			taskID:               1,
			requestPath:          "/api/1",
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not update task in database"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTaskList(c)
			tt.mockBehavior(repo, tt.userID, tt.taskID, tt.inputTask.Description)

			service := &service.Service{TaskList: repo}
			handler := &Handler{service}

			r := gin.New()
			r.PUT("/api/:id", func(c *gin.Context) {
				c.Set(userCtx, tt.userID)
				handler.updateTask(c)
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, tt.requestPath, bytes.NewBufferString(tt.inputBody))

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}

func TestHandler_deleteTask(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockTaskList, userID, taskID int)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		taskID               int
		requestPath          string
		userID               int
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name: "Success",
			mockBehavior: func(s *mock_service.MockTaskList, userID, taskID int) {
				s.EXPECT().DeleteTask(userID, taskID).Return(nil)
			},
			userID:               1,
			taskID:               1,
			requestPath:          "/api/1",
			expectedStatusCode:   200,
			expectedBodyResponse: `{"message":"deleted"}`,
		},
		{
			name:                 "Bad Param ID",
			mockBehavior:         func(s *mock_service.MockTaskList, userID, taskID int) {},
			userID:               1,
			requestPath:          "/api/abc",
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid task id"}`,
		},
		{
			name: "Server Error",
			mockBehavior: func(s *mock_service.MockTaskList, userID, taskID int) {
				s.EXPECT().DeleteTask(userID, taskID).Return(errors.New(`"error":"Could not delete task in database"`))
			},
			userID:               1,
			taskID:               1,
			requestPath:          "/api/1",
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not delete task in database"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTaskList(c)
			tt.mockBehavior(repo, tt.userID, tt.taskID)

			service := &service.Service{TaskList: repo}
			handler := &Handler{service}

			r := gin.New()
			r.DELETE("/api/:id", func(c *gin.Context) {
				c.Set(userCtx, tt.userID)
				handler.deleteTask(c)
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, tt.requestPath, nil)

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}
