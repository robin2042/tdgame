package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

type Oper interface {
	RPush(key string, argv ...interface{}) error
	LRange(key string, start, stop int64) ([]string, error)
}

// Storage struct is used to access database
type CloudStore struct {
	Oper
	rds *redis.Client
}

func ExampleNewClient() *CloudStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	return &CloudStore{
		rds: rdb,
	}

	// Output: PONG <nil>
}

func (c *CloudStore) RPush(key string, argv ...interface{}) error {
	return c.rds.RPush(ctx, key, argv...).Err()

}

func (c *CloudStore) LRange(key string, start, stop int64) ([]string, error) {
	return c.rds.LRange(ctx, key, start, stop).Result()
}
