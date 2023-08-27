package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     int16
}

func Connect(dbConfig DatabaseConfig) (*gorm.DB, error) {
	var databaseSchemaConnection string = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%v",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Dbname, dbConfig.Port, "Asia/Kolkata",
	)
	db, err := gorm.Open(postgres.Open(databaseSchemaConnection), &gorm.Config{})

	if err != nil {
		fmt.Errorf(err.Error())
		return nil, errors.New(err.Error())
	}

	return db, nil
}
