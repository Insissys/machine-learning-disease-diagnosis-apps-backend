package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
)

// GET /api/patient/medicalrecord
func GetMedicalRecords(c *gin.Context) {
	encryptedBase64 := c.Query("id")
	id := utils.DecryptToUint64(encryptedBase64)

	database := container.NewContainer()
	data, err := database.MedicalRecords.GetMedicalRecords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went Wrong"})
		return
	}

	response := []model.MedicalRecord{}

	for _, v := range data {
		response = append(response, model.MedicalRecord{
			Base: model.Base{
				ID:        utils.EncryptUint64(uint64(v.ID)),
				CreatedAt: &v.CreatedAt,
				UpdatedAt: &v.UpdatedAt,
			},
			MedicalRecordNumber: v.MedicalRecordNumber,
			Patient: model.Patient{
				Base: model.Base{
					ID:        utils.EncryptUint64(uint64(v.Patient.ID)),
					CreatedAt: &v.Patient.CreatedAt,
					UpdatedAt: &v.Patient.UpdatedAt,
				},
				MedicalRecordNumber: v.Patient.MedicalRecordNumber,
				Name:                v.Patient.Name,
				Gender:              v.Patient.Gender,
				BirthDate:           model.DateOnly{Time: v.Patient.BirthDate},
			},
			Interrogator: &model.User{
				Name:  v.Interrogator.Name,
				Email: v.Interrogator.Email,
				Role: model.Roles{
					Name: v.Interrogator.Role.Name,
				},
				Expired: v.Interrogator.Expired,
			},
			Feedback: &model.DoctorFeedback{
				Base: model.Base{},
				Interrogator: &model.User{
					Name:  v.Feedback.Interrogator.Name,
					Email: v.Feedback.Interrogator.Email,
					Role: model.Roles{
						Name: v.Feedback.Interrogator.Role.Name,
					},
					Expired: v.Feedback.Interrogator.Expired,
				},
				Response: v.Feedback.Response,
				Approved: v.Feedback.Approved,
			},
			Diagnosis:   v.Diagnosis,
			Predictions: v.Predictions,
		})
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Retrive Registration",
		Data:    response,
	})
}

// PATCH /api/patient/medicalrecord/:id
func PatchMedicalRecord(c *gin.Context) {
	encryptedBase64 := c.Param("id")
	id := utils.DecryptToUint64(encryptedBase64)
	var request model.MedicalRecord

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request"})
		return
	}

	database := container.NewContainer()

	// Update Medical Record
	medicalrecord := &migration.MedicalRecord{
		Diagnosis:   request.Diagnosis,
		Predictions: request.Predictions,
	}
	if err := database.MedicalRecords.PatchMedicalRecord(id, medicalrecord); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	// Update Feedback
	feedbackId := utils.DecryptToUint64(request.Feedback.ID)
	feedback := &migration.DoctorFeedback{
		Response: request.Feedback.Response,
		Approved: request.Feedback.Approved,
	}
	if err := database.Feedback.PatchDoctorFeedback(feedbackId, feedback); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{Message: "Medical Record Updated"})
}
