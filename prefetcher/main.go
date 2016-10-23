package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/processor/utils"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
)

func main() {
	log.Info("Digest Prefetcher")
	utils.Redis.Setup()
	shared.Mongo.Setup()

	requestChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.PrefetchRequest{}, requestChannel)
	}()

	for {
		r := (<-requestChannel).(*models.PrefetchRequest)
		fmt.Println(r)
	}
}
