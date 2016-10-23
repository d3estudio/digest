package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/emoji-manager/remote"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"github.com/d3estudio/digest/shared/redis"
	"gopkg.in/mgo.v2/bson"
)

func buildDatabase() {
	log.Info("Build operation started.")
	asset, err := dataDbBytes()
	if err != nil {
		log.Error("Build failed: ", err)
		return
	}
	result, err := remote.BuildEmojiDatabase(asset)
	if err != nil {
		log.Error("Build failed: ", err)
		return
	}
	bulk := shared.Mongo.Emojis().Bulk()
	operations := 0
	for _, i := range result {
		for _, alias := range i.Aliases {
			bulk.Upsert(bson.M{"name": alias}, models.SingleEmoji{Name: alias, URL: i.URL, Unicode: i.Unicode})
			operations++
			if operations == 999 {
				log.Info("Bulk operations limit reached. Flushing...")
				bulkResult, err := bulk.Run()
				if err != nil {
					log.Error("Bulk operation failed: , err")
					operations = 0
					continue
				}
				log.Info("Bulk operation completed. Matched: ", bulkResult.Matched, " Modified: ", bulkResult.Modified)
				bulk = shared.Mongo.Emojis().Bulk()
			}
		}
	}
	if operations != 0 {
		bulkResult, err := bulk.Run()
		if err != nil {
			log.Error("Bulk operation failed: , err")
		}
		log.Info("Bulk operation completed. Matched: ", bulkResult.Matched, " Modified: ", bulkResult.Modified)
	}
	log.Info("Build completed.")
}

func main() {
	log.Info("Digest Emoji Manager")
	shared.Mongo.Setup()
	redis := redis.Client{}
	redis.Setup()
	log.Info("Performing initial build...")
	buildDatabase()
	log.Info("Listening for updates...")
	ch := make(chan interface{})
	go func() {
		redis.SubscribeToQueueOfType(models.DigestEmojiChanged{}, ch)
	}()
	for {
		<-ch
		buildDatabase()
	}
}
