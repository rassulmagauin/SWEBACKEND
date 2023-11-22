package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
)

type AuctionVehicleResponse struct {
	ID        uint            `json:"id"`
	VehicleID *uint           `json:"vehicle_id"`
	Details   *string         `json:"details"`
	Images    []ImageResponse `json:"images"`
	// Include other fields from gorm.Model if needed
}

type ImageResponse struct {
	ID  uint    `json:"id"`
	Url *string `json:"url"`
}

// CreateAuction godoc
// @Summary Create an auction
// @Description Admins can create auctions with vehicle details and images
// @Tags auction
// @Accept multipart/form-data
// @Produce json
// @Param vehicle_id formData uint true "Vehicle ID"
// @Param details formData string false "Details of the auction"
// @Param images formData file false "Images for the auction"
// @Success 200 {object} AuctionVehicleResponse "Successful response with auction details"
// @Router /auction [post]
// @Security ApiKeyAuth
func (s *Server) CreateAuction(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("only admins can create auctions")))
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Manually extract and set fields for AuctionVehicle
	var auction models.AuctionVehicle
	auction.VehicleID = parseUint(c.PostForm("vehicle_id"))
	auction.Details = parseString(c.PostForm("details"))

	// Process image uploads
	files := c.Request.MultipartForm.File["images"] // "images" is the name attribute in the form
	for _, file := range files {
		imagePath, err := saveFile(file, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Create and append the image to the auction.Images slice
		auction.Images = append(auction.Images, models.Image{Url: &imagePath})
	}

	// Create the auction record with images in the database
	if err := s.DB.Create(&auction).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, auction)
}

// Helper function to parse uint from form value
func parseUint(s string) *uint {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	u := uint(v)
	return &u
}

// Helper function to parse string pointer from form value
func parseString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Include your saveFile function here

// GetAuctions godoc
// @Summary Get all auctions
// @Description Get all auctions
// @Tags auction
// @Produce json
// @Success 200 {object} []AuctionVehicleResponse "Successful response with auction details"
// @Router /auction [get]
func (s *Server) GetAuctions(c *gin.Context) {
	var auctions []models.AuctionVehicle
	if err := s.DB.Preload("Images").Find(&auctions).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Convert file paths to URLs
	for i := range auctions {
		for j := range auctions[i].Images {
			temp := convertFilePathToURL(auctions[i].Images[j].Url)
			auctions[i].Images[j].Url = &temp
		}
	}

	c.JSON(http.StatusOK, auctions)
}

// GetAuction godoc
// @Summary Get an auction
// @Description Get an auction
// @Tags auction
// @Produce json
// @Param id path int true "Auction ID"
// @Success 200 {object} AuctionVehicleResponse "Successful response with auction details"
// @Router /auction/{id} [get]
func (s *Server) GetAuction(c *gin.Context) {
	var auction models.AuctionVehicle
	if err := s.DB.Preload("Images").First(&auction, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Convert file paths to URLs
	for i := range auction.Images {
		temp := convertFilePathToURL(auction.Images[i].Url)
		auction.Images[i].Url = &temp
	}

	c.JSON(http.StatusOK, auction)
}

// DeleteAuction godoc
// @Summary Delete an auction
// @Description Delete an auction
// @Tags auction
// @Produce json
// @Param id path int true "Auction ID"
// @Success 200 {object} string "Successful response with message"
// @Router /auction/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteAuction(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("only admins can delete auctions")))
		return
	}

	// Find the auction with its images
	var auction models.AuctionVehicle
	if err := s.DB.Preload("Images").First(&auction, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Delete image files from the server
	for _, image := range auction.Images {
		if image.Url != nil {
			filePath := filepath.Join("uploads", *image.Url) // Construct the full file path
			err := os.Remove(filePath)
			if err != nil {
				// Handle the error but don't stop the process
				log.Printf("Error deleting file: %v", err)
			}
		}
	}

	// Delete the auction record from the database
	if err := s.DB.Delete(&auction).Error; err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auction deleted successfully"})
}
