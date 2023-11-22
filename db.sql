CREATE TYPE task_status AS ENUM ('completed', 'canceled', 'delayed');
CREATE TYPE vehicle_status AS ENUM ('Active', 'Inactive', 'Maintenance');
CREATE TYPE appointment_status AS ENUM ('Pending', 'Confirmed', 'Cancelled');
CREATE TYPE auction_status AS ENUM ('Sold', 'Pending');
CREATE TYPE maintenance_status AS ENUM ('Pending', 'Done');
CREATE TYPE roles_list AS ENUM ('Admin', 'Driver', 'Fueling_person', 'Maintenance_person');



CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "hashed_password" VARCHAR(255) NOT NULL,
  "government_id" VARCHAR(255) NOT NULL,
  "middle_name" VARCHAR(255),
  "address" TEXT,
  "phone_number" VARCHAR(50),
  "driving_license_code" VARCHAR(50),
  "role" roles_list NOT NULL,
  "first_name" VARCHAR(255),
  "last_name" VARCHAR(255),
  "email" VARCHAR(255) UNIQUE,
  "registration_date" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "last_login" TIMESTAMP,
  "status" VARCHAR(50) DEFAULT 'Active',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "roles_permissions" (
  "role" roles_list PRIMARY KEY,
  "can_access_car_info" BOOLEAN DEFAULT FALSE,
  "can_view_own_profile" BOOLEAN DEFAULT FALSE,
  "can_view_driving_history" BOOLEAN DEFAULT FALSE,
  "can_manage_users" BOOLEAN DEFAULT FALSE,
  "can_view_fueling_info" BOOLEAN DEFAULT FALSE,
  "can_update_maintenance_details" BOOLEAN DEFAULT FALSE,
  "can_search_by_license_plate" BOOLEAN DEFAULT FALSE,
  "can_create_auction_vehicles" BOOLEAN DEFAULT FALSE,
  "can_view_auction_page" BOOLEAN DEFAULT FALSE,
  "can_edit_route_details" BOOLEAN DEFAULT FALSE,
  "can_assign_vehicle_to_driver" BOOLEAN DEFAULT FALSE,
  "can_assign_task_to_driver" BOOLEAN DEFAULT FALSE,
  "can_generate_reports" BOOLEAN DEFAULT FALSE
);
INSERT INTO roles_permissions (role, can_access_car_info, can_view_own_profile, can_view_driving_history, can_manage_users, can_view_fueling_info, can_update_maintenance_details, can_search_by_license_plate, can_create_auction_vehicles, can_view_auction_page, can_edit_route_details, can_assign_vehicle_to_driver, can_assign_task_to_driver, can_generate_reports)
VALUES ('Admin', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);

INSERT INTO roles_permissions (role, can_access_car_info, can_view_own_profile, can_view_driving_history, can_manage_users, can_view_fueling_info, can_update_maintenance_details, can_search_by_license_plate, can_create_auction_vehicles, can_view_auction_page, can_edit_route_details, can_assign_vehicle_to_driver, can_assign_task_to_driver, can_generate_reports)
VALUES ('Driver', TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, TRUE, FALSE, TRUE, FALSE);

INSERT INTO roles_permissions (role, can_access_car_info, can_view_own_profile, can_view_driving_history, can_manage_users, can_view_fueling_info, can_update_maintenance_details, can_search_by_license_plate, can_create_auction_vehicles, can_view_auction_page, can_edit_route_details, can_assign_vehicle_to_driver, can_assign_task_to_driver, can_generate_reports)
VALUES ('Fueling_person', FALSE, TRUE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE);

INSERT INTO roles_permissions (role, can_access_car_info, can_view_own_profile, can_view_driving_history, can_manage_users, can_view_fueling_info, can_update_maintenance_details, can_search_by_license_plate, can_create_auction_vehicles, can_view_auction_page, can_edit_route_details, can_assign_vehicle_to_driver, can_assign_task_to_driver, can_generate_reports)
VALUES ('Maintenance_person', TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE);


CREATE TABLE "vehicles" (
  "id" SERIAL PRIMARY KEY,
  "make" VARCHAR(255) NOT NULL,
  "model" VARCHAR(255) NOT NULL,
  "year" INTEGER,
  "license_plate" VARCHAR(255) UNIQUE NOT NULL,
  "sitting_capacity" INTEGER DEFAULT 5,
  "type" VARCHAR(255),
  "color" VARCHAR(50),
  "VIN" VARCHAR(255) UNIQUE,
  "current_mileage" INTEGER DEFAULT 0,
  "last_maintenance_date" DATE,
  "next_scheduled_maintenance_mileage" INTEGER,
  "next_scheduled_maintenance_date" DATE,
  "status" vehicle_status,
  "assigned_driver_id" INTEGER,
  "registration_date" DATE DEFAULT (CURRENT_DATE),
  "notes" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "vehicle_usage" (
  "id" SERIAL PRIMARY KEY,
  "vehicle_id" INTEGER NOT NULL,
  "start_time" TIMESTAMP,
  "end_time" TIMESTAMP,
  "distance_traveled" INTEGER,
  "route_description" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "appointments" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "appointment_date" TIMESTAMP NOT NULL,
  "status" appointment_status,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "auction_vehicles" (
  "id" SERIAL PRIMARY KEY,
  "vehicle_id" INTEGER NOT NULL,
  "images_path" BYTEA[],
  "status" auction_status,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "maintenance_records" (
  "id" SERIAL PRIMARY KEY,
  "vehicle_id" INTEGER NOT NULL,
  "maintenance_person_id" INTEGER NOT NULL,
  "maintenance_date" TIMESTAMP,
  "activity_type" VARCHAR(50),
  "service_type" VARCHAR(255),
  "status" maintenance_status,
  "parts_list" TEXT[],
  "total_cost" DECIMAL(10,2),
  "mileage_at_service" INTEGER,
  "next_scheduled_maintenance_mileage" INTEGER,
  "next_scheduled_maintenance_date" DATE,
  "notes" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "fueling_records" (
  "id" SERIAL PRIMARY KEY,
  "vehicle_id" INTEGER NOT NULL,
  "fueling_person_id" INTEGER NOT NULL,
  "fueling_date" TIMESTAMP,
  "amount_of_fuel" DECIMAL(10,2),
  "total_cost" DECIMAL(10,2),
  "gas_station_name" VARCHAR(255),
  "notes" TEXT,
  "before_fueling_image" BYTEA,
  "after_fueling_image" BYTEA,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "tasks" (
  "id" SERIAL PRIMARY KEY,
  "driver_id" INTEGER NOT NULL,
  "start_latitude" DECIMAL(9,6),
  "start_longitude" DECIMAL(9,6),
  "end_latitude" DECIMAL(9,6),
  "end_longitude" DECIMAL(9,6),
  "start_time" TIMESTAMP,
  "end_time" TIMESTAMP,
  "status" task_status,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);


ALTER TABLE "vehicles" ADD FOREIGN KEY ("assigned_driver_id") REFERENCES "users" ("id");

ALTER TABLE "vehicle_usage" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicles" ("id");

ALTER TABLE "appointments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "auction_vehicles" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicles" ("id");

ALTER TABLE "maintenance_records" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicles" ("id");

ALTER TABLE "maintenance_records" ADD FOREIGN KEY ("maintenance_person_id") REFERENCES "users" ("id");

ALTER TABLE "fueling_records" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicles" ("id");

ALTER TABLE "fueling_records" ADD FOREIGN KEY ("fueling_person_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("driver_id") REFERENCES "users" ("id");
