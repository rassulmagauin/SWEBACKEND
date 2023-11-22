package config

import (
	"github.com/rassulmagauin/VMS_SWE/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open("host=localhost user=root password=sweteam3 dbname=root port=5555 sslmode=disable"), &gorm.Config{})
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
