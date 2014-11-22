package interceptor

import (

)

type ModelListTemplateInterceptor struct{}

func (o ModelListTemplateInterceptor) AfterQueryData(sessionId int, dataSetId string, items []interface{}) []interface{}  {
	if dataSetId == "" {
		return items
	}
	result := []interface{}{}
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		resultMap := itemMap[dataSetId].(map[string]interface{})
		result = append(result, resultMap)
	}
	return result
}
