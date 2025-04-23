package handlers

import (
	"net/http"

	"github.com/AronditFire/todo-app/entity"
	"github.com/gin-gonic/gin"
)

func (h *Handler) registerUser(c *gin.Context) {
	var req entity.UserRegisterRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.services.Authorization.CreateUser(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})
}

func (h *Handler) loginUser(c *gin.Context) {
	var req entity.UserAuthRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.services.Authorization.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not to login user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
