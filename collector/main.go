package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"github.com/d3estudio/digest/shared/redis"
	"github.com/d3estudio/digest/shared/slack"
)

// Collector is a simple dispatcher that parses incoming data from Slack and
// pushes it into a specific Redis queue that will be read by another process
// such as Processor or Prefetch.

func main() {
	log.Info("Digest Collector")
	client := slack.Client{}
	redis := redis.Client{}
	redis.Setup()
	client.Handshake(shared.Settings.Token)
	c := make(chan slack.RTMMessage)
	go client.Listen(c)
	for {
		incomingMessage := <-c
		switch incomingMessage.Type {
		case slack.TypeMessage:
			redis.PushItem(models.DigestMessage{
				Channel:   incomingMessage.Channel,
				User:      models.DigestUserFromSlack(incomingMessage.User),
				Text:      incomingMessage.Message,
				Timestamp: incomingMessage.Timestamp,
			})
		case slack.TypeMessageDeleted:
			redis.PushItem(models.DigestMessageDeleted{
				Timestamp: incomingMessage.DeletionTarget,
			})
		case slack.TypeReactionAdded:
			redis.PushItem(models.DigestMessageReaction{
				Reaction:  incomingMessage.Reaction,
				Timestamp: incomingMessage.Item.Timestamp,
				Added:     true,
			})
		case slack.TypeReactionRemoved:
			redis.PushItem(models.DigestMessageReaction{
				Reaction:  incomingMessage.Reaction,
				Timestamp: incomingMessage.Item.Timestamp,
				Added:     false,
			})
		case slack.TypeEmojiChanged:
			redis.PushItem(models.DigestEmojiChanged{
				Timestamp: incomingMessage.EventTimestamp,
			})
		case slack.TypeMessageChanged:
			redis.PushItem(models.DigestMessageChanged{
				DigestMessage: models.DigestMessage{
					Channel:   incomingMessage.Channel,
					User:      models.DigestUserFromSlack(incomingMessage.User),
					Text:      incomingMessage.Message,
					Timestamp: incomingMessage.Timestamp,
				},
			})
		}
	}
}
