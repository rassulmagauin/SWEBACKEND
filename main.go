package main

import (
	"log"
	"os"

	"github.com/rassulmagauin/VMS_SWE/api"
	"github.com/rassulmagauin/VMS_SWE/config"
	_ "github.com/rassulmagauin/VMS_SWE/docs"
)

func ensureUploadsDir() {
	path := "./uploads"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm) // or os.MkdirAll for nested directories
	}
}

// @title Vehicle Management System API
// @version 1.0
// @description This API serves as a backend for Vehicle Management System
// @host 	https://swebackend-production.up.railway.app
// @BasePath: /api
// @schemes http https
// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.Connect()
	server, err := api.NewServer(config.DB)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}

	ensureUploadsDir()
	// Run the server
	err = server.Run(":" + port)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

}
