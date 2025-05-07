package model

import (
	"database/sql"

	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	// InsertWithToken(name, email, password, token string) error // Add this
	// VerifyEmail(token string) error                            // Add this
}

type Users struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

var ErrInvalidToken = errors.New("invalid verification token")

func (m *UserModel) Insert(name, email, password string) error {
	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?,?,?,UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hassedPassword))
	if err != nil {
		var MysqlError *mysql.MySQLError
		if errors.As(err, &MysqlError) {
			if MysqlError.Number == 1062 && strings.Contains(MysqlError.Message, "users_uc_email") {
				return DuplicateEmail
			}
		}
		return err
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPass []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, InvalidCredientials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, InvalidCredientials
		} else {
			return 0, err
		}
	}
	return id, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	var exist bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exist)
	return exist, err
}

// func (m *UserModel) InsertWithToken(name, email, password, token string) error {
// 	stmt := `INSERT INTO users (name, email, hashed_password, verification_token, created)
//     VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`
// 	_, err := m.DB.Exec(stmt, name, email, password, token)
// 	return err
// }

// func (m *UserModel) VerifyEmail(token string) error {
// 	stmt := `UPDATE users SET verified = TRUE, verification_token = NULL WHERE verification_token = ?`
// 	result, err := m.DB.Exec(stmt, token)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return ErrInvalidToken
// 	}

// 	return nil
// }
