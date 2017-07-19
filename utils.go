package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
)

var (
	redisPool *redis.Pool
)

func initRedisPool() {
	redisPool = &redis.Pool{
		MaxIdle:   *maxIdle,
		MaxActive: *maxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(*redisURI)
			if err != nil {
				log.Panicf("connect to redis(%s) got error: %s", *redisURI, err)
			}
			return c, nil
		},
	}
}

func genUUIDKey(uuid string) string {
	return fmt.Sprintf("obito:uuid:%s", uuid)
}

func genTagKey(tagName string) string {
	return fmt.Sprintf("obito:tag:%s", tagName)
}

func genTagListKey(uuid string) string {
	return fmt.Sprintf("obito:taglist:%s", uuid)
}

func genBadgeKey(uuid string) string {
	return fmt.Sprintf("obito:badge:%s", uuid)
}

// Success 返回成功的json但是要自己return
func Success(c *gin.Context, code int, result map[string]interface{}, message string) {
	if result == nil {
		result = gin.H{}
	}
	c.JSON(
		code,
		gin.H{
			"result":  result,
			"message": message,
		},
	)
}

// Fail 返回失败的json， 但是要自己return
func Fail(c *gin.Context, code int, result map[string]interface{}, message string) {
	if result == nil {
		result = gin.H{}
	}
	c.JSON(
		code,
		gin.H{
			"result":  result,
			"message": message,
		},
	)
}

// false means failed, true in another way
func push(conn redis.Conn, deviceToken string, uuid string, content string) {
	badge, _ := redis.Int(conn.Do("INCR", genBadgeKey(uuid)))
	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken
	notification.Payload = payload.NewPayload().Alert(
		content,
	).Badge(
		badge,
	).Sound("default")
	log.Printf("gonna push to %s", deviceToken)
	go func() {
		daemon <- &Notification{apns: notification}
	}()
}
