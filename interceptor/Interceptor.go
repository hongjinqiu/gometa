package interceptor

import (
	"github.com/hongjinqiu/gometa/mongo"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var rwlock sync.RWMutex = sync.RWMutex{}
var interceptorDict map[string]reflect.Type = map[string]reflect.Type{}

func init() {
	rwlock.Lock()
	defer rwlock.Unlock()
	interceptorDict[reflect.TypeOf(SysUserInterceptor{}).Name()] = reflect.TypeOf(SysUserInterceptor{})
	interceptorDict[reflect.TypeOf(ModelListTemplateInterceptor{}).Name()] = reflect.TypeOf(ModelListTemplateInterceptor{})
	interceptorDict[reflect.TypeOf(PubReferenceLogInterceptor{}).Name()] = reflect.TypeOf(PubReferenceLogInterceptor{})
	interceptorDict[reflect.TypeOf(AccountInOutItemInterceptor{}).Name()] = reflect.TypeOf(AccountInOutItemInterceptor{})
	interceptorDict[reflect.TypeOf(BbsPostReplyInterceptor{}).Name()] = reflect.TypeOf(BbsPostReplyInterceptor{})
	interceptorDict[reflect.TypeOf(BbsPostInterceptor{}).Name()] = reflect.TypeOf(BbsPostInterceptor{})
	interceptorDict[reflect.TypeOf(BbsPostAdminInterceptor{}).Name()] = reflect.TypeOf(BbsPostAdminInterceptor{})
}

type InterceptorCommon struct{}

func (qb InterceptorCommon) GetCreateUnitByUserId(session *mgo.Session, userId int) int {
	collectionName := "SysUser"
	query := map[string]interface{}{
		"_id": userId,
	}
	sysUser := qb.FindByMapWithSessionExact(session, collectionName, query)
	return qb.GetCreateUnitFromSysUser(sysUser)
}

func (qb InterceptorCommon) GetCreateUnitFromSysUser(sysUser map[string]interface{}) int {
	master := sysUser["A"].(map[string]interface{})
	createUnit, err := strconv.Atoi(fmt.Sprint(master["createUnit"]))
	if err != nil {
		panic(err)
	}
	return createUnit
}

func (qb InterceptorCommon) FindByMapWithSessionExact(session *mgo.Session, collection string, query map[string]interface{}) map[string]interface{} {
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

func (o InterceptorCommon) FindByMapWithSession(session *mgo.Session, collection string, query map[string]interface{}) (result map[string]interface{}, found bool) {
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

func (o InterceptorCommon) IndexWithSession(session *mgo.Session, collection string, query map[string]interface{}, pageNo int, pageSize int, orderBy string) (result map[string]interface{}) {
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

func GetInterceptorDict() map[string]reflect.Type {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return interceptorDict
}

type InterceptorManager struct{}

func (o InterceptorManager) ParseBeforeBuildQuery(sessionId int, classMethod string, paramMap map[string]string) map[string]string {
	if classMethod == "" {
		return paramMap
	}

	paramLi := []interface{}{}
	paramLi = append(paramLi, sessionId)
	paramLi = append(paramLi, paramMap)
	values := o.parse(classMethod, paramLi)
	if values != nil {
		return values[0].(map[string]string)
	}
	return paramMap
}

func (o InterceptorManager) ParseAfterBuildQuery(sessionId int, classMethod string, queryLi []map[string]interface{}) []map[string]interface{} {
	if classMethod == "" {
		return queryLi
	}

	paramLi := []interface{}{}
	paramLi = append(paramLi, sessionId)
	paramLi = append(paramLi, queryLi)
	values := o.parse(classMethod, paramLi)
	if values != nil {
		return values[0].([]map[string]interface{})
	}
	return queryLi
}

func (o InterceptorManager) ParseAfterQueryData(sessionId int, classMethod string, dataSetId string, items []interface{}) []interface{} {
	if classMethod == "" {
		return items
	}

	paramLi := []interface{}{}
	paramLi = append(paramLi, sessionId)
	paramLi = append(paramLi, dataSetId)
	paramLi = append(paramLi, items)
	values := o.parse(classMethod, paramLi)
	if values != nil {
		return values[0].([]interface{})
	}
	return items
}

func (o InterceptorManager) parse(classMethod string, param []interface{}) []interface{} {
	if classMethod != "" {
		exprContent := classMethod
		scriptStruct := strings.Split(exprContent, ".")[0]
		scriptStructMethod := strings.Split(exprContent, ".")[1]
		scriptType := GetInterceptorDict()[scriptStruct]
		if scriptType == nil {
			panic(scriptStruct + " is not exist")
		}
		inst := reflect.New(scriptType).Elem().Interface()
		instValue := reflect.ValueOf(inst)
		in := []reflect.Value{}
		for i, _ := range param {
			in = append(in, reflect.ValueOf(param[i]))
		}
		values := instValue.MethodByName(scriptStructMethod).Call(in)
		result := []interface{}{}
		for i, _ := range values {
			result = append(result, values[i].Interface())
		}
		return result
	}
	return nil
}
