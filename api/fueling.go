package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

//	func (s *Server) CreateFuelingRecord(c *gin.Context) {

//		var fueling models.FuelingRecord
//		if err := c.ShouldBindJSON(&fueling); err != nil {
//			c.JSON(400, errorResponse(err))
//			return
//		}
//		if err := s.DB.Create(&fueling).Error; err != nil {
//			c.JSON(400, errorResponse(err))
//			return
//		}
//		c.JSON(200, fueling)
//	}

type FuelingRecordResponse struct {
	ID                 uint      `json:"id"`
	VehicleID          *uint     `json:"vehicle_id"`
	FuelingPersonID    *uint     `json:"fueling_person_id"`
	Amount             *float64  `json:"amount"`
	TotalCost          *float64  `json:"total_cost"`
	GasStation         *string   `json:"gas_station,omitempty"` // omitempty if the field is optional
	Notes              *string   `json:"notes,omitempty"`       // omitempty if the field is optional
	BeforeFuelingImage *string   `json:"before_fueling_image"`
	AfterFuelingImage  *string   `json:"after_fueling_image"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	// Include other fields from gorm.Model if needed
	// If Vehicle is a detailed object and you want to include it in the response, define a corresponding struct
	// Vehicle            *VehicleResponse `json:"vehicle,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateFuelingRecord godoc
// @Summary Create fueling record
// @Description Admins and fueling personnel can create fueling records
// @Tags fueling
// @Accept multipart/form-data
// @Produce json
// @Param vehicle_id formData uint true "Vehicle ID"
// @Param fueling_person_id formData uint true "Fueling Person ID"
// @Param amount formData number true "Amount of fuel"
// @Param total_cost formData number true "Total cost of fueling"
// @Param gas_station formData string false "Gas Station"
// @Param notes formData string false "Additional notes"
// @Param before_fueling_image formData file true "Image before fueling"
// @Param after_fueling_image formData file true "Image after fueling"
// @Success 200 {object} FuelingRecordResponse "Successful response with fueling record details"
// @Failure 400 {object} ErrorResponse "Bad Request with error message"
// @Failure 500 {object} ErrorResponse "Internal Server Error with error message"
// @Router /fueling [post]
// @Security ApiKeyAuth
func (s *Server) CreateFuelingRecord(c *gin.Context) {

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Fueling" {
		c.JSON(400, errorResponse(errors.New("only admins and fueling can create fueling records")))
		return
	}
	// Initialize your FuelingRecord struct
	var fueling models.FuelingRecord

	// Manually handle non-file fields
	// Handling vehicle_id
	vehicleIDStr := c.PostForm("vehicle_id")
	if vehicleIDStr != "" {
		vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id must be a number"})
			return
		}
		vehicleIDUint := uint(vehicleID)
		fueling.VehicleID = &vehicleIDUint
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	// Handling fueling_person_id
	fuelingPersonIDStr := c.PostForm("fueling_person_id")
	if fuelingPersonIDStr != "" {
		fuelingPersonID, err := strconv.ParseUint(fuelingPersonIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fueling_person_id must be a number"})
			return
		}
		fuelingPersonIDUint := uint(fuelingPersonID)
		fueling.FuelingPersonID = &fuelingPersonIDUint
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fueling_person_id is required"})
		return
	}

	// Handling amount
	amountStr := c.PostForm("amount")
	if amountStr != "" {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be a valid number"})
			return
		}
		fueling.Amount = &amount
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is required"})
		return
	}

	// Handling total_cost
	totalCostStr := c.PostForm("total_cost")
	if totalCostStr != "" {
		totalCost, err := strconv.ParseFloat(totalCostStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "total_cost must be a valid number"})
			return
		}
		fueling.TotalCost = &totalCost
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "total_cost is required"})
		return
	}

	// Handling optional string fields: GasStation and Notes
	gasStation := c.PostForm("gas_station")
	if gasStation != "" {
		fueling.GasStation = &gasStation
	}

	notes := c.PostForm("notes")
	if notes != "" {
		fueling.Notes = &notes
	}

	// Handle file upload for BeforeFuelingImage
	beforeFuelingImage, err := c.FormFile("before_fueling_image")
	if err == nil { // If there's a file
		beforeImagePath, err := saveFile(beforeFuelingImage, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fueling.BeforeFuelingImage = &beforeImagePath
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "before_fueling_image is required"})
		return
	}

	// Handle file upload for AfterFuelingImage
	afterFuelingImage, err := c.FormFile("after_fueling_image")
	if err == nil { // If there's a file
		afterImagePath, err := saveFile(afterFuelingImage, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fueling.AfterFuelingImage = &afterImagePath
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "after_fueling_image is required"})
		return
	}

	// Save the record in the database
	if err := s.DB.Create(&fueling).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fueling)
}

func saveFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	// Create a unique filename to avoid conflicts
	newFileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("./uploads", newFileName)

	// Save the file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return newFileName, nil
}

// GetFuelingRecord godoc
// @Summary Get a fueling record
// @Description Get a fueling record
// @Tags fueling
// @Accept  json
// @Produce  json
// @Param id path int true "Fueling Record ID"
// @Success 200 {object} FuelingRecordResponse "Successful response with fueling record details"
// @Failure 400 {object} ErrorResponse "Bad Request with error message"
// @Failure 404 {object} ErrorResponse "Not Found with error message"
// @Router /fueling/{id} [get]
// @Security ApiKeyAuth
func (s *Server) GetFuelingRecord(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var fueling models.FuelingRecord
	if err := s.DB.First(&fueling, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	// Assuming you have a function to convert file paths to URLs
	temp1 := convertFilePathToURL(fueling.BeforeFuelingImage)
	temp2 := convertFilePathToURL(fueling.AfterFuelingImage)
	fueling.BeforeFuelingImage = &temp1
	fueling.AfterFuelingImage = &temp2

	c.JSON(http.StatusOK, fueling)
}
func convertFilePathToURL(filePath *string) string {
	if filePath == nil {
		return ""
	}

	// Use a relative URL for flexibility
	return fmt.Sprintf("/static/%s", *filePath)
}

// GetFuelingRecords godoc
// @Summary Get all fueling records
// @Description Get all fueling records
// @Tags fueling
// @Produce  json
// @Success 200 {object} []FuelingRecordResponse "Successful response with fueling record details"
// @Failure 400 {object} ErrorResponse "Bad Request with error message"
// @Router /fueling [get]
// @Security ApiKeyAuth
func (s *Server) GetFuelingRecords(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" && authPayload.Role != "Fueling" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("only admins and fueling can get fueling records")))
		return
	}

	var fuelings []models.FuelingRecord
	if err := s.DB.Find(&fuelings).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Convert file paths to URLs for each record
	for i := range fuelings {
		if fuelings[i].BeforeFuelingImage != nil {
			temo1 := convertFilePathToURL(fuelings[i].BeforeFuelingImage)
			fuelings[i].BeforeFuelingImage = &temo1
		}
		if fuelings[i].AfterFuelingImage != nil {
			temp2 := convertFilePathToURL(fuelings[i].AfterFuelingImage)
			fuelings[i].AfterFuelingImage = &temp2
		}
	}

	c.JSON(http.StatusOK, fuelings)
}

// Include your saveFile function here

// DeleteFuelingRecord godoc
// @Summary Delete a fueling record
// @Description Delete a fueling record
// @Tags fueling
// @Produce  json
// @Param id path int true "Fueling Record ID"
// @Success 200 {object} string "Successful response with message"
// @Failure 400 {object} ErrorResponse "Bad Request with error message"
// @Failure 404 {object} ErrorResponse "Not Found with error message"
// @Router /fueling/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteFuelingRecord(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Retrieve the record from the database
	var fueling models.FuelingRecord
	if err := s.DB.First(&fueling, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	// Delete the image files if they exist
	deleteFileIfExists(fueling.BeforeFuelingImage)
	deleteFileIfExists(fueling.AfterFuelingImage)

	// Delete the record from the database
	if err := s.DB.Delete(&fueling).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}

func deleteFileIfExists(filePath *string) {
	if filePath != nil {
		fullPath := "./uploads/" + *filePath // Adjust the path according to your directory structure
		if _, err := os.Stat(fullPath); err == nil {
			os.Remove(fullPath)
		}
	}
}

func (s *Server) GetFuelingRecordsOfVehicle(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Fueling" {
		c.JSON(400, errorResponse(errors.New("only admins and fueling can get fueling records")))
		return
	}
	var fueling []models.FuelingRecord
	if err := s.DB.Where("vehicle_id = ?", c.Param("vehicle_id")).Find(&fueling).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, fueling)
}

func (s *Server) GetFuelingRecordsOfUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" || authPayload.Role != "Fueling" {
		c.JSON(400, errorResponse(errors.New("only admins and fueling can get fueling records")))
		return
	}
	var fueling []models.FuelingRecord
	if err := s.DB.Where("user_id = ?", c.Param("user_id")).Find(&fueling).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	c.JSON(200, fueling)
}
