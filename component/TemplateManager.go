package component

import (
	//	"github.com/hongjinqiu/gometa/dictionary"
	"github.com/hongjinqiu/gometa/config"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/interceptor"
	"github.com/hongjinqiu/gometa/layer"
	//	"github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/script"
	//	"github.com/hongjinqiu/gometa/tree"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var templaterwlock sync.RWMutex = sync.RWMutex{}
var gListTemplateDict map[string]ListTemplateInfo = map[string]ListTemplateInfo{}
var gFormTemplateDict map[string]FormTemplateInfo = map[string]FormTemplateInfo{}
var gSelectorTemplateDict map[string]SelectorTemplateInfo = map[string]SelectorTemplateInfo{}

type ListTemplateInfo struct {
	Path         string
	ListTemplate ListTemplate
}

type FormTemplateInfo struct {
	Path         string
	FormTemplate FormTemplate
}

type SelectorTemplateInfo struct {
	Path         string
	ListTemplate ListTemplate
}

type TemplateManager struct{}

func (o TemplateManager) GetListTemplateInfoLi() []ListTemplateInfo {
	listTemplateInfo := []ListTemplateInfo{}
	if len(gListTemplateDict) == 0 {
		o.loadListTemplate()
	}
	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	for _, item := range gListTemplateDict {
		listTemplateInfo = append(listTemplateInfo, item)
	}
	return listTemplateInfo
}

func (o TemplateManager) RefretorListTemplateInfo() []ListTemplateInfo {
	o.clearListTemplate()
	o.loadListTemplate()
	listTemplateInfo := []ListTemplateInfo{}
	for _, item := range gListTemplateDict {
		listTemplateInfo = append(listTemplateInfo, item)
	}
	return listTemplateInfo
}

func (o TemplateManager) GetListTemplate(id string) ListTemplate {
	return o.GetListTemplateInfo(id).ListTemplate
}

func (o TemplateManager) GetListTemplateInfo(id string) ListTemplateInfo {
	if config.String("mode.dev") == "true" {
		listTemplateInfo, found := o.findListTemplateInfo(id)
		if found {
			listTemplateInfo, err := o.loadSingleListTemplateWithLock(listTemplateInfo.Path)
			if err != nil {
				panic(err)
			}
			if listTemplateInfo.ListTemplate.Id == id {
				return listTemplateInfo
			}
		}
		o.clearListTemplate()
		o.loadListTemplate()
		listTemplateInfo, found = o.findListTemplateInfo(id)
		if found {
			return listTemplateInfo
		}
		panic(id + " not exists in ListTemplate list")
	}

	if len(gListTemplateDict) == 0 {
		o.loadListTemplate()
	}
	listTemplateInfo, found := o.findListTemplateInfo(id)
	if found {
		return listTemplateInfo
	}
	panic(id + " not exists in ListTemplate list")
}

func (o TemplateManager) findListTemplateInfo(id string) (ListTemplateInfo, bool) {
	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	if gListTemplateDict[id].Path != "" {
		return gListTemplateDict[id], true
	}
	return ListTemplateInfo{}, false
}

func (o TemplateManager) clearListTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	gListTemplateDict = map[string]ListTemplateInfo{}
}

func (o TemplateManager) loadListTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	path := config.String("LIST_TEMPLATE_PATH")
	if path != "" {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Index(path, "list_") > -1 && strings.Index(path, ".xml") > -1 && !info.IsDir() {
				_, err = o.loadSingleListTemplate(path)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (o TemplateManager) loadSingleListTemplateWithLock(path string) (ListTemplateInfo, error) {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	return o.loadSingleListTemplate(path)
}

func (o TemplateManager) loadSingleListTemplate(path string) (ListTemplateInfo, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return ListTemplateInfo{}, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return ListTemplateInfo{}, err
	}

	listTemplate := ListTemplate{}
	err = xml.Unmarshal(data, &listTemplate)
	if err != nil {
		return ListTemplateInfo{}, err
	}

	if listTemplate.Adapter.Name != "" {
		classMethod := listTemplate.Adapter.Name + ".ApplyAdapter"
		commonMethod := CommonMethod{}
		paramLi := []interface{}{listTemplate}
		values := commonMethod.Parse(classMethod, paramLi)
		listTemplate = values[0].Interface().(ListTemplate)
	}

	listTemplateInfo := ListTemplateInfo{
		Path:         path,
		ListTemplate: listTemplate,
	}
	gListTemplateDict[listTemplate.Id] = listTemplateInfo
	return listTemplateInfo, nil
}

func (o TemplateManager) GetSelectorTemplateInfoLi() []SelectorTemplateInfo {
	selectorTemplateInfo := []SelectorTemplateInfo{}
	if len(gSelectorTemplateDict) == 0 {
		o.loadSelectorTemplate()
	}
	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	for _, item := range gSelectorTemplateDict {
		selectorTemplateInfo = append(selectorTemplateInfo, item)
	}
	return selectorTemplateInfo
}

func (o TemplateManager) RefretorSelectorTemplateInfo() []SelectorTemplateInfo {
	o.clearSelectorTemplate()
	o.loadSelectorTemplate()
	selectorTemplateInfo := []SelectorTemplateInfo{}
	for _, item := range gSelectorTemplateDict {
		selectorTemplateInfo = append(selectorTemplateInfo, item)
	}
	return selectorTemplateInfo
}

func (o TemplateManager) GetSelectorTemplate(id string) ListTemplate {
	return o.GetSelectorTemplateInfo(id).ListTemplate
}

func (o TemplateManager) GetSelectorTemplateInfo(id string) SelectorTemplateInfo {
	if config.String("mode.dev") == "true" {
		selectorTemplateInfo, found := o.findSelectorTemplateInfo(id)
		if found {
			selectorTemplateInfo, err := o.loadSingleSelectorTemplateWithLock(selectorTemplateInfo.Path)
			if err != nil {
				panic(err)
			}
			if strings.Index(selectorTemplateInfo.Path, "list_") == -1 {
				if selectorTemplateInfo.ListTemplate.Id == id {
					return selectorTemplateInfo
				}
			} else {
				if selectorTemplateInfo.ListTemplate.SelectorId == id {
					return selectorTemplateInfo
				}
			}
		}
		o.clearSelectorTemplate()
		o.loadSelectorTemplate()
		selectorTemplateInfo, found = o.findSelectorTemplateInfo(id)
		if found {
			return selectorTemplateInfo
		}
		panic(id + " not exists in ListTemplate list")
	}

	if len(gSelectorTemplateDict) == 0 {
		o.loadSelectorTemplate()
	}
	selectorTemplateInfo, found := o.findSelectorTemplateInfo(id)
	if found {
		return selectorTemplateInfo
	}
	panic(id + " not exists in ListTemplate list")
}

func (o TemplateManager) findSelectorTemplateInfo(id string) (SelectorTemplateInfo, bool) {
	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	if gSelectorTemplateDict[id].Path != "" {
		return gSelectorTemplateDict[id], true
	}
	return SelectorTemplateInfo{}, false
}

func (o TemplateManager) clearSelectorTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	gSelectorTemplateDict = map[string]SelectorTemplateInfo{}
}

func (o TemplateManager) loadSelectorTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	path := config.String("SELECTOR_TEMPLATE_PATH")
	if path != "" {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Index(path, "list_") > -1 && strings.Index(path, ".xml") > -1 && !info.IsDir() {
				_, err = o.loadSingleSelectorTemplate(path)
				if err != nil {
					return err
				}
			}

			return nil
		})
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Index(path, "selector_") > -1 && strings.Index(path, ".xml") > -1 && !info.IsDir() {
				_, err = o.loadSingleSelectorTemplate(path)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (o TemplateManager) loadSingleSelectorTemplateWithLock(path string) (SelectorTemplateInfo, error) {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	return o.loadSingleSelectorTemplate(path)
}

func (o TemplateManager) loadSingleSelectorTemplate(path string) (SelectorTemplateInfo, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return SelectorTemplateInfo{}, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return SelectorTemplateInfo{}, err
	}

	listTemplate := ListTemplate{}
	err = xml.Unmarshal(data, &listTemplate)
	if err != nil {
		return SelectorTemplateInfo{}, err
	}

	isAdd := true
	if strings.Index(path, "list_") > -1 {
		if listTemplate.SelectorId == "" {
			isAdd = false
		}
	}
	if isAdd {
		if listTemplate.Adapter.Name != "" {
			classMethod := listTemplate.Adapter.Name + ".ApplyAdapter"
			commonMethod := CommonMethod{}
			paramLi := []interface{}{listTemplate}
			values := commonMethod.Parse(classMethod, paramLi)
			listTemplate = values[0].Interface().(ListTemplate)
		}

		selectorTemplateInfo := SelectorTemplateInfo{
			Path:         path,
			ListTemplate: listTemplate,
		}
		if strings.Index(path, "list_") > -1 {
			gSelectorTemplateDict[listTemplate.SelectorId] = selectorTemplateInfo
		} else {
			gSelectorTemplateDict[listTemplate.Id] = selectorTemplateInfo
		}
		return selectorTemplateInfo, nil
	}
	return SelectorTemplateInfo{}, nil
}

func (o TemplateManager) GetFormTemplateInfoLi() []FormTemplateInfo {
	formTemplateInfo := []FormTemplateInfo{}
	if len(gFormTemplateDict) == 0 {
		o.loadFormTemplate()
	}

	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	for _, item := range gFormTemplateDict {
		formTemplateInfo = append(formTemplateInfo, item)
	}
	return formTemplateInfo
}

func (o TemplateManager) RefretorFormTemplateInfo() []FormTemplateInfo {
	o.clearFormTemplate()
	o.loadFormTemplate()
	formTemplateInfo := []FormTemplateInfo{}
	for _, item := range gFormTemplateDict {
		formTemplateInfo = append(formTemplateInfo, item)
	}
	return formTemplateInfo
}

func (o TemplateManager) GetFormTemplate(id string) FormTemplate {
	return o.GetFormTemplateInfo(id).FormTemplate
}

func (o TemplateManager) GetFormTemplateInfo(id string) FormTemplateInfo {
	if config.String("mode.dev") == "true" {
		formTemplateInfo, found := o.findFormTemplateInfo(id)
		if found {
			formTemplateInfo, err := o.loadSingleFormTemplateWithLock(formTemplateInfo.Path)
			if err != nil {
				panic(err)
			}
			if formTemplateInfo.FormTemplate.Id == id {
				return formTemplateInfo
			}
		}
		o.clearFormTemplate()
		o.loadFormTemplate()
		formTemplateInfo, found = o.findFormTemplateInfo(id)
		if found {
			return formTemplateInfo
		}
		panic(id + " not exists in FormTemplate list")
	}

	if len(gFormTemplateDict) == 0 {
		o.loadFormTemplate()
	}
	formTemplateInfo, found := o.findFormTemplateInfo(id)
	if found {
		return formTemplateInfo
	}
	panic(id + " not exists in FormTemplate list")
}

func (o TemplateManager) findFormTemplateInfo(id string) (FormTemplateInfo, bool) {
	templaterwlock.RLock()
	defer templaterwlock.RUnlock()

	if gFormTemplateDict[id].Path != "" {
		return gFormTemplateDict[id], true
	}
	return FormTemplateInfo{}, false
}

func (o TemplateManager) clearFormTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	gFormTemplateDict = map[string]FormTemplateInfo{}
}

func (o TemplateManager) loadFormTemplate() {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	path := config.String("FORM_TEMPLATE_PATH")
	if path != "" {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Index(path, "form_") > -1 && strings.Index(path, ".xml") > -1 && !info.IsDir() {
				_, err = o.loadSingleFormTemplate(path)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (o TemplateManager) loadSingleFormTemplateWithLock(path string) (FormTemplateInfo, error) {
	templaterwlock.Lock()
	defer templaterwlock.Unlock()

	return o.loadSingleFormTemplate(path)
}

func (o TemplateManager) loadSingleFormTemplate(path string) (FormTemplateInfo, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return FormTemplateInfo{}, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return FormTemplateInfo{}, err
	}

	formTemplate := FormTemplate{}
	err = xml.Unmarshal(data, &formTemplate)
	if err != nil {
		return FormTemplateInfo{}, err
	}

	for i, _ := range formTemplate.FormElemLi {
		formElem := &formTemplate.FormElemLi[i]
		if formElem.XMLName.Local == "html" {
			formElemXmlData, err := xml.Marshal(formElem)
			if err != nil {
				panic(err)
			}
			err = xml.Unmarshal(formElemXmlData, &(formElem.Html))
			if err != nil {
				panic(err)
			}
		} else if formElem.XMLName.Local == "toolbar" {
			formElemXmlData, err := xml.Marshal(formElem)
			if err != nil {
				panic(err)
			}
			err = xml.Unmarshal(formElemXmlData, &(formElem.Toolbar))
			if err != nil {
				panic(err)
			}
		} else if formElem.XMLName.Local == "column-model" {
			formElemXmlData, err := xml.Marshal(formElem)
			if err != nil {
				panic(err)
			}
			err = xml.Unmarshal(formElemXmlData, &(formElem.ColumnModel))
			if err != nil {
				panic(err)
			}
		}
	}

	if formTemplate.Adapter.Name != "" {
		classMethod := formTemplate.Adapter.Name + ".ApplyAdapter"
		commonMethod := CommonMethod{}
		paramLi := []interface{}{formTemplate}
		values := commonMethod.Parse(classMethod, paramLi)
		formTemplate = values[0].Interface().(FormTemplate)
	}

	formTemplateInfo := FormTemplateInfo{
		Path:         path,
		FormTemplate: formTemplate,
	}
	gFormTemplateDict[formTemplate.Id] = formTemplateInfo
	return formTemplateInfo, nil
}

func (o TemplateManager) QueryDataForListTemplate(sessionId int, listTemplate *ListTemplate, paramMap map[string]string, pageNo int, pageSize int) map[string]interface{} {
	interceptorManager := InterceptorManager{}
	paramMap = interceptorManager.ParseBeforeBuildQuery(sessionId, listTemplate.BeforeBuildQuery, paramMap)

	session, _ := global.GetConnection(sessionId)
	querySupport := QuerySupport{}
	queryMap := map[string]interface{}{}
	permissionSupport := PermissionSupport{}
	permissionQueryDict := permissionSupport.GetPermissionQueryDict(sessionId, listTemplate.Security)
	for k, v := range permissionQueryDict {
		queryMap[k] = v
	}

	queryLi := []map[string]interface{}{}

	collection := listTemplate.DataProvider.Collection
	mapStr := listTemplate.DataProvider.Map
	reduce := listTemplate.DataProvider.Reduce
	fixBsonQuery := listTemplate.DataProvider.FixBsonQuery

	if fixBsonQuery != "" {
		fixBsonQueryMap := map[string]interface{}{}
		err := json.Unmarshal([]byte(fixBsonQuery), &fixBsonQueryMap)
		if err != nil {
			panic(err)
		}

		queryLi = append(queryLi, fixBsonQueryMap)
	}

	queryParameters := listTemplate.QueryParameterGroup.QueryParameterLi
	queryParameterBuilder := QueryParameterBuilder{}
	for _, queryParameter := range queryParameters {
		if queryParameter.Editor != "" && queryParameter.Restriction != "" && queryParameter.UseIn != "none" {
			name := queryParameter.Name
			if paramMap[name] != "" {
				//				if listTemplate.Adapter.Name != "" {
				//					classMethod := listTemplate.Adapter.Name + ".ApplyQueryParameter"
				//					commonMethod := CommonMethod{}
				//					paramLi := []interface{}{}
				//					var listTemplateParam interface{} = listTemplate
				//					paramLi = append(paramLi, listTemplateParam)
				//					var queryParameterParam interface{} = queryParameter
				//					paramLi = append(paramLi, queryParameterParam)
				//					commonMethod.Parse(classMethod, paramLi)
				//				}
				if listTemplate.QueryParameterGroup.DataSetId != "" {
					queryParameter.DataSetId = listTemplate.QueryParameterGroup.DataSetId
				}
				queryParameterMap := queryParameterBuilder.buildQuery(queryParameter, paramMap[name])
				queryLi = append(queryLi, queryParameterMap)
			}
		}
	}

	queryLi = interceptorManager.ParseAfterBuildQuery(sessionId, listTemplate.AfterBuildQuery, queryLi)

	if len(queryLi) == 1 {
		for k, v := range queryLi[0] {
			queryMap[k] = v
		}
	} else if len(queryLi) > 1 {
		queryMap["$and"] = queryLi
	}

	queryByte, err := json.MarshalIndent(queryMap, "", "\t")
	if err != nil {
		panic(err)
	}
	if mapStr == "" {
		orderBy := listTemplate.ColumnModel.BsonOrderBy
		log.Println("QueryDataForListTemplate,collection:" + collection + ",query is:" + string(queryByte) + ",orderBy is:" + orderBy)
		result := querySupport.IndexWithSession(session, collection, queryMap, pageNo, pageSize, orderBy)
		items := result["items"].([]interface{})
		items = interceptorManager.ParseAfterQueryData(sessionId, listTemplate.AfterQueryData, listTemplate.ColumnModel.DataSetId, items)
		result["items"] = items
		return result
	}
	mapReduce := mgo.MapReduce{
		Map:    mapStr,
		Reduce: reduce,
	}

	mapReduceByte, err := json.MarshalIndent(mapReduce, "", "\t")
	if err != nil {
		panic(err)
	}

	log.Println("QueryDataForListTemplate,collection:" + collection + ",query is:" + string(queryByte) + ",mapReduce:" + string(mapReduceByte))
	mapReduceLi := querySupport.MapReduceAll(collection, queryMap, mapReduce)
	items := []interface{}{}
	for _, item := range mapReduceLi {
		item["id"] = item["_id"]
		items = append(items, item)
	}
	items = interceptorManager.ParseAfterQueryData(sessionId, listTemplate.AfterQueryData, listTemplate.ColumnModel.DataSetId, items)
	return map[string]interface{}{
		"totalResults": len(mapReduceLi),
		"items":        items,
	}
}

func (o TemplateManager) GetColumnModelDataForListTemplate(sessionId int, listTemplate ListTemplate, items []interface{}) map[string]interface{} {
	//	o.applyAdapterColumnName(listTemplate)
	return o.GetColumnModelDataForColumnModel(sessionId, listTemplate.ColumnModel, items)
}

func (o TemplateManager) GetColumnModelDataForFormTemplate(sessionId int, formTemplate FormTemplate, bo map[string]interface{}) map[string]interface{} {
	relationBo := map[string]interface{}{}
	formTemplateIterator := FormTemplateIterator{}
	var result interface{} = ""
	formTemplateIterator.IterateAllTemplateColumnModel(formTemplate, &result, func(columnModel ColumnModel, result *interface{}) {
		if bo[columnModel.DataSetId] != nil {
			if columnModel.DataSetId == "A" {
				items := []interface{}{
					map[string]interface{}{
						"A": bo["A"],
					},
				}
				columnModelData := o.GetColumnModelDataForColumnModel(sessionId, columnModel, items)
				columnModelItems := columnModelData["items"].([]interface{})
				if len(columnModelItems) > 0 {
					// 主数据集有多个,因此,需要累加
					if bo["A"] == nil {
						bo["A"] = map[string]interface{}{}
					}
					aDict := bo["A"].(map[string]interface{})
					item0 := columnModelItems[0].(map[string]interface{})
					for key, value := range item0 {
						aDict[key] = value
					}
					bo["A"] = aDict
					columnModelRelationBo := columnModelData["relationBo"].(map[string]interface{})
					o.mergeRelationBo(&relationBo, columnModelRelationBo)
				}
			} else {
				items := bo[columnModel.DataSetId].([]interface{})
				dataSetItems := []interface{}{}
				for _, item := range items {
					dataSetItems = append(dataSetItems, map[string]interface{}{
						columnModel.DataSetId: item,
					})
				}
				columnModelData := o.GetColumnModelDataForColumnModel(sessionId, columnModel, dataSetItems)

				items = columnModelData["items"].([]interface{})

				bo[columnModel.DataSetId] = items
				columnModelRelationBo := columnModelData["relationBo"].(map[string]interface{})
				o.mergeRelationBo(&relationBo, columnModelRelationBo)
			}
		}
	})
	o.mergeSelectorInfo2RelationBo(formTemplate, &relationBo)
	return map[string]interface{}{
		"bo":         bo,
		"relationBo": relationBo,
	}
}

func (o TemplateManager) mergeRelationBo(relationBo *map[string]interface{}, columnModelRelationBo map[string]interface{}) {
	for key, item := range columnModelRelationBo {
		if (*relationBo)[key] == nil {
			(*relationBo)[key] = item
		} else {
			relationBoItem := map[string]interface{}{}
			if (*relationBo)[key] != nil {
				relationBoItem = (*relationBo)[key].(map[string]interface{})
			}
			columnModelRelationBoItem := item.(map[string]interface{})
			for subKey, subItem := range columnModelRelationBoItem {
				if relationBoItem[subKey] == nil {
					relationBoItem[subKey] = subItem
				}
			}

			(*relationBo)[key] = relationBoItem
		}
	}
}

//func (o TemplateManager) applyAdapterColumnName(listTemplate ListTemplate) {
//	if listTemplate.Adapter.Name != "" {
//		//ApplyColumnName(listTemplate *ListTemplate, column *Column) {
//		classMethod := listTemplate.Adapter.Name + ".ApplyColumnName"
//		commonMethod := CommonMethod{}
//		paramLi := []interface{}{}
//		paramLi = append(paramLi, listTemplate)
//		paramLi = append(paramLi, listTemplate.ColumnModel.IdColumn)
//		commonMethod.Parse(classMethod, paramLi)
//		for i, _ := range listTemplate.ColumnModel.ColumnLi {
//			o.recursionApplyAdapterColumnName(listTemplate, listTemplate.ColumnModel.ColumnLi[i])
//		}
//	}
//}

// TODO, bytest
//func (o TemplateManager) recursionApplyAdapterColumnName(listTemplate ListTemplate, columnItem Column) {
//	if columnItem.XMLName.Local != "virtual-column" {
//		if columnItem.ColumnModel.ColumnLi != nil {
//			for i, _ := range columnItem.ColumnModel.ColumnLi {
//				o.recursionApplyAdapterColumnName(listTemplate, columnItem.ColumnModel.ColumnLi[i])
//			}
//		} else {
//			commonMethod := CommonMethod{}
//			classMethod := listTemplate.Adapter.Name + ".ApplyColumnName"
//			paramLi := []interface{}{}
//			var listTemplateParam interface{} = listTemplate
//			paramLi = append(paramLi, listTemplateParam)
//			var columnItemParam interface{} = columnItem
//			paramLi = append(paramLi, columnItemParam)
//			commonMethod.Parse(classMethod, paramLi)
//		}
//	}
//}

func (o TemplateManager) GetSelectorInfoForListTemplate(listTemplate ListTemplate) map[string]interface{} {
	result := map[string]interface{}{}

	listTemplateIterator := ListTemplateIterator{}
	var iterateResult interface{} = ""
	listTemplateIterator.IterateTemplateColumn(listTemplate, &iterateResult, func(column Column, iterateResult *interface{}) {
		if column.CRelationDS.CRelationItemLi != nil {
			for _, relationItem := range column.CRelationDS.CRelationItemLi {
				selectorName := relationItem.CRelationConfig.SelectorName
				listTemplate := o.GetSelectorTemplate(selectorName)
				selectorInfo := map[string]interface{}{
					"Description": listTemplate.Description,
					"url":         o.GetViewUrl(listTemplate),
				}
				result[selectorName] = selectorInfo
			}
		}
	})
	listTemplateIterator.IterateTemplateQueryParameter(listTemplate, &iterateResult, func(queryParameter QueryParameter, iterateResult *interface{}) {
		if queryParameter.CRelationDS.CRelationItemLi != nil {
			for _, relationItem := range queryParameter.CRelationDS.CRelationItemLi {
				selectorName := relationItem.CRelationConfig.SelectorName
				listTemplate := o.GetSelectorTemplate(selectorName)
				selectorInfo := map[string]interface{}{
					"Description": listTemplate.Description,
					"url":         o.GetViewUrl(listTemplate),
				}
				result[selectorName] = selectorInfo
			}
		}
	})

	return result
}

func (o TemplateManager) GetColumnModelDataForColumnModel(sessionId int, columnModel ColumnModel, items []interface{}) map[string]interface{} {
	// set dataSetId to columnItem.DataSetId
	o.recursionSetDefaultDataSetId(columnModel.DataSetId, &columnModel)
	columnDict := map[string]Column{}
	o.recursionGetColumnItem(columnModel, &columnDict)
	columnNodeDict := o.buildColumnNode(columnDict)
	relationBo := map[string]interface{}{}

	columnModelItems := []interface{}{}
	expressionParser := ExpressionParser{}
	for _, item := range items {
		record := item.(map[string]interface{})
		record["pendingTransactions"] = []interface{}{}
		recordJsonByte, err := json.Marshal(record)
		if err != nil {
			panic(err)
		}
		recordJson := string(recordJsonByte)

		loopItem := map[string]interface{}{}
		loopItem[columnModel.CheckboxColumn.Name] = expressionParser.Parse(recordJson, columnModel.CheckboxColumn.Expression)
		idColumnName := columnModel.IdColumn.Name
		if idColumnName != "" {
			if columnModel.IdColumn.DataSetId != "" {
				idColumnName = columnModel.IdColumn.DataSetId + "." + idColumnName
			}
			loopItem[columnModel.IdColumn.Name] = o.getValueBySpot(record, idColumnName)
			loopItem["id"] = o.getValueBySpot(record, idColumnName)
			loopItem["_id"] = o.getValueBySpot(record, idColumnName)
		}
		for _, columnNode := range columnNodeDict {
			o.GetColumnModelDataForColumnItem(sessionId, columnNodeDict, columnNode, record, &relationBo, &loopItem)
		}

		columnModelItems = append(columnModelItems, loopItem)
	}

	return map[string]interface{}{
		"items":      columnModelItems,
		"relationBo": relationBo,
	}
}

func (o TemplateManager) recursionSetDefaultDataSetId(dataSetId string, columnModel *ColumnModel) {
	idColumn := &columnModel.IdColumn
	if idColumn.DataSetId == "" {
		idColumn.DataSetId = dataSetId
	}
	if columnModel.ColumnLi != nil {
		for i, _ := range columnModel.ColumnLi {
			columnItem := &columnModel.ColumnLi[i]
			o.recursionSetDefaultDataSetId(dataSetId, &columnItem.ColumnModel)
			if columnItem.DataSetId == "" {
				columnItem.DataSetId = dataSetId
			}
		}
	}
	if columnModel.DataSetId == "" {
		columnModel.DataSetId = dataSetId
	}
}

type ColumnNode struct {
	preColumnLi []string // 前驱
	column      Column
}

func (o TemplateManager) recursionGetColumnItem(columnModel ColumnModel, columnDict *map[string]Column) {
	for _, columnItem := range columnModel.ColumnLi {
		if columnItem.ColumnModel.ColumnLi != nil {
			o.recursionGetColumnItem(columnItem.ColumnModel, columnDict)
			(*columnDict)[columnItem.Name] = columnItem
		} else {
			(*columnDict)[columnItem.Name] = columnItem
		}
	}
}

func (o TemplateManager) buildColumnNode(columnDict map[string]Column) map[string]ColumnNode {
	columnNodeDict := map[string]ColumnNode{}
	for _, columnItem := range columnDict {
		columnNode := ColumnNode{
			column: columnItem,
		}
		for _, subColumnItem := range columnDict {
			if subColumnItem.XMLName.Local == "select-column" {
				if subColumnItem.CRelationDS.CRelationItemLi != nil {
					for _, relationItem := range subColumnItem.CRelationDS.CRelationItemLi {
						if relationItem.CCopyConfigLi != nil {
							for _, copyConfig := range relationItem.CCopyConfigLi {
								if copyConfig.CopyColumnName == columnItem.Name {
									if columnNode.preColumnLi == nil {
										columnNode.preColumnLi = []string{}
									}
									isIn := false
									for _, preColumnItem := range columnNode.preColumnLi {
										if subColumnItem.Name == preColumnItem {
											isIn = true
											break
										}
									}
									if !isIn {
										columnNode.preColumnLi = append(columnNode.preColumnLi, subColumnItem.Name)
									}
								}
							}
						}
					}
				}
			}
		}
		columnNodeDict[columnItem.Name] = columnNode
	}
	return columnNodeDict
}

func (o TemplateManager) parseModelExpression(bo map[string]interface{}, data map[string]interface{}, mode string, content string) string {
	bo["pendingTransactions"] = []interface{}{}
	expressionParser := ExpressionParser{}
	if mode == "" || mode == "text" {
		return content
	} else if mode == "python" {
		dataJsonData, err := json.Marshal(&data)
		if err != nil {
			panic(err)
		}
		dataJson := string(dataJsonData)
		boJsonData, _ := json.Marshal(&bo)
		boJson := string(boJsonData)
		return expressionParser.ParseModel(boJson, dataJson, content)
	} else if mode == "golang" {
		return expressionParser.ParseGolang(bo, data, content)
	}
	return ""
}

func (o TemplateManager) applyCopyValueField(sessionId int, preColumnName string, columnItem Column, relationItem CRelationItem, record map[string]interface{}, relationBo *map[string]interface{}, loopItem *map[string]interface{}) {
	for _, copyConfig := range relationItem.CCopyConfigLi {
		if copyConfig.CopyColumnName == columnItem.Name {
			selectorName := relationItem.CRelationConfig.SelectorName
			if (*relationBo)[selectorName] != nil {
				selectorDict := (*relationBo)[selectorName].(map[string]interface{})
				if (*loopItem)[preColumnName] != nil {
					id := fmt.Sprint((*loopItem)[preColumnName])
					if selectorDict[id] != nil {
						selectorData := selectorDict[id].(map[string]interface{})
						valueFieldLi := strings.Split(copyConfig.CopyValueField, ",")
						valueLi := []string{}
						for _, valueField := range valueFieldLi {
							if selectorData[valueField] != nil {
								valueLi = append(valueLi, fmt.Sprint(selectorData[valueField]))
							}
						}
						(*loopItem)[columnItem.Name] = strings.Join(valueLi, ",")
						//如果是select-column,取rs中的值出来,放到relationBo里面,
						if columnItem.XMLName.Local == "select-column" {
							if (*loopItem)[columnItem.Name] != nil {
								valueStr := fmt.Sprint((*loopItem)[columnItem.Name])
								if valueStr != "" {
									value, err := strconv.Atoi(valueStr)
									if err != nil {
										panic(err)
									}
									o.applyRelationBoBySelectField(sessionId, columnItem, value, record, relationBo)
								}
							}
						}
					}
				}
			}
			break
		}
	}
}

func (o TemplateManager) getData4Expression(column Column, record map[string]interface{}) map[string]interface{} {
	return record[column.DataSetId].(map[string]interface{})
	//	if column.DataSetId == "A" {
	//		return record["A"].(map[string]interface{})
	//	}
	//	return record
}

func (o TemplateManager) GetColumnModelDataForColumnItem(sessionId int, columnNodeDict map[string]ColumnNode, columnNode ColumnNode, record map[string]interface{}, relationBo *map[string]interface{}, loopItem *map[string]interface{}) {
	columnItem := columnNode.column
	if columnNode.preColumnLi != nil {
		for _, name := range columnNode.preColumnLi {
			o.GetColumnModelDataForColumnItem(sessionId, columnNodeDict, columnNodeDict[name], record, relationBo, loopItem)
		}
		// 算完前驱,算自身,
		for _, name := range columnNode.preColumnLi {
			preColumnNode := columnNodeDict[name]
			relationLi := preColumnNode.column.CRelationDS.CRelationItemLi
			if relationLi != nil {
				for _, relationItem := range relationLi {
					mode := relationItem.CRelationExpr.Mode
					exprContent := relationItem.CRelationExpr.Content
					bo := record
					data := o.getData4Expression(columnItem, record)
					content := o.parseModelExpression(bo, data, mode, exprContent)
					if strings.ToLower(content) == "true" {
						if relationItem.CCopyConfigLi != nil {
							o.applyCopyValueField(sessionId, name, columnItem, relationItem, record, relationBo, loopItem)
						}
						break
					}
				}
			}
		}
	} else {
		if columnItem.XMLName.Local != "virtual-column" {
			columnItemName := columnItem.Name
			if columnItem.DataSetId != "" {
				columnItemName = columnItem.DataSetId + "." + columnItemName
			}
			(*loopItem)[columnItem.Name] = o.getValueBySpot(record, columnItemName)
			//o.ApplyDictionaryColumnData(loopItem, columnItem)
			o.ApplyScriptColumnData(loopItem, record, columnItem)
			// 如果是select-column,取rs中的值出来,放到relationBo里面,多选,按逗号分隔
			if columnItem.XMLName.Local == "select-column" {
				if (*loopItem)[columnItem.Name] != nil {
					valueStr := fmt.Sprint((*loopItem)[columnItem.Name])
					if valueStr != "" {
						for _, item := range strings.Split(valueStr, ",") {
							value, err := strconv.Atoi(item)
							if err != nil {
								panic(err)
							}
							o.applyRelationBoBySelectField(sessionId, columnItem, value, record, relationBo)
						}
					}
				}
			}
		} else {
			if (*loopItem)[columnItem.Name] == nil {
				virtualColumn := map[string]interface{}{}
				buttons := []interface{}{}
				virtualColumn["buttons"] = buttons
				(*loopItem)[columnItem.Name] = virtualColumn
			}

			record["pendingTransactions"] = []interface{}{}
			recordJsonByte, err := json.Marshal(record)
			if err != nil {
				panic(err)
			}
			recordJson := string(recordJsonByte)

			expressionParser := ExpressionParser{}
			for _, buttonItem := range columnItem.Buttons.ButtonLi {
				button := map[string]interface{}{}
				button["isShow"] = expressionParser.Parse(recordJson, buttonItem.Expression)
				virtualColumn := (*loopItem)[columnItem.Name].(map[string]interface{})
				buttons := virtualColumn["buttons"].([]interface{})
				buttons = append(buttons, button)
				virtualColumn["buttons"] = buttons // append will generate a new reference, so must reset value
			}
		}
	}
}

func (o TemplateManager) applyRelationBoBySelectField(sessionId int, column Column, value int, data map[string]interface{}, relationBo *map[string]interface{}) {
	if column.CRelationDS.CRelationItemLi != nil {
		for _, relationItem := range column.CRelationDS.CRelationItemLi {
			mode := relationItem.CRelationExpr.Mode
			exprContent := relationItem.CRelationExpr.Content
			bo := data
			data4Expression := o.getData4Expression(column, data)
			content := o.parseModelExpression(bo, data4Expression, mode, exprContent)
			if strings.ToLower(content) == "true" {
				selectorName := relationItem.CRelationConfig.SelectorName
				if (*relationBo)[selectorName] != nil {
					selectorDict := (*relationBo)[selectorName].(map[string]interface{})
					if selectorDict[fmt.Sprint(value)] != nil {
						continue
					}
				}
				li := []map[string]interface{}{}
				li = append(li, map[string]interface{}{
					"relationId": value,
					"selectorId": selectorName,
				})
				singleRelationBo := o.GetRelationBo(sessionId, li)
				if singleRelationBo[selectorName] != nil {
					singleRelationItem := singleRelationBo[selectorName].(map[string]interface{})
					if (*relationBo)[selectorName] == nil {
						(*relationBo)[selectorName] = map[string]interface{}{}
					}
					relationItem := (*relationBo)[selectorName].(map[string]interface{})
					for k, v := range singleRelationItem {
						relationItem[k] = v
					}
					(*relationBo)[selectorName] = relationItem
				}
			}
		}
	}
}

func (o TemplateManager) getValueBySpot(record map[string]interface{}, name string) interface{} {
	current := record
	nameLi := strings.Split(name, ".")
	for i, _ := range nameLi {
		if i < len(nameLi)-1 {
			if current[nameLi[i]] == nil {
				return nil
			}
			if reflect.ValueOf(current[nameLi[i]]).Kind() == reflect.Map {
				current = current[nameLi[i]].(map[string]interface{})
			} else {
				return nil
			}
		} else {
			return current[nameLi[i]]
		}
	}
	return nil
}

//func (o TemplateManager) ApplyDictionaryColumnData(loopItem *map[string]interface{}, columnItem Column) {
//	dictionaryManager := dictionary.GetInstance()
//	if columnItem.XMLName.Local == "dictionary-column" {
//		dictionaryItem := dictionaryManager.GetDictionary(columnItem.Dictionary)
//		items := dictionaryItem["items"]
//		if items != nil {
//			itemsLi := items.([]map[string]interface{})
//			columnValue := fmt.Sprint((*loopItem)[columnItem.Name])
//			for _, codeNameItem := range itemsLi {
//				code := fmt.Sprint(codeNameItem["code"])
//				if code == columnValue {
//					(*loopItem)[columnItem.Name+"_DICTIONARY_NAME"] = codeNameItem["name"]
//					break
//				}
//			}
//		}
//	}
//}

func (o TemplateManager) ApplyScriptColumnData(loopItem *map[string]interface{}, record map[string]interface{}, columnItem Column) {
	if columnItem.XMLName.Local == "script-column" {
		data, err := json.Marshal(record)
		if err != nil {
			panic(err)
		}

		expressionParser := ExpressionParser{}
		scriptValue := expressionParser.ParseString(string(data), columnItem.Script)
		(*loopItem)[columnItem.Name] = scriptValue
	}
}

func (o TemplateManager) GetToolbarForListTemplate(listTemplate ListTemplate) []interface{} {
	return o.GetToolbarBo(listTemplate.Toolbar)
}

func (o TemplateManager) GetToolbarBo(toolbar Toolbar) []interface{} {
	result := []interface{}{}

	expressionParser := ExpressionParser{}
	for _, buttonItem := range toolbar.ButtonLi {
		button := map[string]interface{}{}
		button["isShow"] = expressionParser.Parse("", buttonItem.Expression)
		result = append(result, button)
	}

	return result
}

/**
 * 获取模版业务对象
 */
func (o TemplateManager) GetBoForListTemplate(sessionId int, listTemplate *ListTemplate, paramMap map[string]string, pageNo int, pageSize int) map[string]interface{} {
	queryResult := o.QueryDataForListTemplate(sessionId, listTemplate, paramMap, pageNo, pageSize)
	items := queryResult["items"].([]interface{})
	itemsDict := o.GetColumnModelDataForListTemplate(sessionId, *listTemplate, items)
	bo := itemsDict["items"].([]interface{})
	relationBo := itemsDict["relationBo"].(map[string]interface{})

	selectorInfo := o.GetSelectorInfoForListTemplate(*listTemplate)
	o.mergeRelationBo(&relationBo, selectorInfo)

	return map[string]interface{}{
		"totalResults": queryResult["totalResults"],
		"items":        bo,
		"relationBo":   relationBo,
	}
}

func (o TemplateManager) GetColumns(listTemplate *ListTemplate) []string {
	fields := []string{}
	//	loopItem["isShowCheckbox"] = expressionParser.Parse(recordJson, listTemplate.ColumnModel.CheckboxColumn.Expression)
	//		loopItem["id"] = record[listTemplate.ColumnModel.IdColumn.Name]
	//		for _, columnItem := range listTemplate.ColumnModel.ColumnLi {
	fields = append(fields, listTemplate.ColumnModel.IdColumn.Name)
	for _, columnItem := range listTemplate.ColumnModel.ColumnLi {
		//		fields = append(fields, columnItem.Name)
		fields = append(fields, columnItem.Name)
	}
	return fields
}

func (o TemplateManager) GetShowParameterLiForListTemplate(listTemplate *ListTemplate) []QueryParameter {
	queryParameterLi := []QueryParameter{}
	for _, item := range listTemplate.QueryParameterGroup.QueryParameterLi {
		if item.Editor != "hiddenfield" {
			queryParameterLi = append(queryParameterLi, item)
		}
	}
	return queryParameterLi
}

func (o TemplateManager) GetHiddenParameterLiForListTemplate(listTemplate *ListTemplate) []QueryParameter {
	queryParameterLi := []QueryParameter{}
	for _, item := range listTemplate.QueryParameterGroup.QueryParameterLi {
		if item.Editor == "hiddenfield" {
			queryParameterLi = append(queryParameterLi, item)
		}
	}
	return queryParameterLi
}

//func (o TemplateManager) ApplyDictionaryForQueryParameter(listTemplate *ListTemplate) {
//	mongoDBFactory := mongo.GetInstance()
//	session, db := mongoDBFactory.GetConnection()
//	defer session.Close()
//
//	dictionaryManager := dictionary.GetInstance()
//	for i, _ := range listTemplate.QueryParameterGroup.QueryParameterLi {
//		item := &(listTemplate.QueryParameterGroup.QueryParameterLi[i])
//		for _, parameterAttribute := range item.ParameterAttributeLi {
//			if parameterAttribute.Name == "dictionary" {
//				item.Dictionary = dictionaryManager.GetDictionaryBySession(db, parameterAttribute.Value)
//				break
//			}
//		}
//	}
//}

//func (o TemplateManager) ApplyTreeForQueryParameter(listTemplate *ListTemplate) {
//	mongoDBFactory := mongo.GetInstance()
//	session, db := mongoDBFactory.GetConnection()
//	defer session.Close()
//
//	treeManager := tree.GetInstance()
//	for i, _ := range listTemplate.QueryParameterGroup.QueryParameterLi {
//		item := &(listTemplate.QueryParameterGroup.QueryParameterLi[i])
//		for _, parameterAttribute := range item.ParameterAttributeLi {
//			if parameterAttribute.Name == "tree" {
//				item.Tree = treeManager.GetTreeBySession(db, parameterAttribute.Value)
//				break
//			}
//		}
//	}
//}

func (o TemplateManager) GetLayerForFormTemplate(sId int, formTemplate FormTemplate) map[string]interface{} {
	_, db := global.GetConnection(sId)

	result := map[string]interface{}{}
	resultLi := map[string]interface{}{}
	layerManager := layer.GetInstance()
	for _, item := range formTemplate.FormElemLi {
		if item.XMLName.Local == "column-model" {
			for _, column := range item.ColumnModel.ColumnLi {
				if column.Dictionary != "" {
					layerMap := layerManager.GetLayerBySession(sId, db, column.Dictionary)
					if layerMap != nil {
						items := layerMap["items"]
						if items != nil {
							dictMap := map[string]interface{}{}
							for _, item := range items.([]map[string]interface{}) {
								dictMap[fmt.Sprint(item["code"])] = item
							}
							result[column.Dictionary] = dictMap
							resultLi[column.Dictionary] = items
						} else {
							result[column.Dictionary] = map[string]interface{}{}
							resultLi[column.Dictionary] = []interface{}{}
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"layerBo":   result,
		"layerBoLi": resultLi,
	}
}

// TODO
func (o TemplateManager) GetLayerForListTemplate(sId int, listTemplate ListTemplate) map[string]interface{} {
	_, db := global.GetConnection(sId)

	result := map[string]interface{}{}
	resultLi := map[string]interface{}{}
	layerManager := layer.GetInstance()

	listTemplateIterator := ListTemplateIterator{}
	var iterateResult interface{} = ""
	listTemplateIterator.IterateTemplateColumn(listTemplate, &iterateResult, func(column Column, iterateResult *interface{}) {
		if column.Dictionary != "" {
			layerMap := layerManager.GetLayerBySession(sId, db, column.Dictionary)
			if layerMap != nil {
				items := layerMap["items"]
				if items != nil {
					dictMap := map[string]interface{}{}
					for _, item := range items.([]map[string]interface{}) {
						dictMap[fmt.Sprint(item["code"])] = item
					}
					result[column.Dictionary] = dictMap
					resultLi[column.Dictionary] = items
				} else {
					result[column.Dictionary] = map[string]interface{}{}
					resultLi[column.Dictionary] = []interface{}{}
				}
			}
		}
	})

	return map[string]interface{}{
		"layerBo":   result,
		"layerBoLi": resultLi,
	}
}

func (o TemplateManager) GetRelationBo(sId int, relationLi []map[string]interface{}) map[string]interface{} {
	_, db := global.GetConnection(sId)
	result := map[string]interface{}{}
	for _, item := range relationLi {
		relationId, err := strconv.Atoi(fmt.Sprint(item["relationId"]))
		if err != nil {
			panic(err)
		}
		if relationId == 0 {
			continue
		}
		selectorId := fmt.Sprint(item["selectorId"])
		listTemplate := o.GetSelectorTemplate(selectorId)
		collection := listTemplate.DataProvider.Collection
		element := map[string]interface{}{}
		queryMap := map[string]interface{}{
			"_id": relationId,
		}
		queryByte, err := json.MarshalIndent(queryMap, "", "\t")
		if err != nil {
			panic(err)
		}
		log.Println("GetRelationBo,collection:" + collection + ",query is:" + string(queryByte))
		err = db.C(collection).Find(queryMap).One(&element)
		if err != nil {
			if err == mgo.ErrNotFound {
				continue
			}
			panic(err)
		}
		items := []interface{}{element}
		interceptorManager := InterceptorManager{}
		items = interceptorManager.ParseAfterQueryData(sId, listTemplate.AfterQueryData, listTemplate.ColumnModel.DataSetId, items)
		if len(items) > 0 {
			itemsDict := o.GetColumnModelDataForListTemplate(sId, listTemplate, items)
			items = itemsDict["items"].([]interface{})
			element = items[0].(map[string]interface{})
		} else {
			continue
		}
		if result[selectorId] == nil {
			result[selectorId] = map[string]interface{}{}
		}
		selectorDict := result[selectorId].(map[string]interface{})
		selectorDict[fmt.Sprint(relationId)] = element
		if selectorDict["url"] == nil {
			selectorDict["url"] = o.GetViewUrl(listTemplate)
		}
		if selectorDict["Description"] == nil {
			selectorDict["Description"] = listTemplate.Description
		}
		result[selectorId] = selectorDict
	}
	return result
}

func (o TemplateManager) mergeSelectorInfo2RelationBoFromCRelationDS(cRelationDS CRelationDS, relationBo *map[string]interface{}) {
	if cRelationDS.CRelationItemLi != nil {
		for _, item := range cRelationDS.CRelationItemLi {
			selectorName := item.CRelationConfig.SelectorName
			if (*relationBo)[selectorName] != nil {
				selectorDict := (*relationBo)[selectorName].(map[string]interface{})
				if selectorDict["Description"] == nil {
					selectorTemplate := o.GetSelectorTemplate(selectorName)
					selectorDict["Description"] = selectorTemplate.Description
				}
				if selectorDict["url"] == nil {
					selectorTemplate := o.GetSelectorTemplate(selectorName)
					selectorDict["url"] = o.GetViewUrl(selectorTemplate)
				}
				(*relationBo)[selectorName] = selectorDict
			} else {
				selectorDict := map[string]interface{}{}
				selectorTemplate := o.GetSelectorTemplate(selectorName)
				selectorDict["Description"] = selectorTemplate.Description
				selectorDict["url"] = o.GetViewUrl(selectorTemplate)
				(*relationBo)[selectorName] = selectorDict
			}
		}
	}
}

func (o TemplateManager) mergeSelectorInfo2RelationBo(formTemplate FormTemplate, relationBo *map[string]interface{}) {
	formTemplateIterator := FormTemplateIterator{}
	var result interface{} = ""
	formTemplateIterator.IterateTemplateColumn(formTemplate, &result, func(column Column, result *interface{}) {
		o.mergeSelectorInfo2RelationBoFromCRelationDS(column.CRelationDS, relationBo)
	})
	formTemplateIterator.IterateTemplateButton(formTemplate, &result, func(toolbar Toolbar, column ColumnModel, isToolbarBtn bool, button Button, result *interface{}) {
		o.mergeSelectorInfo2RelationBoFromCRelationDS(button.CRelationDS, relationBo)
	})
}

func (o TemplateManager) GetViewUrl(listTemplate ListTemplate) string {
	for _, item := range listTemplate.ColumnModel.ColumnLi {
		if item.Buttons.ButtonLi != nil {
			for _, buttonItem := range item.Buttons.ButtonLi {
				if buttonItem.Name == "btn_view" {
					return fmt.Sprint(buttonItem.Handler)
				}
			}
		}
	}
	return ""
}

func (o TemplateManager) GetQueryDefaultValue(listTemplate ListTemplate) map[string]string {
	defaultDict := map[string]string{}
	listTemplateIterator := ListTemplateIterator{}
	var result interface{} = ""
	listTemplateIterator.IterateTemplateQueryParameter(listTemplate, &result, func(queryParameter QueryParameter, result *interface{}) {
		mode := queryParameter.CDefaultValueExpr.Mode
		content := queryParameter.CDefaultValueExpr.Content
		value := o.parseQueryParameterExpression(mode, content)
		if value != "" {
			defaultDict[queryParameter.Name] = value
		}
	})
	return defaultDict
}

func (o TemplateManager) parseQueryParameterExpression(mode string, content string) string {
	expressionParser := ExpressionParser{}
	if mode == "" || mode == "text" {
		return content
	} else if mode == "python" {
		return expressionParser.ParseString("{}", content)
	} else if mode == "golang" {
		bo := map[string]interface{}{}
		data := map[string]interface{}{}
		return expressionParser.ParseGolang(bo, data, content)
	}
	return ""
}
