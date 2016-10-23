package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/processor/processors"
	"github.com/d3estudio/digest/processor/utils"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"github.com/d3estudio/digest/shared/redis"
)

func main() {
	log.Info("Digest Processor")
	utils.Redis = redis.Client{}
	utils.Redis.Setup()
	shared.Mongo.Setup()

	// FIXME: We are using this dummy channel to keep everything up and running. An alternative is a waitgroup.
	dummyChannel := make(chan interface{})

	incomingMessageChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.DigestMessage{}, incomingMessageChannel)
	}()
	go func() {
		for {
			processors.ProcessIncomingMessage((<-incomingMessageChannel).(*models.DigestMessage))
		}
	}()

	deletedMessageChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.DigestMessageDeleted{}, deletedMessageChannel)
	}()
	go func() {
		for {
			processors.ProcessDeletedMessage((<-deletedMessageChannel).(*models.DigestMessageDeleted))
		}
	}()

	reactionChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.DigestMessageReaction{}, reactionChannel)
	}()
	go func() {
		for {
			processors.ProcessReaction((<-reactionChannel).(*models.DigestMessageReaction))
		}
	}()

	updatedChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.DigestMessageChanged{}, updatedChannel)
	}()

	go func() {
		for {
			processors.ProcessUpdatedMessage((<-updatedChannel).(*models.DigestMessageChanged))
		}
	}()

	_ = <-dummyChannel
}
