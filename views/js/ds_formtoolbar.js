function editData() {//修改
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/EditData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"]
		},
		callback: function(o) {
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("edit");
		}
	});
}

function saveData() {//保存
	var formManager = new FormManager();
	var bo = formManager.getBo();
	var validateResult = formManager.dsFormValidator(g_dataSourceJson, bo);
	
	if (!validateResult.result) {
		showError(validateResult.message);
	} else {
		ajaxRequest({
			url: "/" + g_dataSourceJson.Id + "/SaveData?format=json"
			,params: {
				"dataSourceModelId": g_dataSourceJson.Id,
				"formTemplateId": g_formTemplateJsonData.Id,
				"jsonData": bo
			},
			callback: function(o) {
				showSuccess("保存数据成功");
				formManager.setFormStatus("view");
				formManager.applyGlobalParamFromAjaxData(o);
				formManager.loadData2Form(g_dataSourceJson, o.bo);
			}
		});
	}
}

function newData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/NewData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id
		},
		callback: function(o) {
			formManager.setDetailIncId(g_dataSourceJson, o.bo);
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("edit");
		}
	});
}

function copyData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/CopyData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"]
		},
		callback: function(o) {
			formManager.setDetailIncId(g_dataSourceJson, o.bo);
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("edit");
		}
	});
}

function giveUpData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	showConfirm("您确定要放弃吗？", function(){
		if (!bo["id"] || bo["id"] == "0") {
			location.href = "/console/listschema?@name=" + g_dataSourceJson.Id;
		} else {
			ajaxRequest({
				url: "/" + g_dataSourceJson.Id + "/GiveUpData?format=json"
				,params: {
					"dataSourceModelId": g_dataSourceJson.Id,
					"formTemplateId": g_formTemplateJsonData.Id,
					"id": bo["id"]
				},
				callback: function(o) {
					formManager.applyGlobalParamFromAjaxData(o);
					formManager.loadData2Form(g_dataSourceJson, o.bo);
					formManager.setFormStatus("view");
				}
			});
		}
	});
}

function deleteData() {
	showConfirm("您确定要删除吗？", function(){
		var formManager = new FormManager();
		var bo = formManager.getBo();
		ajaxRequest({
			url: "/" + g_dataSourceJson.Id + "/DeleteData?format=json"
			,params: {
				"dataSourceModelId": g_dataSourceJson.Id,
				"formTemplateId": g_formTemplateJsonData.Id,
				"id": bo["id"]
			},
			callback: function(o) {
				location.href = "/console/listschema?@name=" + g_dataSourceJson.Id;
			}
		});
	})
}

function refreshData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/RefreshData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"]
		},
		callback: function(o) {
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("view");
		}
	});
}

function logList() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	var dialog = showModalDialog({
		"title": "被用查询",
		"url": "/console/listschema?@name=PubReferenceLog&beReferenceDataSourceModelId=" + g_dataSourceJson.Id + "&beReferenceId=" + bo["id"] + "&date=" + new Date()
	});
}

function cancelData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/CancelData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"]
		},
		callback: function(o) {
			showSuccess("作废数据成功");
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("view");
		}
	});
}

function unCancelData() {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/UnCancelData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"]
		},
		callback: function(o) {
			showSuccess("反作废数据成功");
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("view");
		}
	});
}

/*
function getData() {
	ajaxRequest({
		url: "/ActionTest/GetData?format=json"
		,params: {
			"id": 26,
			"dataSourceModelId": "ActionTest"
		},
		callback: function(o) {
			console.log(o);
			showSuccess(o.responseText);
		}
	});
}
*/
function test() {
	var relationManager = new RelationManager();
	relationManager.getRelationBo("SysUserSelector", 16);
	return;
}

function ToolbarManager(){}

function setBorderTmp(btn, status) {
//	if (true) {
//		return;
//	}
	if (status == "disabled") {
		if (!btn.origClassName) {
			btn.origClassName = btn.className;
		}
//		btn.style.border = "1px solid black";
		btn.className = "disable_but_box";
	} else {
//		btn.style.border = "1px solid red";
		if (btn.className == "disable_but_box") {
			btn.className = btn.origClassName;
		}
	}
}

ToolbarManager.prototype.enableDisableToolbarBtn = function() {
	if (g_formStatus == "view") {
		var viewEnableBtnLi = ["listBtn","newBtn","copyBtn","cancelBtn","unCancelBtn","refreshBtn","usedQueryBtn"];
		var viewDisableBtnLi = ["saveBtn","giveUpBtn"];
		
		// cancelBtn,
		if (g_masterFormFieldDict["billStatus"]) {
			if (g_masterFormFieldDict["billStatus"].get("value") == "1") {// 正常
				viewEnableBtnLi.push("cancelBtn");
				viewEnableBtnLi.push("editBtn");
			} else {
				viewDisableBtnLi.push("cancelBtn");
				viewDisableBtnLi.push("editBtn");
			}
		} else {
			viewEnableBtnLi.push("editBtn");
		}
		// unCancelBtn,
		if (g_masterFormFieldDict["billStatus"]) {
			if (g_masterFormFieldDict["billStatus"].get("value") == "4") {// 作废
				viewEnableBtnLi.push("unCancelBtn");
			} else {
				viewDisableBtnLi.push("unCancelBtn");
			}
		}
		// delBtn,
		var isUsed = false;
		if (g_usedCheck) {
			if (g_usedCheck["A"]) {
				var id = g_masterFormFieldDict["id"].get("value");
				if (g_usedCheck["A"][id]) {
					isUsed = true;
				}
			}
		}
		if (isUsed) {
			viewDisableBtnLi.push("delBtn");
		} else {
			viewEnableBtnLi.push("delBtn");
		}
		
		
		for (var i = 0; i < viewEnableBtnLi.length; i++) {
			var btn = document.getElementById(viewEnableBtnLi[i]);
			if (btn) {
				btn.disabled = "";
				setBorderTmp(btn, "");
			}
		}
		/*
		var cancelBtn = document.getElementById("cancelBtn");
		if (cancelBtn && g_masterFormFieldDict["billStatus"]) {
			if (g_masterFormFieldDict["billStatus"].get("value") == "1") {
				cancelBtn.disabled = "";
				setBorderTmp(cancelBtn, "");
			} else {
				cancelBtn.disabled = "disabled";
				setBorderTmp(cancelBtn, "disabled");
			}
		}
		var unCancelBtn = document.getElementById("unCancelBtn");
		if (unCancelBtn && g_masterFormFieldDict["billStatus"]) {
			if (g_masterFormFieldDict["billStatus"].get("value") == "4") {
				unCancelBtn.disabled = "";
				setBorderTmp(unCancelBtn, "");
			} else {
				unCancelBtn.disabled = "disabled";
				setBorderTmp(unCancelBtn, "disabled");
			}
		}
		*/
		
		for (var i = 0; i < viewDisableBtnLi.length; i++) {
			var btn = document.getElementById(viewDisableBtnLi[i]);
			if (btn) {
				btn.disabled = "disabled";
				setBorderTmp(btn, "disabled");
			}
		}
	} else {
		var editEnableBtnLi = ["listBtn","saveBtn","giveUpBtn"];
		var editDisableBtnLi = ["newBtn","copyBtn","editBtn","delBtn","cancelBtn","unCancelBtn","refreshBtn","usedQueryBtn"];
		for (var i = 0; i < editEnableBtnLi.length; i++) {
			var btn = document.getElementById(editEnableBtnLi[i]);
			if (btn) {
				btn.disabled = "";
				setBorderTmp(btn, "");
			}
		}
		for (var i = 0; i < editDisableBtnLi.length; i++) {
			var btn = document.getElementById(editDisableBtnLi[i]);
			if (btn) {
				btn.disabled = "disabled";
				setBorderTmp(btn, "disabled");
			}
		}
	}
}





