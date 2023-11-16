package storage

import "database/sql"

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(u *sql.DB) *UserStorage {
	return &UserStorage{
		db: u,
	}
}
func (u *UserStorage) Save(user User) (int64, error) {
	query := "INSERT INTO users (login, password) VALUES ($1,$2) RETURNING id"
	var id int64
	err := u.db.QueryRow(query, user.Login, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserStorage) GetUserByID(id int64) (User, error) {
	query := "SELECT id, login FROM users WHERE id = $1 "

	var user User

	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Login)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserStorage) Login(login string, password string) (int64, error) {
	query := "SELECT id FROM users WHERE login=$1 AND password = $2"

	var id int64
	err := u.db.QueryRow(query, login, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
