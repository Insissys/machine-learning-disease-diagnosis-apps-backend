package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
	"gorm.io/gorm"
)

// GET /api/patient/registration
func GetRegistrations(c *gin.Context) {
	groupID := c.MustGet("groupId").(uint64)

	database := container.NewContainer()
	data, err := database.Registrations.GetRegistrations(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went Wrong"})
		return
	}

	response := []model.Registration{}

	for _, v := range data {
		response = append(response, model.Registration{
			Base: model.Base{
				ID:        utils.EncryptUint64(uint64(v.ID)),
				CreatedAt: &v.CreatedAt,
				UpdatedAt: &v.UpdatedAt,
			},
			RegistrationNumber: v.RegistrationNumber,
			MedicalRecord: model.MedicalRecord{
				Base: model.Base{
					ID:        utils.EncryptUint64(uint64(v.MedicalRecord.ID)),
					CreatedAt: &v.MedicalRecord.CreatedAt,
					UpdatedAt: &v.MedicalRecord.UpdatedAt,
				},
				Patient: model.Patient{
					Base: model.Base{
						ID:        utils.EncryptUint64(uint64(v.MedicalRecord.Patient.ID)),
						CreatedAt: &v.MedicalRecord.Patient.CreatedAt,
						UpdatedAt: &v.MedicalRecord.Patient.UpdatedAt,
					},
					MedicalRecordNumber: v.MedicalRecord.Patient.MedicalRecordNumber,
					Name:                v.MedicalRecord.Patient.Name,
					Gender:              v.MedicalRecord.Patient.Gender,
					BirthDate:           model.DateOnly{Time: v.MedicalRecord.Patient.BirthDate},
				},
				Interrogator: &model.User{
					Name:  v.MedicalRecord.Interrogator.Name,
					Email: v.MedicalRecord.Interrogator.Email,
					Role: model.Roles{
						Name: v.MedicalRecord.Interrogator.Role.Name,
					},
					IsActive: v.MedicalRecord.Interrogator.IsActive,
					Expired:  v.MedicalRecord.Interrogator.Expired,
				},
				MedicalRecordNumber: v.MedicalRecord.MedicalRecordNumber,
			},
			Group: &model.Group{
				Base: model.Base{
					ID:        utils.EncryptUint64(uint64(v.Group.ID)),
					CreatedAt: &v.Group.CreatedAt,
					UpdatedAt: &v.Group.UpdatedAt,
				},
				Name:    v.Group.Name,
				Address: v.Group.Address,
			},
		})
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Retrive Registration",
		Data:    response,
	})
}

// POST /api/patient/registration
func StoreRegistration(c *gin.Context) {
	var input model.Registration
	var patient *migration.Patient
	var err error
	groupId := c.MustGet("groupId").(uint64)

	// 1. Validate input binding
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid input"})
		return
	}

	database := container.NewContainer()

	// 2. Check if patient exist
	if input.MedicalRecord.Patient.ID != "" {
		patient, err = database.Patients.GetPatientById(utils.DecryptToUint64(input.MedicalRecord.Patient.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Patient not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, model.ApiResponse{
				Message:  "Something went wrong",
				Error:    err,
				Data:     nil,
				Metadata: nil,
			})
			return
		}
	} else {
		patient = &migration.Patient{
			MedicalRecordNumber: input.MedicalRecord.MedicalRecordNumber,
			Name:                input.MedicalRecord.Patient.Name,
			Gender:              input.MedicalRecord.Patient.Gender,
			BirthDate:           input.MedicalRecord.Patient.BirthDate.Time,
			GroupID:             uint(groupId),
		}

		err := database.Patients.StorePatient(patient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ApiResponse{
				Message:  "Something went wrong",
				Error:    err,
				Data:     nil,
				Metadata: nil,
			})
			return
		}
	}

	// 3. Check if user exist
	userdoctor, err := database.Users.GetUserById(&migration.User{
		Model: gorm.Model{ID: uint(utils.DecryptToUint64(input.MedicalRecord.Interrogator.ID))},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			database.Patients.DestroyPatient(uint64(patient.ID))
			c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "User not found"})
			return
		}
		database.Patients.DestroyPatient(uint64(patient.ID))
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	// 4. Input New Feedback
	feedback := &migration.DoctorFeedback{
		InterrogatorID: userdoctor.ID,
		Approved:       false,
	}

	err = database.Feedback.StoreDoctorFeedback(feedback)
	if err != nil {
		database.Patients.DestroyPatient(uint64(patient.ID))
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	// 5. Input New Medical Record
	medicalrecord := &migration.MedicalRecord{
		MedicalRecordNumber: input.MedicalRecord.MedicalRecordNumber,
		PatientID:           patient.ID,
		InterrogatorID:      userdoctor.ID,
		FeedbackID:          feedback.ID,
	}

	err = database.MedicalRecords.StoreMedicalRecord(medicalrecord)
	if err != nil {
		database.Patients.DestroyPatient(uint64(patient.ID))
		database.Feedback.DestroyDoctorFeedback(uint64(feedback.ID))
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	// 6. Input Registration
	group, err := database.Users.GetUserGroup(&migration.User{GroupID: uint(groupId)})
	if err != nil {
		database.Patients.DestroyPatient(uint64(patient.ID))
		database.Feedback.DestroyDoctorFeedback(uint64(feedback.ID))
		database.MedicalRecords.DestroyMedicalRecord(uint64(medicalrecord.ID))
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	regisnumber := utils.GenerateRegisterNumber(group.Name)

	registration := &migration.Registration{
		RegistrationNumber: regisnumber,
		MedicalRecordID:    medicalrecord.ID,
		GroupID:            uint(groupId),
	}

	err = database.Registrations.StoreRegistration(registration)
	if err != nil {
		database.Patients.DestroyPatient(uint64(patient.ID))
		database.Feedback.DestroyDoctorFeedback(uint64(feedback.ID))
		database.MedicalRecords.DestroyMedicalRecord(uint64(medicalrecord.ID))
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	// 7. Return the created registration or status
	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Add Registration",
	})
}

// PATCH /api/patient/registration/:id
func PatchRegistration(c *gin.Context) {
	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Add Registration",
	})
}

// DELETE /api/patient/registration/:id
func DestroyRegistration(c *gin.Context) {
	encryptedBase64 := c.Param("id")
	id := utils.DecryptToUint64(encryptedBase64)

	database := container.NewContainer()
	if err := database.Registrations.DestroyRegistration(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Something went wrong",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{Message: "Registration deleted"})
}
