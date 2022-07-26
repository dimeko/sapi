package rdbms

import (
	"database/sql"
	"os"
	"strings"

	"github.com/dimeko/sapi/app"
	"github.com/dimeko/sapi/store"
)

type RDBS struct {
	store.Store
	db *sql.DB
}

func env() map[string]string {
	return map[string]string{
		"username": os.Getenv("POSTGRES_USER"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"db_name":  os.Getenv("POSTGRES_DB"),
		"port":     os.Getenv("POSTGRES_PORT"),
	}
}

func New() *RDBS {
	dsn := strings.Join([]string{"postgres://",
		env()["username"], ":",
		env()["password"], "@localhost:",
		env()["port"], "/", env()["db_name"], "/sslmode=disable"}, "")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("Could not initialize database")
	}

	store := &RDBS{
		db: db,
	}

	return store
}

func (s *RDBS) Create(payload app.CreatePayload) (string, error) {
	var id string
	err := s.db.QueryRow("INSERT INTO users(username, firstname, lastname) VALUES($1, $2, $3)",
		payload.Username, payload.Firstname, payload.Lastname).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil

}

func (s *RDBS) Update(id string, payload app.UpdatePayload) (store.User, error) {
	var user store.User
	err := s.db.QueryRow("UPDATE users SET username=$1, firstname=$2, lastname=$3 WHERE id=$4",
		payload.Username, payload.Firstname, payload.Lastname, id).Scan(&user)
	if err != nil {
		return store.User{
			Id:        "",
			Username:  "",
			Firstname: "",
			Lastname:  "",
		}, err
	}
	return user, err
}

func (s *RDBS) Get(id string) (store.User, error) {
	var user store.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&user)
	if err != nil {
		return store.User{
			Id:        "",
			Username:  "",
			Firstname: "",
			Lastname:  "",
		}, err
	}
	return user, err
}

func (s *RDBS) List(limit string, offset string) ([]store.User, error) {
	rows, err := s.db.Query("SELECT * FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []store.User{}
	for rows.Next() {
		var u store.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *RDBS) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
