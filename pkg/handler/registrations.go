package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

// GET /api/patients
func GetRegistrations(c *gin.Context) {
	groupID := c.MustGet("groupId").(uint64)

	database := container.NewContainer()
	data, err := database.Registrations.GetRegistrations(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Retrive Registration",
		Data:    data,
	})
}

// POST /api/patients
func StoreRegistration(c *gin.Context) {
	var input model.Registration

	// 1. Validate input binding
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// input.GroupID = c.MustGet("groupId").(uint64)

	// // 2. Call repository/service to add the patient
	// database := container.NewContainer()
	// if err := database.Patients.StorePatient(input); err != nil {
	// 	if strings.Contains(strings.ToLower(err.Error()), "duplicate entry") {
	// 		c.JSON(http.StatusConflict, gin.H{"message": "Medical record number already exist"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create patient", "error": err.Error()})
	// 	return
	// }

	// 3. Return the created patient or status
	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Add Registration",
	})
}

// // PATCH /api/patients/:id
// func PatchRegistration(c *gin.Context) {
// 	id := c.Param("id")
// 	var request model.Patient

// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
// 		return
// 	}

// 	database := container.NewContainer()
// 	if err := database.Patients.PatchPatient(id, request); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Patient updated"})
// }

// // DELETE /api/patients/:id
// func DestroyRegistration(c *gin.Context) {
// 	id := c.Param("id")

// 	database := container.NewContainer()
// 	if err := database.Patients.DestroyPatient(id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted"})
// }
