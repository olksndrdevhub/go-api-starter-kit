package db

const (
	CreateUserQuery = `
    insert into users (
      email, password, first_name, last_name
    ) values (
      ?, ?, ?, ?, ?
    )
  `
	UpdateUserQuery = `
    update users set
      first_name = ?, last_name = ?, updated_at = current_timestamp
    where id = ?
  `
	ChangeUserPasswordQuery = `
    update users set
      password = ?, updated_at = current_timestamp
    where id = ?
  `
	CheckUserExistsByEmailQuery = `
    select count(*) from users where email = ?
  `
	GetUserByIDQuery = `
    select * from users where id = ?
  `
	GetUserByEmailQuery = `
    select * from users where email = ?
  `
)
