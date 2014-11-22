package component

import (
	. "github.com/hongjinqiu/gometa/model"
	"encoding/xml"
)

type ModelFormTemplateAdapter struct{}

func (o ModelFormTemplateAdapter) ApplyAdapter(iFormTemplate interface{}) FormTemplate {
	formTemplate := (iFormTemplate).(FormTemplate)
	if formTemplate.DataSourceModelId != "" {
		modelTemplateFactory := ModelTemplateFactory{}
		dataSource := modelTemplateFactory.GetDataSource(formTemplate.DataSourceModelId)
		o.applyDetailDataSet(dataSource, &formTemplate)
	}
	return formTemplate
}

func (o ModelFormTemplateAdapter) applyDetailDataSet(dataSource DataSource, formTemplate *FormTemplate) {
	commonMethod := CommonMethod{}
	var result interface{} = ""
	for i, _ := range formTemplate.FormElemLi {
		if formTemplate.FormElemLi[i].XMLName.Local == "column-model" {
			// 应用title,
			for _, detailDataItem := range dataSource.DetailDataLi {
				if formTemplate.FormElemLi[i].DataSetId == detailDataItem.Id {
					formTemplate.FormElemLi[i].Text = detailDataItem.DisplayName
					formTemplate.FormElemLi[i].ColumnModel.Text = detailDataItem.DisplayName
				}
			}
		
			commonMethod.recursionApplyColumnModel(dataSource, &formTemplate.FormElemLi[i].ColumnModel, &result)
			data, err := xml.Marshal(&formTemplate.FormElemLi[i].ColumnModel)
			if err != nil {
				panic(err)
			}
			columnModel := formTemplate.FormElemLi[i].ColumnModel
			err = xml.Unmarshal(data, &formTemplate.FormElemLi[i])
			if err != nil {
				panic(err)
			}
			formTemplate.FormElemLi[i].ColumnModel = columnModel
		}
	}
}
