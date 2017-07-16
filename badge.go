package main

import (
	"github.com/gin-gonic/gin"
)

type setAsReadJSON struct {
	UUID string `json:"uuid" binding:"required"`
}

// ClearBadge 清除badge
func ClearBadge(c *gin.Context) {
	conn := redisPool.Get()
	defer conn.Close()

	var jsonDict setAsReadJSON
	if c.BindJSON(&jsonDict) != nil {
		Fail(c, 400, nil, "Invalid json object.")
		return
	}
	conn.Do("DEL", genBadgeKey(jsonDict.UUID))
	Success(c, 200, nil, "")
}
