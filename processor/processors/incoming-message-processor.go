package processors

import (
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/processor/utils"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"github.com/mvdan/xurls"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ProcessIncomingMessage deals with incoming messages from the remote RTM
// server. It is also called by ProcessUpdatedMessage when an updated message
// is not on our records.
func ProcessIncomingMessage(m *models.DigestMessage) {
	if m.Text == "" {
		return
	}
	logger := log.WithField("timestamp", m.Timestamp)
	url := xurls.Strict.FindString(m.Text)
	if len(url) > 0 {
		url = attemptURLNormalization(url)
		var item *models.DigestedMessage
		collection := shared.Mongo.Items()
		err := collection.Find(bson.M{"ts": m.Timestamp}).One(&item)
		if err != nil && err != mgo.ErrNotFound {
			logger.Error("Error obtaining item from collection: ", err)
			return
		}
		if item == nil {
			digested := m.Digest()
			digested.DetectedURL = url
			err = collection.Insert(digested)
			if err != nil {
				logger.Error("Error inserting item into collection: ", err)
				return
			}
			logger.Debug("Requesting prefetch")
			utils.Redis.PushItem(models.PrefetchRequest{Timestamp: m.Timestamp})
		}
	}
}
