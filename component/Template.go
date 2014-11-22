package component

import (
	"encoding/xml"
	"html/template"
)

type ListTemplate struct {
	XMLName             xml.Name            `xml:"list-template"`
	Id                  string              `xml:"id"`
	SelectorId          string              `xml:"selector-id"`
	DataSourceModelId   string              `xml:"data-source-model-id"`
	Adapter             Adapter             `xml:"adapter"`
	Description         string              `xml:"description"`
	Cookie              Cookie              `xml:"cookie"`
	Scripts             template.URL        `xml:"scripts"`
	ViewTemplate        ViewTemplate        `xml:"view-template"`
	Toolbar             Toolbar             `xml:"toolbar"`
	Security            Security            `xml:"security"`
	BeforeBuildQuery    string              `xml:"before-build-query"`
	AfterBuildQuery     string              `xml:"after-build-query"`
	AfterQueryData      string              `xml:"after-query-data"`
	DataProvider        DataProvider        `xml:"data-provider"`
	ColumnModel         ColumnModel         `xml:"column-model"`
	QueryParameterGroup QueryParameterGroup `xml:"query-parameters"`
}

type FormTemplate struct {
	XMLName           xml.Name     `xml:"form-template"`
	Id                string       `xml:"id"`
	DataSourceModelId string       `xml:"data-source-model-id"`
	Adapter           Adapter      `xml:"adapter"`
	Description       string       `xml:"description"`
	Scripts           template.URL `xml:"scripts"`
	ViewTemplate      ViewTemplate `xml:"view-template"`
	Security          Security     `xml:"security"`
	FormElemLi        []FormElem   `xml:",any"`
}

type FormElem struct {
	XMLName     xml.Name    `xml:""`
	InnerHTML   string      `xml:",innerxml"`
	Html        Html        `xml:"-"`
	Toolbar     Toolbar     `xml:"-"`
	ColumnModel ColumnModel `xml:"-"`
	RenderTag   string      `xml:"-"` // 页面上渲染时用,主要是golang html的渲染能力较弱,

	ColumnModelAttributeGroup
}

type Adapter struct {
	XMLName xml.Name `xml:"adapter"`
	Name    string   `xml:"name,attr,omitempty"`
}

type Cookie struct {
	XMLName xml.Name `xml:"cookie"`
	Name    string   `xml:"name,attr,omitempty"`
}

type ViewTemplate struct {
	XMLName         xml.Name     `xml:"view-template"`
	View            string       `xml:"view,attr,omitempty"`
	SelectorView    string       `xml:"selectorView,attr,omitempty"`
	SelectorScripts template.URL `xml:"selectorScripts,attr,omitempty"`
}

type Html struct {
	XMLName xml.Name      `xml:"html"`
	Value   template.HTML `xml:",chardata"`
	ColSpan string        `xml:"colSpan,attr,omitempty"`
}

type Toolbar struct {
	XMLName xml.Name `xml:"toolbar"`
	ToolbarCommon
}

type EditorToolbar struct {
	XMLName xml.Name `xml:"editor-toolbar"`
	ToolbarCommon
}

type ToolbarCommon struct {
	ButtonGroup ButtonGroup `xml:"button-group"`
	ButtonLi    []Button    `xml:",any"`

	Name           string `xml:"name,attr,omitempty"`
	Export         string `xml:"export,attr,omitempty"`
	Exporter       string `xml:"exporter,attr,omitempty"`
	ExportParam    string `xml:"exportParam,attr,omitempty"`
	FreezedHeader  string `xml:"freezedHeader,attr,omitempty"`
	ExportChart    string `xml:"exportChart,attr,omitempty"`
	ExcelChart     string `xml:"excelChart,attr,omitempty"`
	ExcelChartType string `xml:"excelChartType,attr,omitempty"`
	ExportTitle    string `xml:"exportTitle,attr,omitempty"`
	ExportSuffix   string `xml:"exportSuffix,attr,omitempty"`
}

type Security struct {
	XMLName xml.Name `xml:"security"`

	ByUnit                string `xml:"byUnit,attr,omitempty"`
	ByAdmin               string `xml:"byAdmin,attr,omitempty"`
	FunctionId            string `xml:"functionId,attr,omitempty"`
	Override              string `xml:"override,attr,omitempty"`
	DEFAULT_RESOURCE_CODE string `xml:"DEFAULT_RESOURCE_CODE,attr,omitempty"`
}

type DataProvider struct {
	XMLName xml.Name `xml:"data-provider"`

	Collection    string `xml:"collection"`
	FixBsonQuery  string `xml:"fix-bson-query"`
	Map           string `xml:"map"`
	Reduce        string `xml:"reduce"`
	Size          string `xml:"size,attr,omitempty"`
	BsonIntercept string `xml:"bsonIntercept,attr,omitempty"`
}

type ColumnModel struct {
	XMLName xml.Name `xml:"column-model"`

	CheckboxColumn CheckboxColumn `xml:"checkbox-column"`
	IdColumn       IdColumn       `xml:"id-column"`
	Toolbar        Toolbar        `xml:"toolbar"`
	EditorToolbar  EditorToolbar  `xml:"editor-toolbar"`
	ColumnLi       []Column       `xml:",any"`

	ColumnModelAttributeGroup
}

type CheckboxColumn struct {
	XMLName xml.Name `xml:"checkbox-column"`

	ColumnAttributeLi []ColumnAttribute `xml:"column-attribute"`
	Expression        string            `xml:"expression"`
	Hideable          string            `xml:"hideable,attr,omitempty"`
	Name              string            `xml:"name,attr,omitempty"`
}

type ColumnAttribute struct {
	XMLName xml.Name `xml:"column-attribute"`

	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type IdColumn struct {
	XMLName xml.Name `xml:"id-column"`

	ColumnAttributeGroup
}

type ColumnModelAttributeGroup struct {
	Name                  string `xml:"name,attr,omitempty"`
	AutoLoad              string `xml:"autoLoad,attr,omitempty"`
	SummaryLoad           string `xml:"summaryLoad,attr,omitempty"`
	SummaryStat           string `xml:"summaryStat,attr,omitempty"`
	Summation             string `xml:"summation,attr,omitempty"`
	GroupSummation        string `xml:"groupSummation,attr,omitempty"`
	GroupMerge            string `xml:"groupMerge,attr,omitempty"`
	ShowGroupFilter       string `xml:"showGroupFilter,attr,omitempty"`
	ShowAggregationFilter string `xml:"showAggregationFilter,attr,omitempty"`
	AutoRowHeight         string `xml:"autoRowHeight,attr,omitempty"`
	Nowrap                string `xml:"nowrap,attr,omitempty"`
	Rownumber             string `xml:"rownumber,attr,omitempty"`
	SelectionMode         string `xml:"selectionMode,attr,omitempty"`
	ShowClearBtn          string `xml:"showClearBtn,attr,omitempty"`
	SelectionSupport      string `xml:"selectionSupport,attr,omitempty"`
	GroupField            string `xml:"groupField,attr,omitempty"`
	BsonOrderBy           string `xml:"bsonOrderBy,attr,omitempty"`
	SaveUrl               string `xml:"saveUrl,attr,omitempty"`
	DeleteUrl             string `xml:"deleteUrl,attr,omitempty"`
	StoreIntercept        string `xml:"storeIntercept,attr,omitempty"`
	RecordIntercept       string `xml:"recordIntercept,attr,omitempty"`
	SelectionTemplate     string `xml:"selectionTemplate,attr,omitempty"`
	SelectionTitle        string `xml:"selectionTitle,attr,omitempty"`
	DisplayMode           string `xml:"displayMode,attr,omitempty"`
	DataSetId             string `xml:"dataSetId,attr,omitempty"`
	ColSpan               string `xml:"colSpan,attr,omitempty"`
	Text                  string `xml:"text,attr,omitempty"`
}

type ColumnAttributeGroup struct {
	Name             string `xml:"name,attr,omitempty"`
	Bson             string `xml:"bson,attr,omitempty"`
	Text             string `xml:"text,attr,omitempty"`
	Align            string `xml:"align,attr,omitempty"`
	Graggable        string `xml:"graggable,attr,omitempty"`
	Groupable        string `xml:"groupable,attr,omitempty"`
	Hideable         string `xml:"hideable,attr,omitempty"`
	Editable         string `xml:"editable,attr,omitempty"`
	MenuDisabled     string `xml:"menuDisabled,attr,omitempty"`
	Sortable         string `xml:"sortable,attr,omitempty"`
	Comparable       string `xml:"comparable,attr,omitempty"`
	Locked           string `xml:"locked,attr,omitempty"`
	Auto             string `xml:"auto,attr,omitempty"`
	Width            string `xml:"width,attr,omitempty"`
	FieldWidth       string `xml:"fieldWidth,attr,omitempty"`
	FieldHeight      string `xml:"fieldHeight,attr,omitempty"`
	FieldCls         string `xml:"fieldCls,attr,omitempty"`
	ExcelWidth       string `xml:"excelWidth,attr,omitempty"`
	Renderer         string `xml:"renderer,attr,omitempty"`
	RendererTemplate string `xml:"rendererTemplate,attr,omitempty"`
	SummaryText      string `xml:"summaryText,attr,omitempty"`
	SummaryType      string `xml:"summaryType,attr,omitempty"`
	Cycle            string `xml:"cycle,attr,omitempty"`
	Exported         string `xml:"exported,attr,omitempty"`
	ColSpan          string `xml:"colSpan,attr,omitempty"`
	ColumnWidth      string `xml:"columnWidth,attr,omitempty"`
	LabelWidth       string `xml:"labelWidth,attr,omitempty"`
	DsFieldMap       string `xml:"dsFieldMap,attr,omitempty"`
	FixReadOnly      string `xml:"fixReadOnly,attr,omitempty"`
	ReadOnly         string `xml:"readOnly,attr,omitempty"`
	ZeroShowEmpty    string `xml:"zeroShowEmpty,attr,omitempty"`
	ManualRender     string `xml:"manualRender,attr,omitempty"`
	DataSetId        string `xml:"-"`
}

type Column struct {
	XMLName           xml.Name          `xml:""` // 有可能是string-column,number-column,date-column,boolean-column,dictionary-column,virtual-column,script-column,select-column
	Html              template.HTML     `xml:",chardata"`
	Name              string            `xml:"name,attr,omitempty"`
	ColumnAttributeLi []ColumnAttribute `xml:"column-attribute"`
	Editor            Editor            `xml:"editor"`
	Listeners         Listeners         `xml:"listeners"`
	CRelationDS       CRelationDS       `xml:"relationDS"`
	ColumnAttributeGroup
	ColumnModel ColumnModel `xml:"column-model"`

	//	Format         string `xml:"format,attr,omitempty"`
	DisplayPattern string `xml:"displayPattern,attr,omitempty"`
	DbPattern      string `xml:"dbPattern,attr,omitempty"`
	BooleanColumnAttributeGroup

	Dictionary string `xml:"dictionary,attr,omitempty"`
	Complex    string `xml:"complex,attr,omitempty"`

	Buttons            Buttons `xml:"buttons"`
	Prefix             string  `xml:"prefix,attr,omitempty"`
	DecimalPlaces      string  `xml:"decimalPlaces,attr,omitempty"`
	DecimalSeparator   string  `xml:"decimalSeparator,attr,omitempty"`
	ThousandsSeparator string  `xml:"thousandsSeparator,attr,omitempty"`
	Suffix             string  `xml:"suffix,attr,omitempty"`

	// 财务格式化字段
	CurrencyField string `xml:"currencyField,attr,omitempty"` // 币别字段
	IsMoney       string `xml:"isMoney,attr,omitempty"`       // 是否金额
	IsUnitPrice   string `xml:"isUnitPrice,attr,omitempty"`   // 是否单价
	IsCost        string `xml:"isCost,attr,omitempty"`        // 是否成本
	IsPercent     string `xml:"isPercent,attr,omitempty"`     // 是否百分比
	IsQuantity    string `xml:"isQuantity,attr,omitempty"`    // 是否数量

	Script string `xml:"script,attr,omitempty"`
	UseIn  string `xml:"use-in,attr,omitempty"`

	// select-column的内容
	//	DisplayField  string `xml:"displayField,attr,omitempty"`
	//	ValueField    string `xml:"valueField,attr,omitempty"`
	//	SelectorName  string `xml:"selectorName,attr,omitempty"`
	//	SelectionMode string `xml:"selectionMode,attr,omitempty"`
}

type Editor struct {
	XMLName         xml.Name          `xml:"editor"`
	EditorAttribute []EditorAttribute `xml:"editor_attribute"`
	Name            string            `xml:"name,attr,omitempty"`
}

type CRelationDS struct {
	XMLName         xml.Name        `xml:"relationDS"`
	CRelationItemLi []CRelationItem `xml:"relationItem"`
}

type CRelationItem struct {
	XMLName         xml.Name        `xml:"relationItem"`
	Name            string          `xml:"name,attr,omitempty"`
	CRelationExpr   CRelationExpr   `xml:"relationExpr"`
	CJsRelationExpr CJsRelationExpr `xml:"jsRelationExpr"`
	CRelationConfig CRelationConfig `xml:"relationConfig"`
	CCopyConfigLi   []CCopyConfig   `xml:"copyConfig"`
}

type CRelationConfig struct {
	XMLName       xml.Name `xml:"relationConfig"`
	SelectorName  string   `xml:"selectorName,attr,omitempty"`
	DisplayField  string   `xml:"displayField,attr,omitempty"`
	ValueField    string   `xml:"valueField,attr,omitempty"`
	SelectionMode string   `xml:"selectionMode,attr,omitempty"`
}

type CCopyConfig struct {
	XMLName        xml.Name `xml:"copyConfig"`
	CopyColumnName string   `xml:"copyColumnName,attr,omitempty"`
	CopyValueField string   `xml:"copyValueField,attr,omitempty"`
}

type CRelationExpr struct {
	XMLName xml.Name `xml:"relationExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type CJsRelationExpr struct {
	XMLName xml.Name `xml:"jsRelationExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type Listeners struct {
	XMLName     xml.Name `xml:"listeners"`
	Change      string   `xml:"change,attr,omitempty"`
	Selection   string   `xml:"selection,attr,omitempty"`
	UnSelection string   `xml:"unSelection,attr,omitempty"`
}

type EditorAttribute struct {
	XMLName xml.Name `xml:"editor_attribute"`

	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type BooleanColumnAttributeGroup struct {
	TrueText      string `xml:"trueText,attr,omitempty"`
	FalseText     string `xml:"falseText,attr,omitempty"`
	UndefinedText string `xml:"undefinedText,attr,omitempty"`
}

type QueryParameterGroup struct {
	XMLName xml.Name `xml:"query-parameters"`

	FixedParameterLi []FixedParameter `xml:"fixed-parameter"`
	QueryParameterLi []QueryParameter `xml:"query-parameter"`

	FormColumns      string `xml:"formColumns,attr,omitempty"`
	EnableEnterParam string `xml:"enableEnterParam,attr,omitempty"`
	DataSetId        string `xml:"dataSetId,attr,omitempty"`
}

type FixedParameter struct {
	XMLName xml.Name `xml:"fixed-parameter"`

	QueryParamAttributeGroup
}

type QueryParameter struct {
	XMLName xml.Name `xml:"query-parameter"`

	ParameterAttributeLi []ParameterAttribute `xml:"parameter-attribute"`
	CRelationDS          CRelationDS          `xml:"relationDS"`
	CDefaultValueExpr    CDefaultValueExpr    `xml:"defaultValueExpr"`
	QueryParamAttributeGroup
	Dictionary map[string]interface{} `xml:"-"`
	Tree       map[string]interface{} `xml:"-"`
	DataSetId  string                 `xml:"-"`
}

type CDefaultValueExpr struct {
	XMLName xml.Name `xml:"defaultValueExpr"`
	Mode    string   `xml:"mode,attr"`
	Content string   `xml:",chardata"`
}

type ParameterAttribute struct {
	XMLName xml.Name `xml:"parameter-attribute"`

	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type QueryParamAttributeGroup struct {
	Name        string `xml:"name,attr,omitempty"`
	Text        string `xml:"text,attr,omitempty"`
	ColumnName  string `xml:"columnName,attr,omitempty"`
	EnterParam  string `xml:"enterParam,attr,omitempty"`
	Auto        string `xml:"auto,attr,omitempty"`
	Editor      string `xml:"editor,attr,omitempty"`
	Restriction string `xml:"restriction,attr,omitempty"`
	ColSpan     string `xml:"colSpan,attr,omitempty"`
	RowSpan     string `xml:"rowSpan,attr,omitempty"`
	Value       string `xml:"value,attr,omitempty"`
	OtherName   string `xml:"otherName,attr,omitempty"`
	Having      string `xml:"having,attr,omitempty"`
	Required    string `xml:"required,attr,omitempty"`
	UseIn       string `xml:"use-in,attr,omitempty"`
	ReadOnly    string `xml:"readOnly,attr,omitempty"`
}

type ButtonGroup struct {
	XMLName  xml.Name `xml:"button-group"`
	ButtonLi []Button `xml:",any"`
}

type Buttons struct {
	XMLName xml.Name `xml:"buttons"`

	ButtonLi []Button `xml:"button"`
}

type Button struct {
	XMLName         xml.Name        `xml:""` // 有可能是button,也有可能是split-button
	Expression      string          `xml:"expression"`
	ButtonAttribute ButtonAttribute `xml:"button-attribute"`
	CRelationDS     CRelationDS     `xml:"relationDS"`
	ButtonAttributeGroup
}

type ButtonAttribute struct {
	XMLName xml.Name `xml:"button-attribute"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:"value,attr,omitempty"`
}

type ButtonAttributeGroup struct {
	Xtype      string      `xml:"xtype,attr,omitempty"`
	Name       string      `xml:"name,attr,omitempty"`
	Text       string      `xml:"text,attr,omitempty"`
	IconCls    string      `xml:"iconCls,attr,omitempty"`
	IconAlign  string      `xml:"iconAlign,attr,omitempty"`
	Disabled   string      `xml:"disabled,attr,omitempty"`
	Hidden     string      `xml:"hidden,attr,omitempty"`
	ArrowAlign string      `xml:"arrowAlign,attr,omitempty"`
	Scale      string      `xml:"scale,attr,omitempty"`
	Rowspan    string      `xml:"rowspan,attr,omitempty"`
	Handler    template.JS `xml:"handler,attr,omitempty"`
	Mode       string      `xml:"mode,attr,omitempty"`
	UseIn      string      `xml:"use-in,attr,omitempty"`
}
