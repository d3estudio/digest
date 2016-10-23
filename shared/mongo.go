package shared

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type mongo struct {
	Session *mgo.Session
}

var Mongo mongo

func init() {
	Mongo = mongo{}
}

func (m *mongo) Setup() (err error) {
	log.Debug("Connecting to mongo at ", Settings.MongoServer)
	sess, err := mgo.Dial(Settings.MongoServer)
	if err != nil {
		log.Error("Cannot connect to mongo server")
	}
	m.Session = sess

	m.Items().EnsureIndex(mgo.Index{
		Key:    []string{"timestamp"},
		Unique: true,
	})

	m.Emojis().EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	return
}

func (m *mongo) collection(name string) *mgo.Collection {
	return m.Session.DB(Settings.MongoDatabase).C(name)
}

func (m *mongo) Items() *mgo.Collection {
	return m.collection("items")
}

func (m *mongo) Emojis() *mgo.Collection {
	return m.collection("emojis")
}
