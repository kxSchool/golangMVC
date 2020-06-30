package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

const (
	RedisURL            = "redis://20.10.1.51:6379"
	redisMaxIdle        = 300 //最大空闲连接数
	redisIdleTimeoutSec = 0   //最大空闲连接时间
	RedisPassword       = ""
)

// NewRedisPool 返回redis连接池
func init() {
	pool = &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(RedisURL)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			//if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
			//	return nil, fmt.Errorf("redis auth password error: %s", authErr)
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func set(k, v string) {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func getStringValue(k string) string {
	c := pool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		fmt.Println("Get Error: ", err.Error())
		return ""
	}
	return username
}

func SetKeyExpire(k string, ex int) {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func CheckKey(k string) bool {
	c := pool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return exist
	}
}

func DelKey(k string) error {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SetJson(k string, data interface{}) error {
	c := pool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func getJsonByte(k string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonGet, nil
}

func main() {
	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/pool", pools)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pools(w http.ResponseWriter, r *http.Request) {
	str := getStringValue("abc")
	fmt.Fprintf(w, "%s", str)
	//fmt.Fprintf(w, "%s", string(b))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
