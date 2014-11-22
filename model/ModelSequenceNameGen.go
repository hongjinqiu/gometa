package model

import (
	"strings"
)

func GetMasterSequenceName(dataSource DataSource) string {
	modelTemplateFactory := ModelTemplateFactory{}
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	byte0 := collectionName[0]
	return strings.ToLower(string(byte0)) + collectionName[1:] + "Id"
}

func GetDetailSequenceName(dataSource DataSource, detailData DetailData) string {
	return GetMasterSequenceName(dataSource)
//	modelTemplateFactory := ModelTemplateFactory{}
//	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
//	byte0 := collectionName[0]
//	return strings.ToLower(string(byte0)) + collectionName[1:] + detailData.Id + "Id"
}
