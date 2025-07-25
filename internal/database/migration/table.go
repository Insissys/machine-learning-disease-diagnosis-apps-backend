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
		&DoctorFeedback{},
		&Registration{},
	)

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

	UserID uint `gorm:"not null"`
	User   User

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

	IsActive *bool     `gorm:"not null"`
	Expired  time.Time `gorm:"not null;size:50"`

	GroupID uint `gorm:"not null"`
	Group   Group
}

type Patient struct {
	gorm.Model

	MedicalRecordNumber *string   `gorm:"uniqueIndex;size:50"`
	Name                string    `gorm:"not null;size:100"`
	Gender              string    `gorm:"not null;size:10"`
	BirthDate           time.Time `gorm:"not null;size:50"`

	GroupID uint `gorm:"not null"`
	Group   Group
}

type MedicalRecord struct {
	gorm.Model

	PatientID uint `gorm:"not null"`
	Patient   Patient

	InterrogatorID uint `gorm:"not null"`
	Interrogator   User `gorm:"foreignKey:InterrogatorID"`

	Diagnosis   string `gorm:"not null;size:1000"`
	Predictions string `gorm:"not null;size:1000"`
}

type DoctorFeedback struct {
	gorm.Model

	MedicalRecordID uint `gorm:"not null;uniqueIndex"`
	MedicalRecord   MedicalRecord

	InterrogatorID uint `gorm:"not null"`
	Interrogator   User `gorm:"foreignKey:InterrogatorID"`

	Response string `gorm:"not null;size:1000"`
	Approved *bool
}

type Registration struct {
	gorm.Model

	RegistrationNumber *string `gorm:"uniqueIndex;size:50"`
	PatientID          uint    `gorm:"not null"`
	Patient            Patient
	GroupID            uint `gorm:"not null"`
	Group              Group
}
