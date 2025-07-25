package database

import (
	"time"

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
	birthdate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
		return err
	}

	data := &migration.Patient{
		MedicalRecordNumber: &request.MedicalRecord,
		Name:                request.Name,
		Gender:              request.Gender,
		BirthDate:           birthdate,
		GroupID:             request.GroupID,
	}

	err = db.Gorm.Debug().Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) PatchPatient(request string, data *migration.Patient) error {
	var patient *migration.Patient

	birthdate, err := time.Parse("2006-01-02", data.BirthDate)
	if err != nil {
		return err
	}

	d := &migration.Patient{
		MedicalRecordNumber: &data.MedicalRecord,
		Name:                data.Name,
		Gender:              data.Gender,
		BirthDate:           birthdate,
		GroupID:             data.GroupID,
	}

	err = db.Gorm.Debug().Where("id = ?", request).First(&patient).Error
	if err != nil {
		return err
	}

	err = db.Gorm.Debug().Model(patient).Updates(d).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) DestroyPatient(request string) error {
	err := db.Gorm.Debug().Delete(&migration.Patient{}, request).Error
	if err != nil {
		return err
	}
	return nil
}
