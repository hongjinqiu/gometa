function QueryParameterManager() {}

/**
 * 设置传入的默认值
 */
QueryParameterManager.prototype.applyQueryDefaultValue = function(Y) {
		if (g_defaultBo) {
			for (var key in g_masterFormFieldDict) {
				if (g_defaultBo[key]) {
					g_masterFormFieldDict[key].set("value", g_defaultBo[key]);
				} else {
					g_masterFormFieldDict[key].set("value", "");
				}
			}
		} else {
			for (var key in g_masterFormFieldDict) {
				g_masterFormFieldDict[key].set("value", "");
			}
		}
}

/**
 * 设置url传入的参数值
 */
QueryParameterManager.prototype.applyFormData = function(Y) {
		if (g_formDataJson) {
			for (var key in g_formDataJson) {
				if (g_masterFormFieldDict[key]) {
					if (g_formDataJson[key]) {
						g_masterFormFieldDict[key].set("value", g_formDataJson[key]);
					} else {
						g_masterFormFieldDict[key].set("value", "");
					}
				}
			}
		}
}

QueryParameterManager.prototype.applyObserveEventBehavior = function() {
	for (var key in g_masterFormFieldDict) {
		var self = g_masterFormFieldDict[key];
		var listTemplateIterator = new ListTemplateIterator();
		var result = "";
		listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(queryParameter, result){
			if (queryParameter.Name == self.get("name")) {
				if (queryParameter.ParameterAttributeLi) {
					for (var j = 0; j < queryParameter.ParameterAttributeLi.length; j++) {
						if (queryParameter.ParameterAttributeLi[j].Name == "observe") {
//							function() {
//								
//							}();
							YUI(g_financeModule).use("finance-module", function(Y){
								self.after('valueChange', Y.bind(function(e) {
									listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(targetQueryParameter, result){
										if (targetQueryParameter.Name == queryParameter.ParameterAttributeLi[j].Value) {
											var queryParameterManager = new QueryParameterManager();
											var treeUrlAttr = queryParameterManager.findQueryParameterAttr(targetQueryParameter, "treeUrl");
											
											if (treeUrlAttr) {
												var uri = "/tree/" + treeUrlAttr.Value;
												if (uri.indexOf("?") > -1) {
													uri += "&parentId=" + g_masterFormFieldDict[queryParameter.Name].get("value");
												} else {
													uri += "?parentId=" + g_masterFormFieldDict[queryParameter.Name].get("value");
												}
												function complete(id, o, args) {
													var id = id; // Transaction ID.
													var data = Y.JSON.parse(o.responseText);
													var choicesLi = [];
													for (var k = 0; k < data.length; k++) {
														choicesLi.push({
															"label": data[k].name,
															"value": data[k].code
														});
													}
													g_masterFormFieldDict[targetQueryParameter.Name].set("choices", choicesLi);
													g_masterFormFieldDict[targetQueryParameter.Name].set("value", "");
												};
												Y.on('io:complete', complete, Y, []);
												var request = Y.io(uri);
											} else {
												if (g_masterFormFieldDict[targetQueryParameter.Name]) {
													g_masterFormFieldDict[targetQueryParameter.Name].set("value", "");
												}
											}
											
											return true;
										}
										return false;
									});
								},
								this));
							// 弄到这里,
							});
							
							break;
						}
					}
				}
				return true;
			}
			return false;
		});
	}
}

QueryParameterManager.prototype.findQueryParameterAttr = function(queryParameter, name) {
	if (queryParameter.ParameterAttributeLi) {
		for (var i = 0; i < queryParameter.ParameterAttributeLi.length; i++) {
			if (queryParameter.ParameterAttributeLi[i].Name == name) {
				return queryParameter.ParameterAttributeLi[i];
			}
		}
	}
	return null;
}

QueryParameterManager.prototype.getQueryField = function(Y, name) {
	var listTemplateIterator = new ListTemplateIterator();
	var result = "";
	var field = null;
	listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(queryParameter, result){
		if (queryParameter.Name == name) {
			var readOnly = queryParameter.ReadOnly === "true";
			if (queryParameter.Editor == "hiddenfield") {
				field = new Y.LHiddenField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "textfield") {
				field = new Y.LTextField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "textareafield") {
				field = new Y.LTextareaField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "numberfield") {
				field = new Y.LNumberField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "datefield") {
				field = new Y.LDateField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "combofield") {
				field = new Y.LSelectField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "displayfield") {
				field = new Y.LDisplayField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "checkboxfield") {
				field = new Y.LChoiceField({
					name : name,
					validateInline: true,
					multi: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "radiofield") {
				field = new Y.LChoiceField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			} else if (queryParameter.Editor == "triggerfield") {
				field = new Y.LTriggerField({
					name : name,
					validateInline: true,
					readonly: readOnly
				});
			}
			return true;
		}
		return false;
	});
	return field;
}

QueryParameterManager.prototype.getQueryFormData = function() {
	var result = {};
	for (var key in g_masterFormFieldDict) {
		result[key] = g_masterFormFieldDict[key].get("value");
	}
	return result;
}

