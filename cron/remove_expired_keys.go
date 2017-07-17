package main

// iterating keys with pattern `itachi:taglist:*` in redis, and check if the
// key exists, if not, then delete it from the corresponding tag, and finally,
// delete the key itself.

import (
	"flag"
	"log"

	"github.com/garyburd/redigo/redis"
)

var (
    redisURI = flag.String("redisURI", "redis://127.0.0.1:6379/0", "redis URI")
)

func main() {
	conn, err := redis.DialURL(*redisURI)
	if err != nil {
		log.Panicf("failed to connect to redis: %s", err)
	}

	// iterating keys
	tagListKeys, err := redis.Strings(conn.Do("KEYS", "obito:taglist:*"))
	if err != nil {
		log.Panicf("failed to iterate with error: %s", err)
	}

	for _, tagList := range tagListKeys {
		uuid := tagList[14:] // "obito:taglist:"
		exists, _ := redis.Bool(conn.Do("EXISTS", "obito:uuid:"+uuid))
		if !exists {
			log.Printf("key(%s) not exist", tagList)
			// delete from tag and finally delete itself
			joinedTagList, _ := redis.Strings(conn.Do("ZRANGE", tagList, 0, -1))
			for _, joinedTag := range joinedTagList {
				log.Printf("gonna remove %s in tag %s", uuid, joinedTag)
				conn.Do("ZREM", "obito:tag:"+joinedTag, uuid)
			}
			log.Printf("gonna remove key %s", tagList)
			conn.Do("DEL", tagList)
		} else {
			log.Printf("key(%s) exist", tagList)
		}
	}
}
