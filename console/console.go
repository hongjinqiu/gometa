package console

import (
	"encoding/json"
	"encoding/xml"
	//	"fmt"
	"bufio"
	"bytes"
	"fmt"
	. "github.com/hongjinqiu/gometa/common"
	. "github.com/hongjinqiu/gometa/component"
	"github.com/hongjinqiu/gometa/config"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	. "github.com/hongjinqiu/gometa/model/handler"
	"github.com/hongjinqiu/gometa/session"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Console struct {
}

//func (c Console) Summary(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("^^^^^^ Summary")
//}
//
//func (c Console) summary(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("______ Summary")
//}

func (self Console) Summary(w http.ResponseWriter, r *http.Request) {
	sessionId := global.GetSessionId()
	global.SetGlobalAttr(sessionId, "userId", session.GetFromSession(w, r, "userId"))
	global.SetGlobalAttr(sessionId, "adminUserId", session.GetFromSession(w, r, "adminUserId"))
	defer global.CloseSession(sessionId)

	templateManager := TemplateManager{}
	formTemplate := templateManager.GetFormTemplate("Console")

	//	if true {
	//		xmlDataArray, err := xml.Marshal(&formTemplate)
	//		if err != nil {
	//			panic(err)
	//		}
	//		return self.RenderXml(&formTemplate)
	//	}

	formTemplateJsonDataArray, err := json.Marshal(&formTemplate)
	if err != nil {
		panic(err)
	}

	toolbarBo := map[string]interface{}{}

	dataBo := map[string]interface{}{}
	{
		listTemplateInfoLi := templateManager.GetListTemplateInfoLi()
		dataBo["Component"] = getSummaryListTemplateInfoLi(listTemplateInfoLi)
	}
	{
		selectorTemplateInfoLi := templateManager.GetSelectorTemplateInfoLi()
		dataBo["Selector"] = getSummarySelectorTemplateInfoLi(selectorTemplateInfoLi)
	}
	{
		formTemplateInfoLi := templateManager.GetFormTemplateInfoLi()
		dataBo["Form"] = getSummaryFormTemplateInfoLi(formTemplateInfoLi)
	}
	{
		modelTemplateFactory := ModelTemplateFactory{}
		dataSourceInfoLi := modelTemplateFactory.GetDataSourceInfoLi()
		dataBo["DataSource"] = getSummaryDataSourceInfoLi(dataSourceInfoLi)
	}
	for _, item := range formTemplate.FormElemLi {
		if item.XMLName.Local == "column-model" {
			if dataBo[item.ColumnModel.Name] == nil {
				dataBo[item.ColumnModel.Name] = []interface{}{}
			}
			items := dataBo[item.ColumnModel.Name].([]interface{})
			itemsDict := templateManager.GetColumnModelDataForColumnModel(sessionId, item.ColumnModel, items)
			items = itemsDict["items"].([]interface{})
			dataBo[item.ColumnModel.Name] = items
		} else if item.XMLName.Local == "toolbar" {
			toolbarBo[item.Toolbar.Name] = templateManager.GetToolbarBo(item.Toolbar)
		}
	}

	dataBoByte, err := json.Marshal(dataBo)
	if err != nil {
		panic(err)
	}

	//	self.Response.Status = http.StatusOK
	//	self.Response.ContentType = "text/plain; charset=utf-8"
	result := map[string]interface{}{
		"formTemplate":         formTemplate,
		"toolbarBo":            toolbarBo,
		"dataBo":               dataBo,
		"formTemplateJsonData": template.JS(string(formTemplateJsonDataArray)),
		"dataBoJson":           template.JS(string(dataBoByte)),
	}
	// formTemplate.ViewTemplate.View
	//	RenderText(text string, objs ...interface{}) Result

	viewPath := config.String("VIEW_PATH")
	file, err := os.Open(viewPath + "/" + formTemplate.ViewTemplate.View)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	//	self.Response.Out
	//	return self.RenderTemplate(string(fileContent))
	//	funcMap := map[string]interface{}{
	//		"eq": func(a, b interface{}) bool {
	//			return a == b
	//		},
	//	}
	funcMap := map[string]interface{}{}
	//self.Response.ContentType = "text/html; charset=utf-8"
	tmpl, err := template.New("summary").Funcs(funcMap).Parse(string(fileContent))
	if err != nil {
		panic(err)
	}
	tmplResult := map[string]interface{}{
		"result": result,
	}
	//tmpl.Execute(self.Response.Out, result)
	err = tmpl.Execute(w, tmplResult)
	if err != nil {
		panic(err)
	}
	//	return self.Render(string(fileContent), result)
}

func getSummaryListTemplateInfoLi(listTemplateInfoLi []ListTemplateInfo) []interface{} {
	componentItems := []interface{}{}
	for _, item := range listTemplateInfoLi {
		module := "组件模型"
		if item.ListTemplate.DataSourceModelId != "" && item.ListTemplate.Adapter.Name != "" {
			module = "数据源模型适配"
		}
		componentItems = append(componentItems, map[string]interface{}{
			"id":     item.ListTemplate.Id,
			"name":   item.ListTemplate.Description,
			"module": module,
			"path":   item.Path,
		})
	}
	return componentItems
}

func getSummarySelectorTemplateInfoLi(selectorTemplateInfoLi []SelectorTemplateInfo) []interface{} {
	componentItems := []interface{}{}
	for _, item := range selectorTemplateInfoLi {
		module := "组件模型选择器"
		if item.ListTemplate.DataSourceModelId != "" && item.ListTemplate.Adapter.Name != "" {
			module = "数据源模型选择器适配"
		}
		id := item.ListTemplate.SelectorId
		if id == "" {
			id = item.ListTemplate.Id
		}
		componentItems = append(componentItems, map[string]interface{}{
			"id":     id,
			"name":   item.ListTemplate.Description,
			"module": module,
			"path":   item.Path,
		})
	}
	return componentItems
}

func getSummaryFormTemplateInfoLi(formTemplateInfoLi []FormTemplateInfo) []interface{} {
	formItems := []interface{}{}
	for _, item := range formTemplateInfoLi {
		module := "form模型"
		if item.FormTemplate.DataSourceModelId != "" && item.FormTemplate.Adapter.Name != "" {
			module = "数据源模型适配"
		}
		formItems = append(formItems, map[string]interface{}{
			"id":     item.FormTemplate.Id,
			"name":   item.FormTemplate.Description,
			"module": module,
			"path":   item.Path,
		})
	}
	return formItems
}

func getSummaryDataSourceInfoLi(dataSourceInfoLi []DataSourceInfo) []interface{} {
	dataSourceItems := []interface{}{}
	for _, item := range dataSourceInfoLi {
		dataSourceItems = append(dataSourceItems, map[string]interface{}{
			"id":     item.DataSource.Id,
			"name":   item.DataSource.DisplayName,
			"module": "数据源模型",
			"path":   item.Path,
		})
	}
	return dataSourceItems
}

func (self Console) Xml(w http.ResponseWriter, r *http.Request) {
	//	refretorType := self.Params.Get("type")
	//	id := self.Params.Get("@name")
	refretorType := r.URL.Query().Get("type")
	id := r.URL.Query().Get("@name")
	templateManager := TemplateManager{}

	if refretorType == "Component" {
		listTemplate := templateManager.GetListTemplate(id)
		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		data, err := xml.MarshalIndent(&listTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	if refretorType == "Selector" {
		selectorTemplate := templateManager.GetSelectorTemplate(id)
		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		data, err := xml.MarshalIndent(&selectorTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	if refretorType == "Form" {
		formTemplate := templateManager.GetFormTemplate(id)
		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		data, err := xml.MarshalIndent(&formTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	if refretorType == "DataSource" {
		modelTemplateFactory := ModelTemplateFactory{}
		dataSourceTemplate := modelTemplateFactory.GetDataSource(id)
		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		data, err := xml.MarshalIndent(&dataSourceTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	data, err := json.MarshalIndent(map[string]interface{}{
		"message": "可能传入了错误的refretorType:" + refretorType,
	}, "", "\t")
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func (c Console) RawXml(w http.ResponseWriter, r *http.Request) {
	refretorType := r.URL.Query().Get("type")
	id := r.URL.Query().Get("@name")
	templateManager := TemplateManager{}

	if refretorType == "Component" {
		listTemplateInfo := templateManager.GetListTemplateInfo(id)
		listTemplate := ListTemplate{}
		file, err := os.Open(listTemplateInfo.Path)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		err = xml.Unmarshal(data, &listTemplate)
		if err != nil {
			panic(err)
		}

		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		xmlData, err := xml.MarshalIndent(&listTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(xmlData)
		return
	}
	if refretorType == "Selector" {
		selectorTemplateInfo := templateManager.GetSelectorTemplateInfo(id)
		selectorTemplate := ListTemplate{}
		file, err := os.Open(selectorTemplateInfo.Path)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		err = xml.Unmarshal(data, &selectorTemplate)
		if err != nil {
			panic(err)
		}

		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		xmlData, err := xml.MarshalIndent(&selectorTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(xmlData)
		return
	}
	if refretorType == "Form" {
		formTemplateInfo := templateManager.GetFormTemplateInfo(id)
		formTemplate := FormTemplate{}
		file, err := os.Open(formTemplateInfo.Path)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		err = xml.Unmarshal(data, &formTemplate)
		if err != nil {
			panic(err)
		}

		w.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		xmlData, err := xml.MarshalIndent(&formTemplate, "", "\t")
		if err != nil {
			panic(err)
		}
		w.Write(xmlData)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	data, err := json.MarshalIndent(map[string]interface{}{
		"message": "可能传入了错误的refretorType:" + refretorType,
	}, "", "\t")
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func (self Console) ListSchema(w http.ResponseWriter, r *http.Request) {
	schemaName := r.URL.Query().Get("@name")

	templateManager := TemplateManager{}
	listTemplate := templateManager.GetListTemplate(schemaName)

	isFromList := true
	result := self.listSelectorCommon(w, r, &listTemplate, true, isFromList)

	format := r.URL.Query().Get("format")
	if strings.ToLower(format) == "json" {
		callback := r.URL.Query().Get("callback")
		if callback == "" {
			dataBo := result["dataBo"]
			w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
			data, err := json.Marshal(&dataBo)
			if err != nil {
				panic(err)
			}
			w.Write(data)
			return
			//			c.Response.ContentType = "application/json; charset=utf-8"
			//			return c.RenderJson(&dataBo)
		}
		dataBoText := result["dataBoText"].(string)
		w.Header()["Content-Type"] = []string{"text/javascript; charset=utf-8"}
		w.Write([]byte(dataBoText))
		return
		//		c.Response.ContentType = "text/javascript; charset=utf-8"
		//		return c.RenderText(callback + "(" + dataBoText + ");")
	} else {
		//c.Response.ContentType = "text/html; charset=utf-8"
		/*
			{
				c.RenderArgs["result"] = result
				c.RenderArgs["flash"] = c.Flash.Out
				c.RenderArgs["session"] = c.Session
				return c.RenderTemplate(listTemplate.ViewTemplate.View)
			}
		*/
		tmplResult := map[string]interface{}{
			"result": result,
		}
		result["ListPageContent"] = template.HTML(self.getListPageContent(tmplResult))
		result["ListQueryParameterContent"] = template.HTML(self.getListQueryParameterContent(tmplResult))

		viewPath := config.String("VIEW_PATH")
		file, err := os.Open(viewPath + "/" + listTemplate.ViewTemplate.View)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		strContent := string(fileContent)
		tmpl, err := template.New("ListSchema").Funcs(self.getFuncMap()).Parse(strContent)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, tmplResult)
		if err != nil {
			panic(err)
		}
	}
}

func (self Console) getFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"residue": func(a int, b int, c int) bool {
			return a%b == c
		},
		"last": func(a int, b int) bool {
			return a-1 == b
		},
	}
}

func (self Console) getListPageContent(result map[string]interface{}) string {
	viewPath := config.String("VIEW_PATH")
	file, err := os.Open(viewPath + "/Console/ListPage.html")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	strContent := string(fileContent)
	tmpl, err := template.New("ListPage").Funcs(self.getFuncMap()).Parse(strContent)
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)
	tmpl.Execute(bw, result)
	bw.Flush()
	return b.String()
}

func (self Console) getListQueryParameterContent(result map[string]interface{}) string {
	viewPath := config.String("VIEW_PATH")
	file, err := os.Open(viewPath + "/Console/ListQueryParameter.html")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	strContent := string(fileContent)
	tmpl, err := template.New("ListQueryParameter").Funcs(self.getFuncMap()).Parse(strContent)
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)
	tmpl.Execute(bw, result)
	bw.Flush()
	return b.String()
}

func (self Console) listSelectorCommon(w http.ResponseWriter, r *http.Request, listTemplate *ListTemplate, isGetBo bool, isFromList bool) map[string]interface{} {
	sessionId := global.GetSessionId()
	global.SetGlobalAttr(sessionId, "userId", session.GetFromSession(w, r, "userId"))
	global.SetGlobalAttr(sessionId, "adminUserId", session.GetFromSession(w, r, "adminUserId"))
	defer global.CloseSession(sessionId)

	// 1.toolbar bo
	templateManager := TemplateManager{}
	//templateManager.ApplyDictionaryForQueryParameter(listTemplate)
	//templateManager.ApplyTreeForQueryParameter(listTemplate)
	toolbarBo := templateManager.GetToolbarForListTemplate(*listTemplate)
	paramMap := map[string]string{}

	defaultBo := templateManager.GetQueryDefaultValue(*listTemplate)
	defaultBoByte, err := json.Marshal(&defaultBo)
	if err != nil {
		panic(err)
	}
	for key, value := range defaultBo {
		paramMap[key] = value
	}
	paramMap, _ = self.getCookieDataAndParamMap(sessionId, w, r, *listTemplate, isFromList, paramMap)

	formDataByte, err := json.Marshal(&paramMap)
	if err != nil {
		panic(err)
	}

	//	}
	pageNo := 1
	pageSize := 10
	if listTemplate.DataProvider.Size != "" {
		pageSizeInt, err := strconv.Atoi(listTemplate.DataProvider.Size)
		if err != nil {
			panic(err)
		}
		pageSize = pageSizeInt
	}
	if r.FormValue("pageNo") != "" {
		pageNoInt, _ := strconv.ParseInt(r.FormValue("pageNo"), 10, 0)
		if pageNoInt > 1 {
			pageNo = int(pageNoInt)
		}
	}
	if r.FormValue("pageSize") != "" {
		pageSizeInt, _ := strconv.ParseInt(r.FormValue("pageSize"), 10, 0)
		if pageSizeInt >= 10 {
			pageSize = int(pageSizeInt)
		}
	}
	dataBo := map[string]interface{}{
		"totalResults": 0,
		"items":        []interface{}{},
	}
	relationBo := map[string]interface{}{}
	usedCheckBo := map[string]interface{}{}
	//if self.Params.Get("@entrance") != "true" {
	if isGetBo {
		dataBo = templateManager.GetBoForListTemplate(sessionId, listTemplate, paramMap, pageNo, pageSize)
		relationBo = dataBo["relationBo"].(map[string]interface{})
		//delete(dataBo, "relationBo")

		// usedCheck的修改,
		if listTemplate.DataSourceModelId != "" {
			modelTemplateFactory := ModelTemplateFactory{}
			dataSource := modelTemplateFactory.GetDataSource(listTemplate.DataSourceModelId)
			items := dataBo["items"].([]interface{})
			usedCheck := UsedCheck{}

			usedCheckBo = usedCheck.GetListUsedCheck(sessionId, dataSource, items, listTemplate.ColumnModel.DataSetId)
		}
	}
	dataBo["usedCheckBo"] = usedCheckBo

	dataBoByte, err := json.Marshal(&dataBo)
	if err != nil {
		panic(err)
	}

	relationBoByte, err := json.Marshal(&relationBo)
	if err != nil {
		panic(err)
	}

	listTemplateByte, err := json.Marshal(listTemplate)
	if err != nil {
		panic(err)
	}

	usedCheckByte, err := json.Marshal(&usedCheckBo)
	if err != nil {
		panic(err)
	}

	// 系统参数的获取
	commonUtil := CommonUtil{}
	userId := commonUtil.GetIntFromString(session.GetFromSession(w, r, "userId"))
	sysParam := self.getSysParam(sessionId, userId)
	sysParamJson, err := json.Marshal(&sysParam)
	if err != nil {
		panic(err)
	}

	queryParameterRenderLi := self.getQueryParameterRenderLi(*listTemplate)

	//showParameterLi := templateManager.GetShowParameterLiForListTemplate(listTemplate)
	showParameterLi := []QueryParameter{}
	hiddenParameterLi := templateManager.GetHiddenParameterLiForListTemplate(listTemplate)

	layerBo := templateManager.GetLayerForListTemplate(sessionId, *listTemplate)
	iLayerBo := layerBo["layerBo"]
	layerBoByte, err := json.Marshal(&iLayerBo)
	if err != nil {
		panic(err)
	}
	iLayerBoLi := layerBo["layerBoLi"]
	layerBoLiByte, err := json.Marshal(&iLayerBoLi)
	if err != nil {
		panic(err)
	}
	layerBoJson := string(layerBoByte)
	layerBoJson = commonUtil.FilterJsonEmptyAttr(layerBoJson)
	layerBoLiJson := string(layerBoLiByte)
	layerBoLiJson = commonUtil.FilterJsonEmptyAttr(layerBoLiJson)

	result := map[string]interface{}{
		"pageSize":               pageSize,
		"listTemplate":           listTemplate,
		"toolbarBo":              toolbarBo,
		"showParameterLi":        showParameterLi,
		"hiddenParameterLi":      hiddenParameterLi,
		"queryParameterRenderLi": queryParameterRenderLi,
		"dataBo":                 dataBo,
		//		"columns":       columns,
		"dataBoText":       string(dataBoByte),
		"dataBoJson":       template.JS(string(dataBoByte)),
		"relationBoJson":   template.JS(string(relationBoByte)),
		"listTemplateJson": template.JS(string(listTemplateByte)),
		"layerBoJson":      template.JS(layerBoJson),
		"layerBoLiJson":    template.JS(layerBoLiJson),
		"defaultBoJson":    template.JS(string(defaultBoByte)),
		"formDataJson":     template.JS(string(formDataByte)),
		"usedCheckJson":    template.JS(string(usedCheckByte)),
		"sysParamJson":     template.JS(string(sysParamJson)),
		//		"columnsJson":   string(columnsByte),
	}
	return result
}

func (c Console) getSysParam(sessionId int, userId int) map[string]interface{} {
	commonUtil := CommonUtil{}
	systemParameter := c.getSystemParameter(sessionId, userId)
	systemParameterMain := systemParameter["A"].(map[string]interface{})
	currencyTypeId := commonUtil.GetIntFromMap(systemParameterMain, "currencyTypeId")
	currencyType := c.getCurrencyType(sessionId, currencyTypeId)
	currencyTypeMain := currencyType["A"].(map[string]interface{})
	thousandsSeparator := ","
	if fmt.Sprint(systemParameterMain["thousandDecimals"]) == "1" {
		thousandsSeparator = ""
	}
	amtDecimals := commonUtil.GetIntFromMap(currencyTypeMain, "amtDecimals")
	upDecimals := commonUtil.GetIntFromMap(currencyTypeMain, "upDecimals")
	costDecimals := commonUtil.GetIntFromMap(systemParameterMain, "costDecimals")
	percentDecimals := commonUtil.GetIntFromMap(systemParameterMain, "percentDecimals")
	return map[string]interface{}{
		"localCurrency": map[string]interface{}{
			"prefix":                 currencyTypeMain["currencyTypeSign"],
			"decimalPlaces":          amtDecimals - 1,
			"unitPriceDecimalPlaces": upDecimals - 1,
		},
		"unitCostDecimalPlaces": costDecimals - 1,
		"percentDecimalPlaces":  percentDecimals - 1,
		"thousandsSeparator":    thousandsSeparator,
	}
}

func (self Console) getSystemParameter(sessionId int, userId int) map[string]interface{} {
	session, _ := global.GetConnection(sessionId)
	querySupport := QuerySupport{}
	user := querySupport.FindByMapWithSessionExact(session, "SysUser", map[string]interface{}{
		"_id": userId,
	})
	userMain := user["A"].(map[string]interface{})
	systemParameter := querySupport.FindByMapWithSessionExact(session, "SystemParameter", map[string]interface{}{
		"A.createUnit": userMain["createUnit"],
	})
	return systemParameter
}

func (self Console) getCurrencyType(sessionId int, currencyTypeId int) map[string]interface{} {
	session, _ := global.GetConnection(sessionId)
	querySupport := QuerySupport{}
	currencyType := querySupport.FindByMapWithSessionExact(session, "CurrencyType", map[string]interface{}{
		"_id": currencyTypeId,
	})
	return currencyType
}

func (self Console) getQueryParameterRenderLi(listTemplate ListTemplate) [][]map[string]interface{} {
	lineColSpan := 6
	container := [][]map[string]interface{}{}
	containerItem := []map[string]interface{}{}
	lineColSpanSum := 0
	listTemplateIterator := ListTemplateIterator{}
	var result interface{} = ""
	listTemplateIterator.IterateTemplateQueryParameter(listTemplate, &result, func(queryParameter QueryParameter, result *interface{}) {
		if queryParameter.Editor != "hiddenfield" {
			columnColSpan := 2
			containerItem = append(containerItem, map[string]interface{}{
				"label": queryParameter.Text,
				"name":  queryParameter.Name,
			})
			lineColSpanSum += columnColSpan
			if lineColSpanSum >= lineColSpan {
				container = append(container, containerItem)
				containerItem = []map[string]interface{}{}
				lineColSpanSum = lineColSpanSum - lineColSpan
			}
		}
	})
	if 0 < lineColSpanSum && lineColSpanSum < lineColSpan {
		container = append(container, containerItem)
	}
	return container
}

func (self Console) getCookieDataAndParamMap(sessionId int, w http.ResponseWriter, r *http.Request, listTemplate ListTemplate, isFromList bool, paramMap map[string]string) (map[string]string, map[string]string) {
	isHasCookie := false
	if r.FormValue("cookie") != "false" {
		isHasCookie = true
	}
	isConfigCookie := false
	if listTemplate.Cookie.Name != "" {
		isConfigCookie = true
	}
	cookieData := map[string]string{}
	if isFromList && isHasCookie && isConfigCookie {
		cookie, err := r.Cookie(listTemplate.Cookie.Name)
		if err != nil {
			if err != http.ErrNoCookie {
				panic(err)
			}
		} else {
			cookieStr := cookie.Value
			if cookieStr != "" {
				cookieStr = strings.Replace(cookieStr, "&quote", "\"", -1)
				err = json.Unmarshal([]byte(cookieStr), &cookieData)
				if err != nil {
					panic(err)
				}
				for k, v := range cookieData {
					paramMap[k] = v
				}
			}
		}
	}
	formQueryData := map[string]string{}
	//	self.Request.URL
	for k, v := range r.Form {
		value := strings.Join(v, ",")
		paramMap[k] = value
		formQueryData[k] = value
	}
	/*
		for k, v := range self.Params.Query {
			value := strings.Join(v, ",")
			paramMap[k] = value
			formQueryData[k] = value
		}
	*/

	if isFromList && isConfigCookie && !isHasCookie {
		cookie := http.Cookie{
			Name:  listTemplate.Cookie.Name,
			Value: "",
		}
		http.SetCookie(w, &cookie)
	} else if isFromList && isConfigCookie && isHasCookie {
		cookieFormQueryData := map[string]string{}
		for k, v := range cookieData {
			cookieFormQueryData[k] = v
		}
		for k, v := range formQueryData {
			cookieFormQueryData[k] = v
		}
		cookieStr, err := json.Marshal(&cookieFormQueryData)
		if err != nil {
			panic(err)
		}
		cookie := http.Cookie{
			Name:  listTemplate.Cookie.Name,
			Value: strings.Replace(string(cookieStr), "\"", "&quote", -1),
		}
		http.SetCookie(w, &cookie)
	}
	cookieData = map[string]string{}
	cookie, err := r.Cookie(listTemplate.Cookie.Name)
	if err != nil {
		if err != http.ErrNoCookie {
			panic(err)
		}
	} else {
		cookieStr := cookie.Value
		if cookieStr != "" {
			cookieStr = strings.Replace(cookieStr, "&quote", "\"", -1)
			err := json.Unmarshal([]byte(cookieStr), &cookieData)
			if err != nil {
				panic(err)
			}
		}
	}
	return paramMap, cookieData
}

func (self Console) Relation(w http.ResponseWriter, r *http.Request) {
	sessionId := global.GetSessionId()
	global.SetGlobalAttr(sessionId, "userId", session.GetFromSession(w, r, "userId"))
	global.SetGlobalAttr(sessionId, "adminUserId", session.GetFromSession(w, r, "adminUserId"))
	defer global.CloseSession(sessionId)

	selectorId := r.FormValue("selectorId")
	id := r.FormValue("id")

	templateManager := TemplateManager{}
	relationLi := []map[string]interface{}{
		map[string]interface{}{
			"selectorId": selectorId,
			"relationId": id,
		},
	}
	relationBo := templateManager.GetRelationBo(sessionId, relationLi)
	var result interface{} = nil
	var url interface{} = nil
	if relationBo[selectorId] != nil {
		selRelationBo := relationBo[selectorId].(map[string]interface{})
		if selRelationBo[id] != nil {
			result = selRelationBo[id]
			url = selRelationBo["url"]
		}
	}
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	data, err := json.Marshal(map[string]interface{}{
		"result": result,
		"url":    url,
	})
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
