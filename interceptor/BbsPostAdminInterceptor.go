package interceptor

import (

)

type BbsPostAdminInterceptor struct{}

func (o BbsPostAdminInterceptor) AfterBuildQuery(sessionId int, queryLi []map[string]interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, item := range queryLi {
		isCreateUnit := false
		for k, _ := range item {
			if k == "A.createUnit" {
				isCreateUnit = true
				break
			}
		}
		if !isCreateUnit {
			result = append(result, item)
		}
	}
	return result
}
