package handler

import (
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	. "github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/common"
	"fmt"
	"labix.org/v2/mgo"
	"log"
	"strconv"
	"encoding/json"
)

type UsedCheck struct{}

func (o UsedCheck) CheckUsed(sessionId int, dataSource DataSource, bo map[string]interface{}) bool {
	_, db := global.GetConnection(sessionId)
	masterData := bo["A"].(map[string]interface{})
	beReferenceQuery := []interface{}{
		dataSource.Id,
		"A",
		"id",
		masterData["id"],
	}
	count, err := db.C("PubReferenceLog").Find(map[string]interface{}{
		"beReference": beReferenceQuery,
	}).Limit(1).Count()
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (o UsedCheck) CheckDeleteDetailRecordUsed(sessionId int, dataSource DataSource, bo map[string]interface{}, diffDataRow DiffDataRow) bool {
	_, db := global.GetConnection(sessionId)
	fieldGroupLi := diffDataRow.FieldGroupLi
	destData := diffDataRow.DestData
	srcData := diffDataRow.SrcData
	if destData == nil && srcData != nil {// 删除的分录
		beReferenceQuery := []interface{}{
			dataSource.Id,
			fieldGroupLi[0].GetDataSetId(),
			"id",
			srcData["id"],
		}
		count, err := db.C("PubReferenceLog").Find(map[string]interface{}{
			"beReference": beReferenceQuery,
		}).Limit(1).Count()
		if err != nil {
			panic(err)
		}
		return count > 0
	}
	return false
}

func (o UsedCheck) Insert(sessionId int, fieldGroupLi []FieldGroup, bo *map[string]interface{}, data *map[string]interface{}) {
	_, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
	createTime := DateUtil{}.GetCurrentYyyyMMddHHmmss()
	for _, fieldGroup := range fieldGroupLi {
		if fieldGroup.IsRelationField() {
			modelTemplateFactory := ModelTemplateFactory{}
			relationItem, found := modelTemplateFactory.ParseRelationExpr(fieldGroup, *bo, *data)
			if !found {
				panic("数据源:" + fieldGroup.GetDataSource().Id + ",数据集:" + fieldGroup.GetDataSetId() + ",字段:" + fieldGroup.Id + ",配置的关联模型列表,不存在返回true的记录")
			}
			referenceData := map[string]interface{}{
				"A": map[string]interface{}{
					"createBy": (*data)["createBy"],
					"createTime": createTime,
					"createUnit": (*data)["createUnit"],
				},
				"reference":   o.GetSourceReferenceLi(db, fieldGroup, bo, data),
				"beReference": o.GetBeReferenceLi(db, fieldGroup, relationItem, data),
			}
			txnManager.Insert(txnId, "PubReferenceLog", referenceData)
		}
	}
}

func (o UsedCheck) Update(sessionId int, fieldGroupLi []FieldGroup, bo *map[string]interface{}, destData *map[string]interface{}, srcData map[string]interface{}) {
	if destData != nil && srcData == nil {
		o.Insert(sessionId, fieldGroupLi, bo, destData)
	} else if destData == nil && srcData != nil {
		o.Delete(sessionId, fieldGroupLi, srcData)
	} else if destData != nil && srcData != nil {
		modelTemplateFactory := ModelTemplateFactory{}
		fmt.Println(modelTemplateFactory.IsDataDifferent(fieldGroupLi, *destData, srcData))
		// 分析字段,如果字段都相等,不过帐,
		if modelTemplateFactory.IsDataDifferent(fieldGroupLi, *destData, srcData) {
			o.Delete(sessionId, fieldGroupLi, srcData)
			o.Insert(sessionId, fieldGroupLi, bo, destData)
		}
	}
}

/**
 * 不能直接用主数据集["ds", "A", "id", id]来直接删除,
 * 因为会连同分录的一起被删,而分录差异行数据做了different的判断,未修改的分录不会再补上被用记录,就漏掉了
*/
func (o UsedCheck) Delete(sessionId int, fieldGroupLi []FieldGroup, data map[string]interface{}) {
	for _, fieldGroup := range fieldGroupLi {
		if fieldGroup.IsRelationField() {
			srcDataSourceId := fieldGroup.GetDataSource().Id
			srcDataSetId := fieldGroup.GetDataSetId()
			srcFieldName := fieldGroup.Id
			fieldValue := data[srcFieldName]
			if fieldValue != nil {
				referenceQueryLi := []interface{}{}
				referenceQueryLi = append(referenceQueryLi, []interface{}{srcDataSourceId, srcDataSetId, "id", data["id"]})
				referenceQueryLi = append(referenceQueryLi, []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, fieldValue})
//				referenceQuery := []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, fieldValue}
				o.deleteReference(sessionId, referenceQueryLi)
			}
		}
	}
}

func (o UsedCheck) DeleteAll(sessionId int, fieldGroupLi []FieldGroup, data map[string]interface{}) {
	dataSource := fieldGroupLi[0].GetDataSource()
	id, err := strconv.Atoi(fmt.Sprint(data["id"]))
	if err != nil {
		panic(err)
	}
	referenceQuery := []interface{}{
		dataSource.Id,
		fieldGroupLi[0].GetDataSetId(),
		"id",
		id,
	}
	referenceQueryLi := []interface{}{referenceQuery}
	o.deleteReference(sessionId, referenceQueryLi)
}

func (o UsedCheck) deleteReference(sessionId int, referenceQueryLi []interface{}) {
	_, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
//	deleteQuery := map[string]interface{}{
//		"reference": referenceQueryLi,
//	}
	deleteQuery := map[string]interface{}{
		//"$and": referenceQueryLi,
	}
	andQuery := []interface{}{}
	for _, item := range referenceQueryLi {
		andQuery = append(andQuery, map[string]interface{}{
			"reference": item,
		})
	}
	deleteQuery["$and"] = andQuery
	deleteByte, err := json.MarshalIndent(&deleteQuery, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Println("deleteReference,collection:PubReferenceLog,query is:" + string(deleteByte))
	count, err := db.C("PubReferenceLog").Find(deleteQuery).Limit(1).Count()
	if err != nil {
		panic(err)
	}
	if count > 0 {
		_, result := txnManager.RemoveAll(txnId, "PubReferenceLog", deleteQuery)
		if !result {
			panic("删除失败")
		}
	}
}

//reference:[[dataSource, dataSet, fieldName, id], [dataSource, dataSet, fieldName, id]]
//beReference:[[dataSource, dataSet, fieldName, id], [dataSource, dataSet, fieldName, id]]
func (o UsedCheck) GetSourceReferenceLi(db *mgo.Database, fieldGroup FieldGroup, bo *map[string]interface{}, data *map[string]interface{}) []interface{} {
	masterData := (*bo)["A"].(map[string]interface{})
	sourceLi := []interface{}{}

	srcDataSourceId := fieldGroup.GetDataSource().Id
	//srcDataSetId := fieldGroup.GetDataSetId()
	srcDataSetId := "A"
	srcFieldName := "id"
	iId := fmt.Sprint(masterData["id"])
	id, err := strconv.Atoi(iId)
	if err != nil {
		panic(err)
	}
	refLi := []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, id}
	sourceLi = append(sourceLi, refLi)
	if fieldGroup.IsMasterField() {
		srcDataSourceId = fieldGroup.GetDataSource().Id
		srcDataSetId = fieldGroup.GetDataSetId()
		srcFieldName = fieldGroup.Id
		iId := fmt.Sprint(masterData[srcFieldName])
		id, err := strconv.Atoi(iId)
		if err != nil {
			panic(err)
		}
		refLi2 := []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, id}
		sourceLi = append(sourceLi, refLi2)
	} else {
		srcDataSourceId = fieldGroup.GetDataSource().Id
		srcDataSetId = fieldGroup.GetDataSetId()
		//dataSetData := (*bo)[srcDataSetId].(map[string]interface{})
		dataSetData := (*data)
		srcFieldName = "id"
		iId := fmt.Sprint(dataSetData["id"])
		id, err := strconv.Atoi(iId)
		if err != nil {
			panic(err)
		}
		refLi2 := []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, id}
		sourceLi = append(sourceLi, refLi2)

		srcDataSourceId = fieldGroup.GetDataSource().Id
		srcDataSetId = fieldGroup.GetDataSetId()
		//		dataSetData = (*bo)[srcDataSetId].(map[string]interface{})
		dataSetData = (*data)
		srcFieldName = fieldGroup.Id
		iId = fmt.Sprint(dataSetData[srcFieldName])
		id, err = strconv.Atoi(iId)
		if err != nil {
			panic(err)
		}
		refLi3 := []interface{}{srcDataSourceId, srcDataSetId, srcFieldName, id}
		sourceLi = append(sourceLi, refLi3)
	}
	return sourceLi
}

//func (o UsedCheck) getDataSetData() {
//
//}

//reference:[[dataSource, dataSet, fieldName, id], [dataSource, dataSet, fieldName, id]]
//beReference:[[dataSource, dataSet, fieldName, id], [dataSource, dataSet, fieldName, id]]
func (o UsedCheck) GetBeReferenceLi(db *mgo.Database, fieldGroup FieldGroup, relationItem RelationItem, data *map[string]interface{}) []interface{} {
	sourceLi := []interface{}{}
	relationId, err := strconv.Atoi(fmt.Sprint((*data)[fieldGroup.Id]))
	if err != nil {
		panic(err)
	}
	if relationItem.RelationDataSetId == "A" {
		sourceLi = append(sourceLi, []interface{}{relationItem.RelationModelId, "A", "id", relationId})
		return sourceLi
	}
	
	refData := map[string]interface{}{}
	query := map[string]interface{}{
		relationItem.RelationDataSetId + ".id": relationId,
	}
	//{"B.id": 2}
	err = db.C(relationItem.RelationModelId).Find(query).One(&refData)
	if err != nil {
		panic(err)
	}
	masterData := refData["A"].(map[string]interface{})
	masterDataId, err := strconv.Atoi(fmt.Sprint(masterData["id"]))
	if err != nil {
		panic(err)
	}
	sourceLi = append(sourceLi, []interface{}{relationItem.RelationModelId, "A", "id", masterDataId})

	sourceLi = append(sourceLi, []interface{}{relationItem.RelationModelId, relationItem.RelationDataSetId, "id", relationId})
	return sourceLi
}

func (o UsedCheck) GetFormUsedCheck(sessionId int, dataSource DataSource, bo map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	_, db := global.GetConnection(sessionId)
	modelIterator := ModelIterator{}
	var iteResult interface{} = ""
	modelIterator.IterateAllFieldBo(dataSource, &bo, &iteResult, func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, iteResult *interface{}) {
		if fieldGroup.Id == "id" {
			dataSetId := fieldGroup.GetDataSetId()
			referenceQuery := []interface{}{
				dataSource.Id,
				dataSetId,
				fieldGroup.Id,
				(*data)[fieldGroup.Id],
			}
			query := map[string]interface{}{
				"beReference": referenceQuery,
			}
			queryByte, err := json.MarshalIndent(query, "", "\t")
			if err != nil {
				panic(err)
			}
			log.Println("GetFormUsedCheck,collection:PubReferenceLog,query is:" + string(queryByte))
			count, err := db.C("PubReferenceLog").Find(query).Limit(1).Count()
			if err != nil {
				panic(err)
			}
			isUsed := count > 0
			if result[dataSetId] == nil {
				result[dataSetId] = map[string]interface{}{}
			}
			dataSetUsedMap := result[dataSetId].(map[string]interface{})
			dataSetUsedMap[fmt.Sprint((*data)[fieldGroup.Id])] = isUsed
			result[dataSetId] = dataSetUsedMap
		}
	})

	return result
}

func (o UsedCheck) GetListUsedCheck(sessionId int, dataSource DataSource, items []interface{}, dataSetId string) map[string]interface{} {
	result := map[string]interface{}{}
	_, db := global.GetConnection(sessionId)
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		dataSetData := itemMap
		referenceQuery := []interface{}{
			dataSource.Id,
			dataSetId,
			"id",
			dataSetData["id"],
		}
		queryMap := map[string]interface{}{
			"beReference": referenceQuery,
		}
		queryByte, err := json.MarshalIndent(&queryMap, "", "\t")
		if err != nil {
			panic(err)
		}
		log.Println("GetListUsedCheck,collection:PubReferenceLog,query is:" + string(queryByte))
		count, err := db.C("PubReferenceLog").Find(queryMap).Limit(1).Count()
		if err != nil {
			panic(err)
		}
		isUsed := count > 0
		if result[dataSetId] == nil {
			result[dataSetId] = map[string]interface{}{}
		}
		dataSetUsedMap := result[dataSetId].(map[string]interface{})
		dataSetUsedMap[fmt.Sprint(dataSetData["id"])] = isUsed
		result[dataSetId] = dataSetUsedMap
		/*
		if itemMap[dataSetId] != nil {
			dataSetData := itemMap[dataSetId].(map[string]interface{})
			referenceQuery := []interface{}{
				dataSource.Id,
				dataSetId,
				"id",
				dataSetData["id"],
			}
			count, err := db.C("PubReferenceLog").Find(map[string]interface{}{
				"reference": referenceQuery,
			}).Limit(1).Count()
			if err != nil {
				panic(err)
			}
			isUsed := count > 0
			if result[dataSetId] == nil {
				result[dataSetId] = map[string]interface{}{}
			}
			dataSetUsedMap := result[dataSetId].(map[string]interface{})
			dataSetUsedMap[fmt.Sprint(dataSetData["id"])] = isUsed
			result[dataSetId] = dataSetUsedMap
		}
		*/
	}
	return result
}
