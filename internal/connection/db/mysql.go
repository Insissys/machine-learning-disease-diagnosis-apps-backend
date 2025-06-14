package db

import (
	"fmt"
	"log"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Gorm *gorm.DB

func MysqlConnection(cfg *config.Database) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)

	Gorm, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // users, patients, medical_records
		},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
