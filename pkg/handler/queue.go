package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
)

// GET /api/patient/queue
func GetQueue(c *gin.Context) {
	groupId := c.MustGet("groupId").(uint64)
	encodedUserId := c.Query("userId")
	userRole := c.MustGet("role").(string)
	var userIdParam *uint64

	if userRole != "superadmin" {
		id := utils.DecryptToUint64(encodedUserId)
		userIdParam = &id
	}

	database := container.NewContainer()
	data, err := database.Queue.GetQueue(userIdParam, &groupId)
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
				Feedback: &model.DoctorFeedback{
					Base: model.Base{
						ID:        utils.EncryptUint64(uint64(v.MedicalRecord.Feedback.ID)),
						CreatedAt: &v.MedicalRecord.Feedback.CreatedAt,
						UpdatedAt: &v.MedicalRecord.Feedback.UpdatedAt,
					},
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
