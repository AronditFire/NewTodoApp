package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/AronditFire/todo-app/internal/service"
	mock_service "github.com/AronditFire/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentify(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehaivor         mockBehavior
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:        "Success",
			headerName:  authorizationHeader,
			headerValue: "Bearer token",
			token:       "token",
			mockBehaivor: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(&service.TokenClaims{UserID: 1, IsAdmin: false}, nil)
			},
			expectedStatusCode:   200,
			expectedBodyResponse: "1, false",
		},
		{
			name:                 "Wrong Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehaivor:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid header value",
			headerName:           authorizationHeader,
			headerValue:          "Bearr token",
			token:                "token",
			mockBehaivor:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"invalid header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           authorizationHeader,
			headerValue:          "Bearer ",
			token:                "token",
			mockBehaivor:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"empty token"}`,
		},
		{
			name:                 "Empty token",
			headerName:           authorizationHeader,
			headerValue:          "Bearer ",
			token:                "token",
			mockBehaivor:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"empty token"}`,
		},
		{
			name:        "Invalid token/Server Error",
			headerName:  authorizationHeader,
			headerValue: "Bearer token",
			token:       "token",
			mockBehaivor: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(nil, errors.New(`"error":"could not parse access token"`))
			},
			expectedStatusCode:   500,
			expectedBodyResponse: `{"error":"could not parse access token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			tt.mockBehaivor(repo, tt.token)

			service := &service.Service{Authorization: /*repo*/ nil}
			handler := Handler{service}

			r := gin.New()
			r.GET("/identify", handler.userIdentify, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				is_admin, _ := c.Get(adminCtx)
				c.String(200, "%d, %v", id, is_admin)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identify", nil)
			req.Header.Set(tt.headerName, tt.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedBodyResponse)
		})
	}
}

func TestHandler_adminIdentify(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	tests := []struct {
		name                 string
		isAdmin              bool
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			name:                 "Success",
			isAdmin:              true,
			expectedStatusCode:   200,
			expectedBodyResponse: "admin access granted",
		},
		{
			name:                 "Not Admin",
			isAdmin:              false,
			expectedStatusCode:   401,
			expectedBodyResponse: `{"error":"admin only"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{}

			r := gin.New()
			r.GET("/admin", func(c *gin.Context) {
				c.Set(adminCtx, tt.isAdmin)
				handler.adminIdentify(c)
			}, func(c *gin.Context) {
				c.String(200, "admin access granted")
			})

			// Arrange
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin", nil)

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedBodyResponse)
		})
	}
}

func TestHandler_getUserId(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	var getContext = func(id any) *gin.Context {
		c := &gin.Context{}
		c.Set(userCtx, id)
		return c
	}

	tests := []struct {
		name string
		ctx  *gin.Context
		id   int
		err  error
	}{
		{
			name: "Success",
			ctx:  getContext(1),
			id:   1,
			err:  nil,
		},
		{
			name: "No User ID",
			ctx:  &gin.Context{},
			id:   0,
			err:  errors.New("user id not found"),
		},
		{
			name: "Invalid user id",
			ctx:  getContext("invalid"),
			id:   0,
			err:  errors.New("user id is of invalid type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := getUserId(tt.ctx)

			assert.Equal(t, id, tt.id)

			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err, tt.err)
			}
		})
	}
}
