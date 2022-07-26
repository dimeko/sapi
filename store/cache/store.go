package cache

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io/ioutil"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dimeko/sapi/app"
	"github.com/dimeko/sapi/store"
)

type Cache struct {
	store.Store
	mc *memcache.Client
}

// func env() map[string]string {
// 	return map[string]string{
// 		"username": os.Getenv("POSTGRES_USER"),
// 		"password": os.Getenv("POSTGRES_PASSWORD"),
// 		"db_name":  os.Getenv("POSTGRES_DB"),
// 		"port":     os.Getenv("POSTGRES_PORT"),
// 	}
// }

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return string(h.Sum64())
}

func New(childStore store.Store) *Cache {
	mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")

	cache := &Cache{
		Store: childStore,
		mc:    mc,
	}

	return cache
}

func (c *Cache) Create(id string, payload app.CreatePayload) (string, error) {
	cacheKey := hash(id + payload.Username)
	cacheValue, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}
	c.mc.Set(&memcache.Item{Key: cacheKey, Value: cacheValue})

	return cacheKey, nil
}

func (c *Cache) Update(id string, payload app.UpdatePayload) (string, error) {
	cacheKey := hash(id + payload.Username)
	cacheValue, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}
	c.mc.Set(&memcache.Item{Key: cacheKey, Value: cacheValue})

	return cacheKey, nil
}

func (c *Cache) Get(id string) (WalletDetails, error) {
	if _, ok := s.Wallets[id]; !ok {
		return WalletDetails{
			Id:      "",
			Name:    "",
			Balance: 0,
		}, errors.New("notFound")
	}

	return s.Wallets[id].Details, nil
}
func (c *Cache) List() ([]WalletDetails, error) {
	var wlts []WalletDetails

	for _, w := range s.Wallets {
		wlts = append(wlts, w.Details)
	}
	return wlts, nil
}

func writeToFile(fileName string, newContent WalletDetails) error {
	newContentNormalized, _ := json.MarshalIndent(newContent, "", " ")
	err := ioutil.WriteFile(fileName, newContentNormalized, 0644)
	if err != nil {
		return errors.New("databaseError")
	}
	return nil
}
