package service

import (
	"testing"

	mock_repository "github.com/AronditFire/todo-app/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetJsonTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockParsingJSON(ctrl)
	service := NewParseService(mockRepo)

	expectedData := []map[string]any{
		{"key1": "value1", "key2": "value2"},
	}

	mockRepo.EXPECT().GetJsonTable().Return(expectedData, nil)

	data, err := service.GetJsonTable()
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
}
