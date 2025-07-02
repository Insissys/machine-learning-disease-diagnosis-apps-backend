package repository

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

type DatabaseUsersRepository interface {
	GetUser(request *model.LoginRequest) (*model.User, error)
	GetUsers(request string) ([]model.User, error)
	RegisterUser(request *model.RegisterRequest) error
	StoreUser(request model.User) error
	PatchUser(request string, data model.User) error
	DestroyUser(request string) error
	ActivateUser(request string, isActive bool) error
}

type DatabaseAuthRepository interface {
	GetRefreshToken(id string) (*migration.Token, error)
	StoreRefreshToken(token model.CustomClaims) error
	UpdateRefreshToken(token *migration.Token) error
}

type DatabasePatientsRepository interface {
	GetPatients(groupID uint) ([]model.Patient, error)
	StorePatient(request model.Patient) error
	PatchPatient(request string, data model.Patient) error
	DestroyPatient(request string) error
}
