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
	gin.SetMode(gin.ReleaseMode)
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

func TestHandler_loginUser(t *testing.T) {
	type mockBehaivor func(s *mock_service.MockAuthorization, user entity.UserAuthRequest)
	gin.SetMode(gin.ReleaseMode)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            entity.UserAuthRequest
		mockBehaivor         mockBehaivor
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:      "Success login",
			inputBody: `{"username": "test", "pass": "123456"}`,
			inputUser: entity.UserAuthRequest{
				Username: "test",
				Password: "123456",
			},
			mockBehaivor: func(s *mock_service.MockAuthorization, user entity.UserAuthRequest) {
				s.EXPECT().LoginUser(user).Return("abc", "abc", nil)
			},
			expectedStatusCode:   200,
			expectedBodyResponse: `{"accessToken":"abc","refreshToken":"abc"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            entity.UserAuthRequest{},
			mockBehaivor:         func(s *mock_service.MockAuthorization, user entity.UserAuthRequest) {},
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid request body"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"username": "test", "pass": "123456"}`,
			inputUser: entity.UserAuthRequest{
				Username: "test",
				Password: "123456",
			},
			mockBehaivor: func(s *mock_service.MockAuthorization, user entity.UserAuthRequest) {
				s.EXPECT().LoginUser(user).Return("", "", errors.New(`"error":"Could not to create token"`))
			},
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"Could not to create token"}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockAuthorization(ctrl)

			tt.mockBehaivor(repo, tt.inputUser)

			service := &service.Service{Authorization: repo}
			handler := &Handler{service}

			r := gin.New()
			r.POST("/auth/sign-in", handler.loginUser)

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(tt.inputBody))

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedBodyResponse)
		})
	}
}

func TestHandler_refreshTokens(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	type mockRefreshBehaivor func(s *mock_service.MockAuthorization, refreshToken string)
	type mockRenewBehaivor func(s *mock_service.MockAuthorization, id int)

	tests := []struct {
		name                 string
		inputBody            string
		inputToken           entity.RefreshRequest
		inputID              int
		mockRefreshBehaivor  mockRefreshBehaivor
		mockRenewBehaivor    mockRenewBehaivor
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:      "Success",
			inputBody: `{"refreshToken": "abc"}`,
			inputToken: entity.RefreshRequest{
				RefreshToken: "abc",
			},
			inputID: 1,
			mockRefreshBehaivor: func(s *mock_service.MockAuthorization, refreshToken string) {
				s.EXPECT().ParseRefreshToken(refreshToken).Return(1, nil)
			},
			mockRenewBehaivor: func(s *mock_service.MockAuthorization, id int) {
				s.EXPECT().RenewTokens(id).Return("abc", "abc", nil)
			},
			expectedStatusCode:   200,
			expectedBodyResponse: `{"accessToken":"abc","refreshToken":"abc"}`,
		},
		{
			name:                 "Invalid body",
			inputBody:            `{"refreshToken": 1}`,
			inputToken:           entity.RefreshRequest{},
			inputID:              1,
			mockRefreshBehaivor:  func(s *mock_service.MockAuthorization, refreshToken string) {},
			mockRenewBehaivor:    func(s *mock_service.MockAuthorization, id int) {},
			expectedStatusCode:   400,
			expectedBodyResponse: `{"error":"invalid input body"}`,
		},
		{
			name:      "Expired refresh token",
			inputBody: `{"refreshToken": "abc"}`,
			inputToken: entity.RefreshRequest{
				RefreshToken: "abc",
			},
			inputID: 1,
			mockRefreshBehaivor: func(s *mock_service.MockAuthorization, refreshToken string) {
				s.EXPECT().ParseRefreshToken(refreshToken).Return(0, errors.New(`"error":"invalid refresh token"`))
			},
			mockRenewBehaivor:    func(s *mock_service.MockAuthorization, id int) {},
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"invalid refresh token"}`,
		},
		{
			name:      "Server user ID Error",
			inputBody: `{"refreshToken": "abc"}`,
			inputToken: entity.RefreshRequest{
				RefreshToken: "abc",
			},
			inputID: 1,
			mockRefreshBehaivor: func(s *mock_service.MockAuthorization, refreshToken string) {
				s.EXPECT().ParseRefreshToken(refreshToken).Return(1, nil)
			},
			mockRenewBehaivor: func(s *mock_service.MockAuthorization, id int) {
				s.EXPECT().RenewTokens(id).Return("", "", errors.New(`"error":"could not generate tokens"`))
			},
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"could not generate tokens"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			tt.mockRefreshBehaivor(repo, tt.inputToken.RefreshToken)
			tt.mockRenewBehaivor(repo, tt.inputID)

			service := &service.Service{Authorization: repo}
			handler := &Handler{service}

			r := gin.New()
			r.POST("/auth/refresh", handler.refreshTokens)

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(tt.inputBody))

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedBodyResponse)
		})
	}
}
