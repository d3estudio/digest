package processors

import (
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
	"gopkg.in/mgo.v2/bson"
)

// ProcessDeletedMessage is responsible for dealing with deletion events sent
// from the remote RTM server
func ProcessDeletedMessage(m *models.DigestMessageDeleted) {
	logger := log.WithField("timestamp", m.Timestamp)
	change, err := shared.Mongo.Items().RemoveAll(bson.M{"timestamp": m.Timestamp})
	if err != nil {
		logger.Error("Error removing item: ", err)
		return
	}
	logger.Debug("Removed ", change.Removed, " item(s)")
}
