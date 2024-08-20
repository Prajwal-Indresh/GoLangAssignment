package database

import (
	"go-students-api/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB(cfg *config.Config) error {
	var err error
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetDB() *sqlx.DB {
	return db
}
