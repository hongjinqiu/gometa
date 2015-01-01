if (typeof(console) == "undefined") {
	console = {
		log: function() {},
		info: function() {},
		error: function() {},
		fatal: function() {}
	};
}

function CommonUtil() {
}

function executeGYUI(func) {
	if (typeof(g_Y) != "undefined" && g_Y) {
		func(g_Y);
	} else {
		YUI(g_financeModule).use("finance-module", func);
	}
}

CommonUtil.prototype.getFuncOrString = function(text) {
	if (/^[a-zA-Z\d_]*$/.test(text)) {
		if (eval("typeof(" + text + ")") == "function") {
			return eval(text);
		}
	}
	return text;
}

CommonUtil.prototype.getCRelationItem = function(cRelationDS, bo, formData) {
	for (var i = 0; i < cRelationDS.CRelationItemLi.length; i++) {
		var relationItem = cRelationDS.CRelationItemLi[i];
		var mode = relationItem.CJsRelationExpr.Mode;
		var content = relationItem.CJsRelationExpr.Content;
		if (mode == undefined || mode == "" || mode == "text") {
			if (content == "true") {
				return relationItem;
			}
		} else if (mode == "js") {
			var data = formData;
			if (eval(content) === true) {
				return relationItem;
			}
		} else if (mode == "function") {
			var data = formData;
			eval("var f=" + content);
			if (f(data) === true) {
				return relationItem;
			}
		} else if (mode == "functionName") {
			var data = formData;
			if (eval(content + "(data)") === true) {
				return relationItem;
			}
		}
	}
	return null;
}

var panelZIndex = 6;

/**
 * config: title,url
 */
function showModalDialog(config) {
	var moduleName = config.moduleName;
	if (!moduleName) {
		moduleName = "finance-module";
	}
	var func = function(Y) {
		var title = config["title"];
		var url = config["url"];
	//	node.getComputedStyle("width")
//		var width = 700;
//		var height = 500;
		var node = Y.one("window");
	//		width = parseInt(node.getComputedStyle("width"));
	//		height = parseInt(node.getComputedStyle("height"));
		var width = parseInt(node.get("winWidth"), 10);
		var height = parseInt(node.get("winHeight"), 10);
		var edge = 100;
		var frameWidth = width - edge;
		if (frameWidth <= 0) {
			frameWidth = 100;
		}
		var frameHeight = height - edge;
		if (frameHeight <= 0) {
			frameHeight = 100;
		}
//	var bodyContent = null;
		var bodyContent = "<iframe src='{src}' frameborder='0' style='width:100%;height:99%;overflow: auto;'></iframe>";
		bodyContent = Y.Lang.sub(bodyContent, {
			src: url
//			,width: frameWidth
//			,height: frameHeight
		});
	    var dialog = new Y.Panel({
	        contentBox : Y.Node.create('<div id="dialog" />'),
	        headerContent: title,
	        bodyContent: bodyContent,
	        width      : frameWidth,
	        height: frameHeight,
	        zIndex     : (++panelZIndex),
	        centered   : true,
	        modal      : true, // modal behavior
	        render     : '.popupModelDialog',
	        visible    : false, // make visible explicitly with .show()
	        plugins      : [Y.Plugin.Drag],
	        buttons: [
	                  {
	                      value: "",// string or html string
	                      action: function(e) {
	                          e.preventDefault();
	                          dialog.hide();
	                      },
	                      section: Y.WidgetStdMod.HEADER,
	                      classNames: "closeBtn"
	                  }
	              ]
	    });

	    dialog.hide = function() {
	    	window.s_dialog = null;
			return this.destroy();
		}
	    
	    dialog.dd.addHandle('.yui3-widget-hd');
	    dialog.show();
	    window.s_dialog = dialog;
	};
	if (moduleName == "finance-module") {
		executeGYUI(func);
	} else {
		YUI(g_financeModule).use(moduleName, func);
	}
}

function triggerShowModalDialog(config) {
	if (top && top.putTabIfAbsent) {
		top.putTabIfAbsent(config.title.replace("列表", ""), config.url);
	} else {
		showModalDialog(config);
	}
}

/**
 * infoType:info,error,question,warn
 */
function showDialog(config){
	var infoType = config["infoType"];
	var title = config["title"];
	var msg = config["msg"];
	var callback = config["callback"];
	var width = config["width"] || 410;
	var height = config["height"] || 130;
	var bodyHeight = height - 23 - 40 - 50;
	var bodyContent = null;
	var footer = [];
	if (infoType == "info") {
//		bodyContent = '<div class="message icon-info overflowAuto" style="height:' + bodyHeight + 'px;">' + msg + '</div>';
		bodyContent = '<div class="message icon-info overflowAuto" style="">' + msg + '</div>';
		footer = [{
            name     : 'proceed',
            label    : '确定',
            action   : 'onOK',
            classNames: 'message_bt1'
        }];
	} else if (infoType == "success") {
//		bodyContent = '<div class="message icon-success overflowAuto" style="height:' + bodyHeight + 'px;">' + msg + '</div>';
		bodyContent = '<div class="message icon-success overflowAuto" style="">' + msg + '</div>';
		footer = [{
            name     : 'proceed',
            label    : '确定',
            action   : 'onOK',
            classNames: 'message_bt1'
        }];
	} else if (infoType == "warn") {
//		bodyContent = '<div class="message icon-warn overflowAuto" style="height:' + bodyHeight + 'px;">' + msg + '</div>';
		bodyContent = '<div class="message icon-warn overflowAuto" style="">' + msg + '</div>';
		footer = [{
            name     : 'proceed',
            label    : '确定',
            action   : 'onOK',
            classNames: 'message_bt1'
        }];
	} else if (infoType == "question") {
//		bodyContent = '<div class="message icon-question overflowAuto" style="height:' + bodyHeight + 'px;">' + msg + '</div>';
		bodyContent = '<div class="message icon-question overflowAuto" style="">' + msg + '</div>';
		footer = [{
            name  : 'cancel',
            label : '取消',
            action: 'onCancel',
            classNames: 'message_bt1'
        }, {
            name     : 'proceed',
            label    : '确定',
            action   : 'onOK',
            classNames: 'message_bt1'
        }];
	} else if (infoType == "error") {
//		bodyContent = '<div class="message icon-error overflowAuto" style="height:' + bodyHeight + 'px;">' + msg + '</div>';
		bodyContent = '<div class="message icon-error overflowAuto" style="">' + msg + '</div>';
		footer = [{
            name     : 'proceed',
            label    : '确定',
            action   : 'onOK',
            classNames: 'message_bt1'
        }];
	}
	
	executeGYUI(function(Y) {
	    var dialog = new Y.Panel({
	        contentBox : Y.Node.create('<div id="dialog" />'),
	        headerContent: title,
	        bodyContent: bodyContent,
	        width      : width,
	        height: height,
	        zIndex     : (++panelZIndex),
	        centered   : true,
	        modal      : true, // modal behavior
	        render     : '.popupDialog',
	        visible    : false, // make visible explicitly with .show()
	        plugins      : [Y.Plugin.Drag],
	        buttons    : {
	        	footer: footer
	        }
	    });

	    dialog.onCancel = function (e) {
	        e.preventDefault();
	        this.hide();
	    }

	    dialog.onOK = function (e) {
	        e.preventDefault();
	        this.hide();
	        if (callback) {
	        	callback();
	        }
	    }
	    
	    dialog.hide = function() {
			return this.destroy();
		}
	    
	    dialog.dd.addHandle('.yui3-widget-hd');
	    dialog.show();
	    if (infoType == "info" || infoType == "success" || infoType == "warn" || infoType == "error") {
	    	dialog.getButton("proceed").focus();
	    }
	});
}

function showAlert(msg, callback, width, height){
	showDialog({
		"infoType": "info",
		"title": "提示信息",
		"msg": msg,
		"callback": callback,
		"width": width,
		"height": height
	});
}

function showSuccess(msg, callback, width, height){
	showDialog({
		"infoType": "success",
		"title": "成功信息",
		"msg": msg,
		"callback": callback,
		"width": width,
		"height": height
	});
}

function showError(msg, callback, width, height){
	showDialog({
		"infoType": "error",
		"title": "错误信息",
		"msg": msg,
		"callback": callback,
		"width": width,
		"height": height
	});
}

function showWarning(msg, callback, width, height){
	showDialog({
		"infoType": "warn",
		"title": "警告信息",
		"msg": msg,
		"callback": callback,
		"width": width,
		"height": height
	});
}

function showConfirm(msg, callback, width, height){
	showDialog({
		"infoType": "question",
		"title": "确认信息",
		"msg": msg,
		"callback": callback,
		"width": width,
		"height": height
	});
}

/**
 * 配置demo:
 * {
 * 	sync: true | false,
 * 	method: GET | POST,
 * 	params: post data,
 * 	callback: success callback function,
 *  failCallback: fail callback function,
 * }
 */
function ajaxRequest(option){
	// 有用的配置为 doCallback, 自己对failure,error进行提示即可,
	// url,params,async,scope,
	var moduleName = option.moduleName;
	if (!moduleName) {
		moduleName = "finance-module";
	}
	YUI(g_financeModule).use(moduleName, function(Y){
//		var paramData = Y.JSON.stringify(option["params"]);
		var paramData = {};
		if (option.params) {
			for (var k in option.params) {
				if (typeof(option.params[k]) == "object") {
					paramData[k] = Y.JSON.stringify(option.params[k]);
				} else {
					paramData[k] = option.params[k];
				}
			}
		}
		var cfg = {
			sync: option.sync !== undefined ? option.sync : true,
			method: option.method || 'POST',
			data: Y.QueryString.stringify(paramData),
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded'
			},
			on: {
				start: function(){
//								console.log("start");
				},
				complete: function(){
//								console.log("complete");
				},
				success: function(id, o, args){
					if (option.callback || option.failCallback) {
						try {
							var jsonRsp = {};
							if (o.responseText) {
								jsonRsp = Y.JSON.parse(o.responseText);
							}
							if (jsonRsp.success === false) {
								if (option.failCallback) {
									option.failCallback(jsonRsp);
								} else {
									showError(jsonRsp.message);
								}
							} else {
								if (option.callback) {
									option.callback(jsonRsp);
								}
							}
							
							/*
							if (o.responseText) {
								var jsonRsp = Y.JSON.parse(o.responseText);
								option.callback(jsonRsp);
							} else {
								option.callback({});
							}
							*/
						} catch (e) {
							console.error(option["url"]);
							console.error(e);
						}
					}
				},
				failure: function(id, o, args){// failure调用在complete之前,
					var text = o.responseText;
					var reg = /panic\(&#34;(.*?)&#34;\)/.test(text);
					var msg = RegExp.$1;
					if (msg) {
						showError(msg);
					} else {
						showError(text, null, 600, 400);
					}
				},
				end: function(){
//								console.log("end");
				}
			}
		};
//		console.log(Y.QueryString.stringify(paramData));
		Y.io(option["url"], cfg);
//		function complete(id, o, args) {
//			var id = id; // Transaction ID.
//			var data = Y.JSON.parse(o.responseText);
//			
//		};
//		io:complete
//		io:end
//		io:failure
//		io:progress
//		io:start
//		io:success
//		io:xdrReady
//		Y.on('io:complete', complete, Y, []);
//		var request = Y.io(uri);
//		Y.QueryString.stringify
		// import json
	});
}

function g_setMasterFormFieldStatus(status) {
	if (status == "view") {
		for (var key in g_masterFormFieldDict) {
			
		}
	}
}

function isNumber(value) {
	return /^-?\d*(\.\d*)?$/.test(value);
}

/**
 * @param continueAnyAll "true"|"false" 出现赤字是否继续
 */
function limitControlSaveData(continueAnyAll) {
	var formManager = new FormManager();
	var bo = formManager.getBo();
	if (continueAnyAll == "true" || continueAnyAll == "false") {
		bo.continueAnyAll = continueAnyAll;
	}
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
				// 为现金账户和银行账户添加的自定义函数
				if (typeof(enableQueryParameters) != "undefined") {
					enableQueryParameters();
				}
			},
			failCallback: function(o) {
				if (o.code == "3") {// 赤字警告
					showConfirm(o.message + "<br />是否继续？", function(){
						limitControlSaveData("true");
					});
				} else {
					showError(o.message);
				}
			}
		});
	}
}

function limitControlDeleteDataExecute(continueAnyAll) {
	if (continueAnyAll != "true" && continueAnyAll != "false") {
		continueAnyAll = "false";
	}
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/DeleteData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"],
			"continueAnyAll": continueAnyAll || ""
		},
		callback: function(o) {
			location.href = "/console/listschema?@name=" + g_dataSourceJson.Id;
		},
		failCallback: function(o) {
			if (o.code == "3") {// 赤字警告
				showConfirm(o.message + "<br />是否继续？", function(){
					limitControlDeleteDataExecute("true");
				});
			} else {
				showError(o.message);
			}
		}
	});
}

function limitControlDeleteData() {
	showConfirm("您确定要删除吗？", function(){
		limitControlDeleteDataExecute();
	})
}

function limitControlCancelData(continueAnyAll) {
	if (continueAnyAll != "true" && continueAnyAll != "false") {
		continueAnyAll = "false";
	}
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/CancelData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"],
			"continueAnyAll": continueAnyAll || ""
		},
		callback: function(o) {
			showSuccess("作废数据成功");
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("view");
		},
		failCallback: function(o) {
			if (o.code == "3") {// 赤字警告
				showConfirm(o.message + "<br />是否继续？", function(){
					limitControlCancelData("true");
				});
			} else {
				showError(o.message);
			}
		}
	});
}

function limitControlUnCancelData(continueAnyAll) {
	if (continueAnyAll != "true" && continueAnyAll != "false") {
		continueAnyAll = "false";
	}
	var formManager = new FormManager();
	var bo = formManager.getBo();
	ajaxRequest({
		url: "/" + g_dataSourceJson.Id + "/UnCancelData?format=json"
		,params: {
			"dataSourceModelId": g_dataSourceJson.Id,
			"formTemplateId": g_formTemplateJsonData.Id,
			"id": bo["id"],
			"continueAnyAll": continueAnyAll || ""
		},
		callback: function(o) {
			showSuccess("反作废数据成功");
			formManager.applyGlobalParamFromAjaxData(o);
			formManager.loadData2Form(g_dataSourceJson, o.bo);
			formManager.setFormStatus("view");
		},
		failCallback: function(o) {
			if (o.code == "3") {// 赤字警告
				showConfirm(o.message + "<br />是否继续？", function(){
					limitControlUnCancelData("true");
				});
			} else {
				showError(o.message);
			}
		}
	});
}

function getFormJsonData(formName) {
	var form = document.forms[formName];
	return getChildFormFieldValueMap(form);
}

function getChildFormFieldValueMap(elem, seperator) {
	if (!seperator) {
		seperator = ",";
	}
	var result = {};
	var inputLi = elem.getElementsByTagName("input");
	for (var i = 0; i < inputLi.length; i++) {
		if (inputLi[i].type.toLowerCase() == "text" || inputLi[i].type.toLowerCase() == "hidden") {
			var name = inputLi[i].name;
			putOrAppend(result, inputLi[i].name, inputLi[i].value, seperator);
		} else if (inputLi[i].type.toLowerCase() == "radio") {
			if (inputLi[i].checked) {
				putOrAppend(result, inputLi[i].name, inputLi[i].value, seperator);
			}
		} else if (inputLi[i].type.toLowerCase() == "checkbox") {
			if (inputLi[i].checked) {
				putOrAppend(result, inputLi[i].name, inputLi[i].value, seperator);
			}
		}
	}
	
	var selectLi = elem.getElementsByTagName("select");
	for (var i = 0; i < selectLi.length; i++) {
		var name = selectLi[i].name;
		putOrAppend(result, selectLi[i].name, selectLi[i].value, seperator);
	}
	
	var textareaLi = elem.getElementsByTagName("textarea");
	for (var i = 0; i < textareaLi.length; i++) {
		var name = textareaLi[i].name;
		putOrAppend(result, textareaLi[i].name, textareaLi[i].value, seperator);
	}
	
	return result;
}

function getCheckboxValue(checkboxName) {
	var seperator = ",";
	var result = {};
	var inputLi = document.getElementsByTagName("input");
	for (var i = 0; i < inputLi.length; i++) {
		if (inputLi[i].name == checkboxName && inputLi[i]["type"].toLowerCase() == "radio") {
			if (inputLi[i].checked) {
				putOrAppend(result, inputLi[i].name, inputLi[i].value, seperator);
			}
		} else if (inputLi[i].type.toLowerCase() == "checkbox") {
			if (inputLi[i].checked) {
				putOrAppend(result, inputLi[i].name, inputLi[i].value, seperator);
			}
		}
	}
	return result[checkboxName] || "";
}

function putOrAppend(dictObj, name, value, seperator) {
	if (dictObj[name] !== undefined) {
		if (seperator) {
			dictObj[name] += seperator + value;
		} else {
			dictObj[name] += "," + value;
		}
	} else {
		dictObj[name] = value || "";
	}
}

function openTabOrJump(url) {
	if (top && top.putTabIfAbsent) {
		var name = listTemplate.Description.replace("列表", "");
		var isRefresh = true;
		top.putTabIfAbsent(name, url, isRefresh);
	} else {
		location.href = url;
	}
}
