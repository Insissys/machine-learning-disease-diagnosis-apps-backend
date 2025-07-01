package repository

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

type DatabaseUsersRepository interface {
	GetUser(request *model.LoginRequest) (*model.User, error)
	StoreUser(request *model.RegisterRequest) error
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
