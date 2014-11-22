package interceptor

import (
	"encoding/json"
	. "github.com/hongjinqiu/gometa/model"
	"strconv"
//	"fmt"
)

type PubReferenceLogInterceptor struct{}

func (o PubReferenceLogInterceptor) BeforeBuildQuery(sessionId int, paramMap map[string]string) map[string]string {
	dataSourceModelId := paramMap["beReferenceDataSourceModelId"]
	dataSetId := "A"
	idName := "id"
	beReferenceId := paramMap["beReferenceId"]
	if dataSourceModelId == "" || beReferenceId == "" {
		panic("传入的被引用方数据源模型id或被引用方id为空")
	}
	id, err := strconv.Atoi(beReferenceId)
	if err != nil {
		panic(err)
	}
	beReferenceLi := []interface{}{
		dataSourceModelId,
		dataSetId,
		idName,
		id,
	}
	data, err := json.Marshal(&beReferenceLi)
	if err != nil {
		panic(err)
	}
	paramMap["beReference"] = string(data)
	return paramMap
}

func (o PubReferenceLogInterceptor) AfterBuildQuery(sessionId int, queryLi []map[string]interface{}) []map[string]interface{} {
	for i, item := range queryLi {
		for k, v := range item {
			if k == "beReference" {
				beReferenceQuery := v.(string)
				beReferenceLi := []interface{}{}
				err := json.Unmarshal([]byte(beReferenceQuery), &beReferenceLi)
				if err != nil {
					panic(err)
				}
				item[k] = beReferenceLi
				queryLi[i] = item
			}
		}
	}
	return queryLi
}

func (o PubReferenceLogInterceptor) AfterQueryData(sessionId int, dataSetId string, items []interface{}) []interface{}  {
	result := []interface{}{}
	modelTemplateFactory := ModelTemplateFactory{}
	for _, item := range items {
		itemDict := item.(map[string]interface{})
		referenceLi := itemDict["reference"].([]interface{})
		reference := referenceLi[0].([]interface{})
		dataSourceId := reference[0].(string)
		displayName := modelTemplateFactory.GetDataSource(dataSourceId).DisplayName
		result = append(result, map[string]interface{}{
			"_id": itemDict["_id"],
			"id": itemDict["id"],
			"A": map[string]interface{}{
				"_id": itemDict["_id"],
				"id": itemDict["id"],
				"referenceDataSourceModelId": dataSourceId,
				"referenceDataSourceModelDisplayName": displayName,
				"referenceId": reference[3],
			},
		})
	}
	
	return result
}

