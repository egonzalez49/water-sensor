package subscriber

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	var address = os.Getenv("REDIS_ADDRESS")
	var password = os.Getenv("REDIS_PASSWORD")

	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	fmt.Println("Successfully connected to Redis.")
}
