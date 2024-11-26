package auth

import (
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	DB *sql.DB
}

func (db *AuthRepository) GetUser(data UserModel) (UserModel, error) {
	conn := db.DB
	var user UserModel

	query := `SELECT id, username, password, is_admin FROM user
		WHERE username=? OR email=?;`
	err := conn.QueryRow(query, data.Username, data.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrIncorrectCredentials
		}
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return user, ErrIncorrectCredentials
	}

	return user, nil
}

func (db *AuthRepository) RegisterUser(data UserModel) error {
	conn := db.DB

	bytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO user(
			id, first_name, username, email, password
		) VALUES (?, ?, ?, ?, ?);`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(uuid.NewString(), data.FirstName, data.Username, data.Email, string(bytes))
	if err != nil {
		return err
	}

	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i != 1 {
		return ErrUserNotInserted
	}

	return nil
}
