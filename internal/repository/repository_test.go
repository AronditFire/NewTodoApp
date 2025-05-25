package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	sqlDB, gormDB, _ := DbMock(t)
	defer sqlDB.Close()

	svc := NewRepository(gormDB)

	assert.NotNil(t, svc)

	assert.NotNil(t, svc.TaskList)
	assert.NotNil(t, svc.Authorization)
	assert.NotNil(t, svc.ParsingJSON)

}
