package store

import (
	"encoding/json"

	"github.com/dimeko/sapi/models"
	"github.com/dimeko/sapi/store/cache"
	"github.com/dimeko/sapi/store/rdbms"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Store struct {
	Cache *cache.Cache
	Rdbms models.Actions
}

func New() *Store {
	cache := cache.New()
	rdbms := rdbms.New()
	store := &Store{
		Cache: cache,
		Rdbms: rdbms,
	}
	return store
}

func (s *Store) Create(payload models.UserPayload) (*models.User, error) {
	newUser, err := s.Rdbms.Create(payload)

	if err != nil {
		return nil, err
	}

	unhashedCacheKey := "users:listing"
	err1 := s.Cache.Remove(unhashedCacheKey)
	if err1 != nil {
		log.Error(err1)
	}
	return newUser, nil
}

func (s *Store) Update(id string, payload models.UserPayload) (*models.User, error) {
	unhashedCurrentCacheKey := "user:" + id
	updatedUser, err := s.Rdbms.Update(id, payload)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err1 := s.Cache.Remove(unhashedCurrentCacheKey)
	if err1 != nil {
		log.Error(err1)
	}

	unhashedCacheKey := "users:listing"
	err2 := s.Cache.Remove(unhashedCacheKey)
	if err1 != nil {
		log.Error(err2)
	}

	return updatedUser, nil
}

func (s *Store) List(limit string, offset string) ([]*models.User, error) {
	unhashedCacheKey := "users:listing"
	var users []*models.User
	item, err := s.Cache.Get(unhashedCacheKey)
	if err != nil {
		log.Error("Cache error: ", err)
		users, err1 := s.Rdbms.List(limit, offset)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
		_, err2 := s.Cache.Set(unhashedCacheKey, users)
		if err2 != nil {
			log.Error(err2)
		}
		return users, nil
	}

	log.Info("Cache hit on user listing")
	err2 := json.Unmarshal(item.Value, &users)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}

	return users, nil
}

func (s *Store) Get(id string) (*models.User, error) {
	unhashedCacheKey := "user:" + id
	var user *models.User
	item, err := s.Cache.Get(unhashedCacheKey)
	if err != nil {
		log.Error("Cache error: ", err)
		user, err = s.Rdbms.Get(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		_, err2 := s.Cache.Set(unhashedCacheKey, user)
		if err2 != nil {
			log.Error(err2)
		}
		return user, nil
	}

	log.Info("Cache hit on user get")
	err2 := json.Unmarshal(item.Value, &user)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}

	return user, nil
}

func (s *Store) Delete(id string) error {
	unhashedCacheKey := "user:" + id
	err := s.Rdbms.Delete(id)
	if err != nil {
		log.Error(err)
		return err
	}
	err1 := s.Cache.Remove(unhashedCacheKey)
	if err1 != nil {
		log.Error(err1)
	}

	return nil
}
