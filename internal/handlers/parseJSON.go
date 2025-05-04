package handlers

import (
	"net/http"

	"github.com/AronditFire/todo-app/entity"
	"github.com/gin-gonic/gin"
)

type GetAllJSONResponse struct {
	Data []map[string]any `json:"data"`
}

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

func (h *Handler) getJsonFiles(c *gin.Context) {
	data, err := h.services.ParsingJSON.GetJsonTable()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetAllJSONResponse{
		Data: data,
	})
}
