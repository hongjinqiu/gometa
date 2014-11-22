package layer

import (
	"github.com/hongjinqiu/gometa/dictionary"
	"labix.org/v2/mgo"
	"github.com/hongjinqiu/gometa/mongo"
	"github.com/hongjinqiu/gometa/global"
)

func GetInstance() LayerManager {
	return LayerManager{}
}

type LayerManager struct {
}

func (o LayerManager) GetLayer(code string) map[string]interface{} {
	connectionFactory := mongo.GetInstance()
	session, db := connectionFactory.GetConnection()
	defer session.Close()

	sessionId := global.GetSessionId()
	defer global.CloseSession(sessionId)
	
	return o.GetLayerBySession(sessionId, db, code)
}

func (o LayerManager) GetLayerBySession(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	dictionaryManager := dictionary.GetInstance()
	result := dictionaryManager.GetDictionaryBySession(db, code)
	if result == nil {
		programDictionaryManager := dictionary.GetProgramDictionaryInstance()
		result = programDictionaryManager.GetProgramDictionaryBySession(sessionId, db, code)
	}
	return result
}

