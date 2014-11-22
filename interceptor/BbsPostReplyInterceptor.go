package interceptor

import (
	"time"
	"fmt"
)

type BbsPostReplyInterceptor struct{}

func (o BbsPostReplyInterceptor) AfterBuildQuery(sessionId int, queryLi []map[string]interface{}) []map[string]interface{} {
	orQuery := map[string]interface{}{}
	orLi := []map[string]interface{}{}
	createUnitQuery := map[string]interface{}{}
	for _, item := range queryLi {
		if item["A.createUnit"] == nil {
			orLi = append(orLi, item)
		} else {
			createUnitQuery = item
		}
	}
	orQuery["$or"] = orLi
	if len(orLi) > 0 {
		queryLi = []map[string]interface{}{
			orQuery,
		}
		if len(createUnitQuery) > 0 {
			queryLi = append(queryLi, createUnitQuery)
		}
	}
	return queryLi
}

func (o BbsPostReplyInterceptor) AfterQueryData(sessionId int, dataSetId string, items []interface{}) []interface{}  {
	for i, item := range items {
		data := item.(map[string]interface{})
		dataA := data["A"].(map[string]interface{})
		data["A"] = dataA
		items[i] = data
		
		createTimeStr := fmt.Sprint(dataA["createTime"])
		if createTimeStr != "" && createTimeStr != "0" {
			dDate, err := time.Parse("20060102150405", fmt.Sprint(createTimeStr))
			if err != nil {
				panic(err)
			}
			dataA["createTimeDisplay"] = dDate.Format("2006-01-02 15:04:05")
		} else {
			dataA["createTimeDisplay"] = ""
		}
	}
	return items
}
