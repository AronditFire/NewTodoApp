package service

import (
	"testing"

	"github.com/AronditFire/todo-app/internal/repository"
	mock_repository "github.com/AronditFire/todo-app/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskList := mock_repository.NewMockTaskList(ctrl)
	mockAuthorization := mock_repository.NewMockAuthorization(ctrl)
	mockParsingJSON := mock_repository.NewMockParsingJSON(ctrl)

	repo := &repository.Repository{
		TaskList:      mockTaskList,
		Authorization: mockAuthorization,
		ParsingJSON:   mockParsingJSON,
	}

	svc := NewService(repo)
	assert.NotNil(t, svc)

	assert.NotNil(t, svc.TaskList)
	assert.NotNil(t, svc.Authorization)
	assert.NotNil(t, svc.ParsingJSON)

}
