package main

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
)


const tableCreationQuery = `CREATE TABLE IF NOT EXISTS verifiedusers (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	gtoken TEXT UNIQUE NOT NULL,
	validuntil INT UNIQUE
);
`

func ensureTableExists(a App) {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable(a App) {
	a.DB.Exec("DELETE FROM verifiedusers")
}

func TestSaveUserToDB(t *testing.T) {

	a := App{}
	a.Initialize(
		"postgres",
		"testpass",
		"whattodo")

	ensureTableExists(a)
	clearTable(a)

	email := "test@test.com"
	token := "aVerySecretTestToken"
	validuntil := int64(1234)

	fmt.Println(time.Now().Unix())

	affectedRows, _ := saveUserToDB(a.DB, email, token, validuntil)
	assert.Equal(t, affectedRows, int64(1))

	clearTable(a)
}
