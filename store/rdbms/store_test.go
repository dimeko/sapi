package rdbms

import (
	"os"
	"testing"

	"github.com/dimeko/sapi/app"
	"github.com/stretchr/testify/require"
)

var rdbms *RDBS

func TestMain(m *testing.M) {
	rdbms = New()
	code := m.Run()
	tearDownDatabase()
	os.Exit(code)
}

func tearDownDatabase() {
	rdbms.db.Exec("DELETE FROM users")
}

func TestCreate(t *testing.T) {
	u := app.CreatePayload{Username: "test", Firstname: "test_first", Lastname: "test_last"}
	w, err := rdbms.Create(u)

	require.NoError(t, err)
	require.NotEmpty(t, w)
}

func TestUpdate(t *testing.T) {
	u1 := app.CreatePayload{Username: "test1", Firstname: "test_first1", Lastname: "test_last1"}
	u2 := app.CreatePayload{Username: "test2", Firstname: "test_first2", Lastname: "test_last2"}

	u1_id, err := rdbms.Create(u1)
	u2_id, err := rdbms.Create(u2)

	users, err := rdbms.List("5", "0")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 2, len(users))

	u1_up := app.UpdatePayload{Username: "test1_up", Firstname: "test_first1_up", Lastname: "test_last1_up"}
	u2_up := app.UpdatePayload{Username: "test2_up", Firstname: "test_first2_up", Lastname: "test_last2_up"}

	u1_updated, err2 := rdbms.Update(u1_id, u1_up)
	u2_updated, err2 := rdbms.Update(u2_id, u2_up)

	require.NoError(t, err2)
	require.NoError(t, err2)

	require.Equal(t, "test1_up", u1_updated.Username)
	require.Equal(t, "test_first2_up", u2_updated.Firstname)
}

func TestGet(t *testing.T) {
	u1 := app.CreatePayload{Username: "test", Firstname: "test_first1", Lastname: "test_last1"}

	u1_id, err := rdbms.Create(u1)

	users, err := rdbms.List("5", "0")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 1, len(users))

	user, err2 := rdbms.Get(u1_id)
	require.NoError(t, err2)
	require.NotEmpty(t, user)
	require.Equal(t, "test", user.Username)
}

func TestList(t *testing.T) {
	u1 := app.CreatePayload{Username: "test1", Firstname: "test_first1", Lastname: "test_last1"}
	u2 := app.CreatePayload{Username: "test2", Firstname: "test_first2", Lastname: "test_last2"}

	_, err := rdbms.Create(u1)
	_, err = rdbms.Create(u2)

	ws, err := rdbms.List("5", "0")

	require.NoError(t, err)
	require.NotEmpty(t, ws)
	require.Equal(t, 2, len(ws))
}
