package mongo

import (
	"fmt"
	"github.com/hongjinqiu/gometa/config"
	"github.com/hongjinqiu/gometa/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
)

func GetInstance() ConnectionFactory {
	return ConnectionFactory{}
}

type ConnectionFactory struct{}

func (c ConnectionFactory) GetConnection() (*mgo.Session, *mgo.Database) {
	session, err := mgo.Dial(config.String("MONGODB_ADDRESS"))
	if err != nil {
		panic(err)
	}
	//	session.SetMode(mgo.Monotonic, true)
	db := session.DB(config.String("MONGODB_DATABASE_NAME"))
	return session, db
}

func (c ConnectionFactory) GetSession() *mgo.Session {
	session, err := mgo.Dial(config.String("MONGODB_ADDRESS"))
	if err != nil {
		panic(err)
	}
	return session
}

func (c ConnectionFactory) GetDatabase(session *mgo.Session) *mgo.Database {
	//	session.SetMode(mgo.Monotonic, true)
	db := session.DB(config.String("MONGODB_DATABASE_NAME"))
	return db
}

func GetCollectionSequenceName(collection string) string {
	byte0 := collection[0]
	return strings.ToLower(string(byte0)) + collection[1:] + "Id"
}

func GetSequenceNo(db *mgo.Database, sequenceName string) int {
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"c": 1}},
		ReturnNew: true,
	}
	doc := map[string]interface{}{}
	_, err := db.C("counters").Find(bson.M{"_id": sequenceName}).Apply(change, &doc)
	if err != nil {
		log.Error("^^^^sequenceName is:" + sequenceName)
		panic(err)
	}
	idText := fmt.Sprint(doc["c"])
	id, err := strconv.Atoi(idText)
	if err != nil {
		panic(err)
	}
	return id
}
