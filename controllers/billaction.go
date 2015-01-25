package controllers

import (
	"encoding/json"
	"fmt"
	. "github.com/hongjinqiu/gometa/component"
	. "github.com/hongjinqiu/gometa/error"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	. "github.com/hongjinqiu/gometa/model/handler"
	. "github.com/hongjinqiu/gometa/mongo"
	"github.com/hongjinqiu/gometa/session"
	"net/http"
	"strconv"
	"strings"
)

func init() {
}

type BillAction struct {
	BaseDataAction
}

func (c BillAction) SaveData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.SaveCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) DeleteData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}

	modelRenderVO := c.DeleteDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) EditData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.EditDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) NewData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.RNewDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) GetData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.GetDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

/**
 * 复制
 */
func (c BillAction) CopyData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.CopyDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

/**
 * 放弃保存,回到浏览状态
 */
func (c BillAction) GiveUpData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.GiveUpDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

/**
 * 刷新
 */
func (c BillAction) RefreshData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.RefreshDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) LogList(w http.ResponseWriter, r *http.Request) {
	result := c.LogListCommon(w, r)

	format := r.FormValue("format")
	if strings.ToLower(format) == "json" {
		w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
		data, err := json.Marshal(&result)
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	//c.Response.ContentType = "text/html; charset=utf-8"
	return
}

/**
 * 作废
 */
func (c BillAction) CancelData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.CancelDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) CancelDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
	sessionId := global.GetSessionId()
	global.SetGlobalAttr(sessionId, "userId", session.GetFromSession(w, r, "userId"))
	global.SetGlobalAttr(sessionId, "adminUserId", session.GetFromSession(w, r, "adminUserId"))
	defer global.CloseSession(sessionId)
	defer c.RollbackTxn(sessionId)

	userId, err := strconv.Atoi(session.GetFromSession(w, r, "userId"))
	if err != nil {
		panic(err)
	}

	dataSourceModelId := r.FormValue("dataSourceModelId")
	formTemplateId := r.FormValue("formTemplateId")
	strId := r.FormValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		panic(err)
	}
	querySupport := QuerySupport{}
	queryMap := map[string]interface{}{
		"_id": id,
	}
	templateManager := TemplateManager{}
	formTemplate := templateManager.GetFormTemplate(formTemplateId)
	permissionSupport := PermissionSupport{}
	permissionQueryDict := permissionSupport.GetPermissionQueryDict(sessionId, formTemplate.Security)
	for k, v := range permissionQueryDict {
		queryMap[k] = v
	}

	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(dataSourceModelId)
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	bo, found := querySupport.FindByMap(collectionName, queryMap)
	if !found {
		panic("CancelData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	c.setRequestParameterToBo(r, &bo)

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	c.SetModifyFixFieldValue(sessionId, dataSource, &bo)
	c.RActionSupport.BeforeCancelData(sessionId, dataSource, formTemplate, &bo)
	mainData := bo["A"].(map[string]interface{})
	bo["A"] = mainData
	if fmt.Sprint(mainData["billStatus"]) == "4" {
		panic(BusinessError{Message: "单据已作废，不可再次作废"})
	}
	mainData["billStatus"] = 4

	_, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
	_, updateResult := txnManager.Update(txnId, dataSourceModelId, bo)
	if !updateResult {
		panic("作废失败")
	}

	c.RActionSupport.AfterCancelData(sessionId, dataSource, formTemplate, &bo)

	bo, _ = querySupport.FindByMap(collectionName, queryMap)

	usedCheck := UsedCheck{}
	usedCheckBo := usedCheck.GetFormUsedCheck(sessionId, dataSource, bo)

	columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
	bo = columnModelData["bo"].(map[string]interface{})
	relationBo := columnModelData["relationBo"].(map[string]interface{})

	modelTemplateFactory.ClearReverseRelation(&dataSource)
	c.CommitTxn(sessionId)
	return ModelRenderVO{
		UserId:       userId,
		Bo:           bo,
		RelationBo:   relationBo,
		UsedCheckBo:  usedCheckBo,
		DataSource:   dataSource,
		FormTemplate: formTemplate,
	}
}

/**
 * 反作废
 */
func (c BillAction) UnCancelData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.UnCancelDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BillAction) UnCancelDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
	sessionId := global.GetSessionId()
	global.SetGlobalAttr(sessionId, "userId", session.GetFromSession(w, r, "userId"))
	global.SetGlobalAttr(sessionId, "adminUserId", session.GetFromSession(w, r, "adminUserId"))
	defer global.CloseSession(sessionId)
	defer c.RollbackTxn(sessionId)

	userId, err := strconv.Atoi(session.GetFromSession(w, r, "userId"))
	if err != nil {
		panic(err)
	}

	dataSourceModelId := r.FormValue("dataSourceModelId")
	formTemplateId := r.FormValue("formTemplateId")
	strId := r.FormValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		panic(err)
	}
	querySupport := QuerySupport{}
	queryMap := map[string]interface{}{
		"_id": id,
	}
	templateManager := TemplateManager{}
	formTemplate := templateManager.GetFormTemplate(formTemplateId)
	permissionSupport := PermissionSupport{}
	permissionQueryDict := permissionSupport.GetPermissionQueryDict(sessionId, formTemplate.Security)
	for k, v := range permissionQueryDict {
		queryMap[k] = v
	}

	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(dataSourceModelId)
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	bo, found := querySupport.FindByMap(collectionName, queryMap)
	if !found {
		panic("UnCancelData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	c.setRequestParameterToBo(r, &bo)

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	c.SetModifyFixFieldValue(sessionId, dataSource, &bo)
	c.RActionSupport.BeforeUnCancelData(sessionId, dataSource, formTemplate, &bo)
	mainData := bo["A"].(map[string]interface{})
	if fmt.Sprint(mainData["billStatus"]) == "1" {
		panic(BusinessError{Message: "单据已经反作废，不可再次反作废"})
	}
	mainData["billStatus"] = 1

	_, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
	_, updateResult := txnManager.Update(txnId, dataSourceModelId, bo)
	if !updateResult {
		panic("反作废失败")
	}

	c.RActionSupport.AfterUnCancelData(sessionId, dataSource, formTemplate, &bo)

	bo, _ = querySupport.FindByMap(collectionName, queryMap)

	usedCheck := UsedCheck{}
	usedCheckBo := usedCheck.GetFormUsedCheck(sessionId, dataSource, bo)

	columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
	bo = columnModelData["bo"].(map[string]interface{})
	relationBo := columnModelData["relationBo"].(map[string]interface{})

	modelTemplateFactory.ClearReverseRelation(&dataSource)
	c.CommitTxn(sessionId)
	return ModelRenderVO{
		UserId:       userId,
		Bo:           bo,
		RelationBo:   relationBo,
		UsedCheckBo:  usedCheckBo,
		DataSource:   dataSource,
		FormTemplate: formTemplate,
	}
}
