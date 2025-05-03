package handlers

import (
	"net/http"

	"github.com/AronditFire/todo-app/entity"
	"github.com/gin-gonic/gin"
)

func (h *Handler) parseJsonFile(c *gin.Context) {
	var bindfile entity.BindFile
	if err := c.ShouldBind(&bindfile); err != nil { // get json file
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "could not bind uploaded file",
		})
		return
	}

	if err := h.services.ParsingJSON.ParseJSON(bindfile); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file successfully parsed",
	})
}
