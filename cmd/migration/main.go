package main

import (
	"log"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	model "github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()

	db.MysqlConnection(&cfg.Config.Database)
	database := db.Gorm.Migrator()

	// Auto migrate models
	err := database.AutoMigrate(
		&model.ActivityLog{},
		&model.User{},
		&model.Patient{},
		&model.MedicalRecord{},
		&model.DoctorFeedback{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration success")
}
