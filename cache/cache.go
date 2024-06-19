package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Printf("SET key: (%s) to val: (%s)", key, value)
	c.data[string(key)] = value

	// It's a cache, will not hold data forever so delete the data
	go func() {
		// blocked till `ttl` amnt of time passes
		<-time.After(ttl)
		delete(c.data, string(key))
	}()

	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keyString := string(key)
	val, ok := c.data[keyString]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyString)
	}

	log.Printf("GET key: (%s), val = (%s)", keyString, string(val))
	return val, nil

}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(key))

	return nil
}
