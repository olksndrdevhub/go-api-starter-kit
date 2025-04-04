package db

import (
	"time"

	"github.com/oleksandrdevhub/go-api-starter-kit/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // password will not be returned
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func CreateUser(email, password, first_name, last_name string) (int64, error) {
	query := `insert into users (email, password, first_name, last_name) values (?, ?, ?, ?)`
	result, err := DB.Exec(query, email, password, first_name, last_name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateUser(userID int64, first_name, last_name string) error {
	query := `update users set first_name = ?, last_name = ?, updated_at = current_timestamp where id = ?`
	_, err := DB.Exec(query, first_name, last_name, userID)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExistsByEmail(email string) (bool, error) {
	query := `select count(*) from users where email = ?`
	var count int
	err := DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `select id, email, password, first_name, last_name, created_at, updated_at from users where email = ?`
	row := DB.QueryRow(query, email)

	var user User
	var createdAt string
	var updatedAt string
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, err = utils.ParseTime(createdAt)
	user.UpdatedAt, err = utils.ParseTime(updatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id int64) (*User, error) {
	query := `select id, email, password, first_name, last_name, created_at, updated_at from users where id = ?`
	row := DB.QueryRow(query, id)

	var user User
	var createdAt string
	var updatedAt string
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, err = utils.ParseTime(createdAt)
	user.UpdatedAt, err = utils.ParseTime(updatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
