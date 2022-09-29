package common

import (
	"cc_DavidGayle_BackendAPI/internal/app/model"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type db struct {
	dbConfig model.Database
}

type DB interface {
	GetDatabase() *gorm.DB
	InitialMigration() error
	CloseDatabase(connection *gorm.DB)
}

func NewDb(config *model.Config) DB {
	return &db{
		dbConfig: config.Db,
	}
}

func (d *db) GetDatabase() *gorm.DB {
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		d.dbConfig.Host, d.dbConfig.User, d.dbConfig.Password, d.dbConfig.Name, d.dbConfig.Port)

	connection, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sqldb, err := connection.DB()
	if err != nil {
		log.Fatal("Database connection error", err)
	}

	if err = sqldb.Ping(); err != nil {
		log.Fatal("Database connection error", err)
	}

	fmt.Println("Database connection successful.")
	return connection
}

func (d *db) InitialMigration() error {
	connection := d.GetDatabase()
	defer d.CloseDatabase(connection)
	return connection.AutoMigrate(model.User{})
}

func (d *db) CloseDatabase(connection *gorm.DB) {
	var sqldb *sql.DB
	var err error

	if sqldb, err = connection.DB(); err != nil {
		fmt.Println("Error connection to DB", err)
		return
	}

	if err = sqldb.Close(); err != nil {
		fmt.Println("Error closing DB", err)
	}
}
