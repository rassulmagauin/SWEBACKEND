package models

import (
	"time"

	"gorm.io/gorm"
)

// Define enums as custom Go types
type TaskStatus string
type VehicleStatus string
type AppointmentStatus string
type MaintenanceStatus string
type RolesList string

// Constants for the enum values
const (
	TaskStatusCompleted        TaskStatus        = "completed"
	TaskStatusCanceled         TaskStatus        = "canceled"
	TaskStatusDelayed          TaskStatus        = "delayed"
	VehicleStatusActive        VehicleStatus     = "Active"
	VehicleStatusInactive      VehicleStatus     = "Inactive"
	VehicleStatusMaintenance   VehicleStatus     = "Maintenance"
	AppointmentStatusPending   AppointmentStatus = "Pending"
	AppointmentStatusConfirmed AppointmentStatus = "Confirmed"
	AppointmentStatusCancelled AppointmentStatus = "Cancelled"
	MaintenanceStatusPending   MaintenanceStatus = "Pending"
	MaintenanceStatusDone      MaintenanceStatus = "Done"
	RolesListAdmin             RolesList         = "Admin"
	RolesListDriver            RolesList         = "Driver"
	RolesListFuelingPerson     RolesList         = "Fueling_person"
	RolesListMaintenancePerson RolesList         = "Maintenance_person"
)

type User struct {
	ID                   uint       `gorm:"not null" json:"ID"`
	Username             string     `gorm:"not null;unique" json:"username"`
	HashedPassword       *string    `gorm:"not null" json:"hashed_password"`
	GovermentID          *string    `json:"goverment_id"`
	MiddleName           *string    `json:"middle_name"`
	Address              *string    `json:"address"`
	PhoneNumber          *string    `json:"phone_number"`
	DrivingLicenseNumber *string    `json:"driving_license_number"`
	Role                 *string    `gorm:"not null" json:"role"`
	FirstName            *string    `gorm:"not null" json:"first_name"`
	LastName             *string    `gorm:"not null" json:"last_name"`
	Email                *string    `json:"email"`
	LastLogin            *time.Time `json:"last_login"`
	Status               *string    `json:"status"`
	Tasks                []Task     `gorm:"foreignKey:DriverID"`
	Vehicles             []Vehicle  `gorm:"foreignKey:AssignedDriver"`
	gorm.Model
}

type Vehicle struct {
	ID                 uint                `gorm:"not null" json:"ID"`
	Make               *string             `json:"make"`
	CarModel           *string             `json:"car_model"`
	Year               *int                `json:"year"`
	LicensePlate       *string             `gorm:"not null" json:"license_plate"`
	SittingCapacity    *int                `json:"sitting_capacity"`
	Type               *string             `json:"type"`
	Color              *string             `json:"color"`
	VIN                *string             `gorm:"not null" json:"vin"`
	CurrentMileage     *int                `json:"current_mileage"`
	LastMaintenance    *time.Time          `json:"last_maintenance"`
	NextMaintenance    *time.Time          `json:"next_maintenance"`
	Status             *string             `gorm:"not null" json:"status"`
	AssignedDriver     *uint               `json:"assigned_driver"`
	Notes              *string             `json:"notes"`
	Driver             *User               `gorm:"foreignKey:AssignedDriver;references:ID"`
	AuctionVehicles    []AuctionVehicle    `gorm:"foreignKey:VehicleID"`
	MaintenanceRecords []MaintenanceRecord `gorm:"foreignKey:VehicleID"`
	FuelingRecords     []FuelingRecord     `gorm:"foreignKey:VehicleID"`
	VehicleUsages      []VehicleUsage      `gorm:"foreignKey:VehicleID"`
	gorm.Model
}

type Task struct {
	ID             uint       `gorm:"not null" json:"ID"`
	DriverID       *uint      `gorm:"not null;onDelete:CASCADE" json:"driver_id"`
	StartLatitude  *float64   `gorm:"not null" json:"start_latitude"`
	StartLongitude *float64   `gorm:"not null" json:"start_longitude"`
	EndLatitude    *float64   `gorm:"not null" json:"end_latitude"`
	EndLongitude   *float64   `gorm:"not null" json:"end_longitude"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	Status         *string    `gorm:"not null" json:"status"`
	Notes          *string    `json:"notes"`
	Driver         *User      `gorm:"foreignKey:DriverID;references:ID"`
	gorm.Model
}
type Appointment struct {
	gorm.Model
	AppointmentDate *time.Time         `json:"appointment_date"`
	Status          *AppointmentStatus `gorm:"not null" json:"status"`
	Notes           *string            `json:"notes"`
}

type AuctionVehicle struct {
	gorm.Model
	VehicleID *uint    `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	Images    []Image  `gorm:"foreignKey:ID" json:"images"`
	Details   *string  `json:"details"`
	Vehicle   *Vehicle `gorm:"foreignKey:VehicleID;references:ID"`
}

type MaintenanceRecord struct {
	gorm.Model
	VehicleID           *uint              `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	MaintenancePersonID *uint              `gorm:"not null" json:"maintenance_person_id"`
	MaintenanceDate     *time.Time         `json:"maintenance_date"`
	ServiceType         *string            `json:"service_type"`
	Status              *MaintenanceStatus `gorm:"not null" json:"status"`
	Parts               []Part             `gorm:"foreignKey:ID" json:"parts"`
	TotalCost           *float64           `gorm:"not null" json:"total_cost"`
	MileageAtService    *int               `json:"mileage_at_service"`
	Notes               *string            `json:"notes"`
	Vehicle             *Vehicle           `gorm:"foreignKey:VehicleID;references:ID"`
}

type FuelingRecord struct {
	gorm.Model
	VehicleID          *uint    `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	FuelingPersonID    *uint    `gorm:"not null" json:"fueling_person_id"`
	Amount             *float64 `gorm:"not null" json:"amount"`
	TotalCost          *float64 `gorm:"not null" json:"total_cost"`
	GasStation         *string  `json:"gas_station"`
	Notes              *string  `json:"notes"`
	BeforeFuelingImage *string  `gorm:"not null" json:"before_fueling_image"`
	AfterFuelingImage  *string  `gorm:"not null" json:"after_fueling_image"`
	Vehicle            *Vehicle `gorm:"foreignKey:VehicleID;references:ID"`
}
type VehicleUsage struct {
	gorm.Model
	VehicleID *uint      `gorm:"not null;onDelete:CASCADE" json:"vehicle_id"`
	StartTime *time.Time `gorm:"not null" json:"start_time"`
	EndTime   *time.Time `gorm:"not null" json:"end_time"`
	Distance  *float64   `gorm:"not null" json:"distance"`
	Vehicle   *Vehicle   `gorm:"foreignKey:VehicleID;references:ID"`
}

type RolePermission struct {
	Role                     RolesList `gorm:"primaryKey" json:"role"`
	CanAccessCarInfo         bool      `json:"can_access_car_info"`
	CanViewProfile           bool      `json:"can_view_profile"`
	CanViewDrivingHistory    bool      `json:"can_view_driving_history"`
	CanMenageUsers           bool      `json:"can_menage_users"`
	CanViewFuelingInfo       bool      `json:"can_view_fueling_info"`
	CanUpdateMaintenanceInfo bool      `json:"can_update_maintenance_info"`
	CanCreateAuction         bool      `json:"can_create_auction"`
	CanEditRouteDetails      bool      `json:"can_edit_route_details"`
	CanAssignVehicle         bool      `json:"can_assign_vehicle"`
	CanAssignTask            bool      `json:"can_assign_task"`
	CanGenerateReport        bool      `json:"can_generate_report"`
}

type Image struct {
	ID  uint    `gorm:"not null" json:"ID"`
	Url *string `gorm:"not null" json:"url"`
	gorm.Model
}
type Part struct {
	ID   uint    `gorm:"not null" json:"ID"`
	Name *string `gorm:"not null" json:"name"`
	gorm.Model
}
