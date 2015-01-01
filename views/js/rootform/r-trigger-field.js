/**
 * @class RTriggerField
 * @extends RFormField
 * @param config {Object} Configuration object
 * @constructor
 * @description A form field which allows one or multiple values from a 
 * selection of choices
 */
Y.RTriggerField = Y.Base.create('r-trigger-field', Y.RFormField, [Y.WidgetParent, Y.WidgetChild], {

	FIELD_CLASS : 'table-layout-cell trigger_input',
	INPUT_TYPE: "text",
	SELECT_TEMPLATE: '<a></a>',
	SELECT_CLASS: 'trigger_select',
	VIEW_TEMPLATE: '<a></a>',
	VIEW_CLASS: 'trigger_view',
	DELETE_TEMPLATE: '<a></a>',
	DELETE_CLASS: 'trigger_delete',
	INPUT_NODE_CLASS: 'trigger_input',
	_selectNode: null,
	_viewNode: null,
	_deleteNode: null,
	
	_renderSelectNode : function () {
        this._selectNode = this._renderNode(this.SELECT_TEMPLATE, this.SELECT_CLASS);
    },
    
    _renderViewNode : function () {
        this._viewNode = this._renderNode(this.VIEW_TEMPLATE, this.VIEW_CLASS);
    },
    
    _renderDeleteNode : function () {
        this._deleteNode = this._renderNode(this.DELETE_TEMPLATE, this.DELETE_CLASS);
    },
	
	renderUI: function() {
		Y.RTriggerField.superclass.renderUI.apply(this, arguments);
		
		this._renderSelectNode();
		var multi = this._getBooleanOrFunctionResult(this.get("multi"));
		if (!multi) {
			this._renderViewNode();
		}
		this._renderDeleteNode();
    },
    
    _getStringOrFunctionResult: function(val){
    	if (val) {
    		if (Y.Lang.isFunction(val)) {
    			return val();
    		}
    		return val;
    	}
    	return "";
    },
    
    _getBooleanOrFunctionResult: function(val){
		if (Y.Lang.isFunction(val)) {
			return val();
		}
		return val;
    },
    
    bindUI: function() {
    	Y.RTriggerField.superclass.bindUI.apply(this, arguments);
    	
    	this._fieldNode.detach('change');
    	this._fieldNode.detach('blur');
    	
    	this.detach('valueChange');
    	
    	this.on('valueChange', Y.bind(function(e) {
            if (e.src != 'ui') {
            	var newValue = e.newVal + "";
            	if (!newValue || parseInt(newValue, 10) <= 0) {
            		this._fieldNode.set('value', "");
            		return;
            	}
                //this._fieldNode.set('value', e.newVal + "_测试值");
            	var selectorName = this._getStringOrFunctionResult(this.get("selectorName"));
                var relationManager = new RelationManager();
                var li = newValue.split(",");
                var valueLi = [];
                for (var i = 0; i < li.length; i++) {
                	var value = "";
                	var relationItem = relationManager.getRelationBo(selectorName, li[i]);
                	if (!relationItem) {
                		Y.log("selectorName:" + selectorName + ", id:" + li[i] + " can't found relationItem");
                	}
                	var displayField = this._getStringOrFunctionResult(this.get("displayField"));
                	if (displayField.indexOf("{") > -1) {
                		value = Y.Lang.sub(displayField, relationItem);
                	} else {
                		var keyLi = displayField.split(',');
                		for (var j = 0; j < keyLi.length; j++) {
                			if (relationItem[keyLi[j]]) {
                				value += relationItem[keyLi[j]] + ",";
                			}
                		}
                		if (value) {
                			value = value.substr(0, value.length - 1);
                		}
                	}
                	valueLi.push(value);
                }
            	this._fieldNode.set('value', valueLi.join(";"));
            }
        },
        this));
    	
    	this._selectNode.on("click", Y.bind(function(e) {
    		var self = this;
        	window.s_selectFunc = function(selectValueLi) {
        		var selection = self.get("selectFunc");
        		if (selection) {
        			selection(selectValueLi, self);
        		}
        	};
        	window.s_queryFunc = function() {
        		var queryFunc = self.get("queryFunc");
        		if (queryFunc) {
        			return queryFunc();
        		}
        		return {};
        	};
        	
            var url = "/console/selectorschema?@name={NAME_VALUE}&@id={ID_VALUE}&@multi={MULTI_VALUE}&@displayField={DISPLAY_FIELD_VALUE}&date=" + new Date();
            var selectorName = this._getStringOrFunctionResult(this.get("selectorName"));
            if (!selectorName || selectorName == "NullSelector") {
            	showAlert("无法打开选择器");
            } else {
            	url = url.replace("{NAME_VALUE}", selectorName);
            	url = url.replace("{ID_VALUE}", this.get('value'));
            	var multi = this._getBooleanOrFunctionResult(this.get("multi"));
            	url = url.replace("{MULTI_VALUE}", multi);
            	var displayField = this._getStringOrFunctionResult(this.get("displayField"));
            	url = url.replace("{DISPLAY_FIELD_VALUE}", displayField);
            	var selectorTitle = this._getStringOrFunctionResult(this.get("selectorTitle"));
            	var dialog = showModalDialog({
            		"title": selectorTitle,
            		"url": url
            	});
            	window.s_closeDialog = function() {
            		if (window.s_dialog) {
            			window.s_dialog.hide();
            		}
            		window.s_dialog = null;
            		window.s_selectFunc = null;
            		window.s_queryFunc = null;
            	}
            }
    	}, this));
    	
    	var multi = this._getBooleanOrFunctionResult(this.get("multi"));
    	if (!multi) {
    		this._viewNode.on("click", Y.bind(function(e) {
    			var value = this.get("value");
    			if (!value || value == "0") {
    				showAlert("没有数据，无法查看详情");
    			} else {
    				var selectorName = this._getStringOrFunctionResult(this.get("selectorName"));
    				var relationManager = new RelationManager();
    				var relationItem = relationManager.getRelationBo(selectorName, value);
    				var url = g_relationBo[selectorName]["url"] || "";
    				url = Y.Lang.sub(url, relationItem);
    				if (url) {
    					var selectorTitle = this._getStringOrFunctionResult(this.get("selectorTitle"));
    					triggerShowModalDialog({
    						"title": selectorTitle,
    						"url": url
    					});
    				} else {
    					showAlert("url为空，无法打开详情页面");
    				}
    			}
    		}, this));
    	}
    	
    	this._deleteNode.on("click", Y.bind(function(e) {
    		var self = this;
    		var unSelectFunc = self.get("unSelectFunc");
    		if (unSelectFunc) {
    			unSelectFunc(self);
    		}
    	}, this));
    },
    
    _syncReadonly: function(e) {
    	//Y.PDateField.superclass._syncReadonly.apply(this, arguments);
    	
    	var value = this.get('readonly');
        if (value === true) {
        	this._selectNode.setStyle("display", "none");
        	this._deleteNode.setStyle("display", "none");
        	this._fieldNode.addClass('readonly');
        } else {
        	this._selectNode.setStyle("display", "");
        	this._deleteNode.setStyle("display", "");
        	this._fieldNode.removeClass('readonly');
        }
    },
    
    _syncSelectNode: function() {
    	if (this._selectNode) {
    		this._selectNode.setAttrs({
    			href: "javascript:void(0);",
    			title: "多选",
    			id: Y.guid() + Y.RFormField.FIELD_ID_SUFFIX
    		});
    	}
    },
    
    _syncViewNode: function() {
    	if (this._viewNode) {
    		this._viewNode.setAttrs({
    			href: "javascript:void(0);",
    			title: "查看",
    			id: Y.guid() + Y.RFormField.FIELD_ID_SUFFIX
    		});
    	}
    },
    
    _syncDeleteNode: function() {
    	if (this._deleteNode) {
    		this._deleteNode.setAttrs({
    			href: "javascript:void(0);",
    			title: "删除",
    			id: Y.guid() + Y.RFormField.FIELD_ID_SUFFIX
    		});
    	}
    },
    
    _syncFieldNode: function() {
    	Y.RTriggerField.superclass._syncFieldNode.apply(this, arguments);
    	
    	this._fieldNode.setAttribute("readonly", "readonly");
    },
    
    syncUI: function() {
    	Y.RTriggerField.superclass.syncUI.apply(this, arguments);
    	
    	this._syncSelectNode();
    	this._syncViewNode();
    	this._syncDeleteNode();
    },
    
    initializer: function() {
    	Y.RTriggerField.superclass.initializer.apply(this, arguments);
    }

},
{
    ATTRS: {
        /** 
         * @attribute multi
         * @type Boolean
         * @default false
         * @description Set to true to allow multiple values to be selected
         */
        multi: {
        	value: false,
        	validator: function(val) {
        		if (!(Y.Lang.isBoolean(val) || Y.Lang.isFunction(val))){
        			return false;
        		}
        		return true;
        	}
        },
        value: {
            value: '',
            validator: function(val) {
                if (!(Y.Lang.isString(val) || Y.Lang.isNumber(val))){
                	return false;
                }
                if (!val) {
                	return true;
                }
                val = val + "";
                var selectorName = this._getStringOrFunctionResult(this.get("selectorName"));
                var relationManager = new RelationManager();
                var li = val.split(",");
                for (var i = 0; i < li.length; i++) {
                	if (li[i] != "0") {
                		var g_relationBo = relationManager.getRelationBo(selectorName, li[i]);
                		if (!g_relationBo) {
                			return false;
                		}
                	}
                }
                return true;
            }
        },
        selectorName: {
        	value: '',
        	validator: function(val) {
        		if (!(Y.Lang.isString(val) || Y.Lang.isFunction(val))){
        			return false;
        		}
        		return true;
        	}
        },
        displayField: {
        	value: '',
        	validator: function(val) {
        		if (!(Y.Lang.isString(val) || Y.Lang.isFunction(val))){
        			return false;
        		}
        		return true;
        	}
        },
        valueField: {
        	value: 'id',
        	validator: function(val) {
        		if (!(Y.Lang.isString(val) || Y.Lang.isFunction(val))){
        			return false;
        		}
        		return true;
        	}
        },
        selectFunc: {
        	value: null,
        	validator: Y.Lang.isFunction
        },
        unSelectFunc: {
        	value: null,
        	validator: Y.Lang.isFunction
        },
        queryFunc: {
        	value: null,
        	validator: Y.Lang.isFunction
        },
        selectorTitle: {
        	value: '',
        	validator: function(val) {
        		if (!(Y.Lang.isString(val) || Y.Lang.isFunction(val))){
        			return false;
        		}
        		return true;
        	}
        }
    }
});
