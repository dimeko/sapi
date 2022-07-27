package api

import (
	"encoding/json"
	"net/http"

	"github.com/dimeko/sapi/app"
	"github.com/dimeko/sapi/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Api struct {
	Router *mux.Router
	App    *app.App
}

func New(app *app.App) *Api {
	mux := mux.NewRouter()
	api := &Api{
		Router: mux,
		App:    app,
	}

	mux.HandleFunc("/create", api.create).Methods("POST")
	mux.HandleFunc("/update/{id}", api.update).Methods("PUT")
	mux.HandleFunc("/delete/{id}", api.delete).Methods("DELETE")
	mux.HandleFunc("/get/{id}", api.get).Methods("GET")
	mux.HandleFunc("/list", api.list).Methods("GET")

	return api
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
	jsonResponse(w, code, map[string]interface{}{"result": "SUCCESS", "body": body})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (api *Api) create(w http.ResponseWriter, r *http.Request) {
	var payload models.UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if payload.Username == "" {
		errorResponse(w, http.StatusBadRequest, "Property 'username' was not set")
		return
	}

	createdUser, err := api.App.Create(payload)

	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusBadRequest, "Server error")
		return
	}
	successResponse(w, http.StatusCreated, createdUser)
}

func (api *Api) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload models.UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	updatedUser, err2 := api.App.Update(id, payload)

	if err2 != nil {
		log.Error(err2)
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
	user, err := api.App.Get(id)

	if err != nil {
		log.Error(err)
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
	limit := mux.Vars(r)["limit"]
	offset := mux.Vars(r)["offset"]
	if limit == "" {
		limit = "100"
	}

	if offset == "" {
		offset = "0"
	}

	users, err := api.App.List(limit, offset)
	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusInternalServerError, "Server error")
		return
	}
	successResponse(w, http.StatusOK, users)
}

func (api *Api) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := api.App.Delete(id)

	if err != nil {
		log.Error(err)
		if err.Error() == "notFound" {
			errorResponse(w, http.StatusNotFound, "Not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Server error")
		}
		return
	}
	successResponse(w, http.StatusOK, "OK")
}
