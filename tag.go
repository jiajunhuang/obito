package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type tagJSON struct {
	UUID    string   `json:"uuid" binding:"required"`
	TagList []string `json:"tag_list" binding:"required"`
}

// SetTag 把uuid和tag绑定起来
func SetTag(c *gin.Context) {
	var jsonDict tagJSON

	if c.BindJSON(&jsonDict) != nil {
		Fail(c, 400, nil, "Invalid json object.")
		return
	}

	conn := redisPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", genUUIDKey(jsonDict.UUID)))
	if !exists {
		Fail(c, 400, nil, "The given UUID not valid.")
		return
	}

	conn.Send("MULTI")
	for _, tag := range jsonDict.TagList {
		conn.Send("ZADD", genTagListKey(jsonDict.UUID), "NX", 1, tag)
		conn.Send("ZADD", genTagKey(tag), "NX", 1, jsonDict.UUID)
	}
	conn.Send("EXEC")
	err = conn.Flush()
	if err != nil {
		log.Printf("failed to update tag: %s", err)
		Fail(c, 500, nil, "Failed to update tag.")
	}

	Success(c, 200, nil, "")
}

type pushByTagJSON struct {
	Tag     string `json:"tag" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// PushByTag 按tag推送
func PushByTag(c *gin.Context) {
	var jsonDict pushByTagJSON

	if c.BindJSON(&jsonDict) != nil {
		Fail(c, 400, nil, "Invalid json object.")
		return
	}

	conn := redisPool.Get()
	defer conn.Close()

	// 分批次推送，不过目前暂时先一次性取出来
	count, err := redis.Int(conn.Do("ZCARD", genTagKey(jsonDict.Tag)))
	if err != nil {
		log.Printf("failed to get by tag(%s) with error: %s", jsonDict.Tag, err)
		return
	}

	for i := 0; i < count; i += (*step) {
		// get uuid list
		uuidList, err := redis.Strings(conn.Do("ZRANGE", genTagKey(jsonDict.Tag), i, i+(*step)))
		if err != nil {
			log.Printf("failed to get by tag(%s) with error: %s", jsonDict.Tag, err)
			continue
		}
		// get device token list
		uuidTagList := make([]interface{}, *step)
		for i, uuid := range uuidList {
			uuidTagList[i] = genUUIDKey(uuid)
		}
		deviceTokenList, err := redis.Strings(conn.Do("MGET", uuidTagList...))
		// sent it to push daemon
		for i, uuid := range uuidList {
			push(conn, deviceTokenList[i], uuid, jsonDict.Content)
		}
	}

	Success(c, 200, nil, "")
}
