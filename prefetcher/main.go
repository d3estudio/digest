package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/prefetcher/sources"
	"github.com/d3estudio/digest/processor/utils"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"github.com/victorgama/go-unfurl"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	log.Info("Digest Prefetcher")
	utils.Redis.Setup()
	shared.Mongo.Setup()

	processors := []sources.Source{
		sources.YouTube{},
		sources.Vimeo{},
		sources.XKCD{},
		sources.Twitter{},
		sources.Spotify{},
		sources.PoorLink{},
	}

	requestChannel := make(chan interface{})
	go func() {
		utils.Redis.SubscribeToQueueOfType(models.PrefetchRequest{}, requestChannel)
	}()
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:51.0) Gecko/20100101 Firefox/51.0"
	unf := unfurl.NewClientWithOptions(unfurl.Options{
		UserAgent: &ua,
		MaxHops:   20,
	})
	for {
		r := (<-requestChannel).(*models.PrefetchRequest)
		var item models.DigestedMessage
		err := shared.Mongo.Items().Find(bson.M{"timestamp": r.Timestamp}).One(&item)
		if err != nil && err != mgo.ErrNotFound {
			log.WithField("facility", "unf").Error(err)
			continue
		}
		url, err := unf.Process(item.DetectedURL)
		if err != nil {
			log.WithField("facility", "unf").Error(err)
			continue
		}
		for _, proc := range processors {
			if proc.CanHandle(url) {
				result := proc.Process(url)
				if result == nil {
					continue
				}
				res, err := json.Marshal(result)
				if err != nil {
					log.WithField("facility", "marshaller").Error("Error encoding result: ", err)
					continue
				}
				err = shared.Mongo.Items().Update(bson.M{"timestamp": item.Timestamp}, bson.M{"$set": bson.M{"embededcontent": string(res)}})
				if err != nil {
					log.WithField("facility", "mongo").Error("Error updating item: ", err)
					continue
				}
				break
			}
		}
	}
}
