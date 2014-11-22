package model

import (
	"encoding/xml"
)

type DataSourceTest struct {
	XMLName xml.Name `xml:"datasource"`
	Id      string   `xml:"id"`
}

type DataSource struct {
	XMLName                 xml.Name     `xml:"datasource"`
	Id                      string       `xml:"id"`
	DisplayName             string       `xml:"displayName"`
	SystemId                string       `xml:"systemId"`
	CodeFieldName           string       `xml:"codeFieldName"`
	BusinessDateField       string       `xml:"businessDateField"`
	ModelType               string       `xml:"modelType"`
	InUsedDenyEdit          string       `xml:"inUsedDenyEdit"`
	ActionNameSpace         string       `xml:"actionNameSpace"`
	ListUrl                 string       `xml:"listUrl"`
	CollectionName          string       `xml:"collectionName"`
	BillTypeField           string       `xml:"billTypeField"`
	BillTypeParamDataSource string       `xml:"billTypeParamDataSource"`
	HasCheckField           string       `xml:"hasCheckField"`
	ListSortFields          string       `xml:"listSortFields"`
	MasterData              MasterData   `xml:"masterData"`
	DetailDataLi            []DetailData `xml:"detailData"`
}

type MasterData struct {
	XMLName     xml.Name    `xml:"masterData"`
	Id          string      `xml:"id"`
	DisplayName string      `xml:"displayName"`
	AllowCopy   string      `xml:"allowCopy"`
	PrimaryKey  string      `xml:"primaryKey"`
	FixField    FixField    `xml:"fixField"`
	BizField    BizField    `xml:"bizField"`
}

type DetailData struct {
	XMLName     xml.Name `xml:"detailData"`
	Id          string   `xml:"id"`
	DisplayName string   `xml:"displayName"`
	//	ParentId      string      `xml:"parentId"`
	AllowEmpty string `xml:"allowEmpty"`
	//AllowEmptyRow string      `xml:"allowEmptyRow"`
	AllowCopy string `xml:"allowCopy"`
	//Readonly      string      `xml:"readonly"`
	PrimaryKey string      `xml:"primaryKey"`
	FixField   FixField    `xml:"fixField"`
	BizField   BizField    `xml:"bizField"`
}

type FixField struct {
	XMLName     xml.Name    `xml:"fixField"`
	PrimaryKey  PrimaryKey  `xml:"primaryKey"`
	CreateBy    CreateBy    `xml:"createBy"`
	CreateTime  CreateTime  `xml:"createTime"`
	CreateUnit  CreateUnit  `xml:"createUnit"`
	ModifyUnit  ModifyUnit  `xml:"modifyUnit"`
	ModifyBy    ModifyBy    `xml:"modifyBy"`
	ModifyTime  ModifyTime  `xml:"modifyTime"`
	BillStatus  BillStatus  `xml:"billStatus"`
	AttachCount AttachCount `xml:"attachCount"`
	Remark      Remark      `xml:"remark"`
}

type BizField struct {
	XMLName xml.Name    `xml:"bizField"`
	FieldLi []Field     `xml:"field"`
}

type PrimaryKey struct {
	XMLName xml.Name `xml:"primaryKey"`
	FieldGroup
}
type CreateBy struct {
	XMLName xml.Name `xml:"createBy"`
	FieldGroup
}
type CreateTime struct {
	XMLName xml.Name `xml:"createTime"`
	FieldGroup
}
type CreateUnit struct {
	XMLName xml.Name `xml:"createUnit"`
	FieldGroup
}
type ModifyBy struct {
	XMLName xml.Name `xml:"modifyBy"`
	FieldGroup
}
type ModifyUnit struct {
	XMLName xml.Name `xml:"modifyUnit"`
	FieldGroup
}
type ModifyTime struct {
	XMLName xml.Name `xml:"modifyTime"`
	FieldGroup
}
type BillStatus struct { // 0:正常,1:作废
	XMLName xml.Name `xml:"billStatus"`
	FieldGroup
}
type AttachCount struct {
	XMLName xml.Name `xml:"attachCount"`
	FieldGroup
}
type Remark struct {
	XMLName xml.Name `xml:"remark"`
	FieldGroup
}

type Fields struct {
	XMLName xml.Name `xml:"fields"`
	FieldLi []Field  `xml:"field"`
}

type Field struct {
	XMLName xml.Name `xml:"field"`
	FieldGroup
}

type FieldGroup struct {
	Id      string `xml:"id,attr"`
	Extends string `xml:"extends,attr"`
	//	FieldName         string     `xml:"fieldName"`
	DisplayName       string           `xml:"displayName"`
	FieldDataType     string           `xml:"fieldDataType"`
	FieldNumberType   string           `xml:"fieldNumberType"`
	FieldLength       string           `xml:"fieldLength"`
	DefaultValueExpr  DefaultValueExpr `xml:"defaultValueExpr"`
	CheckInUsed       string           `xml:"checkInUsed"`
	FixHide           string           `xml:"fixHide"`
	FixReadOnly       string           `xml:"fixReadOnly"`
	AllowDuplicate    string           `xml:"allowDuplicate"`
	AllowCopy         string           `xml:"allowCopy"`
	DenyEditInUsed    string           `xml:"denyEditInUsed"`
	AllowEmpty        string           `xml:"allowEmpty"`
	LimitOption       string           `xml:"limitOption"`
	LimitMax          string           `xml:"limitMax"`
	LimitMin          string           `xml:"limitMin"`
	ValidateExpr      string           `xml:"validateExpr"`
	ValidateMessage   string           `xml:"validateMessage"`
	Dictionary        string           `xml:"dictionary"`
	DictionaryWhere   string           `xml:"dictionaryWhere"`
	CalcValueExpr     CalcValueExpr    `xml:"calcValueExpr"`
	Virtual           string           `xml:"virtual"`
	ZeroShowEmpty     string           `xml:"zeroShowEmpty"`
	LocalCurrencyency string           `xml:"localCurrencyency"`
	FieldInList       string           `xml:"fieldInList"`
	ListWhereField    string           `xml:"listWhereField"`
	FormatExpr        string           `xml:"formatExpr"`
	RelationDS        RelationDS       `xml:"relationDS"`
	DataSourceId      string           `xml:"-"`
	DataSetId         string           `xml:"-"`
}

type DefaultValueExpr struct {
	XMLName xml.Name `xml:"defaultValueExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type CalcValueExpr struct {
	XMLName xml.Name `xml:"calcValueExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type RelationDS struct {
	XMLName        xml.Name       `xml:"relationDS"`
	RelationItemLi []RelationItem `xml:"relationItem"`
}

type RelationItem struct {
	XMLName           xml.Name       `xml:"relationItem"`
	Name              string         `xml:"name,attr,omitempty"`
	Id                string         `xml:"id"`
	RelationExpr      RelationExpr   `xml:"relationExpr"`
	JsRelationExpr    JsRelationExpr `xml:"jsRelationExpr"`
	RelationModelId   string         `xml:"relationModelId"`
	RelationDataSetId string         `xml:"relationDataSetId"`
	DisplayField      string         `xml:"displayField"`
	ValueField        string         `xml:"valueField"`
}

type RelationExpr struct {
	XMLName xml.Name `xml:"relationExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type JsRelationExpr struct {
	XMLName xml.Name `xml:"jsRelationExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}
