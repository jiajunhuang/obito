package main

import (
	"log"
	"time"
)

func startRetryQueue() {
	for i := range retryQueue {
		go func(n *Notification) {
			n.mu.Lock()
			defer i.mu.Unlock()

			n.times++
			if n.times > *maxRetries {
				log.Printf("The given notify(%v) has exceed max retries, give up", n)
				return
			}

			time.Sleep(time.Duration(*sleepInterval*n.times) * time.Second)
			daemon <- n
		}(i)
	}
}
