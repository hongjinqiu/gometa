package component

import (
	
)

type ListTemplateIterator struct {}

func (o ListTemplateIterator) recursionGetColumnItem(columnModel ColumnModel, columnLi *[]Column) {
	for _, columnItem := range columnModel.ColumnLi {
		if columnItem.ColumnModel.ColumnLi != nil {
			o.recursionGetColumnItem(columnItem.ColumnModel, columnLi)
		}
		*columnLi = append(*columnLi, columnItem)
	}
}

type IterateTemplateColumnFunc func(column Column, result *interface{})

func (o ListTemplateIterator) IterateTemplateColumn(listTemplate ListTemplate, result *interface{}, iterateFunc IterateTemplateColumnFunc) {
	columnLi := []Column{}
	o.recursionGetColumnItem(listTemplate.ColumnModel, &columnLi)
	for _, item := range columnLi {
		iterateFunc(item, result)
	}
}

type IterateTemplateQueryParameterFunc func(queryParameter QueryParameter, result *interface{})

func (o ListTemplateIterator) IterateTemplateQueryParameter(listTemplate ListTemplate, result *interface{}, iterateFunc IterateTemplateQueryParameterFunc) {
	for _, item := range listTemplate.QueryParameterGroup.QueryParameterLi {
		iterateFunc(item, result)
	}
}
