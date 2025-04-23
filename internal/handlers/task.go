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

func (h *Handler) getAllTasks(c *gin.Context) {
	tasks, err := h.services.TaskList.GetAllTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get tasks",
		})
		return
	}
	c.JSON(http.StatusOK, GetAllTaskResponse{
		Data: tasks,
	})
}

func (h *Handler) getTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	task, err := h.services.TaskList.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get task by ID",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) createTask(c *gin.Context) {
	var req entity.Task
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not unbind request while creating task",
		})
		return
	}

	err := h.services.TaskList.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create task in database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
	})
}

// TODO: При попытке обновления не существующего таска, вместо ошибки создаёт новый

func (h *Handler) updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	var updatedDesc entity.TaskRequest
	if err := c.BindJSON(&updatedDesc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not unbind request while updating task",
		})
		return
	}

	err = h.services.TaskList.UpdateTask(id, updatedDesc.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update task in database",
		})
		return
	}
}

// TODO: При удалении несуществующего такса код 200
func (h *Handler) deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	err = h.services.TaskList.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete task in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
