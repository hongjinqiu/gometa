package controllers

import (
	. "github.com/hongjinqiu/gometa/component"
	. "github.com/hongjinqiu/gometa/error"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/model"
	. "github.com/hongjinqiu/gometa/model/handler"
	. "github.com/hongjinqiu/gometa/mongo"
	"github.com/hongjinqiu/gometa/session"
	//	"github.com/hongjinqiu/gometa/mongo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
}

type IActionSupport interface {
	BeforeNewData(sessionId int, dataSource DataSource, formTemplate FormTemplate)
	AfterNewData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeGetData(sessionId int, dataSource DataSource, formTemplate FormTemplate)
	AfterGetData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeCopyData(sessionId int, dataSource DataSource, formTemplate FormTemplate, srcBo map[string]interface{})
	AfterCopyData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	EditValidate(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo map[string]interface{}) (string, bool)
	BeforeEditData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterEditData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeSaveData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterSaveData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}, diffDateRowLi *[]DiffDataRow)
	BeforeGiveUpData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterGiveUpData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeDeleteData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterDeleteData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeRefreshData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterRefreshData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	BeforeUnCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
	AfterUnCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{})
}

type ModelRenderVO struct {
	UserId       int
	Bo           map[string]interface{}
	RelationBo   map[string]interface{}
	UsedCheckBo  map[string]interface{}
	DataSource   DataSource
	FormTemplate FormTemplate
}

type ActionSupport struct {
}

func (o ActionSupport) BeforeNewData(sessionId int, dataSource DataSource, formTemplate FormTemplate) {
}
func (o ActionSupport) AfterNewData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeGetData(sessionId int, dataSource DataSource, formTemplate FormTemplate) {
}
func (o ActionSupport) AfterGetData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeCopyData(sessionId int, dataSource DataSource, formTemplate FormTemplate, srcBo map[string]interface{}) {
}
func (o ActionSupport) AfterCopyData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) EditValidate(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo map[string]interface{}) (string, bool) {
	return "", true
}
func (o ActionSupport) BeforeEditData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterEditData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeSaveData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterSaveData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}, diffDateRowLi *[]DiffDataRow) {
}
func (o ActionSupport) BeforeGiveUpData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterGiveUpData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeDeleteData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterDeleteData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeRefreshData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterRefreshData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) BeforeUnCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}
func (o ActionSupport) AfterUnCancelData(sessionId int, dataSource DataSource, formTemplate FormTemplate, bo *map[string]interface{}) {
}

type BaseDataAction struct {
	RActionSupport IActionSupport
}

func (c BaseDataAction) SetCreateFixFieldValue(sessionId int, dataSource DataSource, bo *map[string]interface{}) {
	var result interface{} = ""
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	createTime, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 64)
	if err != nil {
		panic(err)
	}
	_, db := global.GetConnection(sessionId)
	sysUser := map[string]interface{}{}
	query := map[string]interface{}{
		"_id": userId,
	}
	err = db.C("SysUser").Find(query).One(&sysUser)
	if err != nil {
		panic(err)
	}
	sysUserMaster := sysUser["A"].(map[string]interface{})
	modelIterator := ModelIterator{}
	modelIterator.IterateDataBo(dataSource, bo, &result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		(*data)["createBy"] = userId
		(*data)["createTime"] = createTime
		(*data)["createUnit"] = sysUserMaster["createUnit"]
	})
}

func (c BaseDataAction) SetModifyFixFieldValue(sessionId int, dataSource DataSource, bo *map[string]interface{}) {
	var result interface{} = ""
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	modifyTime, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 64)
	if err != nil {
		panic(err)
	}
	_, db := global.GetConnection(sessionId)
	sysUser := map[string]interface{}{}
	query := map[string]interface{}{
		"_id": userId,
	}
	err = db.C("SysUser").Find(query).One(&sysUser)
	if err != nil {
		panic(err)
	}
	sysUserMaster := sysUser["A"].(map[string]interface{})

	srcBo := map[string]interface{}{}
	srcQuery := map[string]interface{}{
		"_id": (*bo)["id"],
	}
	// log
	modelTemplateFactory := ModelTemplateFactory{}
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	srcQueryByte, err := json.Marshal(&srcQuery)
	if err != nil {
		panic(err)
	}
	log.Println("SetModifyFixFieldValue,collectionName:" + collectionName + ", query:" + string(srcQueryByte))
	db.C(collectionName).Find(srcQuery).One(&srcBo)
	modelIterator := ModelIterator{}
	modelIterator.IterateDiffBo(dataSource, bo, srcBo, &result, func(fieldGroupLi []FieldGroup, destData *map[string]interface{}, srcData map[string]interface{}, result *interface{}) {
		if destData != nil && srcData == nil {
			(*destData)["createBy"] = userId
			(*destData)["createTime"] = modifyTime
			(*destData)["createUnit"] = sysUserMaster["createUnit"]
		} else if destData == nil && srcData != nil {
			// 删除,不处理
		} else if destData != nil && srcData != nil {
			isMasterData := fieldGroupLi[0].IsMasterField()
			isDetailDataDiff := (!fieldGroupLi[0].IsMasterField()) && modelTemplateFactory.IsDataDifferent(fieldGroupLi, *destData, srcData)
			if isMasterData || isDetailDataDiff {
				(*destData)["createBy"] = srcData["createBy"]
				(*destData)["createTime"] = srcData["createTime"]
				(*destData)["createUnit"] = srcData["createUnit"]

				(*destData)["modifyBy"] = userId
				(*destData)["modifyTime"] = modifyTime
				(*destData)["modifyUnit"] = sysUserMaster["createUnit"]
			}
		}
	})
}

func (c BaseDataAction) RollbackTxn(sessionId int) {
	txnId := global.GetGlobalAttr(sessionId, "txnId")
	if txnId != nil {
		if x := recover(); x != nil {
			_, db := global.GetConnection(sessionId)
			txnManager := TxnManager{db}
			txnManager.Rollback(txnId.(int))
			panic(x)
		}
	}
}

func (c BaseDataAction) CommitTxn(sessionId int) {
	txnId := global.GetGlobalAttr(sessionId, "txnId")
	if txnId != nil {
		_, db := global.GetConnection(sessionId)
		txnManager := TxnManager{db}
		txnManager.Commit(txnId.(int))
	}
}

func (c BaseDataAction) RenderCommon(w http.ResponseWriter, r *http.Request, modelRenderVO ModelRenderVO) {
	bo := modelRenderVO.Bo
	relationBo := modelRenderVO.RelationBo
	dataSource := modelRenderVO.DataSource
	usedCheckBo := modelRenderVO.UsedCheckBo

	modelIterator := ModelIterator{}
	var result interface{} = ""
	modelIterator.IterateAllFieldBo(dataSource, &bo, &result, func(fieldGroup FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		if (*data)[fieldGroup.Id] != nil {
			(*data)[fieldGroup.Id] = fmt.Sprint((*data)[fieldGroup.Id])
		}
	})
	format := r.FormValue("format")
	if strings.ToLower(format) == "json" {
		w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
		data, err := json.Marshal(map[string]interface{}{
			"bo":          bo,
			"relationBo":  relationBo,
			"usedCheckBo": usedCheckBo,
		})
		if err != nil {
			panic(err)
		}
		w.Write(data)
		return
	}
	//	return c.Render()
}

/**
 * 列表页
 */
//func (baseData BaseDataAction) ListData() revel.Result {
//
//}

/**
 * 新增
 */
func (c BaseDataAction) NewData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}

	modelRenderVO := c.RNewDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) RNewDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(dataSourceModelId)
	templateManager := TemplateManager{}
	formTemplate := templateManager.GetFormTemplate(formTemplateId)
	c.RActionSupport.BeforeNewData(sessionId, dataSource, formTemplate)
	bo := modelTemplateFactory.GetInstanceByDS(dataSource)
	c.RActionSupport.AfterNewData(sessionId, dataSource, formTemplate, &bo)

	columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
	bo = columnModelData["bo"].(map[string]interface{})
	relationBo := columnModelData["relationBo"].(map[string]interface{})

	modelTemplateFactory.ClearReverseRelation(&dataSource)

	c.CommitTxn(sessionId)
	return ModelRenderVO{
		UserId:       userId,
		Bo:           bo,
		RelationBo:   relationBo,
		DataSource:   dataSource,
		FormTemplate: formTemplate,
	}
}

func (c BaseDataAction) GetData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.GetDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) GetDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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

	session, _ := global.GetConnection(sessionId)
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
	c.RActionSupport.BeforeGetData(sessionId, dataSource, formTemplate)
	bo, found := querySupport.FindByMapWithSession(session, collectionName, queryMap)
	if !found {
		panic("GetData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}
	c.RActionSupport.AfterGetData(sessionId, dataSource, formTemplate, &bo)

	usedCheck := UsedCheck{}
	usedCheckBo := usedCheck.GetFormUsedCheck(sessionId, dataSource, bo)

	columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
	bo = columnModelData["bo"].(map[string]interface{})
	relationBo := columnModelData["relationBo"].(map[string]interface{})

	modelTemplateFactory.ConvertDataType(dataSource, &bo)

	modelTemplateFactory.ClearReverseRelation(&dataSource)
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
 * 复制
 */
func (c BaseDataAction) CopyData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.CopyDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) CopyDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
	srcBo, found := querySupport.FindByMap(collectionName, queryMap)
	if !found {
		panic("CopyData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	modelTemplateFactory.ConvertDataType(dataSource, &srcBo)
	c.RActionSupport.BeforeCopyData(sessionId, dataSource, formTemplate, srcBo)
	dataSource, bo := modelTemplateFactory.GetCopyInstance(dataSourceModelId, srcBo)
	c.RActionSupport.AfterCopyData(sessionId, dataSource, formTemplate, &bo)

	columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
	bo = columnModelData["bo"].(map[string]interface{})
	relationBo := columnModelData["relationBo"].(map[string]interface{})

	modelTemplateFactory.ClearReverseRelation(&dataSource)
	c.CommitTxn(sessionId)
	return ModelRenderVO{
		UserId:       userId,
		Bo:           bo,
		RelationBo:   relationBo,
		DataSource:   dataSource,
		FormTemplate: formTemplate,
	}
}

/**
 * 修改,只取数据,没涉及到数据库保存,涉及到数据库保存的方法是SaveData,
 */
func (c BaseDataAction) EditData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}

	modelRenderVO := c.EditDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) EditDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
		panic("RefreshData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	editMessage, isValid := c.RActionSupport.EditValidate(sessionId, dataSource, formTemplate, bo)
	if !isValid {
		panic(editMessage)
	}

	c.RActionSupport.BeforeEditData(sessionId, dataSource, formTemplate, &bo)
	c.RActionSupport.AfterEditData(sessionId, dataSource, formTemplate, &bo)

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
 * 保存
 */
func (c BaseDataAction) SaveData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.SaveCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) SaveCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
	jsonBo := r.FormValue("jsonData")

	bo := map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonBo), &bo)
	if err != nil {
		panic(err)
	}

	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(dataSourceModelId)
	templateManager := TemplateManager{}
	formTemplate := templateManager.GetFormTemplate(formTemplateId)
	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	strId := modelTemplateFactory.GetStrId(bo)
	if strId == "" || strId == "0" {
		c.SetCreateFixFieldValue(sessionId, dataSource, &bo)
	} else {
		c.SetModifyFixFieldValue(sessionId, dataSource, &bo)
		editMessage, isValid := c.RActionSupport.EditValidate(sessionId, dataSource, formTemplate, bo)
		if !isValid {
			panic(editMessage)
		}
	}

	c.RActionSupport.BeforeSaveData(sessionId, dataSource, formTemplate, &bo)
	financeService := FinanceService{}

	diffDataRowLi := financeService.SaveData(sessionId, dataSource, &bo)

	c.RActionSupport.AfterSaveData(sessionId, dataSource, formTemplate, &bo, diffDataRowLi)

	querySupport := QuerySupport{}
	queryMap := map[string]interface{}{
		"_id": bo["id"],
	}
	permissionSupport := PermissionSupport{}
	permissionQueryDict := permissionSupport.GetPermissionQueryDict(sessionId, formTemplate.Security)
	for k, v := range permissionQueryDict {
		queryMap[k] = v
	}

	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
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
 * 放弃保存,回到浏览状态
 */
func (c BaseDataAction) GiveUpData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.GiveUpDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) GiveUpDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
		panic("giveUpData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	c.RActionSupport.BeforeGiveUpData(sessionId, dataSource, formTemplate, &bo)
	c.RActionSupport.AfterGiveUpData(sessionId, dataSource, formTemplate, &bo)

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
 * 删除
 */
func (c BaseDataAction) DeleteData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}

	modelRenderVO := c.DeleteDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) setRequestParameterToBo(r *http.Request, bo *map[string]interface{}) {
	keyLi := []string{"dataSourceModelId", "formTemplateId", "id"}
	for k, v := range r.Form {
		isIn := false
		for _, item := range keyLi {
			if item == k {
				isIn = true
				break
			}
		}
		if !isIn {
			(*bo)[k] = strings.Join(v, ",")
		}
	}
}

func (c BaseDataAction) DeleteDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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

	_, db := global.GetConnection(sessionId)
	querySupport := QuerySupport{}
	queryMap := map[string]interface{}{
		"_id": id,
	}
	formTemplate := FormTemplate{}
	if formTemplateId != "" {
		templateManager := TemplateManager{}
		formTemplate = templateManager.GetFormTemplate(formTemplateId)
	}
	permissionSupport := PermissionSupport{}
	permissionQueryDict := permissionSupport.GetPermissionQueryDict(sessionId, formTemplate.Security)
	for k, v := range permissionQueryDict {
		queryMap[k] = v
	}

	modelTemplateFactory := ModelTemplateFactory{}
	dataSource := modelTemplateFactory.GetDataSource(dataSourceModelId)
	// 列表页也调用这个删除方法,但是列表页又没有传递formTemplateId,只有 gatheringBill等要做赤字判断,走与form相同的逻辑,才会传 formTemplateId,
	collectionName := modelTemplateFactory.GetCollectionName(dataSource)
	bo, found := querySupport.FindByMap(collectionName, queryMap)
	if !found {
		panic("DeleteData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}
	// 将客户端传入的各种参数写入,程序在业务方法before,after中有可能会用到
	c.setRequestParameterToBo(r, &bo)

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	c.RActionSupport.BeforeDeleteData(sessionId, dataSource, formTemplate, &bo)

	usedCheck := UsedCheck{}
	if usedCheck.CheckUsed(sessionId, dataSource, bo) {
		panic(BusinessError{Message: "已被用，不能删除"})
	}

	modelIterator := ModelIterator{}
	var result interface{} = ""
	modelIterator.IterateDataBo(dataSource, &bo, &result, func(fieldGroupLi []FieldGroup, data *map[string]interface{}, rowIndex int, result *interface{}) {
		//		if fieldGroupLi[0].IsMasterField() {
		usedCheck.DeleteAll(sessionId, fieldGroupLi, *data)
		//		}
	})

	txnManager := TxnManager{db}
	txnId := global.GetTxnId(sessionId)
	_, removeResult := txnManager.Remove(txnId, dataSourceModelId, bo)
	if !removeResult {
		panic("删除失败")
	}

	c.RActionSupport.AfterDeleteData(sessionId, dataSource, formTemplate, &bo)

	// 列表页也调用这个删除方法,但是列表页又没有传递formTemplateId
	relationBo := map[string]interface{}{}
	if formTemplateId != "" {
		templateManager := TemplateManager{}
		columnModelData := templateManager.GetColumnModelDataForFormTemplate(sessionId, formTemplate, bo)
		bo = columnModelData["bo"].(map[string]interface{})
		relationBo = columnModelData["relationBo"].(map[string]interface{})
	}

	modelTemplateFactory.ClearReverseRelation(&dataSource)
	c.CommitTxn(sessionId)
	return ModelRenderVO{
		UserId:       userId,
		Bo:           bo,
		RelationBo:   relationBo,
		DataSource:   dataSource,
		FormTemplate: formTemplate,
	}
}

/**
 * 刷新
 */
func (c BaseDataAction) RefreshData(w http.ResponseWriter, r *http.Request) {
	c.RActionSupport = ActionSupport{}
	modelRenderVO := c.RefreshDataCommon(w, r)
	c.RenderCommon(w, r, modelRenderVO)
}

func (c BaseDataAction) RefreshDataCommon(w http.ResponseWriter, r *http.Request) ModelRenderVO {
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
		panic("RefreshData, dataSouceModelId=" + dataSourceModelId + ", id=" + strId + " not found")
	}

	modelTemplateFactory.ConvertDataType(dataSource, &bo)
	c.RActionSupport.BeforeRefreshData(sessionId, dataSource, formTemplate, &bo)
	c.RActionSupport.AfterRefreshData(sessionId, dataSource, formTemplate, &bo)

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
 * 被用查询
 */
func (c BaseDataAction) LogList(w http.ResponseWriter, r *http.Request) {
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
	return
//	return c.Render()
}

func (c BaseDataAction) LogListCommon(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	dataSourceModelId := r.FormValue("dataSourceModelId")
	//formTemplateId := c.Params.Get("formTemplateId")
	strId := r.FormValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		panic(err)
	}

	collectionName := "PubReferenceLog"
	// reference,beReference
	querySupport := QuerySupport{}
	query := map[string]interface{}{
		"beReference": []interface{}{dataSourceModelId, "A", "id", id},
	}
	pageNo := 1
	pageSize := 10
	orderBy := ""
	return querySupport.Index(collectionName, query, pageNo, pageSize, orderBy)
}
