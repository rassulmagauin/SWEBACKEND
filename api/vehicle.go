package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

type createVehicleRequest struct {
	Make            *string    `json:"make"`
	CarModel        *string    `json:"car_model"`
	Year            *int       `json:"year"`
	LicensePlate    *string    `gorm:"not null" json:"license_plate"`
	SittingCapacity *int       `json:"sitting_capacity"`
	Type            *string    `json:"type"`
	Color           *string    `json:"color"`
	VIN             *string    `gorm:"not null" json:"vin"`
	CurrentMileage  *int       `json:"current_mileage"`
	LastMaintenance *time.Time `json:"last_maintenance"`
	NextMaintenance *time.Time `json:"next_maintenance"`
	Status          *string    `gorm:"not null" json:"status"`
	AssignedDriver  *uint      `json:"assigned_driver"`
	Notes           *string    `json:"notes"`
}

type createVehicleResponse struct {
	ID              uint       `gorm:"not null" json:"ID"`
	Make            *string    `json:"make"`
	CarModel        *string    `json:"car_model"`
	Year            *int       `json:"year"`
	LicensePlate    *string    `gorm:"not null" json:"license_plate"`
	SittingCapacity *int       `json:"sitting_capacity"`
	Type            *string    `json:"type"`
	Color           *string    `json:"color"`
	VIN             *string    `gorm:"not null" json:"vin"`
	CurrentMileage  *int       `json:"current_mileage"`
	LastMaintenance *time.Time `json:"last_maintenance"`
	NextMaintenance *time.Time `json:"next_maintenance"`
	Status          *string    `gorm:"not null" json:"status"`
	AssignedDriver  *uint      `json:"assigned_driver"`
	Notes           *string    `json:"notes"`
}

// CreateVehicle godoc
// @Summary Create a vehicle
// @Description Create a vehicle
// @Tags vehicle
// @Accept  json
// @Produce  json
// @Param vehicle body createVehicleRequest true "Vehicle"
// @Success 200 {object} createVehicleResponse{}
// @Router /vehicle [post]
// @Security ApiKeyAuth
func (s *Server) CreateVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can create vehicles")))
		return
	}
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	temp := "Active"
	vehicle.Status = &temp
	if err := s.DB.Create(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	createVehicleResponse := createVehicleResponse{
		ID:              vehicle.ID,
		Make:            vehicle.Make,
		CarModel:        vehicle.CarModel,
		Year:            vehicle.Year,
		LicensePlate:    vehicle.LicensePlate,
		SittingCapacity: vehicle.SittingCapacity,
		Type:            vehicle.Type,
		Color:           vehicle.Color,
		VIN:             vehicle.VIN,
		CurrentMileage:  vehicle.CurrentMileage,
		LastMaintenance: vehicle.LastMaintenance,
		NextMaintenance: vehicle.NextMaintenance,
		Status:          vehicle.Status,
		AssignedDriver:  vehicle.AssignedDriver,
		Notes:           vehicle.Notes,
	}
	c.JSON(200, createVehicleResponse)
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Get all vehicles
// @Tags vehicle
// @Produce  json
// @Success 200 {array} []createVehicleResponse{}
// @Router /vehicle [get]
// @Security ApiKeyAuth
func (s *Server) GetVehicles(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	var vehicles []models.Vehicle
	if authPayload.Role == "Driver" {
		fmt.Println("Here")
		var user models.User
		username := authPayload.Username
		if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		userID := user.ID
		if err := s.DB.Where("assigned_driver = ?", userID).Find(&vehicles).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		c.JSON(200, vehicles)
		return
	}

	if err := s.DB.Find(&vehicles).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, vehicles)
}

// GetVehicle godoc
// @Summary Get a vehicle
// @Description Get a vehicle
// @Tags vehicle
// @Produce  json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} createVehicleResponse{}
// @Router /vehicle/{id} [get]
// @Security ApiKeyAuth
func (s *Server) GetVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	var vehicle models.Vehicle
	if authPayload.Role == "Driver" {
		var user models.User
		username := authPayload.Username
		if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		if err := s.DB.First(&vehicle, c.Param("id")).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		if vehicle.AssignedDriver == nil || vehicle.AssignedDriver != &user.ID {
			c.JSON(404, errorResponse(errors.New("driver has no assigned vehicles")))
		}

		c.JSON(200, vehicle)
		return
	}

	if err := s.DB.First(&vehicle, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, vehicle)
}

// UpdateVehicle godoc
// @Summary Update a vehicle
// @Description Update a vehicle
// @Tags vehicle
// @Accept  json
// @Produce  json
// @Param id path int true "Vehicle ID"
// @Param vehicle body createVehicleRequest true "Vehicle"
// @Success 200 {object} createVehicleResponse{}
// @Router /vehicle/{id} [put]
// @Security ApiKeyAuth
func (s *Server) UpdateVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can update vehicles")))
		return
	}
	var vehicle models.Vehicle
	if err := s.DB.First(&vehicle, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Save(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, vehicle)
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Delete a vehicle
// @Tags vehicle
// @Produce  json
// @Param id path int true "Vehicle ID"
// @Success 200
// @Router /vehicle/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can delete vehicles")))
		return
	}
	if err := s.DB.First(&vehicle, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Delete(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, gin.H{})
}

type registerVehicleRequest struct {
	Make            *string    `json:"make"`
	CarModel        *string    `json:"car_model"`
	Year            *int       `json:"year"`
	LicensePlate    *string    `gorm:"not null" json:"license_plate"`
	SittingCapacity *int       `json:"sitting_capacity"`
	Type            *string    `json:"type"`
	Color           *string    `json:"color"`
	VIN             *string    `gorm:"not null" json:"vin"`
	CurrentMileage  *int       `json:"current_mileage"`
	LastMaintenance *time.Time `json:"last_maintenance"`
	NextMaintenance *time.Time `json:"next_maintenance"`
	AssignedDriver  *uint      `json:"assigned_driver"`
	Notes           *string    `json:"notes"`
}

// RegisterVehicle godoc
// @Summary Register a vehicle
// @Description Register a vehicle
// @Tags vehicle
// @Accept  json
// @Produce  json
// @Param vehicle body registerVehicleRequest true "Vehicle"
// @Success 200 {object} createVehicleResponse{}
// @Router /vehicle/register [post]
// @Security ApiKeyAuth
func (s *Server) RegisterVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	temp := "Pending"
	vehicle.Status = &temp
	if err := s.DB.Create(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, vehicle)
}

// ActivateVehicle godoc
// @Summary Activate pending vehicle
// @Description Changes status of pending vehicles to Active
// @Tags vehicle
// @Produce  json
// @Param id path int true "Vehicle ID"
// @Success 200 {array} createVehicleResponse{}
// @Router /vehicle/{id} [post]
// @Security ApiKeyAuth
func (s *Server) ActivateVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can activate vehicles")))
		return
	}
	var vehicle models.Vehicle
	if err := s.DB.First(&vehicle, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	temp := "Active"
	vehicle.Status = &temp
	if err := s.DB.Save(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, vehicle)
}

type assignVehicleRequest struct {
	UserID    uint `json:"user_id"`
	VehicleID uint `json:"vehicle_id"`
}

// AssignVehicle godoc
// @Summary Assign a vehicle to a driver
// @Description Assign a vehicle to a driver
// @Tags vehicle
// @Accept  json
// @Produce  json
// @Param vehicle body assignVehicleRequest true "Vehicle"
// @Success 200 {array} userResponse{}
// @Router /vehicle/assign [post]
// @Security ApiKeyAuth
func (s *Server) AssignVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can assign vehicles")))
		return
	}
	var req assignVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var user models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var vehicle models.Vehicle
	if err := s.DB.Where("id = ?", req.VehicleID).First(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if vehicle.AssignedDriver != nil {
		c.JSON(400, errorResponse(errors.New("vehicle already assigned to a driver")))
		return
	}
	vehicle.AssignedDriver = &user.ID
	if err := s.DB.Save(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	user.Vehicles = append(user.Vehicles, vehicle)
	if err := s.DB.Save(&user).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	userResponse := newUserResonse(user)
	c.JSON(200, userResponse)
}

// UnassignVehicle godoc
// @Summary Unassign a vehicle from a driver
// @Description Unassign a vehicle from a driver
// @Tags vehicle
// @Accept  json
// @Produce  json
// @Param vehicle body assignVehicleRequest true "Vehicle"
// @Success 200 {array} userResponse{}
// @Router /vehicle/unassign [post]
// @Security ApiKeyAuth
func (s *Server) UnassignVehicle(c *gin.Context) {

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can unassign vehicles")))
		return
	}
	var req assignVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var user models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var vehicle models.Vehicle
	if err := s.DB.Where("id = ?", req.VehicleID).First(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if vehicle.AssignedDriver == nil {
		c.JSON(400, errorResponse(errors.New("vehicle is not assigned to a driver")))
		return
	}
	vehicle.AssignedDriver = nil
	if err := s.DB.Save(&vehicle).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var vehicles []models.Vehicle
	for _, v := range user.Vehicles {
		if v.ID != vehicle.ID {
			vehicles = append(vehicles, v)
		}
	}
	user.Vehicles = vehicles
	if err := s.DB.Save(&user).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	response := newUserResonse(user)
	c.JSON(200, response)
}
