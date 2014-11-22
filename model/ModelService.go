package model

import (
	"fmt"
	"strconv"
)

type ModelIterator struct {}

type IterateFunc func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{})

func (o ModelIterator) IterateAllFieldBo(dataSource DataSource, bo *map[string]interface{}, result *interface{}, iterateFunc IterateFunc) {
	o.IterateDataBo(dataSource, bo, result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}){
		for _, item := range fieldGroupLi {
			iterateFunc(item, data, rowIndex, result)
		}
	})
}

func (o ModelIterator) GetFixFieldLi(fixField *FixField) *[]*FieldGroup {
	fixFieldLi := []*FieldGroup{}
	fixFieldLi = append(fixFieldLi, &fixField.PrimaryKey.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.CreateBy.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.CreateTime.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.CreateUnit.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.ModifyBy.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.ModifyTime.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.ModifyUnit.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.BillStatus.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.AttachCount.FieldGroup)
	fixFieldLi = append(fixFieldLi, &fixField.Remark.FieldGroup)
	return &fixFieldLi
}

type IterateTwoBoFunc func(fieldGroup *FieldGroup, destData *map[string]interface{}, srcData map[string]interface{}, result *interface{})

func (o ModelIterator) IterateAllFieldTwoBo(dataSource *DataSource, destBo *map[string]interface{}, srcBo map[string]interface{}, result *interface{}, iterateFunc IterateTwoBoFunc) {
	destData := (*destBo)["A"].(map[string]interface{})
	srcData := srcBo["A"].(map[string]interface{})
	fieldGroupLi := o.getDataSetFieldGroupLi(&dataSource.MasterData.FixField, &dataSource.MasterData.BizField)
	for i, _ := range *fieldGroupLi {
		iterateFunc((*fieldGroupLi)[i], &destData, srcData, result)
	}
	
	for i, _ := range dataSource.DetailDataLi {
		fieldGroupLi = o.getDataSetFieldGroupLi(&dataSource.DetailDataLi[i].FixField, &dataSource.DetailDataLi[i].BizField)
		destDataLi := (*destBo)[dataSource.DetailDataLi[i].Id].([]interface{})
		srcDataLi := srcBo[dataSource.DetailDataLi[i].Id].([]interface{})
		for subIndex, _ := range destDataLi {
			destDetailData := destDataLi[subIndex].(map[string]interface{})
			srcDetailData := srcDataLi[subIndex].(map[string]interface{})
			for j, _ := range *fieldGroupLi {
				iterateFunc((*fieldGroupLi)[j], &destDetailData, srcDetailData, result)
			}
		}
	}
}

type IterateFieldFunc func(fieldGroup *FieldGroup, result *interface{})

func (o ModelIterator) IterateAllField(dataSource *DataSource, result *interface{}, iterateFunc IterateFieldFunc) {
	fieldGroupLi := o.getDataSetFieldGroupLi(&dataSource.MasterData.FixField, &dataSource.MasterData.BizField)
	for i, _ := range *fieldGroupLi {
		iterateFunc((*fieldGroupLi)[i], result)
	}

	for i, _ := range dataSource.DetailDataLi {
		fieldGroupLi = o.getDataSetFieldGroupLi(&dataSource.DetailDataLi[i].FixField, &dataSource.DetailDataLi[i].BizField)
		for j, _ := range *fieldGroupLi {
			iterateFunc((*fieldGroupLi)[j], result)
		}
	}
}

type IterateDiffFunc func(fieldGroupLi []FieldGroup, destData *map[string]interface{}, srcData map[string]interface{}, result *interface{})

func (o ModelIterator) IterateDiffBo(dataSource DataSource, destBo *map[string]interface{}, srcBo map[string]interface{}, result *interface{}, iterateFunc IterateDiffFunc) {
	o.iterateDiffMasterDataBo(dataSource, destBo, srcBo, result, iterateFunc)
	o.iterateDiffDetailDataBo(dataSource, destBo, srcBo, result, iterateFunc)
}

func (o ModelIterator) iterateDiffMasterDataBo(dataSource DataSource, destBo *map[string]interface{}, srcBo map[string]interface{}, result *interface{}, iterateFunc IterateDiffFunc) {
	masterFieldGroupLi := o.getDataSetFieldGroupLi(&dataSource.MasterData.FixField, &dataSource.MasterData.BizField)

	destData := (*destBo)["A"].(map[string]interface{})
	srcData := srcBo["A"].(map[string]interface{})
	
	fieldGroupLi := []FieldGroup{}
	for _, item := range *masterFieldGroupLi {
		fieldGroupLi = append(fieldGroupLi, *item)
	}
	iterateFunc(fieldGroupLi, &destData, srcData, result)
}

func (o ModelIterator) iterateDiffDetailDataBo(dataSource DataSource, destBo *map[string]interface{}, srcBo map[string]interface{}, result *interface{}, iterateFunc IterateDiffFunc) {
	for i, _ := range dataSource.DetailDataLi {
		detailFieldGroupLi := o.getDataSetFieldGroupLi(&dataSource.DetailDataLi[i].FixField, &dataSource.DetailDataLi[i].BizField)
		fieldGroupLi := []FieldGroup{}
		for _, item := range *detailFieldGroupLi {
			fieldGroupLi = append(fieldGroupLi, *item)
		}
		
		iDestDataLi := (*destBo)[dataSource.DetailDataLi[i].Id].([]interface{})
		iSrcDataLi := srcBo[dataSource.DetailDataLi[i].Id].([]interface{})
		
		destDataLi := []map[string]interface{}{}
		srcDataLi := []map[string]interface{}{}
		
		for j, _ := range iDestDataLi {
			item := iDestDataLi[j].(map[string]interface{})
			destDataLi = append(destDataLi, item)
		}
		
		for j, _ := range iSrcDataLi {
			item := iSrcDataLi[j].(map[string]interface{})
			srcDataLi = append(srcDataLi, item)
		}
		
		destDataIdDict := o.getDataIdDict(&destDataLi)
		srcDataIdDict := o.getDataIdDict(&srcDataLi)
		
		// delete,destData is nil
		for j, _ := range srcDataLi {
			dataItem := srcDataLi[j]
			id, _ := strconv.Atoi(fmt.Sprint(dataItem["id"]))
			if (*destDataIdDict)[id] == nil {
				destData := (*map[string]interface{})(nil)
				srcData := dataItem
				iterateFunc(fieldGroupLi, destData, srcData, result)
			}
		}
		
		// insert, destData is nil
		for j, _ := range destDataLi {
			dataItem := destDataLi[j]
			idStr := fmt.Sprint(dataItem["id"])
			if idStr == "" || idStr == "0" {
				destData := dataItem
				srcData := (map[string]interface{})(nil)
				iterateFunc(fieldGroupLi, &destData, srcData, result)
			}
		}
		
		// modify,
		for j, _ := range destDataLi {
			dataItem := destDataLi[j]
			idStr := fmt.Sprint(dataItem["id"])
			if idStr != "" && idStr != "0" {
				id, _ := strconv.Atoi(idStr)
				if (*srcDataIdDict)[id] != nil {// insert时会加id,因此,需要做srcDataIdDict的空判断
					destData := (*destDataIdDict)[id].(map[string]interface{})
					srcData := (*srcDataIdDict)[id].(map[string]interface{})
					iterateFunc(fieldGroupLi, &destData, srcData, result)
				}
			}
		}
	}
}

func (o ModelIterator) getDataIdDict(dataLi *[]map[string]interface{}) *map[int]interface{} {
	destDataIdDict := map[int]interface{}{}
	for i, _ := range *dataLi {
		dataItem := (*dataLi)[i]
		strId := fmt.Sprint(dataItem["id"])
		if strId != "" {
			id, err := strconv.Atoi(strId)
			if err != nil {
				panic(err)
			}
			if id > 0 {
				destDataIdDict[id] = dataItem
			}
		}
	}
	return &destDataIdDict
}

func (o ModelIterator) getDataSetFieldGroupLi(fixField *FixField, bizField *BizField) *[]*FieldGroup {
	fieldGroupLi := []*FieldGroup{}
	fixFieldLi := o.GetFixFieldLi(fixField)
	for i, _ := range *fixFieldLi {
		fieldGroupLi = append(fieldGroupLi, (*fixFieldLi)[i])
	}
	for i, _ := range bizField.FieldLi {
		fieldGroupLi = append(fieldGroupLi, &bizField.FieldLi[i].FieldGroup)
	}
	return &fieldGroupLi
}

type IterateDataFunc func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{})

func (o ModelIterator) IterateDataBo(dataSource DataSource, bo *map[string]interface{}, result *interface{}, iterateFunc IterateDataFunc) {
	o.iterateMasterDataBo(dataSource, bo, result, iterateFunc)
	o.iterateDetailDataBo(dataSource, bo, result, iterateFunc)
}

func (o ModelIterator) iterateMasterDataBo(dataSource DataSource, bo *map[string]interface{}, result *interface{}, iterateFunc IterateDataFunc) {
	masterFieldGroupLi := o.getDataSetFieldGroupLi(&dataSource.MasterData.FixField, &dataSource.MasterData.BizField)
	data := (*bo)["A"].(map[string]interface{})
	fieldGroupLi := []FieldGroup{}
	for _, item := range *masterFieldGroupLi {
		fieldGroupLi = append(fieldGroupLi, *item)
	}
	rowIndex := 0
	iterateFunc(fieldGroupLi, &data, rowIndex, result)
}

func (o ModelIterator) iterateDetailDataBo(dataSource DataSource, bo *map[string]interface{}, result *interface{}, iterateFunc IterateDataFunc) {
	for i, _ := range dataSource.DetailDataLi {
		item := dataSource.DetailDataLi[i]
		detailFieldGroupLi := o.getDataSetFieldGroupLi(&item.FixField, &item.BizField)
		fieldGroupLi := []FieldGroup{}
		for _, item := range *detailFieldGroupLi {
			fieldGroupLi = append(fieldGroupLi, *item)
		}
		dataLi := (*bo)[item.Id].([]interface{})
		for j, _ := range dataLi {
			data := dataLi[j].(map[string]interface{})
			iterateFunc(fieldGroupLi, &data, j, result)
		}
	}
}

