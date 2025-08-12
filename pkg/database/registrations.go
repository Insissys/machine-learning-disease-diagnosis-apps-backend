package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseRegistrations struct{}

func NewDatabaseRegistrations() *DatabaseRegistrations {
	return &DatabaseRegistrations{}
}

func (*DatabaseRegistrations) GetRegistrations(groupID uint64) ([]migration.Registration, error) {
	var patients []migration.Registration

	err := db.Gorm.
		Joins("JOIN medical_records ON medical_records.id = registrations.medical_record_id AND medical_records.diagnosis = ?", "").
		Where("group_id = ?", groupID).
		Preload("MedicalRecord.Patient").
		Preload("MedicalRecord.Interrogator.Role").
		Preload("Group").
		Order("created_at DESC").
		Find(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (*DatabaseRegistrations) StoreRegistration(request *migration.Registration) error {
	err := db.Gorm.Create(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabaseRegistrations) PatchRegistration(id uint64, data *migration.Registration) error {
	// TODO: Coming Soon
	return nil
}

func (*DatabaseRegistrations) DestroyRegistration(id uint64) error {
	return db.Gorm.Delete(&migration.Registration{}, id).Error
}
