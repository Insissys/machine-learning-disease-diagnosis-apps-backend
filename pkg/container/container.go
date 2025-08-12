package container

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/database"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/repository"
)

type Container struct {
	Users          repository.DatabaseUsersRepository
	Auth           repository.DatabaseAuthRepository
	Patients       repository.DatabasePatientsRepository
	Registrations  repository.DatabaseRegistrations
	MedicalRecords repository.DatabaseMedicalRecordRepository
	Feedback       repository.DatabaseDoctorFeedbackRepository
	Queue          repository.DatabaseQueueRepository
}

func NewContainer() *Container {
	users := database.NewDatabaseUsers()
	auth := database.NewDatabaseAuth()
	patients := database.NewDatabasePatients()
	registrations := database.NewDatabaseRegistrations()
	medicalrecords := database.NewDatabaseMedicalRecord()
	feedback := database.NewDatabaseDoctorFeedback()
	queue := database.NewDatabaseQueue()

	return &Container{
		Users:          users,
		Auth:           auth,
		Patients:       patients,
		Registrations:  registrations,
		MedicalRecords: medicalrecords,
		Feedback:       feedback,
		Queue:          queue,
	}
}
