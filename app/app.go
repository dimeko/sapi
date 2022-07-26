package app

import (
	"github.com/dimeko/sapi/store"
)

type App struct {
	Actions
	store *store.Store
}

func New() *App {
	store := store.New()
	return &App{
		store: store,
	}
}

func List(params ...string) {

}

func Get(id string, params ...string) {

}

func Create(payload CreatePayload) {

}

func Update(payload UpdatePayload) {

}

func Delete(id string) {

}
