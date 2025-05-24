package handlers

import (
	"net/http"
	"strconv"

	"github.com/AronditFire/todo-app/entity"
	"github.com/gin-gonic/gin"
)

type GetAllTaskResponse struct {
	Data []entity.Task `json:"data"`
}

// @Summary Get All tasks
// @Security ApiKeyAuth
// @Tags tasks
// @Description get all tasks
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} GetAllTaskResponse
// @Failure 400,404 {string} string "error"
// @Failure 500 {string} string "error"
// @Failure default {string} string "error"
// @Router /api/ [get]
func (h *Handler) getAllTasks(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	tasks, err := h.services.TaskList.GetAllTask(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get tasks for this user",
		})
		return
	}
	c.JSON(http.StatusOK, GetAllTaskResponse{
		Data: tasks,
	})
}

func (h *Handler) getTaskByID(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	task, err := h.services.TaskList.GetTaskByID(userID, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get task by ID",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) createTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	var req entity.Task
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Could not unbind request while creating task",
		})
		return
	}

	err = h.services.TaskList.CreateTask(userID, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create task in database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
	})
}

func (h *Handler) updateTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	var updatedDesc entity.TaskRequest
	if err := c.BindJSON(&updatedDesc); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Could not unbind request while updating task",
		})
		return
	}

	err = h.services.TaskList.UpdateTask(userID, id, updatedDesc.Description)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update task in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
	})
}

func (h *Handler) deleteTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	err = h.services.TaskList.DeleteTask(userID, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete task in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
