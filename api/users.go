package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rassulmagauin/VMS_SWE/models"
	"github.com/rassulmagauin/VMS_SWE/token"
	"github.com/rassulmagauin/VMS_SWE/utils"
	"gorm.io/gorm"
)

type createUserRequest struct {
	Username             string  `gorm:"not null" json:"username"`
	Password             *string `gorm:"not null" json:"password"`
	GovermentID          *string `json:"goverment_id"`
	MiddleName           *string `json:"middle_name"`
	Address              *string `json:"address"`
	PhoneNumber          *string `json:"phone_number"`
	DrivingLicenseNumber *string `json:"driving_license_number"`
	Role                 *string `gorm:"not null" json:"role"`
	FirstName            *string `gorm:"not null" json:"first_name"`
	LastName             *string `gorm:"not null" json:"last_name"`
	Email                *string `json:"email"`
	Status               *string `json:"status"`
}

type userResponse struct {
	ID                   uint    `json:"ID"`
	FirstName            *string `json:"first_name"`
	LastName             *string `json:"last_name"`
	Username             string  `json:"username"`
	GovermentID          *string `json:"goverment_id"`
	MiddleName           *string `json:"middle_name"`
	Address              *string `json:"address"`
	PhoneNumber          *string `json:"phone_number"`
	DrivingLicenseNumber *string `json:"driving_license_number"`
	Role                 *string `json:"role"`
	Email                *string `json:"email"`
	Status               *string `json:"status"`
}

func newUserResonse(user models.User) userResponse {
	return userResponse{
		ID:                   user.ID,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Username:             user.Username,
		GovermentID:          user.GovermentID,
		MiddleName:           user.MiddleName,
		Address:              user.Address,
		PhoneNumber:          user.PhoneNumber,
		DrivingLicenseNumber: user.DrivingLicenseNumber,
		Role:                 user.Role,
		Email:                user.Email,
		Status:               user.Status,
	}
}

// CreateUser godoc
// @Summary Create user
// @Description Creates and saves user to database
// @Param user body createUserRequest true "User"
// @Produce application/json
// @Tags user
// @Success 200 {object} userResponse{}
// @Router /user [post]
func (s *Server) CreateUser(c *gin.Context) {
	// authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	// if authPayload.Role != "Admin" {
	// 	c.JSON(400, errorResponse(errors.New("only admins can create users")))
	// 	return
	// }
	var userReq createUserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(*userReq.Password)
	if err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	user := models.User{
		Username:             userReq.Username,
		HashedPassword:       &hashedPassword,
		GovermentID:          userReq.GovermentID,
		MiddleName:           userReq.MiddleName,
		Address:              userReq.Address,
		PhoneNumber:          userReq.PhoneNumber,
		DrivingLicenseNumber: userReq.DrivingLicenseNumber,
		Role:                 userReq.Role,
		FirstName:            userReq.FirstName,
		LastName:             userReq.LastName,
		Email:                userReq.Email,
		Status:               userReq.Status,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	response := newUserResonse(user)
	c.JSON(200, response)
}

type getUserResponse struct {
	ID                   uint    `gorm:"not null" json:"ID"`
	FirstName            *string `gorm:"not null" json:"first_name"`
	LastName             *string `gorm:"not null" json:"last_name"`
	Username             string  `gorm:"not null;unique" json:"username"`
	GovermentID          *string `json:"goverment_id"`
	MiddleName           *string `json:"middle_name"`
	Address              *string `json:"address"`
	PhoneNumber          *string `json:"phone_number"`
	DrivingLicenseNumber *string `json:"driving_license_number"`
	Role                 *string `gorm:"not null" json:"role"`
	Email                *string `json:"email"`
	Status               *string `json:"status"`
}

// GetUsers godoc
// @Summary Get users
// @Description Gets all users from database
// @Produce application/json
// @Tags user
// @Success 200 {object} []getUserResponse{}
// @Router /user [get]
// @Security ApiKeyAuth
func (s *Server) GetUsers(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		var user models.User
		username := authPayload.Username
		if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(400, errorResponse(err))
			return
		}
		response := getUserResponse{
			ID:                   user.ID,
			FirstName:            user.FirstName,
			LastName:             user.LastName,
			Username:             user.Username,
			GovermentID:          user.GovermentID,
			MiddleName:           user.MiddleName,
			Address:              user.Address,
			PhoneNumber:          user.PhoneNumber,
			DrivingLicenseNumber: user.DrivingLicenseNumber,
			Role:                 user.Role,
			Email:                user.Email,
			Status:               user.Status,
		}
		c.JSON(200, response)
		return
	}
	var users []models.User
	if err := s.DB.Find(&users).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	response := make([]getUserResponse, len(users))
	for i, user := range users {
		response[i] = getUserResponse{
			ID:                   user.ID,
			FirstName:            user.FirstName,
			LastName:             user.LastName,
			Username:             user.Username,
			GovermentID:          user.GovermentID,
			MiddleName:           user.MiddleName,
			Address:              user.Address,
			PhoneNumber:          user.PhoneNumber,
			DrivingLicenseNumber: user.DrivingLicenseNumber,
			Role:                 user.Role,
			Email:                user.Email,
			Status:               user.Status,
		}
	}
	c.JSON(200, response)
}

// GetUser godoc
// @Summary Get user
// @Description Gets user from database
// @Param id path string true "User ID"
// @Produce application/json
// @Tags user
// @Success 200 {object} getUserResponse{}
// @Router /user/{id} [get]
// @Security ApiKeyAuth
func (s *Server) GetUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can get users")))
		return
	}
	var user models.User
	if err := s.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	response := getUserResponse{
		ID:                   user.ID,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Username:             user.Username,
		GovermentID:          user.GovermentID,
		MiddleName:           user.MiddleName,
		Address:              user.Address,
		PhoneNumber:          user.PhoneNumber,
		DrivingLicenseNumber: user.DrivingLicenseNumber,
		Role:                 user.Role,
		Email:                user.Email,
		Status:               user.Status,
	}
	c.JSON(200, response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates user in database
// @Param id path string true "User ID"
// @Param user body createUserRequest true "User"
// @Produce application/json
// @Tags user
// @Success 200 {object} userResponse{}
// @Router /user/{id} [put]
// @Security ApiKeyAuth
func (s *Server) UpdateUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can update users")))
		return
	}
	var user models.User
	if err := s.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var userReq createUserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(*userReq.Password)
	if err != nil {
		c.JSON(400, errorResponse(err))
		return

	}
	user.Username = userReq.Username
	user.HashedPassword = &hashedPassword
	user.GovermentID = userReq.GovermentID
	user.MiddleName = userReq.MiddleName
	user.Address = userReq.Address
	user.PhoneNumber = userReq.PhoneNumber
	user.DrivingLicenseNumber = userReq.DrivingLicenseNumber
	user.Role = userReq.Role
	user.FirstName = userReq.FirstName
	user.LastName = userReq.LastName
	user.Email = userReq.Email
	user.Status = userReq.Status

	if err := s.DB.Save(&user).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	response := newUserResonse(user)
	c.JSON(200, response)
}

type deleteUserResponse struct {
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes user from database
// @Param id path string true "User ID"
// @Produce application/json
// @Tags user
// @Success 200 {object} deleteUserResponse{}
// @Router /user/{id} [delete]
// @Security ApiKeyAuth
func (s *Server) DeleteUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "Admin" {
		c.JSON(400, errorResponse(errors.New("only admins can delete users")))
		return
	}
	var user models.User
	if err := s.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	if err := s.DB.Delete(&user).Error; err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	deleteUserResponse := deleteUserResponse{}
	c.JSON(200, deleteUserResponse)
}

type loginRequest struct {
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}
type loginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

// LoginUser godoc
// @Summary Login user
// @Description Logs user in
// @Param user body loginRequest true "User"
// @Produce application/json
// @Tags user
// @Success 200 {object} loginResponse{}
// @Router /login [post]
func (s *Server) LoginUser(c *gin.Context) {
	var loginReq loginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}
	var user models.User
	if err := s.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(400, errorResponse(err))
		return
	}
	if err := utils.CheckPassword(loginReq.Password, *user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := s.tokenMaker.CreateToken(user.Username, *user.Role, AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := loginResponse{
		AccessToken: accessToken,
		User:        newUserResonse(user),
	}
	c.JSON(http.StatusOK, response)
}
