package interceptor

import (

)

type SysUserInterceptor struct{}

func (o SysUserInterceptor) BeforeBuildQuery(sessionId int, paramMap map[string]string) map[string]string {
	paramMap["nick"] = ""
	return paramMap
}

func (o SysUserInterceptor) AfterBuildQuery(sessionId int, queryLi []map[string]interface{}) []map[string]interface{} {
	queryLi = append(queryLi, map[string]interface{}{
		"age": 20,
	})
	return queryLi
}

func (o SysUserInterceptor) AfterQueryData(sessionId int, dataSetId string, items []interface{}) []interface{}  {
	for i, _ := range items {
		item := items[i].(map[string]interface{})
		item["UNIT_NAME"] = "单位名称aaa"
		item["numTest"] = 1000 * 1000 + 0.12345678901
		item["numTest1"] = 1000 * 1000 + 0.12345678901
		item["numTest2"] = 1000 * 1000 + 0.12345678901
		item["numTest3"] = 1000 * 1000 + 0.12345678901
		item["numTest4"] = 1000 * 1000 + 0.12345678901
		item["numTest5"] = 1000 * 1000 + 0.12345678901
		item["numTest6"] = 1000 * 1000 + 0.12345678901
		item["numTest7"] = 1000 * 1000 + 0.12345678901
		item["numTest8"] = 1000 * 1000 + 0.12345678901
	}
	for i, _ := range items {
		item := items[i].(map[string]interface{})
		if i == 0{
			item["currency"] = map[string]interface{} {
				"prefix": "$",
				"decimalPlaces": 3,
				"unitPriceDecimalPlaces": 6,
				"decimalSeparator": "^",
				"thousandsSeparator": "_",
				"suffix": "&",
			}
		} else {
			item["currency"] = map[string]interface{}{
				"prefix": "*",
				"decimalPlaces": 4,
				"unitPriceDecimalPlaces": 8,
				"decimalSeparator": "@",
				"thousandsSeparator": "=",
				"suffix": "!",
			}
		}
	}
	for i, _ := range items {
		item := items[i].(map[string]interface{})
		item["dateTest"] = 20131020
		item["dateTimeTest"] = int64(20131020114825)
	}
	for i, _ := range items {
		item := items[i].(map[string]interface{})
		item["boolTest"] = true
		item["boolTest2"] = false
	}
	for i, _ := range items {
		item := items[i].(map[string]interface{})
		item["nestDictTest1"] = 1
	}
	return items
}

