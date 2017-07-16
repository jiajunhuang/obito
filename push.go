package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type pushJSON struct {
	UUID    string `json:"uuid" binding:"required"` // 就是我们分配下去的UUID
	Content string `json:"content" binding:"required"`
}

// Push 负责推送的handler
func Push(c *gin.Context) {
	conn := redisPool.Get()
	defer conn.Close()

	var jsonDict pushJSON
	if c.BindJSON(&jsonDict) != nil {
		Fail(c, 400, nil, "Invalid json object.")
		return
	}

	deviceToken, _ := redis.String(conn.Do("GET", genUUIDKey(jsonDict.UUID)))
	if deviceToken == "" {
		Fail(c, 400, nil, fmt.Sprintf("push to %s failed", jsonDict.UUID))
		return
	}

	push(conn, deviceToken, jsonDict.UUID, jsonDict.Content)

	Success(c, 200, nil, "")
}
