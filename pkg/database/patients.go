package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabasePatients struct{}

func NewDatabasePatients() *DatabasePatients {
	return &DatabasePatients{}
}

func (*DatabasePatients) GetPatients(groupID uint64) ([]migration.Patient, error) {
	var patients []migration.Patient

	err := db.Gorm.
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Preload("Group").
		Find(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (*DatabasePatients) GetPatientById(id uint64) (*migration.Patient, error) {
	var patient *migration.Patient

	err := db.Gorm.
		Where("id = ?", id).
		Preload("Group").
		Find(&patient).Error

	if err != nil {
		return nil, err
	}

	return patient, nil
}

func (*DatabasePatients) StorePatient(request *migration.Patient) error {
	err := db.Gorm.Create(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) PatchPatient(id uint64, data *migration.Patient) error {
	updated := &migration.Patient{
		Name:      data.Name,
		Gender:    data.Gender,
		BirthDate: data.BirthDate,
	}

	err := db.Gorm.Model(&migration.Patient{}).Where("id = ?", id).Updates(updated).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) DestroyPatient(id uint64) error {
	return db.Gorm.Delete(&migration.Patient{}, id).Error
}
