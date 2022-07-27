package rdbms

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dimeko/sapi/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type RDBMS struct {
	Db *sql.DB
}

func env() map[string]string {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		panic("Cannot find .env file")
	}
	return map[string]string{
		"username": os.Getenv("POSTGRES_USER"),
		"host":     os.Getenv("POSTGRES_HOST"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"db_name":  os.Getenv("POSTGRES_DB"),
		"port":     os.Getenv("POSTGRES_PORT"),
	}
}

func New() *RDBMS {
	host := env()["host"]
	port := env()["port"]
	user := env()["username"]
	password := env()["password"]
	dbname := env()["db_name"]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	rdbms := &RDBMS{
		Db: db,
	}

	return rdbms
}

func (s *RDBMS) Create(payload models.UserPayload) (*models.User, error) {
	id := 0
	err := s.Db.QueryRow(`INSERT INTO users (username, firstname, lastname) VALUES ($1, $2, $3) RETURNING id`,
		payload.Username, payload.Firstname, payload.Lastname).Scan(&id)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	user := &models.User{}
	row := s.Db.QueryRow("SELECT id, username, firstname, lastname FROM users WHERE id=$1", fmt.Sprint(id))
	if err := row.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname); err != nil { // scan will release the connection
		return nil, err
	}

	return user, nil

}

func (s *RDBMS) Update(id string, payload models.UserPayload) (*models.User, error) {
	updatedId := 0
	err := s.Db.QueryRow("UPDATE users SET username=$1, firstname=$2, lastname=$3 WHERE id=$4 RETURNING id",
		payload.Username, payload.Firstname, payload.Lastname, id).Scan(&updatedId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	user := &models.User{}
	row := s.Db.QueryRow("SELECT id, username, firstname, lastname FROM users WHERE id=$1", updatedId)
	if err := row.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname); err != nil { // scan will release the connection
		return nil, err
	}

	return user, err
}

func (s *RDBMS) Get(id string) (*models.User, error) {
	user := &models.User{}
	row := s.Db.QueryRow("SELECT id, username, firstname, lastname FROM users WHERE id=$1", id)

	if err := row.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname); err != nil { // scan will release the connection
		return nil, err
	}

	return user, nil
}

func (s *RDBMS) List(limit string, offset string) ([]*models.User, error) {
	rows, err := s.Db.Query("SELECT id, username, firstname, lastname FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (s *RDBMS) Delete(id string) error {
	_, err := s.Db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
