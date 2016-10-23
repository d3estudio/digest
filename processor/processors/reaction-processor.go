package processors

import (
	"math"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
)

// ProcessReaction deals with incoming reaction events from the RTM server
// updating the current document. It also performs cleanup operations when
// a given reaction is no longer present.
func ProcessReaction(m *models.DigestMessageReaction) {
	logger := log.WithField("timestamp", m.Timestamp)
	collection := shared.Mongo.Items()
	m.Reaction = strings.Split(m.Reaction, "::")[0]
	delta := 1
	if !m.Added {
		delta = -1
	}
	logger.Debug("Updating reaction ", m.Reaction, " with delta ", delta)

	var item *models.DigestedMessage
	err := collection.Find(bson.M{"timestamp": m.Timestamp}).One(&item)
	if err == mgo.ErrNotFound {
		return
	} else if err != nil {
		logger.Error("Error querying collection: ", err)
		return
	}
	value, ok := item.Reactions[m.Reaction]
	if !ok {
		value = 0
	}
	value = int(math.Max(0, float64(value+delta)))
	item.Reactions[m.Reaction] = value

	var toDelete []string
	for k, v := range item.Reactions {
		if v <= 0 {
			toDelete = append(toDelete, k)
		}
	}

	for _, k := range toDelete {
		delete(item.Reactions, k)
	}

	err = collection.Update(bson.M{"timestamp": m.Timestamp}, item)
	if err != nil {
		logger.Error("Error updating document: ", err)
	}
}
