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
	Del(key string)
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

func (c *CloudStore) Del(key string) {
	c.rds.Del(ctx, key)
}

//获取值
func (c *CloudStore) GetValue(key string) (string, error) {
	val, err := c.rds.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}
	return val, nil
}

//获取list
func (c *CloudStore) GetLrange(key string, start, end int64) ([]string, error) {

	val, err := c.rds.LRange(ctx, key, start, end).Result()
	return val, err
}

//incr 自增
func (c *CloudStore) Incr(key string) {
	c.rds.Incr(ctx, key)
}
