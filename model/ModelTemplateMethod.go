package model

import (
//	"reflect"
	. "github.com/hongjinqiu/gometa/script"
	"encoding/json"
)

func (o FieldGroup) IsMasterField() bool {
	return o.DataSetId == "A"
//	if o.Parent != nil {
//		if reflect.TypeOf(o.Parent).Name() == "FixField" {
//			fixField := o.Parent.(FixField)
//			return reflect.TypeOf(fixField.Parent).Name() == "MasterData"
//		}
//		if reflect.TypeOf(o.Parent).Name() == "BizField" {
//			bizField := o.Parent.(BizField)
//			return reflect.TypeOf(bizField.Parent).Name() == "MasterData"
//		}
//	}
//
//	return false
}

func (o FieldGroup) IsRelationField() bool {
	return len(o.RelationDS.RelationItemLi) > 0
}

func (o FieldGroup) GetRelationItem(bo map[string]interface{}, data map[string]interface{}) (RelationItem, bool) {
	expressionParser := ExpressionParser{}
	tmpBo := bo
	tmpBo["pendingTransactions"] = []interface{}{}
	boJsonData, err := json.Marshal(&tmpBo)
	if err != nil {
		panic(err)
	}
	boJson := string(boJsonData)
	for _, item := range o.RelationDS.RelationItemLi {
		if item.RelationExpr.Mode == "python" {
			dataJsonData, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}
			dataJson := string(dataJsonData)
			content := expressionParser.ParseModel(boJson, dataJson, item.RelationExpr.Content)
			if content == "true" {
				return item, true
			}
		} else if item.RelationExpr.Mode == "" || item.RelationExpr.Mode == "text" {
			if item.RelationExpr.Content == "true" {
				return item, true
			}
		} else if item.RelationExpr.Mode == "golang" {
			content := expressionParser.ParseGolang(bo, data, item.RelationExpr.Content)
			if content == "true" {
				return item, true
			}
		}
	}
	return RelationItem{}, false
}

func (o FieldGroup) GetMasterData() (MasterData, bool) {
	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(o.DataSourceId)
	if o.IsMasterField() {
		return dataSource.MasterData, true
//		if reflect.TypeOf(o.Parent).Name() == "FixField" {
//			fixField := o.Parent.(FixField)
//			return fixField.Parent.(MasterData), true
//		}
//		if reflect.TypeOf(o.Parent).Name() == "BizField" {
//			bizField := o.Parent.(BizField)
//			return bizField.Parent.(MasterData), true
//		}
	}
	return MasterData{}, false
}

func (o FieldGroup) GetDetailData() (DetailData, bool) {
	if o.IsMasterField() {
		return DetailData{}, false
	}
	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(o.DataSourceId)
	for _, detailData := range dataSource.DetailDataLi {
		if detailData.Id == o.DataSetId {
			return detailData, true
		}
	}
//	if reflect.TypeOf(o.Parent).Name() == "FixField" {
//		fixField := o.Parent.(FixField)
//		return fixField.Parent.(DetailData), true
//	}
//	if reflect.TypeOf(o.Parent).Name() == "BizField" {
//		bizField := o.Parent.(BizField)
//		return bizField.Parent.(DetailData), true
//	}
	return DetailData{}, false
}

func (o FieldGroup) GetDataSource() DataSource {
	modelTemplateFactory := ModelTemplateFactory{}
	return modelTemplateFactory.GetDataSource(o.DataSourceId)
//	if o.IsMasterField() {
//		masterData, _ := o.GetMasterData()
//		return masterData.Parent.(DataSource)
//	}
//	detailData, _ := o.GetDetailData()
//	return detailData.Parent.(DataSource)
}

func (o FieldGroup) GetDataSetId() string {
	return o.DataSetId
//	if o.IsMasterField() {
//		masterData, _ := o.GetMasterData()
//		return masterData.Id
//	}
//	detailData, _ := o.GetDetailData()
//	return detailData.Id
}
