package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(addr string) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var err error

	for i := 0; i < 30; i++ {
		if err = db.Ping(context.Background()).Err(); err == nil {
			return db
		}

		time.Sleep(time.Millisecond * 300)
	}

	panic(err)
}
