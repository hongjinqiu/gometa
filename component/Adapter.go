package component

import (
	. "github.com/hongjinqiu/gometa/model"
	"reflect"
	"strings"
	"sync"
	"log"
)

var rwlock sync.RWMutex = sync.RWMutex{}
var adapterDict map[string]reflect.Type = map[string]reflect.Type{}

func init() {
	rwlock.Lock()
	defer rwlock.Unlock()
	adapterDict[reflect.TypeOf(ModelListTemplateAdapter{}).Name()] = reflect.TypeOf(ModelListTemplateAdapter{})
	adapterDict[reflect.TypeOf(ModelFormTemplateAdapter{}).Name()] = reflect.TypeOf(ModelFormTemplateAdapter{})
}

func GetAdapterDict() map[string]reflect.Type {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return adapterDict
}

type CommonMethod struct{}

func (o CommonMethod) Parse(classMethod string, param []interface{}) []reflect.Value {
	exprContent := classMethod
	scriptStruct := strings.Split(exprContent, ".")[0]
	scriptStructMethod := strings.Split(exprContent, ".")[1]
	scriptType := GetAdapterDict()[scriptStruct]
	if scriptType == nil {
		panic("adatpter " + scriptStruct + " not found")
	}
	inst := reflect.New(scriptType).Elem().Interface()
	instValue := reflect.ValueOf(inst)
	in := []reflect.Value{}
	for _, item := range param {
		in = append(in, reflect.ValueOf(item))
	}
	return instValue.MethodByName(scriptStructMethod).Call(in)
}

func (o CommonMethod) recursionApplyColumnModel(dataSource DataSource, columnModel *ColumnModel, result *interface{}) {
	modelIterator := ModelIterator{}
	for i, _ := range columnModel.ColumnLi {
		column := &columnModel.ColumnLi[i]
		if column.ColumnModel.ColumnLi == nil {
			if column.XMLName.Local == "auto-column" || column.Auto == "true" {
				if column.DsFieldMap == "" {
					modelIterator.IterateAllField(&dataSource, result, func(fieldGroup *FieldGroup, result *interface{}) {
						isApplyColumn := false
						columnModelDataSetId := columnModel.DataSetId
						if columnModelDataSetId == "" {
							columnModelDataSetId = "A"
						}
						isApplyColumn = isApplyColumn || (fieldGroup.GetDataSetId() == columnModelDataSetId && column.Name == fieldGroup.Id)
						if isApplyColumn {
							o.applyColumnExtend(*fieldGroup, column)
						}
					})
				} else {
					textLi := strings.Split(column.DsFieldMap, ".")
					if len(textLi) != 3 {
						log.Println("dataSet:" + columnModel.DataSetId + ", column.Name:" + column.Name + ", dsFieldMap:" + column.DsFieldMap + " apply failed, dsFieldMap.len != 3")
					} else {
						dataSourceId := textLi[0]
						dataSetId := textLi[1]
						fieldId := textLi[2]
						modelTemplateFactory := ModelTemplateFactory{}
						outSideDataSource := modelTemplateFactory.GetDataSource(dataSourceId)
						var outSideResult interface{} = ""
						modelIterator.IterateAllField(&outSideDataSource, &outSideResult, func(fieldGroup *FieldGroup, result *interface{}) {
							if fieldGroup.GetDataSetId() == dataSetId && fieldGroup.Id == fieldId {
								o.applyColumnExtend(*fieldGroup, column)
							}
						})
					}
				}
			}
		} else {
			o.recursionApplyColumnModel(dataSource, &column.ColumnModel, result)
		}
	}
}

func (o CommonMethod) applyAutoColumnXMLName(fieldGroup FieldGroup, column *Column) {
	xmlName := o.getColumnXMLName(fieldGroup)
	if xmlName != "" {
		column.XMLName.Local = xmlName
	}
}

func (o CommonMethod) getColumnXMLName(fieldGroup FieldGroup) string {
	isIntField := false
	intArray := []string{"SMALLINT", "INT", "LONGINT"}
	for _, item := range intArray {
		if strings.ToLower(fieldGroup.FieldDataType) == strings.ToLower(item) {
			isIntField = true
			break
		}
	}
	isFloatField := false
	floatArray := []string{"FLOAT", "MONEY", "DECIMAL"}
	for _, item := range floatArray {
		if strings.ToLower(fieldGroup.FieldDataType) == strings.ToLower(item) {
			isFloatField = true
			break
		}
	}
	isDateType := false
	dateArray := []string{"YEAR", "YEARMONTH", "DATE", "TIME", "DATETIME"}
	for _, item := range dateArray {
		if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower(item) {
			isDateType = true
			break
		}
	}
	if fieldGroup.IsRelationField() {
		return "select-column"
	} else if strings.ToLower(fieldGroup.FieldDataType) == "string" || strings.ToLower(fieldGroup.FieldDataType) == "remark" {
		return "string-column"
	} else if (isIntField || isFloatField) && !isDateType && fieldGroup.Dictionary == "" {
		return "number-column"
	} else if (isIntField || isFloatField) && isDateType && fieldGroup.Dictionary == "" {
		return "date-column"
	} else if (isIntField || isFloatField) && !isDateType && fieldGroup.Dictionary != "" {
		return "dictionary-column"
	}
	return ""
}

func (o CommonMethod) getCRelationDSFromModel(relationItem RelationItem) CRelationItem {
	cRelationItem := CRelationItem{}
	cRelationItem.Name = relationItem.Name
	cRelationItem.CRelationExpr.Mode = relationItem.RelationExpr.Mode
	cRelationItem.CRelationExpr.Content = relationItem.RelationExpr.Content
	
	cRelationItem.CJsRelationExpr.Mode = relationItem.JsRelationExpr.Mode
	cRelationItem.CJsRelationExpr.Content = relationItem.JsRelationExpr.Content
	
	cRelationItem.CRelationConfig.SelectorName = relationItem.Id
	cRelationItem.CRelationConfig.DisplayField = relationItem.DisplayField
	cRelationItem.CRelationConfig.ValueField = relationItem.ValueField
	cRelationItem.CRelationConfig.SelectionMode = "single"
	
	return cRelationItem
}

func (o CommonMethod) getCRelationDS(cRelationDS CRelationDS, relationDS RelationDS) CRelationDS {
	resultCRelationDS := CRelationDS{}
	if cRelationDS.CRelationItemLi == nil && relationDS.RelationItemLi != nil {
		cRelationItemLi := []CRelationItem{}
		for _, item := range relationDS.RelationItemLi {
			cRelationItem := o.getCRelationDSFromModel(item)
			cRelationItemLi = append(cRelationItemLi, cRelationItem)
		}
		resultCRelationDS.CRelationItemLi = cRelationItemLi
	} else if cRelationDS.CRelationItemLi != nil && relationDS.RelationItemLi != nil {
//		templateManager := TemplateManager{}
		cRelationItemLi := []CRelationItem{}
		for _, item := range relationDS.RelationItemLi {
			cRelationItem := CRelationItem{}
			
			isInherit := false
			columnRelationItem := CRelationItem{}
//			if len(cRelationDS.CRelationItemLi) > i {
//				for _, subItem := range cRelationDS.CRelationItemLi {
//					if subItem.Name == item.Name {
//						isInherit = true
//						columnRelationItem = subItem
//						break
//					}
//				}
//			}
			for _, subItem := range cRelationDS.CRelationItemLi {
				if subItem.Name == item.Name {
					isInherit = true
					columnRelationItem = subItem
					break
				}
			}
			if isInherit {
				cRelationItem.Name = columnRelationItem.Name
				if columnRelationItem.CRelationExpr.Mode != "" {
					cRelationItem.CRelationExpr.Mode = columnRelationItem.CRelationExpr.Mode
				} else {
					cRelationItem.CRelationExpr.Mode = item.RelationExpr.Mode
				}
				if columnRelationItem.CRelationExpr.Content != "" {
					cRelationItem.CRelationExpr.Content = columnRelationItem.CRelationExpr.Content
				} else {
					cRelationItem.CRelationExpr.Content = item.RelationExpr.Content
				}
				if columnRelationItem.CJsRelationExpr.Mode != "" {
					cRelationItem.CJsRelationExpr.Mode = columnRelationItem.CJsRelationExpr.Mode
				} else {
					cRelationItem.CJsRelationExpr.Mode = item.JsRelationExpr.Mode
				}
				if columnRelationItem.CJsRelationExpr.Content != "" {
					cRelationItem.CJsRelationExpr.Content = columnRelationItem.CJsRelationExpr.Content
				} else {
					cRelationItem.CJsRelationExpr.Content = item.JsRelationExpr.Content
				}
				if columnRelationItem.CRelationConfig.SelectorName != "" {
					cRelationItem.CRelationConfig.SelectorName = columnRelationItem.CRelationConfig.SelectorName
				} else {
					cRelationItem.CRelationConfig.SelectorName = item.Id
				}
				if columnRelationItem.CRelationConfig.DisplayField != "" {
					cRelationItem.CRelationConfig.DisplayField = columnRelationItem.CRelationConfig.DisplayField
				} else {
					cRelationItem.CRelationConfig.DisplayField = item.DisplayField
				}
				if columnRelationItem.CRelationConfig.ValueField != "" {
					cRelationItem.CRelationConfig.ValueField = columnRelationItem.CRelationConfig.ValueField
				} else {
					cRelationItem.CRelationConfig.ValueField = item.ValueField
				}
				if columnRelationItem.CRelationConfig.SelectionMode != "" {
					cRelationItem.CRelationConfig.SelectionMode = columnRelationItem.CRelationConfig.SelectionMode
				} else {
					cRelationItem.CRelationConfig.SelectionMode = "single"
				}
				
				if columnRelationItem.CCopyConfigLi != nil {
					cRelationItem.CCopyConfigLi = columnRelationItem.CCopyConfigLi
				}
			} else {
				cRelationItem = o.getCRelationDSFromModel(item)
			}
			
			cRelationItemLi = append(cRelationItemLi, cRelationItem)
		}
		resultCRelationDS.CRelationItemLi = cRelationItemLi
	}
	return resultCRelationDS
}

func (o CommonMethod) applyColumnExtend(fieldGroup FieldGroup, column *Column) {
	if column.Text == "" {
		column.Text = fieldGroup.DisplayName
	}
	if column.Hideable == "" {
		column.Hideable = fieldGroup.FixHide
	}
	if column.FixReadOnly == "" {
		column.FixReadOnly = fieldGroup.FixReadOnly
	}
	if column.ZeroShowEmpty == "" {
		column.ZeroShowEmpty = fieldGroup.ZeroShowEmpty
	}
	if column.XMLName.Local == "auto-column" {
		o.applyAutoColumnXMLName(fieldGroup, column)
	}

	if column.XMLName.Local == "select-column" {
		column.CRelationDS = o.getCRelationDS(column.CRelationDS, fieldGroup.RelationDS)
	} else if column.XMLName.Local == "string-column" {
		// do nothing
	} else if column.XMLName.Local == "number-column" {
		if column.CurrencyField == "" {
			column.CurrencyField = fieldGroup.FormatExpr
		}
		if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("MONEY") {
			if column.IsMoney == "" {
				column.IsMoney = "true"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("PRICE") {
			if column.IsUnitPrice == "" {
				column.IsUnitPrice = "true"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("UNITCOST") {
			if column.IsCost == "" {
				column.IsCost = "true"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("PERCENT") {
			if column.IsPercent == "" {
				column.IsPercent = "true"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("QUANTITY") {
			if column.IsQuantity == "" {
				column.IsQuantity = "true"
			}
		}
	} else if column.XMLName.Local == "date-column" {
		if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATE") {
			if column.DisplayPattern == "" {
				column.DisplayPattern = "yyyy-MM-dd" //需要从业务中查找,是一个系统配置,TODO,
			}
			if column.DbPattern == "" {
				column.DbPattern = "yyyyMMdd"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("DATETIME") {
			if column.DisplayPattern == "" {
				column.DisplayPattern = "yyyy-MM-dd HH:mm:ss" //需要从业务中查找,是一个系统配置,TODO,
			}
			if column.DbPattern == "" {
				column.DbPattern = "yyyyMMddHHmmss"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEAR") {
			if column.DisplayPattern == "" {
				column.DisplayPattern = "yyyy"
			}
			if column.DbPattern == "" {
				column.DbPattern = "yyyy"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("YEARMONTH") {
			if column.DisplayPattern == "" {
				column.DisplayPattern = "yyyy-MM" //需要从业务中查找,是一个系统配置,TODO,
			}
			if column.DbPattern == "" {
				column.DbPattern = "yyyyMM"
			}
		} else if strings.ToLower(fieldGroup.FieldNumberType) == strings.ToLower("TIME") {
			if column.DisplayPattern == "" {
				column.DisplayPattern = "HH:mm:ss" //需要从业务中查找,是一个系统配置,TODO,
			}
			if column.DbPattern == "" {
				column.DbPattern = "HHmmss"
			}
		}
	} else if column.XMLName.Local == "dictionary-column" {
		if column.Dictionary == "" {
			column.Dictionary = fieldGroup.Dictionary
		}
	}
}
