package model

import (
	"github.com/gmkvaal/wtd/app/shared/database"
)


var db = database.DB

func SaveUserToDB(email string, token string, validUntil int64) error {

	sqlStatement := `
		INSERT INTO verifiedusers (email, gtoken, validuntil)
		VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, email, token, validUntil)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTokenInDB(email string, token string, validUntil int64) error {

	sqlStatement :=
		`UPDATE verifiedusers SET gtoken=$1, validuntil=$2 WHERE email=$3`
	_, err := db.Exec(sqlStatement, token, validUntil, email)
	if err != nil {
		return err
	}

	return nil
}

func CheckIfEmailAlreadyInDB(email string, token string, validUntil int64) (bool, error) {

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

func CheckIfUserIsValidated(email string, token string, validUntil int64) (bool, error) {

	var dbToken string

	sqlStatement := `SELECT gtoken FROM verifiedusers WHERE email=$1`
	err := db.QueryRow(sqlStatement, email).Scan(&dbToken)
	if err != nil {
		return false, err
	}

	if dbToken == token {
		return true, nil
	}

	return false, err
}