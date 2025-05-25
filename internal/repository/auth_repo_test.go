package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	return sqldb, gormdb, mock
}

func TestAuth_CreateUser(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

	tests := []struct {
		name         string
		mock         func()
		inputUserReg entity.UserRegisterRequest
		wantErr      bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"users\"").WithArgs("TestUser", "testpass", false).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUserReg: entity.UserRegisterRequest{
				Username: "TestUser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Error"))
			},
			inputUserReg: entity.UserRegisterRequest{
				Username: "TestUser",
				Password: "testpass",
			},
			wantErr: true,
		},
		{
			name: "InsertError",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"users\"").WithArgs("TestUser", "testpass", false).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
			inputUserReg: entity.UserRegisterRequest{
				Username: "TestUser",
				Password: "testpass",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateUser(tt.inputUserReg)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
func TestAuth_CreateUser_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

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
		err := r.CreateUser(entity.UserRegisterRequest{
			Username: "TestUser",
			Password: "testpass",
		})
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

	tests := []struct {
		name          string
		mock          func()
		inputUsername string
		wantUser      entity.User
		wantErr       bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.
					NewRows([]string{"id", "username", "password", "is_admin"}).
					AddRow(1, "testname", "testpass", false)

					// GORM при вызове First генерит что-то вроде:
					// SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
				)).
					WithArgs("testname", 1).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUsername: "testname",
			wantUser: entity.User{
				ID:       1,
				Username: "testname",
				Password: "testpass",
				IsAdmin:  false,
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUsername: "testname",
			wantUser:      entity.User{},
			wantErr:       true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
				)).
					WithArgs("testname", 1).
					WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUsername: "testname",
			wantUser:      entity.User{},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(tt.inputUsername)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantUser, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_GetUser_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

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
		got, err := r.GetUser("TestUser")
		assert.Equal(t, entity.User{}, got)
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

	tests := []struct {
		name        string
		mock        func()
		inputUserID int
		wantUser    entity.User
		wantErr     bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.
					NewRows([]string{"id", "username", "password", "is_admin"}).
					AddRow(1, "testname", "testpass", false)

					// GORM при вызове First генерит что-то вроде:
					// SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`,
				)).
					WithArgs(1, 1).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			inputUserID: 1,
			wantUser: entity.User{
				ID:       1,
				Username: "testname",
				Password: "testpass",
				IsAdmin:  false,
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("Begin Error"))
			},
			inputUserID: 1,
			wantUser:    entity.User{},
			wantErr:     true,
		},
		{
			name: "Select Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`,
				)).
					WithArgs(1, 1).
					WillReturnError(errors.New("Select Error"))
				mock.ExpectRollback()
			},
			inputUserID: 1,
			wantUser:    entity.User{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByID(tt.inputUserID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantUser, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_GetUserByID_PanicRecovery(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	r := NewAuthRepo(gormDB)

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
		got, err := r.GetUserByID(1)
		assert.Equal(t, entity.User{}, got)
		assert.NoError(t, err)
	})

	// Проверяем, что Rollback действительно вызывался
	assert.NoError(t, mock.ExpectationsWereMet())
}
