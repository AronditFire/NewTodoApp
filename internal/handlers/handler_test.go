package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRoutes(t *testing.T) {
	// Создаём «пустой» handler, если у него нет обязательных зависимостей,
	// или можете передать заглушки так же, как в других тестах.
	h := &Handler{}
	engine := h.InitRoutes()
	routes := engine.Routes()

	// Формируем карту из строк "<METHOD> <PATH>" для быстрого поиска
	have := make(map[string]bool, len(routes))
	for _, r := range routes {
		have[r.Method+" "+r.Path] = true
	}

	expected := []string{
		"GET /swagger/*any",
		"POST /auth/sign-up",
		"POST /auth/sign-in",
		"POST /auth/refresh",
		"GET /api/",
		"GET /api/:id",
		"POST /api/",
		"PUT /api/:id",
		"DELETE /api/:id",
		"POST /api/admin/upload-file",
		"GET /api/admin/get-files",
	}

	for _, exp := range expected {
		assert.True(t, have[exp], "route %q should be registered", exp)
	}
}
