package mongo

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
	"log"
	"encoding/json"
	"github.com/hongjinqiu/gometa/config"
)

const (
	RETRY_COUNT   = 2
	BEGIN_COMMIT  = "begin_commit"
	PREPARE       = "prepare"
	COMMIT        = "commit"
	ABORT         = "abort"
	DONE          = "done"
	CANCELLED     = "cancelled"
	DELETEFLAG    = "deleteFlag"
	SLEEP_TIME    = 500 * time.Millisecond
	TXN_PERIOD    = 1 * time.Minute
	RESUME_PERIOD = 10 * time.Minute
)

type TxnPeriodTask struct{}

func (o TxnPeriodTask) RunTxnPeriod() {
	connectionFactory := ConnectionFactory{}
	session, db := connectionFactory.GetConnection()
	defer session.Close()

	txnManager := TxnManager{db}
	txnManager.ResumePeriod()
	time.AfterFunc(TXN_PERIOD, o.RunTxnPeriod)
}

type TxnManager struct {
	DB *mgo.Database
}

// TODO, byTest
func (o TxnManager) BeginTransaction() int {
	c := o.DB.C("Transactions")
	seqName := GetCollectionSequenceName("Transactions")
	_id := GetSequenceNo(o.DB, seqName)
	transaction := map[string]interface{}{
		"_id":         _id,
		"id":          _id,
		"state":       BEGIN_COMMIT,
		"application": config.String("TXN_NAME"),
		"createTime":  o.GetCurrentDateTime(),
		"updateTime":  o.GetCurrentDateTime(),
	}
	if err := c.Insert(transaction); err != nil {
		panic(err)
	}
	change := map[string]interface{}{
		"$set": map[string]interface{}{
			"state":       PREPARE,
			"collections": []string{},
			"updateTime":  o.GetCurrentDateTime(),
		},
	}
	if err := o.DB.C("Transactions").Update(bson.M{"_id": _id}, change); err != nil {
		panic(err)
	}
	return _id
}

// TODO, byTest
func (o TxnManager) Commit(txnId int) {
	change := map[string]interface{}{
		"$set": map[string]interface{}{
			"state":      COMMIT,
			"updateTime": o.GetCurrentDateTime(),
		},
	}
	if err := o.DB.C("Transactions").Update(bson.M{"_id": txnId}, change); err != nil {
		panic(err)
	}
	doc := map[string]interface{}{}
	if err := o.DB.C("Transactions").Find(bson.M{"_id": txnId}).One(&doc); err != nil {
		panic(err)
	}
	collections := doc["collections"].([]interface{})
	for _, item := range collections {
		collection := item.(string)

		insertUpdateQuery := map[string]interface{}{
			"txnId": txnId,
			"$or": []interface{}{
				bson.M{
					DELETEFLAG: bson.M{
						"$exists": false,
					},
				},
				bson.M{
					DELETEFLAG: bson.M{
						"$ne": 9,
					},
				},
			},
		}
		insertUpdateModify := map[string]interface{}{
			"$set": map[string]interface{}{
				"pendingTransactions": []interface{}{},
			},
			"$unset": map[string]interface{}{
				"txnId":    1,
				DELETEFLAG: 1,
			},
		}
		if _, err := o.DB.C(collection).UpdateAll(insertUpdateQuery, insertUpdateModify); err != nil {
			panic(err)
		}
		deleteQuery := map[string]interface{}{
			"txnId":    txnId,
			DELETEFLAG: 9,
		}
		if _, err := o.DB.C(collection).RemoveAll(deleteQuery); err != nil {
			panic(err)
		}
	}

	change = map[string]interface{}{
		"$set": map[string]interface{}{
			"state":      DONE,
			"updateTime": o.GetCurrentDateTime(),
		},
	}
	if err := o.DB.C("Transactions").Update(bson.M{"_id": txnId}, change); err != nil {
		panic(err)
	}
}

// TODO, byTest
func (o TxnManager) Rollback(txnId int) {
	change := map[string]interface{}{
		"$set": map[string]interface{}{
			"state":      ABORT,
			"updateTime": o.GetCurrentDateTime(),
		},
	}
	if err := o.DB.C("Transactions").Update(bson.M{"_id": txnId}, change); err != nil {
		panic(err)
	}
	doc := map[string]interface{}{}
	if err := o.DB.C("Transactions").Find(bson.M{"_id": txnId}).One(&doc); err != nil {
		panic(err)
	}

	collections := doc["collections"].([]interface{})
	for _, item := range collections {
		collection := item.(string)
		docs := []map[string]interface{}{}
		if err := o.DB.C(collection).Find(bson.M{"txnId": txnId}).All(&docs); err != nil {
			panic(err)
		}
		for _, docItem := range docs {
			pendingTransactions := docItem["pendingTransactions"].([]interface{})
			for i := len(pendingTransactions) - 1; i >= 0; i-- {
				pendingTransaction := pendingTransactions[i].(map[string]interface{})
				var err interface{}
				if pendingTransaction["insert"] != nil {
					err = o.DB.C(collection).RemoveId(docItem["_id"])
				} else if pendingTransaction["update"] != nil {
					err = o.DB.C(collection).Update(bson.M{"_id": docItem["_id"]}, pendingTransaction["update"])
				} else if pendingTransaction["unModify"] != nil {
					lastPendingTransactions := pendingTransactions[0 : i+1]
					modify := map[string]interface{}{
						"$set": map[string]interface{}{
							"pendingTransactions": lastPendingTransactions,
						},
					}
					unModifyMap := pendingTransaction["unModify"].(map[string]interface{})
					for k, v := range unModifyMap {
						if k == "$set" {
							setMap := modify["$set"].(map[string]interface{})
							setMap[k] = v
						} else {
							modify[k] = v
						}
					}
					err = o.DB.C(collection).Update(bson.M{"_id": docItem["_id"]}, modify)
				} else if pendingTransaction["remove"] != nil {
					lastPendingTransactions := pendingTransactions[0 : i+1]
					modify := map[string]interface{}{
						"$set": map[string]interface{}{
							"pendingTransactions": lastPendingTransactions,
						},
						"$unset": bson.M{DELETEFLAG: 1},
					}
					err = o.DB.C(collection).Update(bson.M{"_id": docItem["_id"]}, modify)
				}
				if err != nil {
					panic(err)
				}
			}
		}
		change := bson.M{
			"$unset": bson.M{
				DELETEFLAG: 1,
				"txnId":    txnId,
			},
			"$set": bson.M{
				"pendingTransactions": []interface{}{},
			},
		}
		if _, err := o.DB.C(collection).UpdateAll(bson.M{"txnId": txnId}, change); err != nil {
			panic(err)
		}
	}

	change = map[string]interface{}{
		"$set": map[string]interface{}{
			"state":      CANCELLED,
			"updateTime": o.GetCurrentDateTime(),
		},
	}
	if err := o.DB.C("Transactions").Update(bson.M{"_id": txnId}, change); err != nil {
		panic(err)
	}
}

/*
func throwsPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f() //执行函数f，如果f中出现了panic，那么就可以恢复回来
	return
}
*/

// TODO, byTest,
func (o TxnManager) Insert(txnId int, collection string, doc map[string]interface{}) map[string]interface{} {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	seqName := GetCollectionSequenceName(collection)
	strId := fmt.Sprint(doc["_id"])
	if doc["_id"] == nil || strId == "" || strId == "0" {
		sequenceNo := GetSequenceNo(o.DB, seqName)
		doc["_id"] = sequenceNo
		doc["id"] = sequenceNo
	}
	doc["txnId"] = txnId
	doc["pendingTransactions"] = []interface{}{
		map[string]interface{}{
			//"txnId":  txnId,
			"insert": true,
		},
	}
	if err := o.DB.C(collection).Insert(doc); err != nil {
		panic(err)
	}
	return doc
}

// TODO, byTest
func (o TxnManager) Update(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	for i := 0; i < RETRY_COUNT; i++ {
		if result, found := o.update(txnId, collection, doc); found {
			return result, true
		}
		time.Sleep(SLEEP_TIME)
	}
	return nil, false
}

// TODO, byTest,单条
func (o TxnManager) update(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	if oldDoc, found := o.selectForUpdate(txnId, collection, doc); found {
		query := map[string]interface{}{
			"_id": doc["_id"],
		}
		doc["txnId"] = txnId
		var pendingTransactions []interface{}
		if oldDoc["pendingTransactions"] != nil {
			pendingTransactions = oldDoc["pendingTransactions"].([]interface{})
		} else {
			pendingTransactions = []interface{}{}
		}
		if len(pendingTransactions) == 0 {
			pendingTransactions = append(pendingTransactions, map[string]interface{}{
				"update": oldDoc,
			})
			doc["pendingTransactions"] = pendingTransactions
		}
		if err := o.DB.C(collection).Update(query, doc); err != nil {
			if err == mgo.ErrNotFound {
				return nil, false
			}
			panic(err)
		}
		result := map[string]interface{}{}
		if err := o.DB.C(collection).Find(query).One(&result); err != nil {
			if err == mgo.ErrNotFound {
				return nil, false
			}
			panic(err)
		}
		return result, true
	}
	return nil, false
}

func (o TxnManager) SelectForUpdate(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	return o.selectForUpdate(txnId, collection, doc)
}

// TODO, byTest
func (o TxnManager) selectForUpdate(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	query := map[string]interface{}{
		"_id": doc["_id"],
	}
	return o.selectMultiForUpdate(txnId, collection, query)
}

/**
多条,返回第一条
*/
// TODO, byTest
func (o TxnManager) selectMultiForUpdate(txnId int, collection string, query map[string]interface{}) (map[string]interface{}, bool) {
	change := map[string]interface{}{
		"$set": map[string]interface{}{
			"txnId": txnId,
		},
	}
	txnQuery := bson.M{
		"$or": []interface{}{
			bson.M{
				"txnId": bson.M{
					"$exists": false,
				},
			},
			bson.M{
				"txnId": txnId,
			},
		},
	}
	for k, v := range query {
		txnQuery[k] = v
	}
	changeInfo, err := o.DB.C(collection).UpdateAll(txnQuery, change)
	if err != nil {
		panic(err)
	}
	if changeInfo.Updated > 0 {
		afterUpdateQuery := map[string]interface{}{
			"txnId": txnId,
		}
		for k, v := range query {
			afterUpdateQuery[k] = v
		}
		result := map[string]interface{}{}
		if err = o.DB.C(collection).Find(afterUpdateQuery).One(&result); err != nil {
			panic(err)
		}
		return result, true
	}
	return nil, false
}

// TODO, byTest
/**
atomic for every found document,
but not atomic for all found document
*/
func (o TxnManager) UpdateAll(txnId int, collection string, query map[string]interface{}, update map[string]interface{}, unModify map[string]interface{}) (map[string]interface{}, bool) {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	for i := 0; i < RETRY_COUNT; i++ {
		if result, found := o.updateAll(txnId, collection, query, update, unModify); found {
			return result, true
		}
		time.Sleep(SLEEP_TIME)
	}
	return nil, false
}

// TODO, byTest
func (o TxnManager) updateAll(txnId int, collection string, query map[string]interface{}, update map[string]interface{}, unModify map[string]interface{}) (map[string]interface{}, bool) {
	if _, found := o.selectMultiForUpdate(txnId, collection, query); found {
		changeQuery := map[string]interface{}{
			"txnId": txnId,
		}
		for k, v := range query {
			changeQuery[k] = v
		}

		changeUpdate := map[string]interface{}{
			"$push": map[string]interface{}{
				"pendingTransactions": map[string]interface{}{
					"unModify": unModify,
				},
			},
		}
		for k, v := range update {
			if k == "$push" {
				updateItem := changeUpdate[k].(map[string]interface{})
				updateItem[k] = v
			} else {
				changeUpdate[k] = v
			}
		}
		changeInfo, err := o.DB.C(collection).UpdateAll(changeQuery, changeUpdate)
		if err != nil {
			panic(err)
		}
		if changeInfo.Updated > 0 {
			result := map[string]interface{}{}
			if err = o.DB.C(collection).Find(changeQuery).One(&result); err != nil {
				panic(err)
			}
			return result, true
		}
		return nil, false
	}
	return nil, false
}

//  Update(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
func (o TxnManager) Remove(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	for i := 0; i < RETRY_COUNT; i++ {
		if result, found := o.remove(txnId, collection, doc); found {
			return result, true
		}
		time.Sleep(SLEEP_TIME)
	}
	return nil, false
}

// TODO, byTest,单条
func (o TxnManager) remove(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	if oldDoc, found := o.selectForUpdate(txnId, collection, doc); found {
		query := map[string]interface{}{
			"_id": doc["_id"],
		}
		var pendingTransactions []interface{}
		if oldDoc["pendingTransactions"] != nil {
			pendingTransactions = oldDoc["pendingTransactions"].([]interface{})
		} else {
			pendingTransactions = []interface{}{}
		}
		if len(pendingTransactions) == 0 {
			pendingTransactions = append(pendingTransactions, map[string]interface{}{
				"remove": true,
			})
		}
		update := map[string]interface{}{
			"$set": map[string]interface{}{
				"txnId":               txnId,
				DELETEFLAG:            9,
				"pendingTransactions": pendingTransactions,
			},
		}
		if err := o.DB.C(collection).Update(query, update); err != nil {
			if err == mgo.ErrNotFound {
				println("remove return pos0")
				return nil, false
			}
			panic(err)
		}
		result := map[string]interface{}{}
		if err := o.DB.C(collection).Find(query).One(&result); err != nil {
			if err == mgo.ErrNotFound {
				println("remove return pos1")
				return nil, false
			}
			panic(err)
		}
		println("remove return pos2")
		return result, true
	}
	println("remove return pos3")
	return nil, false
}

// TODO, byTest
/**
atomic for every found document,
but not atomic for all found document
*/
func (o TxnManager) RemoveAll(txnId int, collection string, query map[string]interface{}) (map[string]interface{}, bool) {
	o.pubTxnCollectionIfNotExist(txnId, collection)

	for i := 0; i < RETRY_COUNT; i++ {
		if result, found := o.removeAll(txnId, collection, query); found {
			return result, true
		}
		time.Sleep(SLEEP_TIME)
	}
	return nil, false
}

// TODO, byTest
func (o TxnManager) removeAll(txnId int, collection string, query map[string]interface{}) (map[string]interface{}, bool) {
	if _, found := o.selectMultiForUpdate(txnId, collection, query); found {
		changeQuery := map[string]interface{}{
			"txnId": txnId,
		}
		for k, v := range query {
			changeQuery[k] = v
		}

		changeUpdate := map[string]interface{}{
			"$set": map[string]interface{}{
				DELETEFLAG: 9,
			},
			"$push": map[string]interface{}{
				"pendingTransactions": map[string]interface{}{
					"remove": true,
				},
			},
		}
		changeInfo, err := o.DB.C(collection).UpdateAll(changeQuery, changeUpdate)
		if err != nil {
			panic(err)
		}
		if changeInfo.Updated > 0 {
			result := map[string]interface{}{}
			if err = o.DB.C(collection).Find(changeQuery).One(&result); err != nil {
				panic(err)
			}
			return result, true
		}
		return nil, false
	}
	return nil, false
}

/**
系统启动时运行,从日志中恢复事务
*/
//(begin_commit|prepare|commit|abort|done|cancelled)
func (o TxnManager) Resume() {
	txnLi := []map[string]interface{}{}
	query := bson.M{
		"state": bson.M{
			"$in": []string{BEGIN_COMMIT, PREPARE, COMMIT, ABORT},
		},
		"application": config.String("TXN_NAME"),
	}
	if err := o.DB.C("Transactions").Find(query).All(&txnLi); err != nil {
		panic(err)
	}

	txnParam := []interface{}{}
	for _, item := range txnLi {
		txnParam = append(txnParam, item)
	}
	o.finishTxn(txnParam)
}

/**
程序里面对长时间(10分钟)没反应的事务做的操作
*/
// 协调者的时间更新,
func (o TxnManager) ResumePeriod() {
	t := time.Now()
	t = t.Add(-RESUME_PERIOD)
	txnLi := []map[string]interface{}{}
	query := bson.M{
		"state": bson.M{
			"$in": []string{BEGIN_COMMIT, PREPARE, COMMIT, ABORT},
		},
		"application": config.String("TXN_NAME"),
		"updateTime": bson.M{
			"$lt": o.getYmdhms(t),
		},
	}
	if err := o.DB.C("Transactions").Find(query).All(&txnLi); err != nil {
		panic(err)
	}
	txnParam := []interface{}{}
	for _, item := range txnLi {
		txnParam = append(txnParam, item)
	}
	o.finishTxn(txnParam)
}

func (o TxnManager) finishTxn(txnLi []interface{}) {
	for _, item := range txnLi {
		record := item.(map[string]interface{})
		txnId, err := strconv.Atoi(fmt.Sprint(record["_id"]))
		if err != nil {
			panic(err)
		}
		state := record["state"].(string)
		if state == BEGIN_COMMIT {
			record["state"] = CANCELLED
			if err := o.DB.C("Transactions").Update(bson.M{"_id": record["_id"]}, record); err != nil {
				panic(err)
			}
		} else if state == PREPARE {
			o.Rollback(txnId)
		} else if state == COMMIT {
			o.Commit(txnId)
		} else if state == ABORT {
			o.Rollback(txnId)
		}
	}
}

func (o TxnManager) pubTxnCollectionIfNotExist(txnId int, collection string) {
	query := map[string]interface{}{
		"_id": txnId,
		"collections": map[string]interface{}{
			"$nin": []string{collection},
		},
	}
	update := map[string]interface{}{
		"$push": map[string]interface{}{
			"collections": collection,
		},
	}
	updateByte, err := json.MarshalIndent(&update, "", "\t")
	if err != nil {
		panic(err)
	}
	queryByte, err := json.MarshalIndent(&query, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Println("pubTxnCollectionIfNotExist,update Transactions update:" + string(updateByte) + ", query:" + string(queryByte))
	if err := o.DB.C("Transactions").Update(query, update); err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
	}
}

func (o TxnManager) GetCurrentDateTime() int64 {
	return o.getYmdhms(time.Now())
}

func (o TxnManager) getYmdhms(t time.Time) int64 {
	formatStr := t.Format("20060102150405")
	value, err := strconv.ParseInt(formatStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return value
}
