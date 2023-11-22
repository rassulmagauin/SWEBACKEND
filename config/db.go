package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rassulmagauin/VMS_SWE/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE") // if you're using SSL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	if host == "" {
		host = "localhost"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Exec(`
    CREATE TYPE task_status AS ENUM ('completed', 'canceled', 'delayed'); 
    CREATE TYPE vehicle_status AS ENUM ('Active', 'Inactive', 'Maintenance'); 
    CREATE TYPE appointment_status AS ENUM ('Pending', 'Confirmed', 'Cancelled'); 
    CREATE TYPE auction_status AS ENUM ('Sold', 'Pending'); 
    CREATE TYPE maintenance_status AS ENUM ('Pending', 'Done'); 
    CREATE TYPE roles_list AS ENUM ('Admin', 'Driver', 'Fueling_person', 'Maintenance_person');
`)

	err = db.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.Task{},
		&models.Appointment{},
		&models.AuctionVehicle{},
		&models.MaintenanceRecord{},
		&models.FuelingRecord{},
		&models.VehicleUsage{},
		&models.RolePermission{},
		&models.Image{},
		&models.Part{},
	)
	if err != nil {
		panic(err)
	}

	DB = db
}
