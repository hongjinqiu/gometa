package dictionary

import (
	"github.com/hongjinqiu/gometa/mongo"
	"github.com/hongjinqiu/gometa/global"
	"fmt"
	"strconv"
	"sort"
	"labix.org/v2/mgo"
	"encoding/json"
)

func GetProgramDictionaryInstance() ProgramDictionaryManager {
	return ProgramDictionaryManager{}
}

type ProgramDictionarySort struct {
	objLi []map[string]interface{}
}

func (o ProgramDictionarySort) Len() int {
	return len(o.objLi)
}

func (o ProgramDictionarySort) Less(i, j int) bool {
	orderI := o.objLi[i]["order"]
	if orderI == nil {
		return false
	}
	orderJ := o.objLi[j]["order"]
	if orderJ == nil {
		return false
	}

	order1, err := strconv.Atoi(fmt.Sprint(orderI))
	if err != nil {
		panic(err)
	}
	
	order2, err := strconv.Atoi(fmt.Sprint(orderJ))
	if err != nil {
		panic(err)
	}
	
	return order1 <= order2
}

func (o ProgramDictionarySort) Swap(i, j int) {
	o.objLi[i], o.objLi[j] = o.objLi[j], o.objLi[i]
}


type ProgramDictionaryManager struct {
}

func (o ProgramDictionaryManager) GetProgramDictionary(code string) map[string]interface{} {
	mongoDBFactory := mongo.GetInstance()
	session, db := mongoDBFactory.GetConnection()
	defer session.Close()
	
	sessionId := global.GetSessionId()
	defer global.CloseSession(sessionId)
	
	return o.GetProgramDictionaryBySession(sessionId, db, code)
}

func (o ProgramDictionaryManager) GetProgramDictionaryBySession(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	if code == "SYSUSER_TREE" {
		return o.GetSysUserProgramDictionary(sessionId, db, code)
	}
	if code == "ACCOUNTING_YEAR_START_TREE" {// 会计期年度开始
		return o.GetAccountingYearStartProgramDictionary(sessionId, db, code)
	}
	if code == "ACCOUNTING_YEAR_END_TREE" {// 会计期年度结束
		return o.GetAccountingYearEndProgramDictionary(sessionId, db, code)
	}
	if code == "ACCOUNTING_PERIOD_START_TREE" {// 会计期期间开始
		return o.GetAccountingPeriodStartProgramDictionary(sessionId, db, code)
	}
	if code == "ACCOUNTING_PERIOD_END_TREE" {// 会计期期间结束
		return o.GetAccountingPeriodEndProgramDictionary(sessionId, db, code)
	}
	
	return nil
}

func (o ProgramDictionaryManager) GetAccountingYearStartProgramDictionary(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	collection := "AccountingPeriod"
	c := db.C(collection)
	
	session, _ := global.GetConnection(sessionId)
	queryMap := map[string]interface{}{
		"A.createUnit": o.GetCreateUnitByUserId(session, userId),
	}
	
	itemResult := []map[string]interface{}{}
	err = c.Find(queryMap).All(&itemResult)
	if err != nil {
		panic(err)
	}
	
	result := map[string]interface{}{}
	result["code"] = code
	items := []interface{}{}
	for _, item := range itemResult {
		master := item["A"].(map[string]interface{})
		items = append(items, map[string]interface{}{
			"code": master["accountingYear"],
			"name": master["accountingYear"],
			"order": master["accountingYear"],
		})
	}
	result["items"] = items
	
	// 排序
	o.sortProgramDictionary(&result)
	return result
}

func (o ProgramDictionaryManager) GetAccountingYearEndProgramDictionary(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	return o.GetAccountingYearStartProgramDictionary(sessionId, db, code)
}

func (o ProgramDictionaryManager) GetAccountingPeriodStartProgramDictionary(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	collection := "AccountingPeriod"
	c := db.C(collection)
	
	session, _ := global.GetConnection(sessionId)
	queryMap := map[string]interface{}{
		"A.createUnit": o.GetCreateUnitByUserId(session, userId),
	}
	
	itemResult := []map[string]interface{}{}
	err = c.Find(queryMap).All(&itemResult)
	if err != nil {
		panic(err)
	}
	
	result := map[string]interface{}{}
	result["code"] = code
	items := []interface{}{}
	for _, item := range itemResult {
		detailLi := item["B"].([]interface{})
		for _, detail := range detailLi {
			detailMap := detail.(map[string]interface{})
			isIn := false
			for _, dictItem := range items {
				dictItemMap := dictItem.(map[string]interface{})
				if fmt.Sprint(dictItemMap["code"]) == fmt.Sprint(detailMap["sequenceNo"]) {
					isIn = true
					break
				}
			}
			if !isIn {
				items = append(items, map[string]interface{}{
					"code": detailMap["sequenceNo"],
					"name": detailMap["sequenceNo"],
					"order": detailMap["sequenceNo"],
				})
			}
		}
	}
	result["items"] = items
	
	// 排序
	o.sortProgramDictionary(&result)
	return result
}

func (o ProgramDictionaryManager) GetAccountingPeriodEndProgramDictionary(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	return o.GetAccountingPeriodStartProgramDictionary(sessionId, db, code)
}

func (o ProgramDictionaryManager) GetSysUserProgramDictionary(sessionId int, db *mgo.Database, code string) map[string]interface{} {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	collection := "SysUser"
	c := db.C(collection)
	
	session, _ := global.GetConnection(sessionId)
	queryMap := map[string]interface{}{
		"A.createUnit": o.GetCreateUnitByUserId(session, userId),
	}
	
	sysUserResult := []map[string]interface{}{}
	err = c.Find(queryMap).Limit(10).All(&sysUserResult)
	if err != nil {
		panic(err)
	}
	
	result := map[string]interface{}{}
	result["code"] = code
	//items := []map[string]interface{}{}
	items := []interface{}{}
	for idx, item := range sysUserResult {
		items = append(items, map[string]interface{}{
			"code": item["_id"],
			"name": item["nick"],
			"order": idx,
		})
	}
	result["items"] = items
	
	// 排序
	o.sortProgramDictionary(&result)
	return result
}

func (o ProgramDictionaryManager) sortProgramDictionary(programDictionary *map[string]interface{}) {
	items := (*programDictionary)["items"]
	if items != nil {
		itemsLi := items.([]interface{})
		itemsMapLi := []map[string]interface{}{}
		for i,_ := range itemsLi {
			tmpObj := itemsLi[i].(map[string]interface{})
			itemsMapLi = append(itemsMapLi, tmpObj)
		}
		dSort := ProgramDictionarySort{objLi: itemsMapLi}
		sort.Sort(dSort)
		(*programDictionary)["items"] = itemsMapLi
		
		for i,_ := range itemsMapLi {
			o.sortProgramDictionary(&(itemsMapLi[i]))
		}
	}
}

func (qb ProgramDictionaryManager) GetCreateUnitByUserId(session *mgo.Session, userId int) int {
	collectionName := "SysUser"
	query := map[string]interface{}{
		"_id": userId,
	}
	sysUser := qb.FindByMapWithSessionExact(session, collectionName, query)
	return qb.GetCreateUnitFromSysUser(sysUser)
}

func (qb ProgramDictionaryManager) GetCreateUnitFromSysUser(sysUser map[string]interface{}) int {
	master := sysUser["A"].(map[string]interface{})
	createUnit, err := strconv.Atoi(fmt.Sprint(master["createUnit"]))
	if err != nil {
		panic(err)
	}
	return createUnit
}

func (qb ProgramDictionaryManager) FindByMapWithSessionExact(session *mgo.Session, collection string, query map[string]interface{}) map[string]interface{} {
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

func (qb ProgramDictionaryManager) FindByMapWithSession(session *mgo.Session, collection string, query map[string]interface{}) (result map[string]interface{}, found bool) {
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
