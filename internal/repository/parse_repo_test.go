package repository

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferSQLType(t *testing.T) {
	tests := []struct {
		value    any
		expected string
	}{
		{value: "string", expected: "TEXT"},
		{value: 123.456, expected: "DOUBLE PRECISION"},
		{value: true, expected: "BOOLEAN"},
		{value: nil, expected: "JSONB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := inferSQLType(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseJSON_AnonymousCases(t *testing.T) {
	tests := []struct {
		name        string
		fileData    string
		mockSetup   func(sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name:     "Success",
			fileData: `{"new_column":"some value"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					`ALTER TABLE "files" ADD COLUMN "new_column" TEXT;`,
				)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO "files"`,
				)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: "",
		},
		{
			name:        "Invalid JSON",
			fileData:    `invalid json`,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectedErr: "Could not decode json",
		},
		{
			name:     "Fail begin tx",
			fileData: `{"new_column":"some value"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("error"))
			},
			expectedErr: "error",
		},
		{
			name:     "Alter table error",
			fileData: `{"new_column":"some value"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					`ALTER TABLE "files" ADD COLUMN "new_column" TEXT;`,
				)).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
			expectedErr: "error",
		},
		{
			name:     "Insert error",
			fileData: `{"new_column":"some value"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					`ALTER TABLE "files" ADD COLUMN "new_column" TEXT;`,
				)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO "files"`,
				)).WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
			expectedErr: "Could not create rows with data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB, gormDB, mock := DbMock(t)
			defer sqlDB.Close()
			tc.mockSetup(mock)

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			part, err := writer.CreateFormFile("file", "dummy.json")
			require.NoError(t, err)
			_, err = part.Write([]byte(tc.fileData))
			require.NoError(t, err)
			require.NoError(t, writer.Close())

			req := httptest.NewRequest("POST", "/", &buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			require.NoError(t, req.ParseMultipartForm(int64(buf.Len())))

			_, fh, err := req.FormFile("file")
			require.NoError(t, err)

			bindFile := entity.BindFile{File: fh}
			repo := NewParseRepo(gormDB)

			err = repo.ParseJSON(bindFile)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetJsonTable(t *testing.T) {
	sqlDB, gormDB, mock := DbMock(t)
	defer sqlDB.Close()

	tests := []struct {
		name         string
		mockSetup    func(sqlmock.Sqlmock)
		expectedBody []map[string]any
		wantErr      bool
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "files"`,
				)).WillReturnRows(sqlmock.NewRows([]string{"new_column"}).AddRow("some value"))
				mock.ExpectCommit()
			},
			expectedBody: []map[string]any{
				{"new_column": "some value"},
			},
			wantErr: false,
		},
		{
			name: "Fail begin tx",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "Could Not find rows",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "files"`,
				)).WillReturnError(errors.New("Could not write data in FilesData"))
				mock.ExpectRollback()
			},
			expectedBody: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			repo := NewParseRepo(gormDB)
			body, err := repo.GetJsonTable()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, body)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
