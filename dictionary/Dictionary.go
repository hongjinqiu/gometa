package dictionary

import (
	"github.com/hongjinqiu/gometa/mongo"
	"fmt"
	"strconv"
	"sort"
	"gopkg.in/mgo.v2"
)

func GetInstance() DictionaryManager {
	return DictionaryManager{}
}

type DictionarySort struct {
	objLi []map[string]interface{}
}

func (o DictionarySort) Len() int {
	return len(o.objLi)
}

func (o DictionarySort) Less(i, j int) bool {
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

func (o DictionarySort) Swap(i, j int) {
	o.objLi[i], o.objLi[j] = o.objLi[j], o.objLi[i]
}


type DictionaryManager struct {
}

func (o DictionaryManager) GetDictionary(code string) map[string]interface{} {
	mongoDBFactory := mongo.GetInstance()
	session, db := mongoDBFactory.GetConnection()
	defer session.Close()
	
	return o.GetDictionaryBySession(db, code)
}

func (o DictionaryManager) GetDictionaryBySession(db *mgo.Database, code string) map[string]interface{} {
	collection := "PubDictionary"
	c := db.C(collection)
	
	queryMap := map[string]interface{}{
		"code": code,
	}
	
	result := map[string]interface{}{}
	err := c.Find(queryMap).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		panic(err)
	}
	
	// 排序
	o.sortDictionary(&result)
	return result
}

func (o DictionaryManager) sortDictionary(dictionary *map[string]interface{}) {
	items := (*dictionary)["items"]
	if items != nil {
		itemsLi := items.([]interface{})
		itemsMapLi := []map[string]interface{}{}
		for i,_ := range itemsLi {
			tmpObj := itemsLi[i].(map[string]interface{})
			itemsMapLi = append(itemsMapLi, tmpObj)
		}
		dSort := DictionarySort{objLi: itemsMapLi}
		sort.Sort(dSort)
		(*dictionary)["items"] = itemsMapLi
		
		for i,_ := range itemsMapLi {
			o.sortDictionary(&(itemsMapLi[i]))
		}
	}
}
