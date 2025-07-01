package migration

import (
	"time"

	"gorm.io/gorm"
)

func GetTable() []any {
	var table []any

	table = append(table,
		&Group{},
		&Token{},
		&ActivityLog{},
		&User{},
		&Roles{},
		&Patient{},
		&MedicalRecord{},
		&DoctorFeedback{})

	return table
}

type Token struct {
	ID      string    `gorm:"not null;size:50"`
	User    string    `gorm:"not null;size:100"`
	Expired time.Time `gorm:"not null;size:50"`
	Revoked bool      `gorm:"not null;"`
}

type Group struct {
	gorm.Model

	Name    string `gorm:"not null;size:100"`
	Address string `gorm:"size:255"`
}

type ActivityLog struct {
	gorm.Model

	UserID  uint   `gorm:"not null"`
	Action  string `gorm:"not null;size:50"`
	Details string `gorm:"size:1000"`
}

type Roles struct {
	gorm.Model

	Name string `gorm:"not null;size:50"`
}

type User struct {
	gorm.Model

	Name     string `gorm:"not null;size:100"`
	Email    string `gorm:"uniqueIndex;not null;size:100"`
	Password string `gorm:"not null;size:255"`

	RoleID uint `gorm:"not null"`
	Role   Roles

	IsActive *bool     `gorm:"not null;size:1"`
	Expired  time.Time `gorm:"not null;size:50"`

	GroupID uint `gorm:"not null"`
	Group   Group
}

type Patient struct {
	gorm.Model

	MedicalRecord *string `gorm:"uniqueIndex;size:50"`
	Name          string  `gorm:"not null;size:100"`
	Gender        string  `gorm:"not null;size:10"`
	BirthDate     string  `gorm:"not null;size:10"`

	GroupID uint `gorm:"not null"`
	Group   Group
}

type MedicalRecord struct {
	gorm.Model

	PatientID    uint   `gorm:"not null"`
	Interrogator string `gorm:"not null;size:100"`
	Diagnosis    string `gorm:"not null;size:1000"`
	Predictions  string `gorm:"not null;size:1000"`
}

type DoctorFeedback struct {
	gorm.Model

	MedicalRecordID uint   `gorm:"not null;uniqueIndex"`
	Interrogator    string `gorm:"not null;size:100"`
	Response        string `gorm:"not null;size:1000"`
	Approved        *bool
}
