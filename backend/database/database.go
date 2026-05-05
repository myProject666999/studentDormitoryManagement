package database

import (
	"log"
	"student-dormitory-management/models"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open("dormitory_management.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Student{},
		&models.MaintenanceWorker{},
		&models.Dormitory{},
		&models.DormitoryAssignment{},
		&models.Notice{},
		&models.RepairRequest{},
		&models.RepairRecord{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	createDefaultAdmin()

	log.Println("Database initialized successfully")
}

func createDefaultAdmin() {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	
	if count == 0 {
		defaultAdmin := &models.Admin{
			User: models.User{
				Username: "admin",
				Password: "admin123",
				Role:     "admin",
				Status:   1,
			},
			Name:  "系统管理员",
			Phone: "13800138000",
			Email: "admin@example.com",
		}
		
		defaultAdmin.BeforeCreate(nil)
		
		if err := DB.Create(defaultAdmin).Error; err != nil {
			log.Printf("Failed to create default admin: %v", err)
		} else {
			log.Println("Default admin created: username=admin, password=admin123")
		}
	}
}
