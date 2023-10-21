package database

import (
	"errors"
	"time"
)

type Token struct {
	Time time.Time `json:"time"`
}

func (db *DB) RevokeToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	token := Token{
		Time: time.Now().UTC(),
	}

	status, err := db.CheckTokenRevokeStatus(tokenString)
	if err != nil {
		return err
	}
	if status {
		return errors.New("Token already revoked")
	}

	dbStructure.Tokens[tokenString] = token

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// Checks if a Token was revoked. Returns true if it was revoked, false if not.
func (db *DB) CheckTokenRevokeStatus(tokenString string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}
	_, ok := dbStructure.Tokens[tokenString]
	if !ok {
		return false, nil
	}
	return true, nil
}
