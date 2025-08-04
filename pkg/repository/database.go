package repository

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseUsersRepository interface {
	GetUserById(request *migration.User) (*migration.User, error)
	GetUserByEmail(request *migration.User) (*migration.User, error)
	GetUsers(request string, roles []string) ([]migration.User, error)
	RegisterUser(request *migration.User) error
	StoreUser(request *migration.User) error
	PatchUser(request string, data *migration.User) error
	DestroyUser(request string) error
	ActivateUser(request string, isActive bool) error
}

type DatabaseAuthRepository interface {
	GetRefreshToken(id string) (*migration.Token, error)
	StoreRefreshToken(token *migration.Token) error
	UpdateRefreshToken(token *migration.Token) error
}

type DatabasePatientsRepository interface {
	GetPatients(groupID uint) ([]migration.Patient, error)
	StorePatient(request *migration.Patient) error
	PatchPatient(request string, data *migration.Patient) error
	DestroyPatient(request string) error
}

type DatabaseRegistrations interface {
	GetRegistrations(groupID uint) ([]migration.Registration, error)
	StoreRegistration(request *migration.Registration) error
	PatchRegistration(request string, data *migration.Registration) error
	DestroyRegistration(request string) error
}
