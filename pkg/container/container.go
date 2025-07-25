package container

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/database"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/repository"
)

type Container struct {
	Users         repository.DatabaseUsersRepository
	Auth          repository.DatabaseAuthRepository
	Patients      repository.DatabasePatientsRepository
	Registrations repository.DatabaseRegistrations
}

func NewContainer() *Container {
	users := database.NewDatabaseUsers()
	auth := database.NewDatabaseAuth()
	patients := database.NewDatabasePatients()
	registrations := database.NewDatabaseRegistrations()

	return &Container{
		Users:         users,
		Auth:          auth,
		Patients:      patients,
		Registrations: registrations,
	}
}
