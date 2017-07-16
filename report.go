package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type reportJSON struct {
	UUID        string `json:"uuid" binding:"required"`
	DeviceToken string `json:"device_token" binding:"required"`
}

// ReportInfo 上报设备信息
func ReportInfo(c *gin.Context) {
	var jsonDict reportJSON

	if c.BindJSON(&jsonDict) != nil {
		Fail(c, 400, nil, "Invalid json object.")
		return
	}
	conn := redisPool.Get()
	defer conn.Close()

	// set uuid and deviceToken with expiration time
	_, err := conn.Do("SETEX", genUUIDKey(jsonDict.UUID), (*lifetime)*3600*24, jsonDict.DeviceToken)
	if err != nil {
		log.Printf("set uuid:device(%s:%s) failed: %s", jsonDict.UUID, jsonDict.DeviceToken, err)
		return
	}
	Success(c, 200, nil, "")
	return
}
