package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dimeko/sapi/api"
	"github.com/dimeko/sapi/store"
	"github.com/stretchr/testify/require"
)

const dbPath = "db/"
const dbWalletPath = "db/wallets/"

type WalletResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

func TestMain(m *testing.M) {
	os.RemoveAll(dbPath)
	tst := m.Run()
	os.Exit(tst)
}

func TestCreate(t *testing.T) {
	server, _ := initServer(t)
	defer server.Close()
	os.RemoveAll(dbPath)
}

func TestAdd(t *testing.T) {
	server, wallets := initServer(t)
	client := &http.Client{}

	defer server.Close()

	t.Run("add amount to wallets", func(t *testing.T) {
		var w1 = []byte(`{"amount": 88}`)
		req1, _ := http.NewRequest("PUT", server.URL+"/add/"+wallets[0].Id, bytes.NewBuffer(w1))
		req1.Header.Set("Content-Type", "application/json")
		resp1, _ := client.Do(req1)
		defer resp1.Body.Close()

		var w1Result WalletResponse
		body, err := ioutil.ReadAll(resp1.Body)
		require.NoError(t, err)
		json.Unmarshal(body, &w1Result)
		require.Equal(t, int64(588), w1Result.Balance)
		require.Equal(t, "wallet_1", w1Result.Name)

		var w2 = []byte(`{"amount": 200}`)
		req2, _ := http.NewRequest("PUT", server.URL+"/add/"+wallets[1].Id, bytes.NewBuffer(w2))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := client.Do(req2)
		defer resp1.Body.Close()

		var w2Result WalletResponse
		body, err = ioutil.ReadAll(resp2.Body)
		require.NoError(t, err)
		json.Unmarshal(body, &w2Result)
		require.Equal(t, int64(800), w2Result.Balance)
		require.Equal(t, "wallet_2", w2Result.Name)
	})
	os.RemoveAll(dbPath)
}

func TestRemove(t *testing.T) {
	server, wallets := initServer(t)
	client := &http.Client{}

	defer server.Close()
	t.Run("remove amount from wallets", func(t *testing.T) {
		var w1 = []byte(`{"amount": 88}`)
		req1, _ := http.NewRequest("PUT", server.URL+"/remove/"+wallets[0].Id, bytes.NewBuffer(w1))
		req1.Header.Set("Content-Type", "application/json")
		resp1, _ := client.Do(req1)
		defer resp1.Body.Close()

		var w1Result WalletResponse
		body, err := ioutil.ReadAll(resp1.Body)
		require.NoError(t, err)
		json.Unmarshal(body, &w1Result)
		require.Equal(t, int64(412), w1Result.Balance)
		require.Equal(t, "wallet_1", w1Result.Name)

		var w2 = []byte(`{"amount": 700}`)
		req2, _ := http.NewRequest("PUT", server.URL+"/remove/"+wallets[1].Id, bytes.NewBuffer(w2))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := client.Do(req2)

		require.Equal(t, 400, resp2.StatusCode)
	})

	t.Run("prevent double spent amount from wallets", func(t *testing.T) {
		var call1 = []byte(`{"amount": 500}`)
		req1, _ := http.NewRequest("PUT", server.URL+"/remove/"+wallets[2].Id, bytes.NewBuffer(call1))
		req1.Header.Set("Content-Type", "application/json")
		var resp1 *http.Response
		go func() {
			resp1, _ = client.Do(req1)
		}()

		var call2 = []byte(`{"amount": 500}`)
		req2, _ := http.NewRequest("PUT", server.URL+"/remove/"+wallets[2].Id, bytes.NewBuffer(call2))
		req2.Header.Set("Content-Type", "application/json")
		var resp2 *http.Response
		go func() {
			resp2, _ = client.Do(req2)
		}()

		time.Sleep(1 * time.Second)

		var call1Result WalletResponse
		body, err := ioutil.ReadAll(resp1.Body)
		require.NoError(t, err)
		json.Unmarshal(body, &call1Result)

		var call2Result WalletResponse
		body, err = ioutil.ReadAll(resp2.Body)
		require.NoError(t, err)
		json.Unmarshal(body, &call2Result)
		finalBalance := call1Result.Balance + call2Result.Balance

		require.Equal(t, int64(200), finalBalance)

		defer resp1.Body.Close()
		defer resp2.Body.Close()
	})

	os.RemoveAll(dbPath)
}

func TestGet(t *testing.T) {
	server, wallets := initServer(t)

	defer server.Close()
	t.Run("found wallet", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/get/" + wallets[0].Id)
		require.NoError(t, err)
		var result WalletResponse
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)
		require.Equal(t, wallets[0].Id, result.Id)
	})
	os.RemoveAll(dbPath)
}

func TestList(t *testing.T) {
	server, _ := initServer(t)

	defer server.Close()
	t.Run("listing wallets", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/list")
		require.NoError(t, err)
		var result []WalletResponse
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)
		require.Equal(t, 4, len(result))
	})
	os.RemoveAll(dbPath)
}

func initServer(t *testing.T) (*httptest.Server, []WalletResponse) {
	store := store.Init()
	r := api.Init(store)
	ts := httptest.NewServer(r)
	client := &http.Client{}

	var w1 = []byte(`{"name":"wallet_1", "balance": 500}`)
	req1, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(w1))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := client.Do(req1)
	defer resp1.Body.Close()

	var w2 = []byte(`{"name":"wallet_2", "balance": 600}`)
	req2, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(w2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := client.Do(req2)
	defer resp2.Body.Close()

	var w3 = []byte(`{"name":"wallet_3", "balance": 700}`)
	req3, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(w3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := client.Do(req3)
	defer resp3.Body.Close()

	var w4 = []byte(`{"name":"wallet_4", "balance": 800}`)
	req4, _ := http.NewRequest("POST", ts.URL+"/create", bytes.NewBuffer(w4))
	req4.Header.Set("Content-Type", "application/json")
	resp4, _ := client.Do(req4)
	defer resp4.Body.Close()

	var w1Result WalletResponse
	body, err := ioutil.ReadAll(resp1.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &w1Result)
	require.Equal(t, int64(500), w1Result.Balance)
	require.Equal(t, "wallet_1", w1Result.Name)

	var w2Result WalletResponse
	body, err = ioutil.ReadAll(resp2.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &w2Result)
	require.Equal(t, int64(600), w2Result.Balance)
	require.Equal(t, "wallet_2", w2Result.Name)

	var w3Result WalletResponse
	body, err = ioutil.ReadAll(resp3.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &w3Result)
	require.Equal(t, int64(700), w3Result.Balance)
	require.Equal(t, "wallet_3", w3Result.Name)

	var w4Result WalletResponse
	body, err = ioutil.ReadAll(resp4.Body)
	require.NoError(t, err)
	json.Unmarshal(body, &w4Result)
	require.Equal(t, int64(800), w4Result.Balance)
	require.Equal(t, "wallet_4", w4Result.Name)

	return ts, []WalletResponse{w1Result, w2Result, w3Result, w4Result}
}
