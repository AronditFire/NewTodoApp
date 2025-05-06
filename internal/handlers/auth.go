package handlers

import (
	"net/http"

	"github.com/AronditFire/todo-app/entity"
	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body entity.UserRegisterRequest true "account info"
// @Success 201 {string} string "message"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
// @Failure default {string} string "error"
// @Router /auth/sign-up [post]
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

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body entity.UserAuthRequest true "credentials"
// @Success 200 {object} object "Успешный ответ"
// @Response 200 {object} map[string]any "Успешный ответ"
// @Failure 400,401 {string} string "error"
// @Failure 500 {string} string "error"
// @Failure default {string} string "error"
// @Router /auth/sign-in [post]
func (h *Handler) loginUser(c *gin.Context) {
	var req entity.UserAuthRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid request body",
		})
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.LoginUser(req) // TODO: ADD REFRESH
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) refreshTokens(c *gin.Context) {
	var req entity.RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := h.services.Authorization.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	newAccess, newRefresh, err := h.services.Authorization.RenewTokens(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  newAccess,
		"refreshToken": newRefresh,
	})
}
