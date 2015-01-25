package controllers

//import "github.com/robfig/revel"
import (
	. "github.com/hongjinqiu/gometa/component"
	. "github.com/hongjinqiu/gometa/error"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	. "github.com/hongjinqiu/gometa/model/handler"
	"github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/script"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

type FinanceService struct{}

func (o FinanceService) SaveData(sessionId int, dataSource DataSource, bo *map[string]interface{}) *[]DiffDataRow {
	modelTemplateFactory := ModelTemplateFactory{}

	strId := modelTemplateFactory.GetStrId(*bo)

	modelTemplateFactory.ConvertDataType(dataSource, bo)
	// 主数据集和分录数据校验
	message := o.validateBO(sessionId, dataSource, (*bo))
	if message != "" {
		panic(BusinessError{Message: message})
	}
	_, db := global.GetConnection(sessionId)

	modelIterator := ModelIterator{}
	var result interface{} = ""

	if strId == "" || strId == "0" {
		// 主数据集和分录id赋值,
		modelIterator.IterateAllFieldBo(dataSource, bo, &result, func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
			o.setDataId(db, dataSource, &fieldGroup, bo, data)
		})
		// 被用过帐
		usedCheck := UsedCheck{}
		result = ""
		diffDataRowLi := []DiffDataRow{}
		modelIterator.IterateDataBo(dataSource, bo, &result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
			diffDataRowLi = append(diffDataRowLi, DiffDataRow{
				FieldGroupLi: fieldGroupLi,
				DestBo:       bo,
				DestData:     data,
				SrcData:      nil,
				SrcBo:        nil,
			})
		})
		for i, _ := range diffDataRowLi {
			fieldGroupLi := diffDataRowLi[i].FieldGroupLi
			data := diffDataRowLi[i].DestData
			usedCheck.Insert(sessionId, fieldGroupLi, bo, data)
		}
		o.deleteExtraField(dataSource, bo)
		txnManager := TxnManager{db}
		txnId := global.GetTxnId(sessionId)
		collectionName := modelTemplateFactory.GetCollectionName(dataSource)
		txnManager.Insert(txnId, collectionName, *bo)

		return &diffDataRowLi
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		panic(err)
	}
	srcBo := map[string]interface{}{}
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	err = db.C(collectionName).Find(map[string]interface{}{"_id": id}).One(&srcBo)
	if err != nil {
		panic(err)
	}

	modelTemplateFactory.ConvertDataType(dataSource, &srcBo)
	diffDataRowLi := []DiffDataRow{}
	modelIterator.IterateDiffBo(dataSource, bo, srcBo, &result, func(fieldGroupLi []FieldGroup, destData *map[string]interface{}, srcData map[string]interface{}, result *interface{}) {
		// 分录+id
		if destData != nil {
			dataStrId := fmt.Sprint((*destData)["id"])
			if dataStrId == "" || dataStrId == "0" {
				for i, _ := range fieldGroupLi {
					o.setDataId(db, dataSource, &fieldGroupLi[i], bo, destData)
				}
			}
		}
		diffDataRowLi = append(diffDataRowLi, DiffDataRow{
			FieldGroupLi: fieldGroupLi,
			DestBo:       bo,
			DestData:     destData,
			SrcData:      srcData,
			SrcBo:        srcBo,
		})
	})

	usedCheck := UsedCheck{}
	// 删除的分录行数据的被用判断
	for _, diffDataRow := range diffDataRowLi {
		if usedCheck.CheckDeleteDetailRecordUsed(sessionId, dataSource, *bo, diffDataRow) {
			panic(BusinessError{Message: "部分分录数据已被用，不可删除"})
		}
	}

	// 被用差异行处理
	for i, _ := range diffDataRowLi {
		fieldGroupLi := diffDataRowLi[i].FieldGroupLi
		destData := diffDataRowLi[i].DestData
		srcData := diffDataRowLi[i].SrcData
		usedCheck.Update(sessionId, fieldGroupLi, bo, destData, srcData)
	}
	o.deleteExtraField(dataSource, bo)
	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
	//	txnManager.Update(txnId int, collection string, doc map[string]interface{}) (map[string]interface{}, bool) {
	if _, updateResult := txnManager.Update(txnId, collectionName, *bo); !updateResult {
		panic("更新失败")
	}
	return &diffDataRowLi
}

func (o FinanceService) deleteExtraField(dataSource DataSource, bo *map[string]interface{}) {
	modelIterator := ModelIterator{}
	var result interface{} = ""
	var notInKeyDict = map[string]interface{}{}
	modelIterator.IterateDataBo(dataSource, bo, &result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		if notInKeyDict[fieldGroupLi[0].GetDataSetId()] == nil {
			notInKeyLi := []string{}
			for key, _ := range *data {
				isIn := false
				for _, fieldGroup := range fieldGroupLi {
					if fieldGroup.Id == key {
						isIn = true
						break
					}
				}
				if !isIn && key != "_id" {
					notInKeyLi = append(notInKeyLi, key)
				}
			}
			notInKeyDict[fieldGroupLi[0].GetDataSetId()] = notInKeyLi
		}
		notInKeyLi := notInKeyDict[fieldGroupLi[0].GetDataSetId()].([]string)
		for _, key := range notInKeyLi {
			delete((*data), key)
		}
	})
}

func (o FinanceService) setDataId(db *mgo.Database, dataSource DataSource, fieldGroup *FieldGroup, bo *map[string]interface{}, data *map[string]interface{}) {
	if fieldGroup.Id == "id" {
		if fieldGroup.IsMasterField() {
			masterSeqName := GetMasterSequenceName((dataSource))
			masterSeqId := mongo.GetSequenceNo(db, masterSeqName)
			//			(*data)["_id"] = masterSeqId
			(*data)["id"] = masterSeqId
			(*bo)["_id"] = masterSeqId
			(*bo)["id"] = masterSeqId
		} else {
			detailData, found := fieldGroup.GetDetailData()
			if found {
				detailSeqName := GetDetailSequenceName((dataSource), detailData)
				detailSeqId := mongo.GetSequenceNo(db, detailSeqName)
				//				(*data)["_id"] = detailSeqId
				(*data)["id"] = detailSeqId
			}
		}
	}
}

func (o FinanceService) validateBO(sessionId int, dataSource DataSource, bo map[string]interface{}) string {
	messageLi := []string{}
	modelIterator := ModelIterator{}
	//	detailIndex := map[string]int{}
	//	for _, item := range dataSource.DetailDataLi {
	//		detailIndex[item.Id] = 0
	//	}
	var result interface{} = messageLi
	modelIterator.IterateAllFieldBo(dataSource, &bo, &result, func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		fieldMessageLi := o.validateFieldGroup(fieldGroup, *data)
		if fieldGroup.IsMasterField() {
			for _, item := range fieldMessageLi {
				messageLi = append(messageLi, item)
			}
		} else {
			detailData, _ := fieldGroup.GetDetailData()
			//			detailIndex[detailData.Id]++
			for _, item := range fieldMessageLi {
				//				messageLi = append(messageLi, "分录:"+detailData.DisplayName+"序号为"+strconv.Itoa(detailIndex[detailData.Id])+"的数据,"+item)
				messageLi = append(messageLi, "分录:"+detailData.DisplayName+"序号为"+strconv.Itoa(rowIndex+1)+"的数据,"+item)
			}
		}
	})
	// validate detailData.allowEmpty
	if dataSource.DetailDataLi != nil {
		for _, item := range dataSource.DetailDataLi {
			if item.AllowEmpty == "false" {
				isEmpty := false
				if bo[item.Id] == nil {
					isEmpty = true
				} else {
					lineData := bo[item.Id].([]interface{})
					if len(lineData) == 0 {
						isEmpty = true
					}
				}
				if isEmpty {
					messageLi = append(messageLi, "分录:"+item.DisplayName+"不允许为空")
				}
			}
		}
	}

	duplicateMessage := o.validateBODuplicate(sessionId, dataSource, bo)
	if duplicateMessage != "" {
		messageLi = append(messageLi, duplicateMessage)
	}

	return strings.Join(messageLi, "<br />")
}

func (o FinanceService) validateBODuplicate(sessionId int, dataSource DataSource, bo map[string]interface{}) string {
	result := ""
	result += o.validateMasterDataDuplicate(sessionId, dataSource, bo)
	if result != "" {
		result += "<br />"
	}
	result += o.validateDetailDataDuplicate(sessionId, dataSource, bo)
	return result
}

func (o FinanceService) validateMasterDataDuplicate(sessionId int, dataSource DataSource, bo map[string]interface{}) string {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	session, _ := global.GetConnection(sessionId)

	message := ""
	modelTemplateFactory := ModelTemplateFactory{}
	strId := modelTemplateFactory.GetStrId(bo)
	andQueryLi := []map[string]interface{}{}
	qb := QuerySupport{}
	andQueryLi = append(andQueryLi, map[string]interface{}{
		"deleteFlag": map[string]interface{}{
			"$ne": 9,
		},
		"A.createUnit": qb.GetCreateUnitByUserId(session, userId),
	})
	andFieldNameLi := []string{}
	modelIterator := ModelIterator{}
	var result interface{} = ""
	modelIterator.IterateAllFieldBo(dataSource, &bo, &result, func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		if fieldGroup.IsMasterField() {
			if fieldGroup.AllowDuplicate == "false" && fieldGroup.Id != "id" {
				andQueryLi = append(andQueryLi, map[string]interface{}{
					"A." + fieldGroup.Id: (*data)[fieldGroup.Id],
				})
				andFieldNameLi = append(andFieldNameLi, fieldGroup.DisplayName)
			}
		}
	})
	if len(andFieldNameLi) > 0 {
		if !(strId == "" || strId == "0") {
			andQueryLi = append(andQueryLi, map[string]interface{}{
				"_id": map[string]interface{}{
					"$ne": bo["id"],
				},
			})
		}
		duplicateQuery := map[string]interface{}{
			"$and": andQueryLi,
		}
		collectionName := modelTemplateFactory.GetCollectionName(dataSource)
		_, db := global.GetConnection(sessionId)
		duplicateQueryByte, err := json.MarshalIndent(duplicateQuery, "", "\t")
		if err != nil {
			panic(err)
		}
		log.Println("validateMasterDataDuplicate,collectionName:" + collectionName + ", query:" + string(duplicateQueryByte))
		count, err := db.C(collectionName).Find(duplicateQuery).Limit(1).Count()
		if err != nil {
			panic(err)
		}
		if count > 0 {
			message = strings.Join(andFieldNameLi, "+") + "不允许重复"
		}
	}

	return message
}

func (o FinanceService) validateDetailDataDuplicate(sessionId int, dataSource DataSource, bo map[string]interface{}) string {
	messageLi := []string{}
	modelIterator := ModelIterator{}
	var result interface{} = ""
	duplicateFieldIdLi := []string{}
	duplicateFieldNameLi := []string{}
	modelIterator.IterateAllField(&dataSource, &result, func(fieldGroup *FieldGroup, result *interface{}) {
		if !fieldGroup.IsMasterField() {
			if fieldGroup.AllowDuplicate == "false" && fieldGroup.Id != "id" {
				duplicateFieldIdLi = append(duplicateFieldIdLi, fieldGroup.Id)
				duplicateFieldNameLi = append(duplicateFieldNameLi, fieldGroup.DisplayName)
			}
		}
	})
	duplicateFieldNameJoin := strings.Join(duplicateFieldNameLi, "+") + "不允许重复"
	modelIterator.IterateDataBo(dataSource, &bo, &result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		modelIterator.IterateDataBo(dataSource, &bo, result, func(innerFieldGroupLi []FieldGroup, innerData *map[string]interface{}, innerRowIndex int, innerResult *interface{}) {
			if innerRowIndex > rowIndex {
				isDuplicate := true
				if len(duplicateFieldIdLi) == 0 {
					isDuplicate = false
				}
				for _, item := range duplicateFieldIdLi {
					if (*innerData)[item] != (*data)[item] {
						isDuplicate = false
						break
					}
				}
				if isDuplicate {
					detailData, _ := fieldGroupLi[0].GetDetailData()
					messageLi = append(messageLi, "分录:"+detailData.DisplayName+"序号为"+strconv.Itoa(rowIndex+1)+","+strconv.Itoa(innerRowIndex+1)+"的数据，"+duplicateFieldNameJoin)
				}
			}
		})
	})
	return strings.Join(messageLi, "<br />")
}

func (o FinanceService) validateFieldGroup(fieldGroup FieldGroup, data map[string]interface{}) []string {
	messageLi := []string{}

	if fieldGroup.AllowEmpty != "true" {
		value := data[fieldGroup.Id]
		if value != nil {
			strValue := fmt.Sprint(value)
			if strValue == "" {
				messageLi = append(messageLi, fieldGroup.Id+","+fieldGroup.DisplayName+"不允许空值")
				return messageLi
			}
		} else {
			messageLi = append(messageLi, fieldGroup.Id+","+fieldGroup.DisplayName+"不允许空值")
			return messageLi
		}
	}
	fieldValue := fmt.Sprint(data[fieldGroup.Id])
	if fieldGroup.ValidateExpr != "" {
		// python and golang validate, TODO miss golang validate
		expressionParser := ExpressionParser{}
		if !expressionParser.Validate(fieldValue, fieldGroup.ValidateExpr) {
			messageLi = append(messageLi, fieldGroup.DisplayName+fieldGroup.ValidateMessage)
			return messageLi
		}
	}
	isDataTypeNumber := false
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "DECIMAL"
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "FLOAT"
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "INT"
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "LONGINT"
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "MONEY"
	isDataTypeNumber = isDataTypeNumber || fieldGroup.FieldDataType == "SMALLINT"
	isUnLimit := fieldGroup.LimitOption != "" && fieldGroup.LimitOption != "unLimit"
	if isDataTypeNumber && isUnLimit {
		fieldValueFloat, err := strconv.ParseFloat(fieldValue, 64)
		if err != nil {
			panic(err)
		}
		if fieldGroup.LimitOption == "limitMax" {
			maxValue, err := strconv.ParseFloat(fieldGroup.LimitMax, 64)
			if err != nil {
				panic(err)
			}
			if maxValue < fieldValueFloat {
				messageLi = append(messageLi, fieldGroup.DisplayName+"超出最大值"+fieldGroup.LimitMax)
			}
		} else if fieldGroup.LimitOption == "limitMin" {
			minValue, err := strconv.ParseFloat(fieldGroup.LimitMin, 64)
			if err != nil {
				panic(err)
			}
			if fieldValueFloat < minValue {
				messageLi = append(messageLi, fieldGroup.DisplayName+"小于最小值"+fieldGroup.LimitMin)
			}
		} else if fieldGroup.LimitOption == "limitRange" {
			minValue, err := strconv.ParseFloat(fieldGroup.LimitMin, 64)
			if err != nil {
				panic(err)
			}
			maxValue, err := strconv.ParseFloat(fieldGroup.LimitMax, 64)
			if err != nil {
				panic(err)
			}
			if fieldValueFloat < minValue || maxValue < fieldValueFloat {
				messageLi = append(messageLi, fieldGroup.DisplayName+"超出范围("+fieldGroup.LimitMin+"~"+fieldGroup.LimitMax+")")
			}
		}
	} else {
		isDataTypeString := false
		isDataTypeString = isDataTypeString || fieldGroup.FieldDataType == "STRING"
		isDataTypeString = isDataTypeString || fieldGroup.FieldDataType == "REMARK"
		isFieldLengthLimit := fieldGroup.FieldLength != ""
		if isDataTypeString && isFieldLengthLimit {
			limit, err := strconv.Atoi(fmt.Sprint(fieldGroup.FieldLength))
			if err != nil {
				panic(err)
			}
			fieldValueLength := utf8.RuneCountInString(fieldValue)
			if fieldValueLength > limit {
				messageLi = append(messageLi, fieldGroup.DisplayName+"长度超出最大值"+fieldGroup.FieldLength)
			}
		}
	}

	return messageLi
}
