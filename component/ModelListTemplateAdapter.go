package component

import (
	. "github.com/hongjinqiu/gometa/model"
	"strings"
//	"encoding/json"
)

type ModelListTemplateAdapter struct{}

// TODO, bytest
func (o ModelListTemplateAdapter) ApplyAdapter(iListTemplate interface{}) ListTemplate {
	listTemplate := iListTemplate.(ListTemplate)
	if listTemplate.DataSourceModelId != "" {
		modelTemplateFactory := ModelTemplateFactory{}
		dataSource := modelTemplateFactory.GetDataSource(listTemplate.DataSourceModelId)
		o.applyDataProvider(dataSource, &listTemplate)
		o.applyColumnModel(dataSource, &listTemplate)
		o.applyQueryParameter(dataSource, &listTemplate)
	}
	return listTemplate
}

// TODO, bytest
//func (o ModelListTemplateAdapter) ApplyQueryParameter(iListTemplate *interface{}, iQueryParameter *interface{}) {
//	listTemplate := (*iListTemplate).(ListTemplate)
//	queryParameter := (*iQueryParameter).(QueryParameter)
//	if listTemplate.DataSourceModelId != "" {
//		if listTemplate.QueryParameterGroup.DataSetId != "" {
//			queryParameter.Name = listTemplate.QueryParameterGroup.DataSetId + "." + queryParameter.Name
//		}
//	}
//}

// TODO, bytest
//func (o ModelListTemplateAdapter) ApplyColumnName(iListTemplate *interface{}, iColumn *interface{}) {
//	listTemplate := (*iListTemplate).(ListTemplate)
//	column := (*iColumn).(Column)
//	if listTemplate.DataSourceModelId != "" {
//		if listTemplate.QueryParameterGroup.DataSetId != "" {
//			column.Name = listTemplate.QueryParameterGroup.DataSetId + "." + column.Name
//		}
//	}
//}

// TODO, bytest
func (o ModelListTemplateAdapter) applyDataProvider(dataSource DataSource, listTemplate *ListTemplate) {
	if listTemplate.DataProvider.Collection == "" {
		modelTemplateFactory := ModelTemplateFactory{}
		listTemplate.DataProvider.Collection = modelTemplateFactory.GetCollectionName(dataSource)
	}
}

func (o ModelListTemplateAdapter) applyColumnModel(dataSource DataSource, listTemplate *ListTemplate) {
	var result interface{} = ""
	commonMethod := CommonMethod{}
	commonMethod.recursionApplyColumnModel(dataSource, &listTemplate.ColumnModel, &result)
}

func (o ModelListTemplateAdapter) applyQueryParameter(dataSource DataSource, listTemplate *ListTemplate) {
	commonMethod := CommonMethod{}
	var result interface{} = ""
	modelIterator := ModelIterator{}
	for i, _ := range listTemplate.QueryParameterGroup.QueryParameterLi {
		queryParameter := &listTemplate.QueryParameterGroup.QueryParameterLi[i]
		queryParameterDataSetId := listTemplate.QueryParameterGroup.DataSetId
		if queryParameterDataSetId == "" {
			queryParameterDataSetId = "A"
		}
		if queryParameter.Auto == "true" {
			modelIterator.IterateAllField(&dataSource, &result, func(fieldGroup *FieldGroup, result *interface{}){
				name := queryParameter.Name
				if queryParameter.ColumnName != "" {
					name = queryParameter.ColumnName
				}
				if fieldGroup.GetDataSetId() == queryParameterDataSetId && name == fieldGroup.Id {
					if queryParameter.Text == "" {
						queryParameter.Text = fieldGroup.DisplayName
					}
					if fieldGroup.FixHide == "true" {
						if queryParameter.Editor == "" {
							queryParameter.Editor = "hiddenfield"
						}
					}
					xmlName := commonMethod.getColumnXMLName(*fieldGroup)
					if xmlName != "" {
						o.applyQueryParameterAttr(xmlName, queryParameter)
						o.applyQueryParameterSubAttr(xmlName, *fieldGroup, queryParameter)
					}
				}
			})
		}
	}
}

func (o ModelListTemplateAdapter) applyQueryParameterAttr(xmlName string, queryParameter *QueryParameter) {
	if xmlName == "select-column" {
		if queryParameter.Editor == "" {
			queryParameter.Editor = "triggerfield"
		}
		if queryParameter.Restriction == "" {
			queryParameter.Restriction = "in"
		}
	} else if xmlName == "string-column" {
		if queryParameter.Editor == "" {
			queryParameter.Editor = "textfield"
		}
		if queryParameter.Restriction == "" {
			queryParameter.Restriction = "like"
		}
	} else if xmlName == "number-column" {
		if queryParameter.Editor == "" {
			queryParameter.Editor = "numberfield"
		}
		if queryParameter.Restriction == "" {
			queryParameter.Restriction = "eq"
		}
	} else if xmlName == "date-column" {
		if queryParameter.Editor == "" {
			queryParameter.Editor = "datefield"
		}
		if queryParameter.Restriction == "" {
			queryParameter.Restriction = "eq"
		}
	} else if xmlName == "dictionary-column" {
		if queryParameter.Editor == "" {
			queryParameter.Editor = "combofield"
		}
		if queryParameter.Restriction == "" {
			queryParameter.Restriction = "eq"
		}
	}
}

func (o ModelListTemplateAdapter) applyQueryParameterSubAttr(xmlName string, fieldGroup FieldGroup, queryParameter *QueryParameter) {
	if xmlName == "select-column" {
		commonMethod := CommonMethod{}
		queryParameter.CRelationDS = commonMethod.getCRelationDS(queryParameter.CRelationDS, fieldGroup.RelationDS)
	} else if xmlName == "string-column" {
		// do nothing
	} else if xmlName == "number-column" {
		// do nothing
	} else if xmlName == "date-column" {
		hasInFormat := false
		hasQueryFormat := false
		if queryParameter.ParameterAttributeLi != nil {
			for _, attrItem := range queryParameter.ParameterAttributeLi {
				if attrItem.Name == "displayPattern" {
					hasInFormat = true
					break
				}
			}
			for _, attrItem := range queryParameter.ParameterAttributeLi {
				if attrItem.Name == "dbPattern" {
					hasQueryFormat = true
					break
				}
			}
		}
		if !hasInFormat || !hasQueryFormat {
			if queryParameter.ParameterAttributeLi == nil {
				queryParameter.ParameterAttributeLi = []ParameterAttribute{}
			}
			if !hasInFormat {
				parameterAttribute := ParameterAttribute{}
				parameterAttribute.Name = "displayPattern"
				if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATE") {
					parameterAttribute.Value = "yyyy-MM-dd"//需要从业务中查找,是一个系统配置,TODO,
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATETIME") {
					parameterAttribute.Value = "yyyy-MM-dd HH:mm:ss"//需要从业务中查找,是一个系统配置,TODO,
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEAR") {
					parameterAttribute.Value = "yyyy"
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEARMONTH") {
					parameterAttribute.Value = "yyyy-MM"//需要从业务中查找,是一个系统配置,TODO,
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("TIME") {
					parameterAttribute.Value = "HH:mm:ss"//需要从业务中查找,是一个系统配置,TODO,
				}
				queryParameter.ParameterAttributeLi = append(queryParameter.ParameterAttributeLi, parameterAttribute)
			}
			if !hasQueryFormat {
				parameterAttribute := ParameterAttribute{}
				parameterAttribute.Name = "dbPattern"
				if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATE") {
					parameterAttribute.Value = "yyyyMMdd"
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATETIME") {
					parameterAttribute.Value = "yyyyMMddHHmmss"
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEAR") {
					parameterAttribute.Value = "yyyy"
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEARMONTH") {
					parameterAttribute.Value = "yyyyMM"
				} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("TIME") {
					parameterAttribute.Value = "HHmmss"
				}
				queryParameter.ParameterAttributeLi = append(queryParameter.ParameterAttributeLi, parameterAttribute)
			}
		}
	} else if xmlName == "dictionary-column" {
		hasDictionary := false
		if queryParameter.ParameterAttributeLi != nil {
			for _, attrItem := range queryParameter.ParameterAttributeLi {
				if attrItem.Name == "dictionary" {
					hasDictionary = true
					break
				}
			}
		}
		if !hasDictionary {
			if queryParameter.ParameterAttributeLi == nil {
				queryParameter.ParameterAttributeLi = []ParameterAttribute{}
			}
			parameterAttribute := ParameterAttribute{}
			parameterAttribute.Name = "dictionary"
			parameterAttribute.Value = fieldGroup.Dictionary
			queryParameter.ParameterAttributeLi = append(queryParameter.ParameterAttributeLi, parameterAttribute)
		}
	}
}

