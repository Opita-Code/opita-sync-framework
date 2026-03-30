package store

import "time"

type Item struct {
	Key       string        `json:"key"`
	Value     []byte        `json:"-"`
	TTL       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"created_at"`
}

type Service interface {
	Set(item Item) error
	Get(key string) (Item, bool, error)
	Delete(key string) error
}
