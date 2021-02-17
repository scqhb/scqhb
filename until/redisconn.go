package until

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Rdb = redis.NewClient(&redis.Options{
	Addr:     "9.9.9.99:6379",
	Password: "", // no password set
	DB:       2,  // use default DB
})
var Filepath = "./log.txt"

var ctx = context.Background()

func ExampleNewClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "9.9.9.99:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "9.9.9.99:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for i := uint(0); i < 10; i++ {
		go func() {
			fmt.Println("gorouting begin:", i)
			for t1 := uint(0); t1 < 10000000; t1++ {
				//生产72位字符串,两个36位拼接在一起
				s1 := BufferUuidv4()
				//写入redis
				err := rdb.Set(ctx, string(t1), s1, 0).Err()
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("gorouting end:", i)

		}()
	}
	fmt.Println("数据加载完成t0")

	///////////////////////

	val2, err := rdb.Get(ctx, "111").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("111", val2)
	}
	// Output: key value
	// key2 does not exist
}
