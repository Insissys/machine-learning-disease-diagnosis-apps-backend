package repository

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseUsersRepository interface {
	GetUserById(request *migration.User) (*migration.User, error)
	GetUserByEmail(request *migration.User) (*migration.User, error)
	GetUsers(groupId uint64, roles []string) ([]migration.User, error)
	RegisterUser(request *migration.User) error
	StoreUser(request *migration.User) error
	PatchUser(id uint64, data *migration.User) error
	DestroyUser(id uint64) error
	ActivateUser(id uint64, isActive bool) error
	GetUserGroup(request *migration.User) (*migration.Group, error)
}

type DatabaseAuthRepository interface {
	GetRefreshToken(id string) (*migration.Token, error)
	StoreRefreshToken(token *migration.Token) error
	UpdateRefreshToken(token *migration.Token) error
}

type DatabasePatientsRepository interface {
	GetPatients(groupID uint64) ([]migration.Patient, error)
	GetPatientById(id uint64) (*migration.Patient, error)
	StorePatient(request *migration.Patient) error
	PatchPatient(id uint64, data *migration.Patient) error
	DestroyPatient(id uint64) error
}

type DatabaseRegistrations interface {
	GetRegistrations(groupID uint64) ([]migration.Registration, error)
	StoreRegistration(request *migration.Registration) error
	PatchRegistration(id uint64, data *migration.Registration) error
	DestroyRegistration(id uint64) error
}

type DatabaseMedicalRecordRepository interface {
	GetMedicalRecords(patientId uint64) ([]migration.MedicalRecord, error)
	StoreMedicalRecord(request *migration.MedicalRecord) error
	PatchMedicalRecord(id uint64, data *migration.MedicalRecord) error
	DestroyMedicalRecord(id uint64) error
}

type DatabaseDoctorFeedbackRepository interface {
	StoreDoctorFeedback(request *migration.DoctorFeedback) error
	PatchDoctorFeedback(id uint64, data *migration.DoctorFeedback) error
	DestroyDoctorFeedback(id uint64) error
}

type DatabaseQueueRepository interface {
	GetQueue(userId, groupID *uint64) ([]migration.Registration, error)
}
