package component

import (
	. "github.com/hongjinqiu/gometa/common"
	"fmt"
	"strconv"
	"strings"
//	"time"
)

type QueryParameterBuilder struct {
}

type RestrictionEditorFunc func(queryParameter QueryParameter, value string) map[string]interface{}

func (o QueryParameterBuilder) buildQuery(queryParameter QueryParameter, value string) map[string]interface{} {
	funcMap := map[string]map[string]RestrictionEditorFunc{
		"textfield":     o.stringCmpMap(),
		"textareafield": o.stringCmpMap(),
		"numberfield":   o.intOrFloatCmpMap(),
		"datefield":     o.dateCmpMap(),
		"combofield":    o.intOrFloatOrStringCmpMap(),
		"combotree":     o.intOrFloatOrStringCmpMap(),
		"displayfield":  nil,
		"hiddenfield":   o.intOrFloatOrStringOrLikeCmpMap(),
		"htmleditor":    o.stringCmpMap(),
		"checkboxfield": o.intOrFloatOrStringCmpMap(),
		"checkboxgroup": o.intOrFloatOrStringCmpMap(),
		"radiofield":    o.intOrFloatOrStringCmpMap(),
		"radiogroup":    o.intOrFloatOrStringCmpMap(),
		"triggerfield":  o.intOrFloatCmpMap(),
	}
	if funcMap[queryParameter.Editor] != nil {
		if funcMap[queryParameter.Editor][queryParameter.Restriction] != nil {
			var restrictionEditorFunc RestrictionEditorFunc
			restrictionEditorFunc = funcMap[queryParameter.Editor][queryParameter.Restriction]
			return restrictionEditorFunc(queryParameter, value)
		}
	}
	if funcMap[queryParameter.Editor] != nil && funcMap[queryParameter.Editor][queryParameter.Restriction] == nil {
		panic(queryParameter.Name + ",editor is:" + queryParameter.Editor + ",restriction is:" + queryParameter.Restriction + ", value is:" + value)
	}
	return map[string]interface{}{}
}

func (o QueryParameterBuilder) stringCmpMap() map[string]RestrictionEditorFunc {
	return map[string]RestrictionEditorFunc{
		"eq":             o.stringCmp("$eq"),
		"nq":             o.stringCmp("$ne"),
		"ge":             o.stringCmp("$gte"),
		"le":             o.stringCmp("$lte"),
		"gt":             o.stringCmp("$gt"),
		"lt":             o.stringCmp("$lt"),
		"null":           o.nullCmp("$eq"),
		"not_null":       o.nullCmp("$ne"),
		"exist":          o.existNotExistCmp(true),
		"not_exist":      o.existNotExistCmp(false),
		"in":             o.stringInCmp("$in"),
		"not_in":         o.stringInCmp("$nin"),
		"like":           o.regexpCmp("$regex"),
		"left_like":      o.regexpCmp("$regex"),
		"right_like":     o.regexpCmp("$regex"),
		"not_like":       o.regexpCmp("$regex"),
		"not_left_like":  o.regexpCmp("$regex"),
		"not_right_like": o.regexpCmp("$regex"),
	}
}

func (o QueryParameterBuilder) intOrFloatCmpMap() map[string]RestrictionEditorFunc {
	return map[string]RestrictionEditorFunc{
		"eq":             o.intOrFloatCmp("$eq"),
		"nq":             o.intOrFloatCmp("$ne"),
		"ge":             o.intOrFloatCmp("$gte"),
		"le":             o.intOrFloatCmp("$lte"),
		"gt":             o.intOrFloatCmp("$gt"),
		"lt":             o.intOrFloatCmp("$lt"),
		"null":           nil,
		"not_null":       nil,
		"exist":          o.existNotExistCmp(true),
		"not_exist":      o.existNotExistCmp(false),
		"in":             o.intOrFloatInCmp("$in"),
		"not_in":         o.intOrFloatInCmp("$nin"),
		"like":           nil,
		"left_like":      nil,
		"right_like":     nil,
		"not_like":       nil,
		"not_left_like":  nil,
		"not_right_like": nil,
	}
}

func (o QueryParameterBuilder) intOrFloatOrStringCmpMap() map[string]RestrictionEditorFunc {
	return map[string]RestrictionEditorFunc{
		"eq":             o.intOrFloatOrStringCmp("$eq"),
		"nq":             o.intOrFloatOrStringCmp("$ne"),
		"ge":             o.intOrFloatOrStringCmp("$gte"),
		"le":             o.intOrFloatOrStringCmp("$lte"),
		"gt":             o.intOrFloatOrStringCmp("$gt"),
		"lt":             o.intOrFloatOrStringCmp("$lt"),
		"null":           nil,
		"not_null":       nil,
		"exist":          o.existNotExistCmp(true),
		"not_exist":      o.existNotExistCmp(false),
		"in":             o.intOrFloatOrStringInCmp("$in"),
		"not_in":         o.intOrFloatOrStringInCmp("$nin"),
		"like":           nil,
		"left_like":      nil,
		"right_like":     nil,
		"not_like":       nil,
		"not_left_like":  nil,
		"not_right_like": nil,
	}
}

func (o QueryParameterBuilder) intOrFloatOrStringOrLikeCmpMap() map[string]RestrictionEditorFunc {
	return map[string]RestrictionEditorFunc{
		"eq":             o.intOrFloatOrStringCmp("$eq"),
		"nq":             o.intOrFloatOrStringCmp("$ne"),
		"ge":             o.intOrFloatOrStringCmp("$gte"),
		"le":             o.intOrFloatOrStringCmp("$lte"),
		"gt":             o.intOrFloatOrStringCmp("$gt"),
		"lt":             o.intOrFloatOrStringCmp("$lt"),
		"null":           o.nullCmp("$eq"),
		"not_null":       o.nullCmp("$ne"),
		"exist":          o.existNotExistCmp(true),
		"not_exist":      o.existNotExistCmp(false),
		"in":             o.intOrFloatOrStringInCmp("$in"),
		"not_in":         o.intOrFloatOrStringInCmp("$nin"),
		"like":           o.regexpCmp("$regex"),
		"left_like":      o.regexpCmp("$regex"),
		"right_like":     o.regexpCmp("$regex"),
		"not_like":       o.regexpCmp("$regex"),
		"not_left_like":  o.regexpCmp("$regex"),
		"not_right_like": o.regexpCmp("$regex"),
	}
}

func (o QueryParameterBuilder) dateCmpMap() map[string]RestrictionEditorFunc {
	return map[string]RestrictionEditorFunc{
		"eq":             o.dateCmp("$eq"),
		"nq":             o.dateCmp("$ne"),
		"ge":             o.dateCmp("$gte"),
		"le":             o.dateCmp("$lte"),
		"gt":             o.dateCmp("$gt"),
		"lt":             o.dateCmp("$lt"),
		"null":           nil,
		"not_null":       nil,
		"exist":          o.existNotExistCmp(true),
		"not_exist":      o.existNotExistCmp(false),
		"in":             nil,
		"not_in":         nil,
		"like":           nil,
		"left_like":      nil,
		"right_like":     nil,
		"not_like":       nil,
		"not_left_like":  nil,
		"not_right_like": nil,
	}
}

func (o QueryParameterBuilder) intOrFloatCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.intOrFloatOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) intOrFloatOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	resultLi := []map[string]interface{}{}
	intValue, err := strconv.ParseInt(value, 10, 0)
	commonUtil := CommonUtil{}
	if err == nil {
		if operator == "$eq" {
			//			resultLi = append(resultLi, map[string]interface{}{
			//				o.GetQueryName(queryParameter): intValue,
			//			})
			resultLi = append(resultLi, map[string]interface{}{
				"$or": []map[string]interface{}{
					map[string]interface{}{
						o.GetQueryName(queryParameter): intValue,
					},
					map[string]interface{}{
						o.GetQueryName(queryParameter): commonUtil.GetFloatFormat(fmt.Sprint(intValue)),
					},
				},
			})
		} else {
			//			resultLi = append(resultLi, map[string]interface{}{
			//				o.GetQueryName(queryParameter): map[string]int{
			//					operator: int(intValue),
			//				},
			//			})
			resultLi = append(resultLi, map[string]interface{}{
				"$or": []map[string]interface{}{
					map[string]interface{}{
						o.GetQueryName(queryParameter): map[string]interface{}{
							operator: int(intValue),
						},
					},
					map[string]interface{}{
						o.GetQueryName(queryParameter): map[string]interface{}{
							operator: commonUtil.GetFloatFormat(fmt.Sprint(intValue)),
						},
					},
				},
			})
		}
	}

	if len(resultLi) == 1 {
		return resultLi[0]
	}

	floatValue, err := strconv.ParseFloat(value, 32)
	if err == nil {
		if operator == "$eq" {
			//			resultLi = append(resultLi, map[string]interface{}{
			//				o.GetQueryName(queryParameter): float32(floatValue),
			//			})
			resultLi = append(resultLi, map[string]interface{}{
				"$or": []map[string]interface{}{
					map[string]interface{}{
						o.GetQueryName(queryParameter): float32(floatValue),
					},
					map[string]interface{}{
						o.GetQueryName(queryParameter): commonUtil.GetFloatFormat(value),
					},
				},
			})
		} else {
			//			resultLi = append(resultLi, map[string]interface{}{
			//				o.GetQueryName(queryParameter): map[string]float32{
			//					operator: float32(floatValue),
			//				},
			//			})
			resultLi = append(resultLi, map[string]interface{}{
				"$or": []map[string]interface{}{
					map[string]interface{}{
						o.GetQueryName(queryParameter): map[string]interface{}{
							operator: int(floatValue),
						},
					},
					map[string]interface{}{
						o.GetQueryName(queryParameter): map[string]interface{}{
							operator: commonUtil.GetFloatFormat(value),
						},
					},
				},
			})
		}
	}

	if len(resultLi) == 0 {
		return map[string]interface{}{}
	}

	if len(resultLi) == 1 {
		return resultLi[0]
	}
	return map[string]interface{}{
		"$or": resultLi,
	}
}

func (o QueryParameterBuilder) intOrFloatOrStringCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.intOrFloatOrStringOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) intOrFloatOrStringOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	resultLi := []map[string]interface{}{}

	resultLi = append(resultLi, o.stringOperator(queryParameter, value, operator))

	intOrFloatMap := o.intOrFloatOperator(queryParameter, value, operator)
	if intOrFloatMap["$or"] != nil {
		intOrFloatResultLi := intOrFloatMap["$or"]
		array := intOrFloatResultLi.([]map[string]interface{})
		for _, intOrFloatResult := range array {
			resultLi = append(resultLi, intOrFloatResult)
		}
	} else {
		if len(intOrFloatMap) > 0 {
			resultLi = append(resultLi, intOrFloatMap)
		}
	}

	if len(resultLi) == 0 {
		return map[string]interface{}{}
	}

	if len(resultLi) == 1 {
		return resultLi[0]
	}
	return map[string]interface{}{
		"$or": resultLi,
	}
}

func (o QueryParameterBuilder) stringCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.stringOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) nullCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.stringOperator(queryParameter, "", operator)
	}
}

func (o QueryParameterBuilder) stringOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	if operator == "$eq" {
		return map[string]interface{}{
			o.GetQueryName(queryParameter): value,
		}
	}
	return map[string]interface{}{
		o.GetQueryName(queryParameter): map[string]string{
			operator: value,
		},
	}
}

func (o QueryParameterBuilder) dateCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.dateOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) dateOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	/*
	displayPattern := "yyyy-MM-dd"
	dbPattern := "yyyyMMdd"
	for _, parameterAttribute := range queryParameter.ParameterAttributeLi {
		if parameterAttribute.Name == "displayPattern" {
			displayPattern = parameterAttribute.Value
		}
		if parameterAttribute.Name == "dbPattern" {
			dbPattern = parameterAttribute.Value
		}
	}
	displayPattern = o.replaceDateFormat(displayPattern)
	dbPattern = o.replaceDateFormat(dbPattern)

	t, err := time.Parse(displayPattern, value)
	if err != nil {
		panic(err)
	}

	queryDataStr := t.Format(dbPattern)
	queryData, err := strconv.ParseInt(queryDataStr, 10, 64)
	if err != nil {
		panic(err)
	}
	*/
	queryData, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}

	if operator == "$eq" {
		return map[string]interface{}{
			o.GetQueryName(queryParameter): queryData,
		}
	}
	return map[string]interface{}{
		o.GetQueryName(queryParameter): map[string]int64{
			operator: queryData,
		},
	}
}

func (o QueryParameterBuilder) replaceDateFormat(dateStr string) string {
	result := dateStr
	result = strings.Replace(result, "yyyy", "2006", -1)
	result = strings.Replace(result, "MM", "01", -1)
	result = strings.Replace(result, "dd", "02", -1)
	result = strings.Replace(result, "HH", "15", -1)
	result = strings.Replace(result, "mm", "04", -1)
	result = strings.Replace(result, "ss", "05", -1)
	return result
}

func (o QueryParameterBuilder) intOrFloatOrStringInCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.intOrFloatOrStringInOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) intOrFloatOrStringInOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	resultLi := []map[string]interface{}{}

	resultLi = append(resultLi, o.stringInOperator(queryParameter, value, operator))

	intOrFloatInMap := o.intOrFloatInOperator(queryParameter, value, operator)
	if intOrFloatInMap["$or"] != nil {
		intOrFloatInResultLi := intOrFloatInMap["$or"]
		array := intOrFloatInResultLi.([]map[string]interface{})
		for _, intOrFloatInResult := range array {
			resultLi = append(resultLi, intOrFloatInResult)
		}
	} else {
		if len(intOrFloatInMap) > 0 {
			resultLi = append(resultLi, intOrFloatInMap)
		}
	}

	if len(resultLi) == 0 {
		return map[string]interface{}{}
	}

	if len(resultLi) == 1 {
		return resultLi[0]
	}

	return map[string]interface{}{
		"$or": resultLi,
	}
}

func (o QueryParameterBuilder) intOrFloatInCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.intOrFloatInOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) intOrFloatInOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	valueLi := strings.Split(value, ",")
	resultLi := []map[string]interface{}{}

	intValueLi := []int{}
	for _, valueItem := range valueLi {
		intValue, err := strconv.ParseInt(valueItem, 10, 0)
		if err == nil {
			intValueLi = append(intValueLi, int(intValue))
		}
	}

	if len(intValueLi) > 0 {
		resultLi = append(resultLi, map[string]interface{}{
			o.GetQueryName(queryParameter): map[string]interface{}{
				operator: intValueLi,
			},
		})
	}

	floatValueLi := []float32{}
	for _, valueItem := range valueLi {
		_, err := strconv.ParseInt(valueItem, 10, 0)
		if err == nil {
			continue
		}

		float32Value, err := strconv.ParseFloat(valueItem, 32)
		if err == nil {
			floatValueLi = append(floatValueLi, float32(float32Value))
		}
	}

	if len(floatValueLi) > 0 {
		resultLi = append(resultLi, map[string]interface{}{
			o.GetQueryName(queryParameter): map[string]interface{}{
				operator: floatValueLi,
			},
		})
	}

	if len(resultLi) == 0 {
		return map[string]interface{}{}
	}

	if len(resultLi) == 1 {
		return resultLi[0]
	}
	return map[string]interface{}{
		"$or": resultLi,
	}
}

func (o QueryParameterBuilder) stringInCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.stringInOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) stringInOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	valueLi := strings.Split(value, ",")

	return map[string]interface{}{
		o.GetQueryName(queryParameter): map[string]interface{}{
			operator: valueLi,
		},
	}
}

func (o QueryParameterBuilder) existNotExistCmp(isExist bool) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.existNotExistOperator(queryParameter, isExist)
	}
}

func (o QueryParameterBuilder) existNotExistOperator(queryParameter QueryParameter, isExist bool) map[string]interface{} {
	return map[string]interface{}{
		o.GetQueryName(queryParameter): map[string]bool{
			"$exists": isExist,
		},
	}
}

func (o QueryParameterBuilder) regexpCmp(operator string) RestrictionEditorFunc {
	return func(queryParameter QueryParameter, value string) map[string]interface{} {
		return o.regexpOperator(queryParameter, value, operator)
	}
}

func (o QueryParameterBuilder) regexpOperator(queryParameter QueryParameter, value string, operator string) map[string]interface{} {
	regex := ""
	switch queryParameter.Restriction {
	case "like":
		regex = "(?i)^.*?" + value + ".*?$"
	case "not_like":
		regex = "(?i)^((?!" + value + ").)*$"
	case "left_like":
		regex = "(?i)^" + value + ".*?$"
	case "not_left_like":
		regex = "(?i)^((?!^" + value + ").)*$"
	case "right_like":
		regex = "(?i)^.*?" + value + "$"
	case "not_right_like":
		regex = "(?i)^((?!" + value + "$).)*$"
	}
	switch queryParameter.Restriction {
	case "like", "left_like", "right_like", "not_like", "not_left_like", "not_right_like":
		return map[string]interface{}{
			o.GetQueryName(queryParameter): map[string]string{
				operator: regex,
			},
		}
	}
	panic("input queryParameter.Editor:" + queryParameter.Editor + ", must be one of 'like', 'left_like', 'right_like', 'not_like', 'not_left_like', 'not_left_like'")
}

func (o QueryParameterBuilder) GetQueryName(queryParameter QueryParameter) string {
	if queryParameter.DataSetId == "" {
		if queryParameter.ColumnName != "" {
			return queryParameter.ColumnName
		}
		return queryParameter.Name
	}
	if queryParameter.ColumnName != "" {
		return queryParameter.DataSetId + "." + queryParameter.ColumnName
	}
	return queryParameter.DataSetId + "." + queryParameter.Name
}
