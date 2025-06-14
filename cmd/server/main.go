package main

import (
	_ "net/http"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/router"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()

	db.MysqlConnection(&cfg.Config.Database)

	r := router.SetupRouter()
	r.Run(cfg.Config.Server.Host + ":" + cfg.Config.Server.Port)
}
