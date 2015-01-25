package component

import (
	"github.com/hongjinqiu/gometa/mongo"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
)

type QuerySupport struct{}

func (qb QuerySupport) FindByMap(collection string, query map[string]interface{}) (result map[string]interface{}, found bool) {
	mongoDBFactory := mongo.GetInstance()
	session := mongoDBFactory.GetSession()
	defer session.Close()

	return qb.FindByMapWithSession(session, collection, query)
}

func (qb QuerySupport) FindByMapWithSession(session *mgo.Session, collection string, query map[string]interface{}) (result map[string]interface{}, found bool) {
	mongoDBFactory := mongo.GetInstance()
	db := mongoDBFactory.GetDatabase(session)
	c := db.C(collection)

	result = make(map[string]interface{})
	err := c.Find(query).One(&result)
	if err != nil {
		return result, false
	}

	return result, true
}

//func (qb QuerySupport) GetFixQueryByUserId(session *mgo.Session, userId int) map[string]interface{} {
//	collectionName := "SysUser"
//	query := map[string]interface{}{
//		"_id": userId,
//	}
//	sysUser := qb.FindByMapWithSessionExact(session, collectionName, query)
//	createUnit := qb.GetCreateUnitFromSysUser(sysUser)
//	return map[string]interface{}{
//		"A.createUnit": createUnit,
//	}
//}

func (qb QuerySupport) GetCreateUnitByUserId(session *mgo.Session, userId int) int {
	collectionName := "SysUser"
	query := map[string]interface{}{
		"_id": userId,
	}
	sysUser := qb.FindByMapWithSessionExact(session, collectionName, query)
	return qb.GetCreateUnitFromSysUser(sysUser)
}

func (qb QuerySupport) GetCreateUnitFromSysUser(sysUser map[string]interface{}) int {
	master := sysUser["A"].(map[string]interface{})
	createUnit, err := strconv.Atoi(fmt.Sprint(master["createUnit"]))
	if err != nil {
		panic(err)
	}
	return createUnit
}

func (qb QuerySupport) FindByMapWithSessionExact(session *mgo.Session, collection string, query map[string]interface{}) map[string]interface{} {
	result, found := qb.FindByMapWithSession(session, collection, query)
	if !found {
		queryByte, err := json.MarshalIndent(&query, "", "\t")
		if err != nil {
			panic(err)
		}
		panic("not found, query is:" + string(queryByte))
	}
	return result
}
/*
*/

func (qb QuerySupport) Find(collection string, query string) (result map[string]interface{}, found bool) {
	queryMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(query), &queryMap)
	if err != nil {
		panic(err)
	}

	return qb.FindByMap(collection, queryMap)
}

func (qb QuerySupport) Index(collection string, query map[string]interface{}, pageNo int, pageSize int, orderBy string) (result map[string]interface{}) {
	mongoDBFactory := mongo.GetInstance()
	session, _ := mongoDBFactory.GetConnection()
	defer session.Close()

	return qb.IndexWithSession(session, collection, query, pageNo, pageSize, orderBy)
}

func (qb QuerySupport) IndexWithSession(session *mgo.Session, collection string, query map[string]interface{}, pageNo int, pageSize int, orderBy string) (result map[string]interface{}) {
	mongoDBFactory := mongo.GetInstance()
	db := mongoDBFactory.GetDatabase(session)

	c := db.C(collection)

	items := []map[string]interface{}{}
	var err error
	if orderBy != "" {
		fieldLi := strings.Split(orderBy, ",")
		err = c.Find(query).Sort(fieldLi...).Limit(pageSize).Skip((pageNo - 1) * pageSize).All(&items)
	} else {
		err = c.Find(query).Limit(pageSize).Skip((pageNo - 1) * pageSize).All(&items)
	}
	if err != nil {
		panic(err)
	}

	totalResults, err := c.Find(query).Count()
	if err != nil {
		panic(err)
	}

	mapItems := []interface{}{}
	for _, item := range items {
		mapItems = append(mapItems, item)
	}
	return map[string]interface{}{
		"totalResults": totalResults,
		"items":        mapItems,
	}
}

func (qb QuerySupport) MapReduceAll(collection string, query map[string]interface{}, mapReduce mgo.MapReduce) (result []map[string]interface{}) {
	mongoDBFactory := mongo.GetInstance()
	session, db := mongoDBFactory.GetConnection()
	defer session.Close()

	result = []map[string]interface{}{}
	_, err := db.C(collection).Find(query).MapReduce(&mapReduce, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (qb QuerySupport) MapReduce(collection string, query map[string]interface{}, mapReduce mgo.MapReduce, pageNo int, pageSize int) (result []map[string]interface{}) {
	mongoDBFactory := mongo.GetInstance()
	session, db := mongoDBFactory.GetConnection()
	defer session.Close()

	result = []map[string]interface{}{}
	_, err := db.C(collection).Find(query).Limit(pageSize).Skip((pageNo-1)*pageSize).MapReduce(&mapReduce, &result)
	if err != nil {
		panic(err)
	}

	return result
}
