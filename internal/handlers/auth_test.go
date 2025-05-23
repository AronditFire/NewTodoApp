package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/service"
	mock_service "github.com/AronditFire/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_registerUser(t *testing.T) {
	type mockBehaivor func(s *mock_service.MockAuthorization, user entity.UserRegisterRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            entity.UserRegisterRequest
		mockBehaivor         mockBehaivor
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:      "Successed Test",
			inputBody: `{"username": "test", "pass": "123456"}`,
			inputUser: entity.UserRegisterRequest{
				Username: "test",
				Password: "123456",
			},
			mockBehaivor: func(s *mock_service.MockAuthorization, user entity.UserRegisterRequest) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode:   201,
			expectedBodyResponse: `{"message":"user created"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            entity.UserRegisterRequest{},
			mockBehaivor:         func(s *mock_service.MockAuthorization, user entity.UserRegisterRequest) {},
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid request body"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"username": "test", "pass": "123456"}`,
			inputUser: entity.UserRegisterRequest{
				Username: "test",
				Password: "123456",
			},
			mockBehaivor: func(s *mock_service.MockAuthorization, user entity.UserRegisterRequest) {
				s.EXPECT().CreateUser(user).Return(errors.New(`"error": "Could not to create user"`))
			},
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not to create user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t) // library controller
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			tt.mockBehaivor(repo, tt.inputUser)

			service := &service.Service{Authorization: repo}
			handler := NewHander(service)

			// Arrange
			r := gin.New()
			r.POST("/auth/sign-up", handler.registerUser)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(tt.inputBody))

			// Assert
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedBodyResponse)
		})
	}
}
