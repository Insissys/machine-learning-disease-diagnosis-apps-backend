package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseRegistrations struct{}

func NewDatabaseRegistrations() *DatabaseRegistrations {
	return &DatabaseRegistrations{}
}

func (*DatabaseRegistrations) GetRegistrations(groupID uint) ([]migration.Registration, error) {
	var patients []migration.Registration

	err := db.Gorm.
		Where("group_id = ?", groupID).
		Preload("Patient").
		Preload("Group").
		Order("created_at DESC").
		Find(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (*DatabaseRegistrations) StoreRegistration(request *migration.Registration) error {
	// birthdate, err := time.Parse("2006-01-02", request.BirthDate)
	// if err != nil {
	// 	return err
	// }

	// data := &migration.Registration{
	// 	MedicalRecordNumber: &request.MedicalRecord,
	// 	Name:                request.Name,
	// 	Gender:              request.Gender,
	// 	BirthDate:           birthdate,
	// 	GroupID:             request.GroupID,
	// }

	// err = db.Gorm.Debug().Create(&data).Error
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (*DatabaseRegistrations) PatchRegistration(request string, data *migration.Registration) error {
	// var patient *migration.Registration

	// birthdate, err := time.Parse("2006-01-02", data.BirthDate)
	// if err != nil {
	// 	return err
	// }

	// d := &migration.Registration{
	// 	MedicalRecordNumber: &data.MedicalRecord,
	// 	Name:                data.Name,
	// 	Gender:              data.Gender,
	// 	BirthDate:           birthdate,
	// 	GroupID:             data.GroupID,
	// }

	// err = db.Gorm.Debug().Where("id = ?", request).First(&patient).Error
	// if err != nil {
	// 	return err
	// }

	// err = db.Gorm.Debug().Model(patient).Updates(d).Error
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (*DatabaseRegistrations) DestroyRegistration(request string) error {
	err := db.Gorm.Debug().Delete(&migration.Registration{}, request).Error
	if err != nil {
		return err
	}
	return nil
}
