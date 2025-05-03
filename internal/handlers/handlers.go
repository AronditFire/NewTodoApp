package handlers

import (
	"github.com/AronditFire/todo-app/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHander(sv *service.Service) *Handler {
	return &Handler{services: sv}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.registerUser)
		auth.POST("/sign-in", h.loginUser)
		auth.POST("/refresh", h.refreshTokens)
	}

	api := router.Group("/api", h.userIdentify)
	{
		api.GET("/", h.getAllTasks)      // get all tasks
		api.GET("/:id", h.getTaskByID)   // get 1 task
		api.POST("/", h.createTask)      // create task
		api.PUT("/:id", h.updateTask)    // update task
		api.DELETE("/:id", h.deleteTask) // delete task

		admin := api.Group("/admin", h.adminIdentify)
		{
			admin.POST("/upload-file", h.parseJsonFile)
		}
	}

	return router
}
