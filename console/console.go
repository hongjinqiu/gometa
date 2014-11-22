package console

import (
	"encoding/json"
	"fmt"
	. "github.com/hongjinqiu/gometa/component"
	"github.com/hongjinqiu/gometa/config"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	"github.com/hongjinqiu/gometa/session"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
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

func (self Console) Gogogogo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "gogogogo")
	println("***** gogogogo")
}

func (self Console) Summary(w http.ResponseWriter, r *http.Request) {
	println("session is:", session.GetFromSession(w, r, "userId"))
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

	viewPath := config.String("REVEL_VIEW_PATH")
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
	tmpl.Execute(w, tmplResult)
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
