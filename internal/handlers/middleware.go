package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
	adminCtx            = "isAdmin"
)

func (h *Handler) userIdentify(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "empty auth header",
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid header",
		})
		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "empty token",
		})
		return
	}

	claims, err := h.services.Authorization.ParseAccessToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "could not parse access token",
		})
		return
	}

	c.Set(userCtx, claims.UserID)
	c.Set(adminCtx, claims.IsAdmin)
	c.Next()
}

func (h *Handler) adminIdentify(c *gin.Context) {
	isAdmin, _ := c.Get(adminCtx)
	if isAdmin == false {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "admin only",
		})
		return
	}
	c.Next()
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
