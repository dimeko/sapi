package api

import (
	"encoding/json"
	"net/http"

	"github.com/dimeko/sapi/app"
	"github.com/gorilla/mux"
)

type Api struct {
	app *app.App
}

func New(app *app.App) *mux.Router {
	api := &Api{
		app: app,
	}
	mux := mux.NewRouter()

	mux.HandleFunc("/create", api.create).Methods("POST")
	mux.HandleFunc("/add/{id}", api.update).Methods("PUT")
	mux.HandleFunc("/delete/{id}", api.delete).Methods("DELETE")
	mux.HandleFunc("/get/{id}", api.get).Methods("GET")
	mux.HandleFunc("/list", api.list).Methods("GET")

	return mux
}

/*
For error response we could use http.Error but we set an
errorResposne function to format the error not only string-wise
but json-wise, too
*/
func errorResponse(w http.ResponseWriter, code int, message string) {
	jsonResponse(w, code, map[string]string{"result": "ERROR", "message": message})
}

func successResponse(w http.ResponseWriter, code int, body interface{}) {
	jsonResponse(w, code, map[string]interface{}{"resulst": "SUCCESS", "body": body})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (api *Api) create(w http.ResponseWriter, r *http.Request) {
	var payload app.CreatePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if payload.Username == "" {
		errorResponse(w, http.StatusBadRequest, "Property 'name' was not set")
		return
	}

	createdUser, err := api.app.Create(payload)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Server error")
		return
	}
	successResponse(w, http.StatusCreated, createdUser)
}

func (api *Api) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload app.UpdatePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	if payload.Amount < 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	updatedUser, err2 := api.app.Update(id, payload)

	if err2 != nil {
		if err2.Error() == "notFound" {
			errorResponse(w, http.StatusNotFound, "Not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Server error")
		}
		return
	}
	successResponse(w, http.StatusOK, updatedUser)
}

func (api *Api) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := api.app.Get(id)

	if err != nil {
		if err.Error() == "notFound" {
			errorResponse(w, http.StatusNotFound, "Not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Server error")
		}
		return
	}

	successResponse(w, http.StatusOK, user)
}

func (api *Api) list(w http.ResponseWriter, r *http.Request) {
	users, err := api.app.List()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Server error")
		return
	}
	successResponse(w, http.StatusOK, users)
}

func (api *Api) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := api.app.Delete(id)

	if err != nil {
		if err.Error() == "notFound" {
			errorResponse(w, http.StatusNotFound, "Not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Server error")
		}
		return
	}
	successResponse(w, http.StatusOK, "OK")
}
