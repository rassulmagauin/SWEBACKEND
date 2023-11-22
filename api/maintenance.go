package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

type createMaintenanceRecordRequest struct {
	VehicleID           *uint      `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	MaintenancePersonID *uint      `gorm:"not null" json:"maintenance_person_id"`
	MaintenanceDate     *time.Time `json:"maintenance_date"`
	ServiceType         *string    `json:"service_type"`
	Status              *string    `gorm:"not null" json:"status"`
	TotalCost           *float64   `gorm:"not null" json:"total_cost"`
	MileageAtService    *int       `json:"mileage_at_service"`
	Notes               *string    `json:"notes"`
}
type createMaintenanceRecordResponse struct {
	VehicleID           *uint      `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	MaintenancePersonID *uint      `gorm:"not null" json:"maintenance_person_id"`
	MaintenanceDate     *time.Time `json:"maintenance_date"`
	ServiceType         *string    `json:"service_type"`
	Status              *string    `gorm:"not null" json:"status"`

	TotalCost        *float64 `gorm:"not null" json:"total_cost"`
	MileageAtService *int     `json:"mileage_at_service"`
	Notes            *string  `json:"notes"`
}

// CreateMaintenanceRecord godoc
// @Summary Create a maintenance record
// @Description Create a maintenance record
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param maintenance body createMaintenanceRecordRequest true "Maintenance"
// @Success 200 {object} createMaintenanceRecordResponse{}
// @Router /maintenance [post]
// @Security ApiKeyAuth
func (s *Server) CreateMaintenanceRecord(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can create maintenance records")))
		return
	}
	var maintenance models.MaintenanceRecord
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Create(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}

	c.JSON(200, maintenance)
}

// GetMaintenanceRecordsOfVehicle godoc
// @Summary Get all maintenance records of particular vehicle
// @Description Get all maintenance records if vehicle
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param vehicle_id path int true "Vehicle ID"
// @Success 200 {object} []createMaintenanceRecordResponse{}
// @Router /maintenance [get]
// @Security ApiKeyAuth
func (s *Server) GetMaintenanceRecordsOfVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can get maintenance records")))
		return
	}
	var maintenance []models.MaintenanceRecord
	if err := s.DB.Where("vehicle_id = ?", c.Param("vehicle_id")).Find(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}

// GetMaintenanceRecords godoc
// @Summary Get all maintenance records
// @Description Get all maintenance records
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Success 200 {object} []createMaintenanceRecordResponse{}
// @Router /maintenance [get]
// @Security ApiKeyAuth
func (s *Server) GetMaintenanceRecords(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can get maintenance records")))
		return
	}
	var maintenance []models.MaintenanceRecord
	if err := s.DB.Find(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}

// GetMaintenanceRecord godoc
// @Summary Get a maintenance record
// @Description Get a maintenance record
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param id path int true "Maintenance ID"
// @Success 200 {object} createMaintenanceRecordResponse{}
// @Router /maintenance/{id} [get]
// @Security ApiKeyAuth
func (s *Server) GetMaintenanceRecord(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can get maintenance records")))
		return
	}
	var maintenance models.MaintenanceRecord
	if err := s.DB.First(&maintenance, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}

// UpdateMaintenanceRecord godoc
// @Summary Update a maintenance record
// @Description Update a maintenance record
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param id path int true "Maintenance ID"
// @Param maintenance body createMaintenanceRecordRequest true "Maintenance"
// @Success 200 {object} createMaintenanceRecordResponse{}
// @Router /maintenance/{id} [put]
// @Security ApiKeyAuth
func (s *Server) UpdateMaintenanceRecord(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can update maintenance records")))
		return
	}
	var maintenance models.MaintenanceRecord
	if err := s.DB.First(&maintenance, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Save(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}

// DeleteMaintenanceRecord godoc
// @Summary Delete a maintenance record
// @Description Delete a maintenance record
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param id path int true "Maintenance ID"
// @Success 200 {object} createMaintenanceRecordResponse{}
// @Router /maintenance/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteMaintenanceRecord(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can delete maintenance records")))
		return
	}
	var maintenance models.MaintenanceRecord
	if err := s.DB.First(&maintenance, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Delete(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}

// GetMaintenanceRecordsOfUser godoc
// @Summary Get all maintenance records of particular user
// @Description Get all maintenance records if user
// @Tags maintenance
// @Accept  json
// @Produce  json
// @Param user_id path int true "User ID"
// @Success 200 {object} []createMaintenanceRecordResponse{}
// @Router /maintenance [get]
// @Security ApiKeyAuth
func (s *Server) GetMaintenanceRecordsOfUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Maintenance" {
		c.JSON(400, errorResponse(errors.New("only admins and maintenance can get maintenance records")))
		return
	}
	var maintenance []models.MaintenanceRecord
	if err := s.DB.Where("user_id = ?", c.Param("user_id")).Find(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, maintenance)
}
