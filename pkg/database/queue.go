package database

import (
	"strings"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseQueue struct{}

func NewDatabaseQueue() *DatabaseQueue {
	return &DatabaseQueue{}
}

func (*DatabaseQueue) GetQueue(userId, groupID *uint64) ([]migration.Registration, error) {
	var patients []migration.Registration

	var builder strings.Builder
	var param []any

	builder.WriteString("JOIN medical_records ON medical_records.id = registrations.medical_record_id ")
	builder.WriteString("AND medical_records.diagnosis = ? ")
	param = append(param, "")

	if userId != nil {
		builder.WriteString("AND medical_records.interrogator_id = ? ")
		param = append(param, *userId)
	}

	query := db.Gorm.
		Joins(builder.String(), param...).
		Where("group_id = ?", groupID).
		Preload("MedicalRecord.Patient").
		Preload("MedicalRecord.Interrogator.Role").
		Preload("MedicalRecord.Feedback").
		Preload("Group").
		Order("created_at DESC")

	err := query.Find(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}
