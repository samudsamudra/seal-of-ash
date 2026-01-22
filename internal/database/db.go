package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"seal-of-ash/internal/config"
)

var ActiveDB *gorm.DB
var ForensicDB *gorm.DB

func Init() {
	dsnActive := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Get("DB_USER"),
		config.Get("DB_PASS"),
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_NAME"),
	)

	dsnForensic := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Get("DB_USER"),
		config.Get("DB_PASS"),
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("FORENSIC_DB_NAME"),
	)

	var err error
	ActiveDB, err = gorm.Open(mysql.Open(dsnActive), &gorm.Config{})
	if err != nil {
		panic("failed connect Active DB")
	}

	ForensicDB, err = gorm.Open(mysql.Open(dsnForensic), &gorm.Config{})
	if err != nil {
		panic("failed connect Forensic DB")
	}
}
