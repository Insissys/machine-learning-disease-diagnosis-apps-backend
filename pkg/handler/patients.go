package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
)

// GET /api/patients
func GetPatients(c *gin.Context) {
	groupID := c.MustGet("groupId").(uint64)

	database := container.NewContainer()
	patients, err := database.Patients.GetPatients(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message:  "Something went wrong",
			Error:    err,
			Data:     nil,
			Metadata: nil,
		})
		return
	}

	response := []model.Patient{}

	for _, v := range patients {
		response = append(response, model.Patient{
			Base: model.Base{
				ID:        utils.EncryptUint64(uint64(v.ID)),
				CreatedAt: &v.CreatedAt,
				UpdatedAt: &v.UpdatedAt,
			},
			MedicalRecordNumber: v.MedicalRecordNumber,
			Name:                v.Name,
			Gender:              v.Gender,
			BirthDate:           model.DateOnly{Time: v.BirthDate},
			Group: model.Group{
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
		Message:  "Success Retrieve Patients",
		Error:    nil,
		Data:     response,
		Metadata: nil,
	})
}

// POST /api/patients
func StorePatient(c *gin.Context) {
	var input model.Patient

	// 1. Validate input binding
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Message: "Invalid input",
			Error:   err,
		})
		return
	}

	groupId := c.MustGet("groupId").(uint64)
	input.Group.ID = utils.EncryptUint64(groupId)

	// 2. Call repository/service to add the patient
	database := container.NewContainer()
	if err := database.Patients.StorePatient(&migration.Patient{
		MedicalRecordNumber: input.MedicalRecordNumber,
		Name:                input.Name,
		Gender:              input.Gender,
		BirthDate:           input.BirthDate.Time,
		GroupID:             uint(utils.DecryptToUint64(input.Group.ID)),
	}); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate entry") {
			c.JSON(http.StatusConflict, model.ApiResponse{
				Message: "Medical record number already exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Failed to create patient",
			Error:   err,
		})
		return
	}

	// 3. Return the created patient or status
	c.JSON(http.StatusCreated, model.ApiResponse{Message: "Patient created successfully"})
}

// PATCH /api/patients/:id
func PatchPatient(c *gin.Context) {
	encryptedBase64 := c.Param("id")
	id := strconv.Itoa(int(utils.DecryptToUint64(encryptedBase64)))
	var request model.Patient

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request"})
		return
	}

	database := container.NewContainer()
	if err := database.Patients.PatchPatient(id, &migration.Patient{
		MedicalRecordNumber: request.MedicalRecordNumber,
		Name:                request.Name,
		Gender:              request.Gender,
		BirthDate:           request.BirthDate.Time,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Something went wrong",
			Error:   err,
		})
		return
	}
	c.JSON(http.StatusOK, model.ApiResponse{Message: "Patient updated"})
}

// DELETE /api/patients/:id
func DestroyPatient(c *gin.Context) {
	encryptedBase64 := c.Param("id")
	id := strconv.Itoa(int(utils.DecryptToUint64(encryptedBase64)))

	database := container.NewContainer()
	if err := database.Patients.DestroyPatient(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Something went wrong",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{Message: "Patient deleted"})
}
