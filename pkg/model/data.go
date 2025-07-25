package model

import "time"

type Base struct {
	ID        uint      `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Group struct {
	Base

	Name    string `json:"name"`
	Address string `json:"address"`
}

type Roles struct {
	Base

	Name string `json:"name"`
}

type User struct {
	Base

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Role Roles `json:"role"`

	IsActive bool      `json:"is_active,omitempty"`
	Expired  time.Time `json:"expired"`

	Group Group `json:"group"`
}

type Patient struct {
	Base

	MedicalRecordNumber string    `json:"medical_record_number,omitempty"`
	Name                string    `json:"name"`
	Gender              string    `json:"gender"`
	BirthDate           time.Time `json:"birth_date"`

	Group Group `json:"group"`
}

type MedicalRecord struct {
	Base

	Patient Patient `json:"patient"`

	Interrogator User `json:"interrogator"`

	Diagnosis   string `json:"diagnosis"`
	Predictions string `json:"predictions"`
}

type DoctorFeedback struct {
	Base

	MedicalRecord MedicalRecord `json:"medical_record"`

	Interrogator User `json:"interrogator"`

	Response string `json:"response"`
	Approved bool   `json:"approved,omitempty"`
}

type Registration struct {
	Base

	RegistrationNumber string `json:"registration_number,omitempty"`

	Patient Patient `json:"patient"`

	Group Group `json:"group"`
}

type ActivityLog struct {
	Base

	User User `json:"user"`

	Action  string `json:"action"`
	Details string `json:"details"`
}
