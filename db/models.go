package db

import (
	"time"
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

func CreateUser(email, password, first_name, last_name string) (*User, error) {
	var user User
	query := `
    insert into users 
      (email, password, first_name, last_name) 
    values 
      ($1, $2, $3, $4)
    returning
      id, email, password, first_name, last_name, created_at, updated_at
  `
	err := DB.QueryRow(query, email, password, first_name, last_name).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

func UpdateUser(userID int64, first_name, last_name string) error {
	query := `
    update users set 
      first_name = $1, last_name = $2, updated_at = current_timestamp
    where id = $3
  `
	_, err := DB.Exec(query, first_name, last_name, userID)
	return err

}

func ChangeUserPassword(userID int64, hashedPassword string) error {
	query := `
    update users set 
      password = $1, updated_at = current_timestamp
    where id = $2
  `
	_, err := DB.Exec(query, hashedPassword, userID)
	return err
}

func CheckUserExistsByEmail(email string) (bool, error) {
	query := `select exists(select 1 from users where email = $1)`
	var exists bool
	err := DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func GetUserByEmail(email string) (*User, error) {
	query := `
    select 
      id, email, password, first_name, last_name, created_at, updated_at
    from users
    where 
      email = $1
  `

	var user User
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id int64) (*User, error) {
	query := `
    select 
      id, email, password, first_name, last_name, created_at, updated_at
    from users
    where 
      id = $1
  `

	var user User
	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
