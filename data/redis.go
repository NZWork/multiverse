package data

import (
	"errors"
	"fmt"
	"time"

	"log"

	"github.com/garyburd/redigo/redis"
)

const (
	// MarkdownContentKeyPrefix for the prefix key of log storage
	MarkdownContentKeyPrefix = "markdown_tiki_file_"
	// RedisHost for setup
	RedisHost = "10.1.1.7"
	// RedisPort for setup
	RedisPort = 30005
	// RedisDB for setup
	RedisDB = 0
)

// GetMarkdownContentKey return markdown content key in redis
func GetMarkdownContentKey(token string) string {
	return MarkdownContentKeyPrefix + token
}

var redisPool *redis.Pool

// SetupRedis pool
func SetupRedis() {
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%v:%d", RedisHost, RedisPort),
			redis.DialConnectTimeout(2*time.Second),
			redis.DialDatabase(RedisDB),
		)

		if err != nil {
			log.Printf("Redis init: %v\n", err.Error())
			return nil, err
		}
		return c, err
	}, 10)
}

// Close redis pool
func Close() {
	redisPool.Close()
}

// SetContent stat
func SetContent(token string, content []byte) (err error) {
	r := redisPool.Get()
	defer r.Close()

	_, err = r.Do("SET", GetMarkdownContentKey(token), content)
	return
}

// GetContent stat of log
func GetContent(token string) (content []byte, err error) {
	r := redisPool.Get()
	defer r.Close()

	key := GetMarkdownContentKey(token)

	if exists, _ := redis.Bool(r.Do("EXISTS", key)); exists {
		content, err = redis.Bytes(r.Do("GET", key))
	} else {
		err = errors.New("token not found")
	}
	return
}
