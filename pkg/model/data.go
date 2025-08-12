package model

import (
	"strings"
	"time"
)

type Base struct {
	ID        string     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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
	Password string `json:"password,omitempty"`

	Role Roles `json:"role"`

	IsActive bool      `json:"is_active,omitempty"`
	Expired  time.Time `json:"expired"`

	Group *Group `json:"group,omitempty"`
}

type Patient struct {
	Base

	MedicalRecordNumber string   `json:"medical_record_number,omitempty"`
	Name                string   `json:"name"`
	Gender              string   `json:"gender"`
	BirthDate           DateOnly `json:"birth_date"`

	Group *Group `json:"group,omitempty"`
}

type MedicalRecord struct {
	Base

	MedicalRecordNumber string  `json:"medical_record_number,omitempty"`
	Patient             Patient `json:"patient"`

	Interrogator *User `json:"interrogator,omitempty"`

	Feedback *DoctorFeedback `json:"feedback,omitempty"`

	Diagnosis   string `json:"diagnosis,omitempty"`
	Predictions string `json:"predictions,omitempty"`
}

type DoctorFeedback struct {
	Base

	Interrogator *User `json:"interrogator,omitempty"`

	Response string `json:"response,omitempty"`
	Approved bool   `json:"approved"`
}

type Registration struct {
	Base

	RegistrationNumber string `json:"registration_number"`

	MedicalRecord MedicalRecord `json:"medical_record"`

	Group *Group `json:"group,omitempty"`
}

type ActivityLog struct {
	Base

	User User `json:"user"`

	Action  string `json:"action"`
	Details string `json:"details"`
}

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
