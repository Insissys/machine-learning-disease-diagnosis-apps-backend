package migration

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model

	UserID  uint   `gorm:"not null"`
	Action  string `gorm:"not null"` // e.g., "login", "update_record"
	Details string // Optional JSON or text
}

type User struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"` // e.g., "admin", "doctor"
}

type Patient struct {
	gorm.Model

	MedicalRecordID *string
	Name            string `gorm:"not null"`
	Gender          string `gorm:"not null"` // e.g., "male", "female", "other"
	Birthday        string `gorm:"not null"`
}

type MedicalRecord struct {
	gorm.Model

	PatientID    uint   `gorm:"not null"`
	Interrogator string `gorm:"not null"`
	Diagnosis    string `gorm:"not null"`
	Predictions  string `gorm:"not null"`
}

type DoctorFeedback struct {
	gorm.Model

	MedicalRecordID uint   `gorm:"not null;uniqueIndex"` // One feedback per record
	Interrogator    string `gorm:"not null"`             // FK to User (Role: doctor)
	Response        string `gorm:"not null"`             // Feedback text
	Approved        *bool  // Optional: doctor confirms system diagnosis
}
