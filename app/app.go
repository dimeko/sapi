package app

import (
	"github.com/dimeko/sapi/models"
	"github.com/dimeko/sapi/store"
)

type App struct {
	Store *store.Store
}

func New() *App {
	store := store.New()
	return &App{
		Store: store,
	}
}

func (a *App) List(limit, offset string) ([]*models.User, error) {
	users, err := a.Store.List(limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *App) Get(id string) (*models.User, error) {
	user, err := a.Store.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *App) Create(payload models.UserPayload) (*models.User, error) {
	user, err := a.Store.Create(payload)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *App) Update(id string, payload models.UserPayload) (*models.User, error) {
	updatedUser, err := a.Store.Update(id, payload)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (a *App) Delete(id string) error {
	err := a.Store.Delete(id)
	if err != nil {
		return err
	}
	return nil

}
