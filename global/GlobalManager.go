package global

import (
	"sync"
	"gopkg.in/mgo.v2"
	"github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/mongo"
)

var sId int = 0
var globalMap map[int]interface{} = map[int]interface{}{}
var sIdRwlock sync.RWMutex = sync.RWMutex{}
var globalMapRwlock sync.RWMutex = sync.RWMutex{}

func GetSessionId() int {
	sIdRwlock.Lock()
	defer sIdRwlock.Unlock()
	
	sId++
	if sId == 2147483647 {
		sId = 0
	}
	return sId
}

func GetGlobalAttr(sId int, attr string) interface{} {
	globalMapRwlock.RLock()
	defer globalMapRwlock.RUnlock()
	
	if globalMap[sId] != nil {
		objMap := globalMap[sId].(map[string]interface{})
		return objMap[attr]
	}
	return nil
}

func SetGlobalAttr(sId int, attr string, value interface{}) {
	globalMapRwlock.Lock()
	defer globalMapRwlock.Unlock()
	
	var objMap map[string]interface{}
	if globalMap[sId] != nil {
		objMap = globalMap[sId].(map[string]interface{})
	} else {
		objMap = map[string]interface{}{}
	}
	objMap[attr] = value
	globalMap[sId] = objMap
}

func GetConnection(sId int) (*mgo.Session, *mgo.Database) {
	if session, db, found := getConnection(sId); found {
		return session, db
	}
	return setConnection(sId)
}

func getConnection(sId int) (*mgo.Session, *mgo.Database, bool) {
	globalMapRwlock.RLock()
	defer globalMapRwlock.RUnlock()
	
	if globalMap[sId] != nil {
		objMap := globalMap[sId].(map[string]interface{})
		session := objMap["session"]
		db := objMap["db"]
		if session != nil && db != nil {
			return session.(*mgo.Session), db.(*mgo.Database), true
		}
	}
	return nil, nil, false
}

func setConnection(sId int) (*mgo.Session, *mgo.Database) {
	connectionFactory := mongo.GetInstance()
	session, db := connectionFactory.GetConnection()
	
	globalMapRwlock.Lock()
	defer globalMapRwlock.Unlock()

	if globalMap[sId] == nil {
		globalMap[sId] = map[string]interface{}{}
	}
	objMap := globalMap[sId].(map[string]interface{})
	objMap["session"] = session
	objMap["db"] = db
	globalMap[sId] = objMap
	return session, db
}

func GetTxnId(sId int) int {
	if txnId, found := getTxnId(sId); found {
		return txnId
	}
	return setTxnId(sId)
}

func getTxnId(sId int) (int, bool) {
	globalMapRwlock.RLock()
	defer globalMapRwlock.RUnlock()
	
	if globalMap[sId] != nil {
		objMap := globalMap[sId].(map[string]interface{})
		txnId := objMap["txnId"]
		if txnId != nil {
			return txnId.(int), true
		}
	}
	return 0, false
}

func setTxnId(sId int) int {
	_, db := GetConnection(sId)
	txnManager := TxnManager{db}
	txnId := txnManager.BeginTransaction()
	
	globalMapRwlock.Lock()
	defer globalMapRwlock.Unlock()

	if globalMap[sId] == nil {
		globalMap[sId] = map[string]interface{}{}
	}
	objMap := globalMap[sId].(map[string]interface{})
	objMap["txnId"] = txnId
	globalMap[sId] = objMap
	return txnId
}

/**
* 同时还要closeTxnId
*/
func CloseSession(sId int) {
	globalMapRwlock.Lock()
	defer globalMapRwlock.Unlock()
	
	if globalMap[sId] != nil {
		objMap := globalMap[sId].(map[string]interface{})
		session := objMap["session"]
		if session != nil {
			pSession := session.(*mgo.Session)
			pSession.Close()
		}
//		globalMap[sId] = nil
		delete(globalMap, sId)
	}
}

func RollbackTxn(sessionId int) {
	txnId := GetGlobalAttr(sessionId, "txnId")
	if txnId != nil {
		if x := recover(); x != nil {
			_, db := GetConnection(sessionId)
			txnManager := TxnManager{db}
			txnManager.Rollback(txnId.(int))
			panic(x)
		}
	}
}

func CommitTxn(sessionId int) {
	txnId := GetGlobalAttr(sessionId, "txnId")
	if txnId != nil {
		_, db := GetConnection(sessionId)
		txnManager := TxnManager{db}
		txnManager.Commit(txnId.(int))
	}
}


