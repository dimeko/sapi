package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
)

type Cache struct {
	mc *memcache.Client
}

func env() map[string]string {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		panic("Cannot find .env file")
	}
	return map[string]string{
		"host": os.Getenv("MEMCACHED_HOST"),
		"port": os.Getenv("MEMCACHED_PORT"),
	}
}

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum64())
}

func New() *Cache {
	connectionString := env()["host"] + ":" + env()["port"]
	log.Println("Connecting to cache:", connectionString)
	mc := memcache.New(connectionString)

	cache := &Cache{
		mc: mc,
	}

	return cache
}

func (c *Cache) Set(key string, payload interface{}) (string, error) {
	cacheKey := hash(key)
	cacheValue, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}
	c.mc.Set(&memcache.Item{Key: cacheKey, Value: cacheValue})

	return cacheKey, nil
}

func (c *Cache) Get(key string) (*memcache.Item, error) {
	cacheKey := hash(key)
	item, err := c.mc.Get(cacheKey)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *Cache) Remove(key string) error {
	cacheKey := hash(key)
	err := c.mc.Delete(cacheKey)
	if err != nil {
		return err
	}
	return nil
}
