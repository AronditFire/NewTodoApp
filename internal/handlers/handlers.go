package handlers

import (
	"net/http"

	_ "github.com/AronditFire/todo-app/docs"
	"github.com/AronditFire/todo-app/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHander(sv *service.Service) *Handler {
	return &Handler{services: sv}
}

func (h *Handler) InitRoutes(healthFunc http.HandlerFunc) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", gin.WrapF(healthFunc))

	oauth := router.Group("/oauth2")
	{
		oauth.GET("/login", h.authLogin)
		oauth.GET("/callback", h.callbackGoogle)
	}

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
			admin.GET("/get-files", h.getJsonFiles)
		}
	}

	return router
}
