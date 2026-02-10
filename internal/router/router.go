package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/handler"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	api := r.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/logout", handler.Logout)
	auth.POST("/refresh-token", handler.Refresh)

	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.GET("/patients", handler.GetPatients)
	protected.POST("/patients", handler.StorePatient)
	protected.PATCH("/patients/:id", handler.PatchPatient)
	protected.DELETE("/patients/:id", handler.DestroyPatient)

	protected.GET("/users/me", handler.UsersMe)

	protected.GET("/users", handler.GetUsers)
	protected.POST("/users", handler.StoreUser)
	protected.PATCH("/users/:id", handler.PatchUser)
	protected.PATCH("/users/activate/:id", handler.ActivateUser)
	protected.DELETE("/users/:id", handler.DestroyUser)

	protected.GET("/patient/registration", handler.GetRegistrations)
	protected.POST("/patient/registration", handler.StoreRegistration)
	// protected.PATCH("/patient/registration/:id", handler.PatchRegistration)
	protected.DELETE("/patient/registration/:id", handler.DestroyRegistration)

	protected.GET("/patient/medicalrecord", handler.GetMedicalRecords)
	protected.PATCH("/patient/medicalrecord/:id", handler.PatchMedicalRecord)

	protected.GET("/patient/queue", handler.GetQueue)

	protected.POST("/predict", handler.Predict)

	return r
}
