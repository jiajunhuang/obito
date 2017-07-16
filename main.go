package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

var (
	daemon       chan *apns2.Notification
	cert         = flag.String("cert", "./cert.p12", "p12 cert path")
	passwd       = flag.String("passwd", "", "p12 cert password")
	redisURI     = flag.String("redisURI", "", "redis URI")
	prod         = flag.Bool("prod", false, "prod or develop env?")
	maxConsumer  = flag.Int("maxConsumer", 5, "how many consumer")
	maxQueueSize = flag.Int("maxQueueSize", 4000, "max size in notification queue")
	maxIdle      = flag.Int("maxIdle", 100, "redis concurrency")
	maxActive    = flag.Int("maxActive", 12000, "redis maximum active connections amount")
	lifetime     = flag.Int("lifetime", 30, "life time of reported device token in days")
	step         = flag.Int("step", 1000, "step in loop")
	retryAfter   = flag.Int("retryAfter", 3, "retry after n seconds")
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
	daemon = make(chan *apns2.Notification, *maxQueueSize)
	for i := 0; i < *maxConsumer; i++ {
		go startConsume(cert)
	}

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
