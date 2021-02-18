package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "9.9.9.40:6379",
		Password: "oracle1234", // no password set
		DB:       0,            // use default DB
	})

	for i := 0; i < 100000000; i++ {
		err := rdb.Set(ctx, strconv.Itoa(i), "value", 0).Err()
		if err != nil {
			panic(err)
		}

	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func main() {
	ExampleClient()
}
