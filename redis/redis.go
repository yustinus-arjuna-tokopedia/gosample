package redis

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
)

var Pool *redis.Pool

func Init() {
	redisHost := "devel-redis.tkpd:6379"
	Pool = newPool(redisHost)
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("Ping")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Expire(key string, time int) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func Set(key string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func Get(key string) (string, error) {

	conn := Pool.Get()
	defer conn.Close()
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return data, err
	}
	return data, err
}
