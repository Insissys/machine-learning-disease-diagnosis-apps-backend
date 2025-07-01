package seed

import (
	"log"
	"time"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
)

func UserGroupSeederDev() {
	// Seed groups (clinics/hospitals)
	groups := []migration.Group{
		{Name: "Mitra Sehat Hospital", Address: "Jl. Kesehatan No.1"},
	}

	if err := db.Gorm.Create(&groups).Error; err != nil {
		log.Println("Failed to seed groups:", err)
		return
	}

	log.Println("Groups seeded successfully")

	roles := []migration.Roles{
		{Name: "superadmin"},
		{Name: "admin"},
		{Name: "doctor"},
	}

	if err := db.Gorm.Create(&roles).Error; err != nil {
		log.Println("Failed to seed roles:", err)
		return
	}

	log.Println("Roles seeded successfully")

	active := true
	users := []migration.User{
		{
			Name:     "Super Admin",
			Email:    "superadmin@example.com",
			Password: utils.HashPassword("supersecret"),
			Role:     roles[0],
			IsActive: &active,
			GroupID:  groups[0].ID,
			Expired:  time.Now().Add(7 * 24 * time.Hour),
		},
		{
			Name:     "Admin One",
			Email:    "admin@example.com",
			Password: utils.HashPassword("admin123"),
			Role:     roles[1],
			IsActive: &active,
			GroupID:  groups[0].ID,
			Expired:  time.Now().Add(7 * 24 * time.Hour),
		},
		{
			Name:     "Doctor Strange",
			Email:    "doctor@example.com",
			Password: utils.HashPassword("doctorpass"),
			Role:     roles[2],
			IsActive: &active,
			GroupID:  groups[0].ID,
			Expired:  time.Now().Add(7 * 24 * time.Hour),
		},
	}

	if err := db.Gorm.Create(&users).Error; err != nil {
		log.Println("Failed to seed users:", err)
	} else {
		log.Println("Users seeded successfully")
	}
}
