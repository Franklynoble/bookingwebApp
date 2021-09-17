package dbrepo

import (
	"database/sql"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewpostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
