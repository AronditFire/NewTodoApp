package service

import (
	"errors"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	mock_repository "github.com/AronditFire/todo-app/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Override generatePasswordHash to return a static value.
	generatePasswordHash = func(_ []byte, _ int) ([]byte, error) {
		return []byte("static_hashed"), nil
	}

	type mockBehavior func(mockAuthRepo *mock_repository.MockAuthorization, userReg entity.UserRegisterRequest)
	tests := []struct {
		name          string
		mockBehavior  mockBehavior
		inputUser     entity.UserRegisterRequest
		expectedError string
	}{
		{
			name: "Success",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, userReg entity.UserRegisterRequest) {
				mockAuthRepo.EXPECT().CreateUser(userReg).Return(nil)
			},
			inputUser: entity.UserRegisterRequest{
				Username: "testuser",
				Password: "static_hashed",
			},
			expectedError: "",
		},
		{
			name:         "Wrong Username Length",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, userReg entity.UserRegisterRequest) {},
			inputUser: entity.UserRegisterRequest{
				Username: "t",
				Password: "static_hashed",
			},
			expectedError: "bad username length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthRepo := mock_repository.NewMockAuthorization(ctrl)
			tt.mockBehavior(mockAuthRepo, tt.inputUser)

			authService := NewAuthService(mockAuthRepo)

			err := authService.CreateUser(tt.inputUser)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err.Error())
			}
		})
	}
}
func TestCreateUser_GeneratePasswordHashError(t *testing.T) {
	generatePasswordHash = func(_ []byte, _ int) ([]byte, error) {
		return nil, errors.New("hash error")
	}
	t.Run("Hash Pass Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthRepo := mock_repository.NewMockAuthorization(ctrl)
		authService := NewAuthService(mockAuthRepo)
		err := authService.CreateUser(entity.UserRegisterRequest{
			Username: "testuser",
			Password: "testpassword",
		})
		assert.Equal(t, "Could not hash user password", err.Error())
	})
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_repository.NewMockAuthorization(ctrl)
	authService := NewAuthService(mockAuthRepo)

	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword"}
	mockAuthRepo.EXPECT().GetUser("testuser").Return(user, nil)

	result, err := authService.GetUser("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_repository.NewMockAuthorization(ctrl)
	authService := NewAuthService(mockAuthRepo)

	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword"}
	mockAuthRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := authService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestLoginUser(t *testing.T) {
	origCompareHashAndPassword := CompareHashAndPassword
	defer func() {
		CompareHashAndPassword = origCompareHashAndPassword
	}()

	type mockBehavior func(mockAuthRepo *mock_repository.MockAuthorization, username string)
	tests := []struct {
		name            string
		mockBehavior    mockBehavior
		inputLogin      entity.UserAuthRequest
		overrideCompare func(hashedPwd, pwd []byte) error
		expectedError   string
	}{
		{
			name: "Success",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, username string) {
				mockAuthRepo.EXPECT().GetUser(username).Return(entity.User{
					ID:       1,
					Username: "testuser",
					Password: "static_hashed",
					IsAdmin:  false,
				}, nil)
			},
			inputLogin: entity.UserAuthRequest{
				Username: "testuser",
				Password: "testpassword",
			},
			overrideCompare: func(hashedPwd, pwd []byte) error {
				if string(hashedPwd) == "static_hashed" && string(pwd) == "testpassword" {
					return nil
				}
				return errors.New("password mismatch")
			},
			expectedError: "",
		},
		{
			name:         "Bad Username Length",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, username string) {},
			inputLogin: entity.UserAuthRequest{
				Username: "t",
				Password: "testpassword",
			},
			expectedError: "bad username length",
		},
		{
			name: "User Not Found",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, username string) {
				mockAuthRepo.EXPECT().GetUser(username).Return(entity.User{}, errors.New("user not found"))
			},
			inputLogin: entity.UserAuthRequest{
				Username: "nonexistent",
				Password: "irrelevant",
			},
			expectedError: "user not found",
		},
		{
			name: "Incorrect Password",
			mockBehavior: func(mockAuthRepo *mock_repository.MockAuthorization, username string) {
				mockAuthRepo.EXPECT().GetUser(username).Return(entity.User{
					ID:       1,
					Username: "testuser",
					Password: "static_hashed",
					IsAdmin:  false,
				}, nil)
			},
			inputLogin: entity.UserAuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			overrideCompare: func(hashedPwd, pwd []byte) error {

				return errors.New("mismatch")
			},
			expectedError: "Incorrect password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CompareHashAndPassword = tt.overrideCompare

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockAuthorization(ctrl)
			tt.mockBehavior(mockRepo, tt.inputLogin.Username)

			authService := NewAuthService(mockRepo)
			accessToken, refreshToken, err := authService.LoginUser(tt.inputLogin)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
			} else {
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
			}
		})
	}
}
