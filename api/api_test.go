package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dimeko/sapi/app"
	"github.com/dimeko/sapi/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var api *Api
var testDbClient *sql.DB

type ErrorResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type SuccessReponse struct {
	Result string      `json:"result"`
	Body   interface{} `json:"body"`
}

type SingleUserSuccessReponse struct {
	Result string      `json:"result"`
	Body   models.User `json:"body"`
}

type MultipleUserSuccessReponse struct {
	Result string        `json:"result"`
	Body   []models.User `json:"body"`
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

func TestMain(m *testing.M) {
	os.Chdir("../")
	host := env()["host"]
	port := env()["port"]
	user := env()["username"]
	password := env()["password"]
	dbname := env()["db_name"]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	testDbClient, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	app := app.New()
	api = New(app)
	tst := m.Run()
	tearDownDatabase()
	os.Exit(tst)
}

func tearDownDatabase() {

	testDbClient.Exec("DELETE FROM users")
}

func TestCreate(t *testing.T) {
	server, _ := initServer(t)
	defer tearDownDatabase()
	defer server.Close()
}

func TestGet(t *testing.T) {
	server, users := initServer(t)
	defer tearDownDatabase()
	defer server.Close()
	t.Run("found user", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/get/" + users[0].Id)
		require.NoError(t, err)
		var result SingleUserSuccessReponse
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)
		require.Equal(t, users[0].Username, result.Body.Username)
	})

}

func TestList(t *testing.T) {
	server, _ := initServer(t)
	defer tearDownDatabase()
	defer server.Close()
	t.Run("listing users", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/list")
		require.NoError(t, err)
		var result MultipleUserSuccessReponse
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)
		require.Equal(t, 4, len(result.Body))
	})
	tearDownDatabase()
}

func initServer(t *testing.T) (*httptest.Server, []models.User) {
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}

	var u1 = []byte(`{"username":"user_one", "firstname": "Nick", "lastname": "Jonas"}`)
	req1, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(u1))

	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := client.Do(req1)
	defer resp1.Body.Close()

	var u2 = []byte(`{"username": "user_two", "firstname": "George", "lastname": "Harrison"}`)
	req2, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(u2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := client.Do(req2)
	defer resp2.Body.Close()

	var u3 = []byte(`{"username": "user_three", "firstname": "Amy", "lastname": "Winehouse"}`)
	req3, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(u3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := client.Do(req3)
	defer resp3.Body.Close()

	var u4 = []byte(`{"username": "user_four", "firstname": "Rihanna", "lastname": "Beyonce"}`)
	req4, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(u4))
	req4.Header.Set("Content-Type", "application/json")
	resp4, _ := client.Do(req4)
	defer resp4.Body.Close()

	var u1r SingleUserSuccessReponse
	body, err := ioutil.ReadAll(resp1.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &u1r)
	require.Equal(t, "user_one", u1r.Body.Username)
	require.Equal(t, "Nick", u1r.Body.Firstname)

	var u2r SingleUserSuccessReponse
	body, err = ioutil.ReadAll(resp2.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &u2r)

	require.Equal(t, "user_two", u2r.Body.Username)
	require.Equal(t, "George", u2r.Body.Firstname)

	var u3r SingleUserSuccessReponse
	body, err = ioutil.ReadAll(resp3.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &u3r)
	require.Equal(t, "user_three", u3r.Body.Username)
	require.Equal(t, "Amy", u3r.Body.Firstname)

	var u4r SingleUserSuccessReponse
	body, err = ioutil.ReadAll(resp4.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &u4r)
	require.Equal(t, "user_four", u4r.Body.Username)
	require.Equal(t, "Rihanna", u4r.Body.Firstname)

	return ts, []models.User{u1r.Body, u2r.Body, u3r.Body, u4r.Body}
}
