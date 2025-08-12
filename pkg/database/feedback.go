package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseDoctorFeedback struct{}

func NewDatabaseDoctorFeedback() *DatabaseDoctorFeedback {
	return &DatabaseDoctorFeedback{}
}

func (*DatabaseDoctorFeedback) StoreDoctorFeedback(request *migration.DoctorFeedback) error {
	err := db.Gorm.Create(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabaseDoctorFeedback) PatchDoctorFeedback(id uint64, data *migration.DoctorFeedback) error {
	updated := &migration.DoctorFeedback{
		Response: data.Response,
		Approved: data.Approved,
	}

	err := db.Gorm.Model(&migration.DoctorFeedback{}).Where("id = ?", id).Updates(updated).Error
	if err != nil {
		return err
	}
	return nil
}

func (*DatabaseDoctorFeedback) DestroyDoctorFeedback(id uint64) error {
	return db.Gorm.Delete(&migration.DoctorFeedback{}, id).Error
}
