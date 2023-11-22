package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/rassulmagauin/VMS_SWE/docs"
	"github.com/rassulmagauin/VMS_SWE/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	Router     *gin.Engine
	tokenMaker token.Maker
	DB         *gorm.DB
}

const (
	tokenSymmetricKey   = "12345678901234567890123456789012"
	AccessTokenDuration = 15 * time.Minute
)

func NewServer(DB *gorm.DB) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(tokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{DB: DB, tokenMaker: tokenMaker}
	server.setupRouter()
	return server, nil
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Static("/static", "./uploads")
	swaggerHost := os.Getenv("HOST")
	if swaggerHost == "" {
		swaggerHost = "localhost:8080" // Default value for local development
	}
	swaggerURL := fmt.Sprintf("https://%s/docs/doc.json", swaggerHost)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerURL)))

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/vehicle", server.CreateVehicle)
	authRoutes.GET("/vehicle", server.GetVehicles)
	authRoutes.GET("/vehicle/:id", server.GetVehicle)
	authRoutes.PUT("/vehicle/:id", server.UpdateVehicle)
	authRoutes.DELETE("/vehicle/:id", server.DeleteVehicle)
	authRoutes.POST("/vehicle/:id", server.ActivateVehicle)
	authRoutes.POST("/vehicle/register", server.RegisterVehicle)

	router.POST("/user", server.CreateUser)
	authRoutes.GET("/user", server.GetUsers)
	authRoutes.GET("/user/:id", server.GetUser)
	authRoutes.PUT("/user/:id", server.UpdateUser)
	authRoutes.DELETE("/user/:id", server.DeleteUser)

	authRoutes.POST("/maintenance", server.CreateMaintenanceRecord)
	authRoutes.GET("/maintenance", server.GetMaintenanceRecords)
	authRoutes.GET("/maintenance/:id", server.GetMaintenanceRecord)
	authRoutes.PUT("/maintenance/:id", server.UpdateMaintenanceRecord)
	authRoutes.DELETE("/maintenance/:id", server.DeleteMaintenanceRecord)
	authRoutes.GET("/maintenances/:vehicle_id", server.GetMaintenanceRecordsOfVehicle)
	authRoutes.GET("/maintenance/user/:user_id", server.GetMaintenanceRecordsOfUser)

	authRoutes.POST("/fueling", server.CreateFuelingRecord)
	authRoutes.GET("/fueling", server.GetFuelingRecords)
	authRoutes.GET("/fueling/:id", server.GetFuelingRecord)
	authRoutes.DELETE("/fueling/:id", server.DeleteFuelingRecord)
	authRoutes.GET("/fuelings/:vehicle_id", server.GetFuelingRecordsOfVehicle)
	authRoutes.GET("/fueling/user/:user_id", server.GetFuelingRecordsOfUser)

	authRoutes.POST("/vehicle/assign", server.AssignVehicle)
	authRoutes.POST("/vehicle/unassign", server.UnassignVehicle)

	authRoutes.POST("/task", server.CreateTask)
	authRoutes.GET("/task", server.GetTasks)
	authRoutes.GET("/task/:id", server.GetTask)
	authRoutes.PUT("/task/:id", server.UpdateTask)
	authRoutes.DELETE("/task/:id", server.DeleteTask)

	authRoutes.GET("/report/:vehicle_id", server.GetReport)

	authRoutes.POST("/auction", server.CreateAuction)
	router.GET("/auction", server.GetAuctions)
	router.GET("/auction/:id", server.GetAuction)
	authRoutes.DELETE("/auction/:id", server.DeleteAuction)

	router.POST("/login", server.LoginUser)
	server.Router = router
}
