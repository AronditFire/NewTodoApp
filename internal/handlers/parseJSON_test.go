package handlers

import (
	//"mime/multipart"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/service"
	mock_service "github.com/AronditFire/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_parseJsonFile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockParsingJSON, bindfile entity.BindFile)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		bindfile             entity.BindFile
		filename             string
		fileBody             string
		expectedCode         int
		expectedBodyResponse string
	}{
		{
			name:                 "bad request",
			mockBehavior:         func(s *mock_service.MockParsingJSON, bindfile entity.BindFile) {},
			bindfile:             entity.BindFile{},
			filename:             "",
			fileBody:             "",
			expectedCode:         400,
			expectedBodyResponse: `{"error":"could not bind uploaded file"}`,
		},
		{
			name: "bad request with error from service",
			mockBehavior: func(s *mock_service.MockParsingJSON, bindfile entity.BindFile) {
				s.EXPECT().ParseJSON(gomock.Any()).Return(errors.New("could not bind uploaded file"))
			},
			bindfile:             entity.BindFile{},
			filename:             "test.xml",
			fileBody:             "<root></root>",
			expectedCode:         500,
			expectedBodyResponse: `{"error":"could not bind uploaded file"}`,
		},
		{
			name: "successful parse",
			mockBehavior: func(s *mock_service.MockParsingJSON, bindfile entity.BindFile) {
				s.EXPECT().ParseJSON(gomock.Any()).Return(nil)
			},
			bindfile:             entity.BindFile{},
			filename:             "test.json",
			fileBody:             "{\"key\":\"value\"}",
			expectedCode:         200,
			expectedBodyResponse: `{"message":"file successfully parsed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParsingJSON := mock_service.NewMockParsingJSON(ctrl)
			tt.mockBehavior(mockParsingJSON, tt.bindfile)

			h := &Handler{
				services: &service.Service{
					ParsingJSON: mockParsingJSON,
				},
			}

			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			if tt.filename != "" {
				part, err := writer.CreateFormFile("file", tt.filename)
				assert.NoError(t, err)
				_, err = part.Write([]byte(tt.fileBody))
				assert.NoError(t, err)
			}

			writer.Close()

			r := gin.New()
			r.POST("/admin/upload-file", h.parseJsonFile)

			w := httptest.NewRecorder()
			// Создаём запрос с корректным Content-Type
			req := httptest.NewRequest("POST", "/admin/upload-file", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			r.ServeHTTP(w, req)

			// Проверяем код ответа и тело
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}

func TestHandler_getJsonFiles(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockParsingJSON)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedBodyResponse string
	}{
		{
			name: "successful get files",
			mockBehavior: func(s *mock_service.MockParsingJSON) {
				s.EXPECT().GetJsonTable().Return([]map[string]any{
					{"key1": "value1", "key2": "value2"},
				}, nil)
			},
			expectedCode:         200,
			expectedBodyResponse: `{"data":[{"key1":"value1","key2":"value2"}]}`,
		},
		{
			name: "Service error",
			mockBehavior: func(s *mock_service.MockParsingJSON) {
				s.EXPECT().GetJsonTable().Return(nil, errors.New("service error"))
			},
			expectedCode:         500,
			expectedBodyResponse: `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParsingJSON := mock_service.NewMockParsingJSON(ctrl)
			tt.mockBehavior(mockParsingJSON)

			h := &Handler{
				services: &service.Service{
					ParsingJSON: mockParsingJSON,
				},
			}

			r := gin.New()
			r.GET("/admin/get-files", h.getJsonFiles)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin/get-files", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBodyResponse, w.Body.String())
		})
	}
}
