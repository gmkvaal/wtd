package main

import (
	"encoding/json"
	"log"
	"database/sql"
	"fmt"
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

func updateTokenInDB(db *sql.DB, email string, token string, validUntil int64) error {

	sqlStatement :=
		`UPDATE verifiedusers SET gtoken=$1, validuntil=$1 WHERE email=$2`
	_, err := db.Exec(sqlStatement, token, validUntil, email)
	if err != nil {
		return err
	}

	return nil
}

func checkIfEmailAlreadyInDB(db *sql.DB, email string, token string, validUntil int64) (bool, error) {

	sqlStatement := `SELECT gtoken FROM verifiedusers WHERE email=$1`
	res, err := db.Exec(sqlStatement, email)
	if err != nil {
		return false, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if affectedRows == 1 {
		return true, nil
	}

	return false, err

}

func checkIfUserIsValidated(db *sql.DB, email string, token string, validUntil int64) (error) {

	var dbToken string

	sqlStatement := `SELECT gtoken FROM verifiedusers WHERE email=$1`
	err := db.QueryRow(sqlStatement, email).Scan(&dbToken)
	if err != nil {
		return err
	}

	fmt.Println("token from DB", dbToken)
	fmt.Println("token from user", token)

	return err
}