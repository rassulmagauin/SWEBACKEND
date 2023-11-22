package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

type Report struct {
	Vehicle            models.Vehicle             `json:"vehicle"`
	FuelingRecords     []models.FuelingRecord     `json:"fueling_records"`
	MaintenanceRecords []models.MaintenanceRecord `json:"maintenance_records"`
}
type ReportForSwagger struct {
	Vehicle            createVehicleResponse `json:"vehicle"`
	FuelingRecords     []myFuelingRecord     `json:"fueling_records"`
	MaintenanceRecords []myMaintenanceRecord `json:"maintenance_records"`
}
type myFuelingRecord struct {
	VehicleID          *uint    `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	FuelingPersonID    *uint    `gorm:"not null" json:"fueling_person_id"`
	Amount             *float64 `gorm:"not null" json:"amount"`
	TotalCost          *float64 `gorm:"not null" json:"total_cost"`
	GasStation         *string  `json:"gas_station"`
	Notes              *string  `json:"notes"`
	BeforeFuelingImage *string  `gorm:"not null" json:"before_fueling_image"`
	AfterFuelingImage  *string  `gorm:"not null" json:"after_fueling_image"`
}
type myMaintenanceRecord struct {
	VehicleID           *uint      `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	MaintenancePersonID *uint      `gorm:"not null" json:"maintenance_person_id"`
	MaintenanceDate     *time.Time `json:"maintenance_date"`
	ServiceType         *string    `json:"service_type"`
	Status              *string    `gorm:"not null" json:"status"`
	//Parts               []Part             `gorm:"foreignKey:ID" json:"parts"`
	TotalCost        *float64 `gorm:"not null" json:"total_cost"`
	MileageAtService *int     `json:"mileage_at_service"`
	Notes            *string  `json:"notes"`
}

// this handler returns all fueiling record, maintenance records for a vehicle
// GetReport godoc
// @Summary Get a report
// @Description Get a report
// @Tags report
// @Accept  json
// @Produce  json
// @Param vehicle_id path int true "Vehicle ID"
// @Success 200 {object} ReportForSwagger{}
// @Router /report/{vehicle_id} [get]
// @Security ApiKeyAuth
func (s *Server) GetReport(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can get reports")))
		return
	}
	var report Report
	var vehicle models.Vehicle
	if err := s.DB.First(&vehicle, c.Param("vehicle_id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	report.Vehicle = vehicle
	var fueling []models.FuelingRecord
	if err := s.DB.Where("vehicle_id = ?", c.Param("vehicle_id")).Find(&fueling).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	report.FuelingRecords = fueling
	var maintenance []models.MaintenanceRecord
	if err := s.DB.Where("vehicle_id = ?", c.Param("vehicle_id")).Find(&maintenance).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	report.MaintenanceRecords = maintenance
	c.JSON(200, report)
}
