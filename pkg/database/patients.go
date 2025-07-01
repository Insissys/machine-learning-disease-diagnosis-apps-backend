package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

type DatabasePatients struct{}

func NewDatabasePatients() *DatabasePatients {
	return &DatabasePatients{}
}

func (*DatabasePatients) GetPatients(groupID uint) ([]model.Patient, error) {
	var patients []migration.Patient

	err := db.Gorm.
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Find(&patients).Error

	if err != nil {
		return nil, err
	}

	var response []model.Patient

	for _, v := range patients {
		response = append(response, model.Patient{
			ID:            v.ID,
			MedicalRecord: *v.MedicalRecord,
			Name:          v.Name,
			Gender:        v.Gender,
			BirthDate:     v.BirthDate,
			GroupID:       groupID,
		})
	}

	return response, nil
}

func (*DatabasePatients) StorePatient(request model.Patient) error {
	data := &migration.Patient{
		MedicalRecord: &request.MedicalRecord,
		Name:          request.Name,
		Gender:        request.Gender,
		BirthDate:     request.BirthDate,
		GroupID:       request.GroupID,
	}

	err := db.Gorm.Debug().Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabasePatients) PatchPatient(request string, data model.Patient) error {
	var patient *migration.Patient

	d := &migration.Patient{
		MedicalRecord: &data.MedicalRecord,
		Name:          data.Name,
		Gender:        data.Gender,
		BirthDate:     data.BirthDate,
		GroupID:       data.GroupID,
	}

	err := db.Gorm.Debug().Where("id = ?", request).First(&patient).Error
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
