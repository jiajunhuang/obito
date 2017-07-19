package main

import (
	"flag"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

// Notification 带重传的消息通知
type Notification struct {
	mu    sync.Mutex
	apns  *apns2.Notification
	times int
}

var (
	daemon        chan *Notification
	retryQueue    chan *Notification
	sleepInterval = flag.Int("sleepInterval", 3, "sleep interval")
	maxRetries    = flag.Int("maxRetries", 3, "maximum retry time")
	cert          = flag.String("cert", "./cert.p12", "p12 cert path")
	passwd        = flag.String("passwd", "", "p12 cert password")
	redisURI      = flag.String("redisURI", "", "redis URI")
	prod          = flag.Bool("prod", false, "prod or develop env?")
	maxConsumer   = flag.Int("maxConsumer", 5, "how many consumer")
	maxQueueSize  = flag.Int("maxQueueSize", 4000, "max size in notification queue")
	maxIdle       = flag.Int("maxIdle", 100, "redis concurrency")
	maxActive     = flag.Int("maxActive", 12000, "redis maximum active connections amount")
	lifetime      = flag.Int("lifetime", 30, "life time of reported device token in days")
	step          = flag.Int("step", 1000, "step in loop")
	retryAfter    = flag.Int("retryAfter", 3, "retry after n seconds")
)

func main() {
	// parse cmd params
	flag.Parse()

	// initial redis
	initRedisPool()
	defer redisPool.Close()

	// load cert
	cert, err := certificate.FromP12File(*cert, *passwd)
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	// start consumers
	daemon = make(chan *Notification, *maxQueueSize)
	for i := 0; i < *maxConsumer; i++ {
		go startConsume(cert)
	}
	// start retry queue
	retryQueue = make(chan *Notification, *maxQueueSize)

	// start web server
	r := gin.Default()
	log.Printf("running gin...")
	r.POST("/push", Push)
	r.POST("/report", ReportInfo)
	r.POST("/tag", SetTag)
	r.POST("/tag/push", PushByTag)
	r.PUT("/badge/clear", ClearBadge)
	log.Printf("start server...")
	r.Run(":9999")
}
