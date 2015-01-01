Y.LTriggerField = Y.Base.create('l-trigger-field', Y.RTriggerField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell trigger_input inputWidth',
	SELECT_CLASS: 'trigger_select ltrigger_select',
	VIEW_CLASS: 'trigger_view ltrigger_select',
	DELETE_CLASS: 'trigger_delete ltrigger_select',
	
	initializer : function () {
		Y.LTriggerField.superclass.initializer.apply(this, arguments);
		var self = this;
		
		new LFormManager().applyEventBehavior(self);
		
		// 需要配置在extraInfo里面,
		var selectFunc = function(selectValueLi, formObj){
			
		}
		var unSelectFunc = function(formObj){
			
		}
		var queryFunc = function() {
			return {};
		}
		var multi = false;
		var selectorName = "";
		var displayField = "";
		var valueField = "id";
		var selectorTitle;
		var listTemplateIterator = new ListTemplateIterator();
		var result = "";
		
		self._setDefaultSelectAction();
		
		listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(queryParameter, result){
			if (queryParameter.Name == self.get("name")) {
				selectFunc = queryParameter.jsConfig.selectFunc;
				unSelectFunc = queryParameter.jsConfig.unSelectFunc;
				queryFunc = queryParameter.jsConfig.queryFunc;
				
				// 从selector里面取,
				selectorName = function() {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					var relationItem = self._relationFuncTemplate(queryParameter, formData);
					if (relationItem) {
						return relationItem.CRelationConfig.SelectorName;
					}
					return "";
				}
				displayField = function() {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					var relationItem = self._relationFuncTemplate(queryParameter, formData);
					if (relationItem) {
						return relationItem.CRelationConfig.DisplayField;
					}
					return "";
				}
				multi = function() {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					var relationItem = self._relationFuncTemplate(queryParameter, formData);
					if (relationItem) {
						return relationItem.CRelationConfig.SelectionMode == "multi";
					}
					return false;
				}
				valueField = function() {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					var relationItem = self._relationFuncTemplate(queryParameter, formData);
					if (relationItem) {
						return relationItem.CRelationConfig.ValueField;
					}
					return "";
				}
				selectorTitle = function() {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					var name = selectorName();
					if (name) {
						return g_relationBo[name].Description;
					}
					return "";
				}
				return true;
			}
			return false;
		});
		
		this.set("multi", multi);
		this.set("selectorName", selectorName);
		this.set("displayField", displayField);
		this.set("valueField", valueField);
		this.set("selectFunc", selectFunc);
		this.set("unSelectFunc", unSelectFunc);
		this.set("queryFunc", queryFunc);
		this.set("selectorTitle", selectorTitle);
    },
    
    _setDefaultSelectAction: function() {
    	var self = this;
    	var listTemplateIterator = new ListTemplateIterator();
		var result = "";
		listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(queryParameter, result){
			if (queryParameter.Name == self.get("name")) {
				if (!queryParameter.jsConfig) {
					queryParameter.jsConfig = {};
				}
				if (!queryParameter.jsConfig.selectFunc) {
					queryParameter.jsConfig.selectFunc = function(selectValueLi, formObj) {
						if (!selectValueLi || selectValueLi.length == 0) {
							queryParameter.jsConfig.unSelectFunc(formObj);
    					} else {
    						formObj.set("value", selectValueLi.join(","));
    					}
					}
				}
				if (!queryParameter.jsConfig.unSelectFunc) {
					queryParameter.jsConfig.unSelectFunc = function(formObj) {
						formObj.set("value", "");
					}
				}
				if (!queryParameter.jsConfig.queryFunc) {
					queryParameter.jsConfig.queryFunc = function() {
						return {};
					}
				}
				
				return true;
			}
			return false;
		});
    },
    
    _relationFuncTemplate: function(queryParameter, formData) {
    	var commonUtil = new CommonUtil();
    	var bo = {"A": formData};
    	return commonUtil.getCRelationItem(queryParameter.CRelationDS, bo, formData);
    },
    
    bindUI: function() {
    	Y.LTriggerField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	this.after('valueChange', Y.bind(function(e) {
			var listTemplateIterator = new ListTemplateIterator();
			var result = "";
			listTemplateIterator.iterateAnyTemplateQueryParameter(result, function(queryParameter, result){
				if (queryParameter.Name == self.get("name")) {
					var queryParameterManager = new QueryParameterManager();
					var formData = queryParameterManager.getQueryFormData();
					
					var relationItem = self._relationFuncTemplate(queryParameter, formData);
					if (relationItem) {
						if (relationItem.CCopyConfigLi) {
							var selectorName = self.get("selectorName")();
							if (self.get("value")) {
								var selectorDict = g_relationManager.getRelationBo(selectorName, self.get("value"));
								if (selectorDict) {
									for (var i = 0; i < relationItem.CCopyConfigLi.length; i++) {
										var copyColumnName = relationItem.CCopyConfigLi[i].CopyColumnName;
										var copyValueField = relationItem.CCopyConfigLi[i].CopyValueField;
										if (g_masterFormFieldDict[copyColumnName]) {
											var valueFieldLi = copyValueField.split(",");
											var valueLi = [];
											for (var j = 0; j < valueFieldLi.length; j++) {
												if (selectorDict[valueFieldLi[j]]) {
													valueLi.push(selectorDict[valueFieldLi[j]]);
												}
											}
											g_masterFormFieldDict[copyColumnName].set("value", valueLi.join(","));
										}
									}
								}
							} else {
								for (var i = 0; i < relationItem.CCopyConfigLi.length; i++) {
									var copyColumnName = relationItem.CCopyConfigLi[i].CopyColumnName;
									if (g_masterFormFieldDict[copyColumnName]) {
										g_masterFormFieldDict[copyColumnName].set("value", "");
									}
								}
							}
						}
					}
					
					return true;
				}
				return false;
			});
    	},
        this));
    	new LFormManager().applyEventBehavior(self);
    }
},
{

    ATTRS: {
    	
    }
});
