DataTableManager.prototype.createAddRowGrid = function(inputDataLi) {
	var self = this;
	executeGYUI(
			 function(Y) {
				var pluginDataTableManager = new DataTableManager();
				var doPopupConfirm = function() {
					var li = pluginDataTableManager.dt.pqe.getRecords();
					// 输入中有,输出中没有,删除
					for (var i = 0; i < inputDataLi.length; i++) {
						var isIn = false;
						for (var j = 0; j < li.length; j++) {
							if (inputDataLi[i].id == li[j].id) {
								isIn = true;
								break;
							}
						}
						if (!isIn) {
							self.dt.removeRow(inputDataLi[i].id);
						}
					}
					
					for (var i = 0; i < li.length; i++) {
						var record = self.dt.getRecord(li[i].id);
						if (record) {
							for (var key in li[i]) {
								record.set(key, li[i][key]);
							}
						}
					}
					self.dt.addRows(li);
				};
				var bodyHtmlLi = [];
				bodyHtmlLi.push("<div class='alignLeft formToolbar'>");
				
				if (self.param.columnModel.EditorToolbar && self.param.columnModel.EditorToolbar.ButtonLi) {
					for (var i = 0; i < self.param.columnModel.EditorToolbar.ButtonLi.length; i++) {
						var btnTemplate = null;
						var replObj = {};
						var button = self.param.columnModel.EditorToolbar.ButtonLi[i];
						if (button.Mode == "fn") {
							btnTemplate = "<input type='button' value='{value}' class='{class}' onclick='{fnName}(\"{columnModelName}\")'/>";
							replObj = {
								value: button.Text,
								"class": button.IconCls,
								fnName: button.Handler,
								columnModelName: self.param.columnModelName
							};
						} else if (button.Mode == "url") {
							btnTemplate = "<input type='button' value='{value}' onclick='location.href=\"{href}\"' class='{class}' />";
							replObj = {
								value: button.Text,
								href: button.Handler,
								"class": button.IconCls
							};
						} else {
							btnTemplate = "<input type='button' value='{value}' onclick='window.open(\"{href}\")' class='{class}' />";
							replObj = {
								value: button.Text,
								href: button.Handler,
								"class": button.IconCls
							};
						}
						if (btnTemplate) {
							btnTemplate = Y.Lang.sub(btnTemplate, replObj);
							bodyHtmlLi.push(btnTemplate);
						}
					}
				}
				bodyHtmlLi.push("</div>");
				bodyHtmlLi.push('<div style="overflow: auto" id="' + self.param.columnModelName + "_addrow" + '"></div>');

				var node = Y.one("window");
				var width = parseInt(node.get("winWidth"), 10);
				var edge = 100;
				var dialogWidth = width - edge;
				if (dialogWidth <= 0) {
					dialogWidth = 100;
				}
				var dialog = new Y.Panel({
					contentBox : Y.Node.create('<div id="detail-grid-addrow-dialog" />'),
					headerContent : "新增" + self.param.columnModel.Text,
					bodyContent : bodyHtmlLi.join(""),
					width : dialogWidth,
					zIndex     : (++panelZIndex),
					centered : true,
					modal : true, // modal behavior
					render : '.example',
					visible : false, // make visible explicitly with .show()
					plugins : [ Y.Plugin.Drag ],
					buttons : {
						footer : [ {
							name : 'cancel',
							label : '取消',
							action : 'onCancel',
							classNames: 'message_bt1'
						}, {
							name : 'proceed',
							label : '确定',
							action : 'onOK',
							classNames: 'message_bt1'
						} ]
					}
				});

				dialog.onCancel = function(e) {
					e.preventDefault();
					this.hide();
					// the callback is not executed, and is
					// callback reference removed, so it won't persist
					this.callback = false;
				}

				dialog.onOK = function(e) {
					e.preventDefault();
					
					// 数据较验
					var formManager = new FormManager();
					var detailDataLi = pluginDataTableManager.dt.pqe.getRecords();
					var dataSetId = self.param.columnModelName;
					var validateResult = formManager.dsDetailValidator(g_dataSourceJson, dataSetId, detailDataLi);
					
					if (!validateResult.result) {
						showError(validateResult.message);
					} else {
						this.hide();
						// code that executes the user confirmed action goes here
						if (this.callback) {
							this.callback();
						}
						// callback reference removed, so it won't persist
						this.callback = false;
					}
				}

				dialog.hide = function() {
					g_gridPanelDict[self.param.columnModelName + "_addrow"] = null;
					return this.destroy();
				}

				dialog.dd.addHandle('.yui3-widget-hd');
				//		Y.one('#dialog .message').setHTML('mnopq Are you sure you want to [take some action]?');
				//		Y.one('#dialog .message').set('className', 'message icon-bubble');
				dialog.callback = doPopupConfirm;
				var data = [{}];
				if (inputDataLi) {
					data = inputDataLi;
				}
				var param = {
					data : data,
					columnModel : self.param.columnModel,
					columnModelName : self.param.columnModelName + "_addrow",// 用于virtualColumn的btn,onclick,回找grid的,暂时没用,
					render : "#" + self.param.columnModelName + "_addrow",// 用panel里面的东东,
					url : "",
					totalResults : g_dataBo.totalResults || 50,
					pageSize : 10000,
					paginatorContainer : null,
					paginatorTemplate : null,
					columnManager : new ColumnDataSourceManager(),
					plugin : Y.Plugin.DataTablePFormQuickEdit
				};

				g_gridPanelDict[self.param.columnModelName + "_addrow"] = pluginDataTableManager;// 这一行要放在createDataGrid之前,因为里面的selectField会触发valueChange,应用到copyField里,需要从全局的g_gridPanelDict里面查找formFieldDict,
				pluginDataTableManager.createDataGrid(Y, param);
				dialog.show();
			});
}

function doPluginVirtualColumnBtnAction(columnModelName, elem, fn){
	var self = g_gridPanelDict[columnModelName + "_addrow"];
	var dt = self.dt;
	var yInst = self.yInst;
	var o = dt.getRecord(yInst.one(elem));
	fn(o, columnModelName);
}

/**
 * 插件表头新增,新增一行
 */
function g_pluginAddRow(dataSetId) {
	var formManager = new FormManager();
	var data = formManager.getDataSetNewData(dataSetId);
	g_gridPanelDict[dataSetId + "_addrow"].dt.addRow(data);
}

/**
 * 插件表头删除,删除多行
 */
function g_pluginRemoveRow(dataSetId) {
	var selectRecordLi = g_gridPanelDict[dataSetId + "_addrow"].getSelectRecordLi();
	if (selectRecordLi.length == 0) {
		showAlert("请先选择");
	} else {
		for (var i = 0; i < selectRecordLi.length; i++) {
			g_gridPanelDict[dataSetId + "_addrow"].dt.removeRow(selectRecordLi[i]);
		}
	}
}

/**
 * 点击删除,删除一行
 */
function g_pluginRemoveSingleRow(o, dataSetId) {
	g_gridPanelDict[dataSetId + "_addrow"].dt.removeRow(o);
}

/**
 * 点击行项复制,复制一行
 */
function g_pluginCopyRow(o, dataSetId) {
//	var inputDataLi = [];
	var formManager = new FormManager();
	var id = o.get("id");
	var li = g_gridPanelDict[dataSetId + "_addrow"].dt.pqe.getRecords();
	var data = {};
	for (var i = 0; i < li.length; i++) {
		if (li[i].id == id) {
			data = li[i];
			break;
		}
	}
	var data = formManager.getDataSetCopyData(dataSetId, data);
//	inputDataLi.push(data);
	g_gridPanelDict[dataSetId + "_addrow"].dt.addRow(data);
}

/*
	"buttonConfig": {
		"selectRowBtn": {
			selectFunc: function(datas){},// 单多选回调
			queryFunc: function(){}// 单多选回调
		}
	}
 */
/**
 * 点击选择,选择回多行,需要从全局配置中读queryFunc,selectFunc是这里面应用上的
 * @param dataSetId
 */
function g_selectRow(dataSetId, btnName) {
//	var formManager = new FormManager();
	var templateIterator = new TemplateIterator();
	var result = "";
	templateIterator.iterateAnyTemplateButton(result, function(toolbarOrColumnModel, button, result) {
		if (toolbarOrColumnModel.DataSetId == dataSetId && button.Name == btnName) {
			window.s_selectFunc = function(selectValueLi) {
				if (button.jsConfig && button.jsConfig.selectFunc) {
					button.jsConfig.selectFunc(selectValueLi);
				} else {
					selectRowBtnDefaultAction(dataSetId, toolbarOrColumnModel, button, selectValueLi);
				}
        	};
        	window.s_queryFunc = function() {
        		var queryFunc = null;
        		if (button.jsConfig && button.jsConfig.queryFunc) {
        			queryFunc = button.jsConfig.queryFunc;
        		}
        		if (queryFunc) {
        			return queryFunc();
        		}
        		return {};
        	};
        	if (button.CRelationDS && button.CRelationDS.CRelationItemLi) {
    			var relationItem = button.CRelationDS.CRelationItemLi[0];
    			var url = "/console/selectorschema?@name={NAME_VALUE}&@multi={MULTI_VALUE}&@displayField={DISPLAY_FIELD_VALUE}&date=" + new Date();
    			var selectorName = relationItem.CRelationConfig.SelectorName;
    			url = url.replace("{NAME_VALUE}", selectorName);
    			var multi = relationItem.CRelationConfig.SelectionMode == "multi";
    			url = url.replace("{MULTI_VALUE}", multi);
    			var displayField = relationItem.CRelationConfig.DisplayField;
    			url = url.replace("{DISPLAY_FIELD_VALUE}", displayField);
    			var selectorTitle = g_relationBo[selectorName].Description;
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
        	
			return true;
		}
		return false;
	});
}

/**
 * 点击新增,新增一行
 */
function g_addRow(dataSetId) {
	var inputDataLi = [];
	var formManager = new FormManager();
	var data = formManager.getDataSetNewData(dataSetId);
	inputDataLi.push(data);
	g_gridPanelDict[dataSetId].createAddRowGrid(inputDataLi);
}

/**
 * 点击删除,删除多行
 */
function g_removeRow(dataSetId) {
	var selectRecordLi = g_gridPanelDict[dataSetId].getSelectRecordLi();
	if (selectRecordLi.length == 0) {
		showAlert("请先选择");
	} else {
		var hasUsed = false;
		for (var i = 0; i < selectRecordLi.length; i++) {
			var isUsed = g_usedCheck && g_usedCheck[dataSetId] && g_usedCheck[dataSetId][selectRecordLi[i].get("id")];
			if (isUsed) {
				hasUsed = true;
			} else {
				g_gridPanelDict[dataSetId].dt.removeRow(selectRecordLi[i]);
			}
		}
		if (hasUsed) {
			showAlert("部门数据已被用，不可删除！");
		}
	}
}

/**
 * 点击删除,删除一行
 */
function g_removeSingleRow(o, dataSetId) {
	g_gridPanelDict[dataSetId].dt.removeRow(o);
}

/**
 * 点击行项编辑,编辑一行
 */
function g_editSingleRow(o, dataSetId) {
	var inputDataLi = [];
	inputDataLi.push(o.toJSON());
	g_gridPanelDict[dataSetId].createAddRowGrid(inputDataLi);
}

/**
 * 点击行项复制,复制一行
 */
function g_copyRow(o, dataSetId) {
	var inputDataLi = [];
	var formManager = new FormManager();
	var data = formManager.getDataSetCopyData(dataSetId, o.toJSON());
	inputDataLi.push(data);
	g_gridPanelDict[dataSetId].createAddRowGrid(inputDataLi);
}

/**
 * 点击表格头编辑,编辑多行
 */
function g_editRow(dataSetId) {
	var selectRecordLi = g_gridPanelDict[dataSetId].getSelectRecordLi();
	if (selectRecordLi.length == 0) {
		showAlert("请先选择");
	} else {
		var inputDataLi = [];
		for (var i = 0; i < selectRecordLi.length; i++) {
			inputDataLi.push(selectRecordLi[i].toJSON());
		}
		g_gridPanelDict[dataSetId].createAddRowGrid(inputDataLi);
	}
}



