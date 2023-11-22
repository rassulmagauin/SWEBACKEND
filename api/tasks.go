package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

type createTaskRequest struct {
	DriverID       *uint      `gorm:"not null;onDelete:CASCADE" json:"driver_id"`
	StartLatitude  *float64   `gorm:"not null" json:"start_latitude"`
	StartLongitude *float64   `gorm:"not null" json:"start_longitude"`
	EndLatitude    *float64   `gorm:"not null" json:"end_latitude"`
	EndLongitude   *float64   `gorm:"not null" json:"end_longitude"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	Status         *string    `gorm:"not null" json:"status"`
	Notes          *string    `json:"notes"`
}

type createTaskResponse struct {
	ID             uint       `gorm:"not null" json:"ID"`
	DriverID       *uint      `gorm:"not null;onDelete:CASCADE" json:"driver_id"`
	StartLatitude  *float64   `gorm:"not null" json:"start_latitude"`
	StartLongitude *float64   `gorm:"not null" json:"start_longitude"`
	EndLatitude    *float64   `gorm:"not null" json:"end_latitude"`
	EndLongitude   *float64   `gorm:"not null" json:"end_longitude"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	Status         *string    `gorm:"not null" json:"status"`
	Notes          *string    `json:"notes"`
}

// CreateTask godoc
// @Summary Create a task
// @Description Create a task
// @Tags task
// @Accept  json
// @Produce  json
// @Param task body createTaskRequest true "Task"
// @Success 200 {object} createTaskResponse{}
// @Router /task [post]
// @Security ApiKeyAuth
func (s *Server) CreateTask(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can create tasks")))
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Create(&task).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, task)
}

// GetTasks godoc
// @Summary Get all tasks
// @Description Get all tasks
// @Tags task
// @Accept  json
// @Produce  json
// @Success 200 {array} []createTaskResponse{}
// @Router /task [get]
// @Security ApiKeyAuth
func (s *Server) GetTasks(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		var user models.User
		username := authPayload.Username
		if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		userID := user.ID
		var tasks []models.Task
		if err := s.DB.Where("driver_id = ?", userID).Find(&tasks).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		c.JSON(200, tasks)
		return
	}
	var tasks []models.Task
	if err := s.DB.Find(&tasks).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, tasks)
}

// GetTask godoc
// @Summary Get a task
// @Description Get a task
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} createTaskResponse{}
// @Router /task/{id} [get]
// @Security ApiKeyAuth
func (s *Server) GetTask(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can get tasks")))
		return
	}
	var task models.Task
	if err := s.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, task)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Param task body createTaskRequest true "Task"
// @Success 200 {object} createTaskResponse{}
// @Router /task/{id} [put]
// @Security ApiKeyAuth
func (s *Server) UpdateTask(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		var user models.User
		username := authPayload.Username
		if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		var task models.Task
		if err := s.DB.First(&task, c.Param("id")).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		if task.DriverID == nil || task.DriverID != &user.ID {
			c.JSON(404, errorResponse(errors.New("driver has no assigned tasks")))
		}
		temp := "Completed"
		task.Status = &temp
		if err := s.DB.Save(&task).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		c.JSON(200, task)
		return
	}
	var task models.Task
	if err := s.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Save(&task).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, task)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} createTaskResponse{}
// @Router /task/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteTask(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can delete tasks")))
		return
	}
	var task models.Task
	if err := s.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Delete(&task).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, task)
}
