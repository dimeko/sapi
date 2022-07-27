package rdbms

import (
	"os"
	"testing"

	"github.com/dimeko/sapi/models"
	"github.com/stretchr/testify/require"
)

var rdbms *RDBMS

func TestMain(m *testing.M) {
	os.Chdir("../../")
	rdbms = New()
	code := m.Run()
	tearDownDatabase()
	os.Exit(code)
}

func tearDownDatabase() {
	rdbms.Db.Exec("DELETE FROM users")
}

func TestCreate(t *testing.T) {
	u := models.UserPayload{Username: "test", Firstname: "test_first", Lastname: "test_last"}
	w, err := rdbms.Create(u)

	require.NoError(t, err)
	require.NotEmpty(t, w)
	tearDownDatabase()
}

func TestUpdate(t *testing.T) {
	u1 := models.UserPayload{Username: "test1", Firstname: "test_first1", Lastname: "test_last1"}
	u2 := models.UserPayload{Username: "test2", Firstname: "test_first2", Lastname: "test_last2"}

	u1_id, err := rdbms.Create(u1)
	u2_id, err := rdbms.Create(u2)

	users, err := rdbms.List("5", "0")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 2, len(users))

	u1_up := models.UserPayload{Username: "test1_up", Firstname: "test_first1_up", Lastname: "test_last1_up"}
	u2_up := models.UserPayload{Username: "test2_up", Firstname: "test_first2_up", Lastname: "test_last2_up"}

	u1_updated, err2 := rdbms.Update(u1_id.Id, u1_up)
	u2_updated, err2 := rdbms.Update(u2_id.Id, u2_up)

	require.NoError(t, err2)
	require.NoError(t, err2)

	require.Equal(t, "test1_up", u1_updated.Username)
	require.Equal(t, "test_first2_up", u2_updated.Firstname)
	tearDownDatabase()
}

func TestGet(t *testing.T) {
	u1 := models.UserPayload{Username: "test", Firstname: "test_first1", Lastname: "test_last1"}

	u1_id, err := rdbms.Create(u1)

	users, err := rdbms.List("5", "0")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 1, len(users))

	user, err2 := rdbms.Get(u1_id.Id)
	require.NoError(t, err2)
	require.NotEmpty(t, user)
	require.Equal(t, "test", user.Username)
	tearDownDatabase()
}

func TestList(t *testing.T) {
	u1 := models.UserPayload{Username: "test1", Firstname: "test_first1", Lastname: "test_last1"}
	u2 := models.UserPayload{Username: "test2", Firstname: "test_first2", Lastname: "test_last2"}

	_, err := rdbms.Create(u1)
	_, err = rdbms.Create(u2)

	ws, err := rdbms.List("5", "0")

	require.NoError(t, err)
	require.NotEmpty(t, ws)
	require.Equal(t, 2, len(ws))
	tearDownDatabase()
}
