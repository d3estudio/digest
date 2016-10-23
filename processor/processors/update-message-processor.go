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

// ProcessUpdatedMessage checks an updated message consistency and updates
// the detected URL, before enqueuing the model for a Prefetch operation. If the
// updated text did not had an URL, it won't be found on the collection, which
// will cause this method to handle the incoming data to ProcessIncomingMessage.
func ProcessUpdatedMessage(m *models.DigestMessageChanged) {
	if m.Text == "" {
		return
	}
	logger := log.WithField("timestamp", m.Timestamp)
	url := xurls.Strict.FindString(m.Text)
	collection := shared.Mongo.Items()
	if len(url) > 0 {

		var item *models.DigestedMessage
		err := collection.Find(bson.M{"ts": m.Timestamp}).One(&item)
		if err != nil && err == mgo.ErrNotFound {
			logger.Debug("Message does not exist. Handling to ProcessIncomingMessage and returning...")
			ProcessIncomingMessage(&m.DigestMessage)
			return
		} else if err != nil && err != mgo.ErrNotFound {
			logger.Error("Error obtaining item from collection: ", err)
			return
		}
		logger.Debug("Message has changed and still contains a link. Processing...")
		url = attemptURLNormalization(url)
		digested := m.Digest()
		digested.EmbededContent = item.EmbededContent
		digested.DetectedURL = url
		err = collection.Update(bson.M{"timestamp": item.Timestamp}, item)
		if err != nil {
			logger.Error("Error updating item: ", err)
			return
		}
		logger.Debug("Requesting prefetch")
		utils.Redis.PushItem(models.PrefetchRequest{Timestamp: m.Timestamp})
	} else {
		logger.Debug("Message has changed and does not contain a link. Removing...")
		err := collection.Remove(bson.M{"timestamp": m.Timestamp})
		if err != nil && err != mgo.ErrNotFound {
			logger.Error("Removal failed: ", err)
			return
		}
		logger.Debug("Removal succeeded.")
	}
}
