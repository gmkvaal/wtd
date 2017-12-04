package main

import (
	"encoding/json"
	"log"
	"database/sql"
)

type UserInfo struct {
	Id string
	Email string
	VerifiedEmail string
	Name string
	GivenName string
	FamilyName string
	Link string
	Picture string
	Gender string
	Locale string
}

func ExtractUserData(body []byte) UserInfo {

	var ui UserInfo
	err := json.Unmarshal(body, &ui)
	if err != nil {
		log.Fatal(err)
	}

	return ui
}


func saveUserToDB(db *sql.DB, email string, token string, validUntil int64) (int64, error) {

	sqlStatement := `
		INSERT INTO verifiedusers (email, gtoken, validuntil)
		VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, email, token, validUntil)
	if err != nil {
		return 0, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}