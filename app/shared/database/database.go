package database

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/lib/pq"

)


var (
	DB *sql.DB
)

type Database struct {
	Postgres Postgres
}

type Postgres struct {
	User string
	Password string
	Name string
	DriverName string
}

func Connect(dbi Postgres) {

	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", dbi.User, dbi.Password, dbi.Name)

	var err error
	DB, err = sql.Open(dbi.DriverName, connectionString)
	if err != nil {
		log.Println("PostgreSQL error", err)
	}
}

func CheckConnection() bool {
	if DB == nil {
		return false
	}

	if DB != nil {
		return true
	}

	return false
}

