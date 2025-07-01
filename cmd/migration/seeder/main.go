package main

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/seed"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()

	db.MysqlConnection(&cfg.Config.Database)
	seed.UserGroupSeederDev()
}
