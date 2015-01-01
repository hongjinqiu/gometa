function g_deleteRecord(o) {
	showConfirm("确认删除？", function(){
		var url = "/" + listTemplate.DataSourceModelId + "/DeleteData?format=json";
		ajaxRequest({
			url: url
			,params: {
				"id": o.get("id"),
				"dataSourceModelId": listTemplate.DataSourceModelId
			},
			callback: function(o) {
//			showSuccess("删除数据成功");
				g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
			}
		});
	});
}

function g_deleteRecords() {
	var selectRecords = g_gridPanelDict["columnModel_1"].getSelectRecordLi();
	if (selectRecords.length > 0) {
		showConfirm("确认删除？", function(){
			var errorMsgLi = [];
			
			var url = "/" + listTemplate.DataSourceModelId + "/DeleteData?format=json";
			for (var i = 0; i < selectRecords.length; i++) {
				ajaxRequest({
					url: url
					,params: {
						"id": selectRecords[i].get("id"),
						"dataSourceModelId": listTemplate.DataSourceModelId
					},
					callback: function(o) {
						//g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
					},
					failCallback: function(o) {
						var message = "记录" + selectRecords[i].get("code");
						message += "：" + o.message+"；";
						errorMsgLi.push(message);
					}
				});
			}
			if (errorMsgLi.length > 0) {
				showError(errorMsgLi.join("<br />"));
			}
			g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
		});
	} else {
		showAlert("请选择记录！");
	}
}

function listLimitControlDeleteDataExecute(record, continueAnyAll) {
	if (continueAnyAll != "true" && continueAnyAll != "false") {
		continueAnyAll = "false";
	}
	ajaxRequest({
		url: "/" + listTemplate.DataSourceModelId + "/DeleteData?format=json"
		,params: {
			"dataSourceModelId": listTemplate.DataSourceModelId,
			"formTemplateId": listTemplate.DataSourceModelId,
			"id": record.get("id"),
			"continueAnyAll": continueAnyAll || ""
		},
		callback: function(o) {
			g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
		},
		failCallback: function(o) {
			if (o.code == "3") {// 赤字警告
				showConfirm(o.message + "<br />是否继续？", function(){
					listLimitControlDeleteDataExecute(record, "true");
				});
			} else {
				showError(o.message);
			}
		}
	});
}

function listLimitControlDeleteData(o) {
	showConfirm("确认删除？", function(){
		listLimitControlDeleteDataExecute(o);
	});
}

function listLimitControlDeleteDataGroupExecute(selectRecords, continueAnyAll) {
	if (continueAnyAll != "true" && continueAnyAll != "false") {
		continueAnyAll = "false";
	}
	var errorMsgLi = [];
	var warnMsgLi = [];
	var warnRecordLi = [];
	
	for (var i = 0; i < selectRecords.length; i++) {
		ajaxRequest({
			url: "/" + listTemplate.DataSourceModelId + "/DeleteData?format=json"
			,params: {
				"dataSourceModelId": listTemplate.DataSourceModelId,
				"formTemplateId": listTemplate.DataSourceModelId,
				"id": selectRecords[i].get("id"),
				"continueAnyAll": continueAnyAll || ""
			},
			callback: function(o) {
				//g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
			},
			failCallback: function(o) {
				if (o.code == "3") {// 赤字警告
					warnMsgLi.push(o.message);
					warnRecordLi.push(selectRecords[i]);
				} else {
					errorMsgLi.push(o.message);
				}
			}
		});
	}
	if (errorMsgLi.length + warnMsgLi.length != selectRecords.length) {
		g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
	}
	if (errorMsgLi.length > 0) {
		showError(errorMsgLi.join("<br />"));
	} else if (warnMsgLi.length > 0) {
		showConfirm(warnMsgLi.join("<br />") + "<br />是否继续？", function(){
			listLimitControlDeleteDataGroupExecute(warnRecordLi, "true");
		});
	}
//	g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
}

function listLimitControlDeleteDataGroup() {
	var selectRecords = g_gridPanelDict["columnModel_1"].getSelectRecordLi();
	if (selectRecords.length > 0) {
		showConfirm("确认删除？", function(){
			listLimitControlDeleteDataGroupExecute(selectRecords);
		});
	} else {
		showAlert("请选择记录！");
	}
}
