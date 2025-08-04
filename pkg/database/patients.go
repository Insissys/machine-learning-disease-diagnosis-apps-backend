package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabasePatients struct{}

func NewDatabasePatients() *DatabasePatients {
	return &DatabasePatients{}
}

func (*DatabasePatients) GetPatients(groupID uint) ([]migration.Patient, error) {
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

func (*DatabasePatients) StorePatient(request *migration.Patient) error {
	data := &migration.Patient{
		Name:      request.Name,
		Gender:    request.Gender,
		BirthDate: request.BirthDate,
		GroupID:   request.GroupID,
	}

	err := db.Gorm.Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) PatchPatient(request string, data *migration.Patient) error {
	var patient *migration.Patient

	d := &migration.Patient{
		Name:      data.Name,
		Gender:    data.Gender,
		BirthDate: data.BirthDate,
		GroupID:   data.GroupID,
	}

	err := db.Gorm.Where("id = ?", request).First(&patient).Error
	if err != nil {
		return err
	}

	err = db.Gorm.Model(patient).Updates(d).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) DestroyPatient(request string) error {
	err := db.Gorm.Delete(&migration.Patient{}, request).Error
	if err != nil {
		return err
	}
	return nil
}
