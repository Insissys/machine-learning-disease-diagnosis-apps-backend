package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseMedicalRecord struct{}

func NewDatabaseMedicalRecord() *DatabaseMedicalRecord {
	return &DatabaseMedicalRecord{}
}

func (*DatabaseMedicalRecord) GetMedicalRecords(patientId uint64) ([]migration.MedicalRecord, error) {
	var medicalrecords []migration.MedicalRecord

	err := db.Gorm.
		Where("patient_id = ?", patientId).
		Preload("Patient").
		Preload("Interrogator").
		Preload("Feedback.Interrogator").
		Order("created_at DESC").
		Find(&medicalrecords).Error

	if err != nil {
		return nil, err
	}

	return medicalrecords, nil
}

func (*DatabaseMedicalRecord) StoreMedicalRecord(request *migration.MedicalRecord) error {
	err := db.Gorm.Create(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabaseMedicalRecord) PatchMedicalRecord(id uint64, data *migration.MedicalRecord) error {
	updated := &migration.MedicalRecord{
		Diagnosis:   data.Diagnosis,
		Predictions: data.Predictions,
	}

	err := db.Gorm.Model(&migration.MedicalRecord{}).Where("id = ?", id).Updates(updated).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabaseMedicalRecord) DestroyMedicalRecord(id uint64) error {
	return db.Gorm.Delete(&migration.MedicalRecord{}, id).Error
}
