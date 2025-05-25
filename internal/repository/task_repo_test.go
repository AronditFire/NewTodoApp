package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		inputTask   entity.Task
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"tasks\"").WithArgs("Test Task", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUserID: 1,
			inputTask: entity.Task{
				Description: "Test Task",
				UserID:      1,
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			inputTask:   entity.Task{},
			wantErr:     true,
		},
		{
			name: "InsertError",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"tasks\"").WithArgs("Test Task", 1).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTask: entity.Task{
				Description: "Test Task",
				UserID:      1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := r.CreateTask(tt.inputUserID, tt.inputTask)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_CreateTask_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	// Зарегистрируем callback, который сработает прямо перед тем, как GORM пошлёт SQL
	gormDB.Callback().Create().Before("gorm:create").
		Register("panic_before_create", func(db *gorm.DB) {
			panic("boom")
		})

	// Мокируем транзакцию и ожидаем откат
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Вызываем — внутри должен произойти panic, но благодаря defer+recover
	// функция вернёт nil и не «упадёт» в тесте.
	assert.NotPanics(t, func() {
		err := r.CreateTask(1, entity.Task{
			Description: "Test Task",
			UserID:      1,
		})
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGellAllTasks(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		wantTasks   []entity.Task
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "description", "user_id"}).AddRow(1, "Test Task 1", 1).
					AddRow(2, "Test Task 2", 1)
				mock.ExpectBegin()
				// GORM при Find генерирует примерно такой запрос:
				// SELECT * FROM "tasks" WHERE user_id = $1 ORDER BY "tasks"."id"
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1`),
				).WithArgs(1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUserID: 1,
			wantTasks: []entity.Task{
				{
					ID:          1,
					Description: "Test Task 1",
					UserID:      1,
				},
				{
					ID:          2,
					Description: "Test Task 2",
					UserID:      1,
				},
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			wantTasks:   nil,
			wantErr:     true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1`)).
					WithArgs(1).WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			wantTasks:   nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAllTask(tt.inputUserID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantTasks, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTasks, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_GetAllTasks_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	// Зарегистрируем callback, который сработает прямо перед тем, как GORM пошлёт SQL
	gormDB.Callback().Query().Before("gorm:query").
		Register("panic_before_query", func(db *gorm.DB) {
			panic("boom")
		})

	// Мокируем транзакцию и ожидаем откат
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Вызываем — внутри должен произойти panic, но благодаря defer+recover
	// функция вернёт nil и не «упадёт» в тесте.
	assert.NotPanics(t, func() {
		got, err := r.GetAllTask(1)
		assert.Equal(t, []entity.Task(nil), got)
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTaskByID(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		inputTaskID int
		wantTask    entity.Task
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.
					NewRows([]string{"id", "description", "user_id"}).
					AddRow(1, "test", 1)
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantTask: entity.Task{
				ID:          1,
				Description: "test",
				UserID:      1,
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantTask:    entity.Task{},
			wantErr:     true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantTask:    entity.Task{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetTaskByID(tt.inputUserID, tt.inputTaskID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, entity.Task{}, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTask, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestTask_GetTaskByID_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	// Зарегистрируем callback, который сработает прямо перед тем, как GORM пошлёт SQL
	gormDB.Callback().Query().Before("gorm:query").
		Register("panic_before_query", func(db *gorm.DB) {
			panic("boom")
		})

	// Мокируем транзакцию и ожидаем откат
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Вызываем — внутри должен произойти panic, но благодаря defer+recover
	// функция вернёт nil и не «упадёт» в тесте.
	assert.NotPanics(t, func() {
		got, err := r.GetTaskByID(1, 1)
		assert.Equal(t, entity.Task{}, got)
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)
	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		inputTaskID int
		inputDesc   string
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				// Expect single transaction lifecycle
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "description", "user_id"}).AddRow(1, "Test Task", 1)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE "tasks" SET "description"=$1,"user_id"=$2 WHERE "id" = $3`,
				)).
					WithArgs("Updated Task", 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			inputUserID: 1,
			inputTaskID: 1,
			inputDesc:   "Updated Task",
			wantErr:     false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			inputTaskID: 1,
			inputDesc:   "Updated Task",
			wantErr:     true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTaskID: 1,
			inputDesc:   "Updated Task",
			wantErr:     true,
		},
		{
			name: "Update Error",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "description", "user_id"}).AddRow(1, "Test Task", 1)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE "tasks" SET "description"=$1,"user_id"=$2 WHERE "id" = $3`,
				)).
					WithArgs("Updated Task", 1, 1).
					WillReturnError(errors.New("Update Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTaskID: 1,
			inputDesc:   "Updated Task",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateTask(tt.inputUserID, tt.inputTaskID, tt.inputDesc)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_UpdateTask_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	// Зарегистрируем callback, который сработает прямо перед тем, как GORM пошлёт SQL
	gormDB.Callback().Query().Before("gorm:query").
		Register("panic_before_query", func(db *gorm.DB) {
			panic("boom")
		})

	// Мокируем транзакцию и ожидаем откат
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Вызываем — внутри должен произойти panic, но благодаря defer+recover
	// функция вернёт nil и не «упадёт» в тесте.
	assert.NotPanics(t, func() {
		err := r.UpdateTask(1, 1, "Updated Task")
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTask(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)
	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		inputTaskID int
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				// Expect single transaction lifecycle
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "description", "user_id"}).AddRow(1, "Test Task", 1)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM "tasks" WHERE "tasks"."id" = $1`,
				)).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantErr:     false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantErr:     true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantErr:     true,
		},
		{
			name: "Update Error",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "description", "user_id"}).AddRow(1, "Test Task", 1)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE user_id = $1 AND id = $2 ORDER BY "tasks"."id" LIMIT $3`,
				)).
					WithArgs(1, 1, 1).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM "tasks" WHERE "tasks"."id" = $1`,
				)).
					WithArgs(1).
					WillReturnError(errors.New("Delete Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			inputTaskID: 1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteTask(tt.inputUserID, tt.inputTaskID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_DeleteTask_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewTaskRepo(gormDB)

	// Зарегистрируем callback, который сработает прямо перед тем, как GORM пошлёт SQL
	gormDB.Callback().Query().Before("gorm:query").
		Register("panic_before_query", func(db *gorm.DB) {
			panic("boom")
		})

	// Мокируем транзакцию и ожидаем откат
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Вызываем — внутри должен произойти panic, но благодаря defer+recover
	// функция вернёт nil и не «упадёт» в тесте.
	assert.NotPanics(t, func() {
		err := r.DeleteTask(1, 1)
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}
