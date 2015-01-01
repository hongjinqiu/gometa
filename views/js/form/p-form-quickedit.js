YUI.add('papersns-form-quickedit', function(Y, NAME) {
	function PFormQuickEdit(config) {
		PFormQuickEdit.superclass.constructor.call(this, config);
	}

	PFormQuickEdit._incId = 0;

	PFormQuickEdit.NAME = "PFormQuickEditPlugin";
	PFormQuickEdit.NS = "pqe";

	PFormQuickEdit.ATTRS = {}

	PFormQuickEdit.getInnerId = function() {
		return ++PFormQuickEdit._incId;
	}

	Y.extend(PFormQuickEdit, Y.Plugin.Base, {
		initializer : function(config) {
			var self = this;
			var host = this.get('host');

			var columns = host.get("columns");
			Y.each(columns, function(rec, i) {
				if (self._isVirtualColumn(host.dataSetId, rec.key)) {
					// 查询出 forEditor的 column,换掉
					var templateIterator = new TemplateIterator();
					var result = "";
					var virtualColumnForEditor = null;
					templateIterator.iterateAnyTemplateColumn(host.dataSetId, result, function IterateFunc(column, result) {
						if (column.XMLName.Local == "virtual-column" && column.UseIn == "editor") {
							virtualColumnForEditor = column;
							return true;
						}
						return false;
					});
					rec.allowHTML = true;
					if (virtualColumnForEditor != null) {
						rec.width = virtualColumnForEditor.Width;
						rec.formatter = function(o) {
							var bodyHtmlLi = [];
							for (var j = 0; j < virtualColumnForEditor.Buttons.ButtonLi.length; j++) {
								var btnTemplate = null;
								if (virtualColumnForEditor.Buttons.ButtonLi[j].Mode == "fn") {
									btnTemplate = "<a title='{value}' onclick='doPluginVirtualColumnBtnAction(\"{columnModelName}\", this, {handler})' class='{class}' href='javascript:void(0);' style='display:block;' />";
								} else if (virtualColumnForEditor.Buttons.ButtonLi[j].Mode == "url") {
									btnTemplate = "<a title='{value}' onclick='location.href=\"{href}\"' class='{class}' href='javascript:void(0);' style='display:block;' />";
								} else {
									btnTemplate = "<a title='{value}' onclick='window.open(\"{href}\")' class='{class}' href='javascript:void(0);' style='display:block;' />";
								}
								var handler = virtualColumnForEditor.Buttons.ButtonLi[j].Handler;
								handler = Y.Lang.sub(handler, o.data);
								bodyHtmlLi.push(Y.Lang.sub(btnTemplate, {
									value: virtualColumnForEditor.Buttons.ButtonLi[j].Text,
									handler: handler,
									"class": virtualColumnForEditor.Buttons.ButtonLi[j].IconCls,
									href: handler,
									columnModelName: host.dataSetId
								}));
							}
							return bodyHtmlLi.join("");
						}
					} else {
						rec.width = "60";
						rec.formatter = function(o) {
							var bodyHtmlLi = [];
							var btnTemplate = "<a title='{value}' onclick='doPluginVirtualColumnBtnAction(\"{columnModelName}\", this, {handler})' class='{class}' href='javascript:void(0);' style='display:block;' />";
							bodyHtmlLi.push(Y.Lang.sub(btnTemplate, {
								value: "复制",
								handler: "g_pluginCopyRow",
								"class": "img_add",
								columnModelName: host.dataSetId
							}));
							bodyHtmlLi.push(Y.Lang.sub(btnTemplate, {
								value: "删除",
								handler: "g_pluginRemoveSingleRow",
								"class": "img_delete",
								columnModelName: host.dataSetId
							}));
							return bodyHtmlLi.join("");
						}
					}
				} else if (!self._isSkip(host.dataSetId, rec.key)) {
					rec.allowHTML = true;
					rec.formatter = function(o) {
						var id = PFormQuickEdit.getInnerId();
						return "<div id='cell_" + id + "' class='pformquickedit-container pformquickedit-key:" + o.column.key + "'></div>";
					}
				}
			});
			host.set("columns", columns);
			var h = this.afterHostEvent('render', function() {
				self._addFormFieldPlugin(host);
				
				// apply model.js beforeEdit 函数
				var modelIterator = new ModelIterator();
				var result = "";
				modelIterator.iterateAllDataSet(g_dataSourceJson, result, function(dataSet, result){
					if (dataSet.Id == host.dataSetId) {
						if (dataSet.jsConfig && dataSet.jsConfig.beforeEdit) {
							// 迭代每一行 record,调用beforeEdit,
							var recordLi = host.get('data');
							recordLi.each(function(rec, recordIndex) {
								dataSet.jsConfig.beforeEdit(recordLi, rec, recordIndex);
							});
						}
					}
				});
			});
			
			self._resetHostAddRow(host);
			self._resetHostAddRows(host);
		},
		_resetHostAddRow: function(host){
			var self = this;
			host.oldHostAddRowFunc = host.addRow;
			host.addRow = function() {
				var data = null;
				if (arguments && arguments.length > 0) {
					data = arguments[0];
				}
				var config = null;
				if (arguments && arguments.length > 1) {
					config = arguments[1];
				}
				var callback = null;
				if (arguments && arguments.length > 2) {
					callback = arguments[2];
				}
				host.oldHostAddRowFunc(data, config, callback);
				
				self._addFormFieldPlugin(host);
			}
		},
		_resetHostAddRows: function(host){
			var self = this;
			host.oldHostAddRowsFunc = host.addRows;
			host.addRows = function() {
				var data = null;
				if (arguments && arguments.length > 0) {
					data = arguments[0];
				}
				var config = null;
				if (arguments && arguments.length > 1) {
					config = arguments[1];
				}
				var callback = null;
				if (arguments && arguments.length > 2) {
					callback = arguments[2];
				}
				host.oldHostAddRowsFunc(data, config, callback);
				
				self._addFormFieldPlugin(host);
			}
		},
		_addFormFieldPlugin: function(host){
			var self = this;
			var dataSetId = host.dataSetId;
			var formFieldFactory = new FormFieldFactory();
			var rows = host._tbodyNode.get('children');
			var templateIterator = new TemplateIterator();
			var formFieldFactory = new FormFieldFactory();
			
			var columnResult = "";
			var columnLi = [];
			templateIterator.iterateAllTemplateColumn(dataSetId, columnResult, function IterateFunc(column, result) {
				columnLi.push(column);
			});
			var columnSequenceService = new ColumnSequenceService();
			var sequenceColumnLi = columnSequenceService.buildSequenceColumnLi(columnLi);
			// 给其加上是否在当前数据集的属性
			var modelIterator = new ModelIterator();
			var modelResult = "";
			var modelFieldLi = [];
			modelIterator.iterateAllField(g_dataSourceJson, modelResult, function(fieldGroup, result){
				if (fieldGroup.getDataSetId() == dataSetId) {
					modelFieldLi.push(fieldGroup);
				}
			});
			for (var i = 0; i < sequenceColumnLi.length; i++) {
				for (var j = 0; j < modelFieldLi.length; j++) {
					if (sequenceColumnLi[i].Name == modelFieldLi[j].Id) {
						sequenceColumnLi[i].isDsField = true;
						break;
					}
				}
			}
			
			host.get('data').each(function(rec, recordIndex) {
				if (!rec.formFieldLi) {
					rec.formFieldLi = [];
					if (!rec.formFieldDict) {
						rec.formFieldDict = {};
					}
					//IdColumn,
					var iterateResult = "";
					templateIterator.iterateAnyTemplateColumnModel(iterateResult, function IterateFunc(columnModel, result) {
						if (columnModel.DataSetId == dataSetId) {
							//columnModel.IdColumn
							var field = formFieldFactory.getFormField(Y, columnModel.IdColumn.Name, dataSetId);
							field.set("value", rec.get(columnModel.IdColumn.Name));
							rec.formFieldLi.push(field);
							rec.formFieldDict[columnModel.IdColumn.Name] = field;
							return true;
						}
						return false;
					});
					// 生成field,再进行赋值,此时,copyConfig才能起作用
					// 生成HiddenColumn,
					var iterateResult = "";
					templateIterator.iterateAllTemplateColumn(dataSetId, iterateResult, function IterateFunc(column, result) {
						if (column.Hideable == "true") {
							var field = formFieldFactory.getFormField(Y, column.Name, dataSetId);
							rec.formFieldLi.push(field);
							rec.formFieldDict[column.Name] = field;
						}
					});
					
					// 生成非HiddenColumn
					var list = rows.item(recordIndex).all('.pformquickedit-container');
					var field_count = list.size();
					for ( var j = 0; j < field_count; j++) {
						var fieldName = self._getColumnKey(list.item(j));
						
						var field = formFieldFactory.getFormField(Y, fieldName, dataSetId);
						field.render("#" + list.item(j).get("id"));
						rec.formFieldLi.push(field);
						rec.formFieldDict[fieldName] = field;
					}
					
					for (var i = 0; i < sequenceColumnLi.length; i++) {
						if (sequenceColumnLi[i].isDsField) {// 对于非数据源模型中的字段,当前只有复制出来的字段,这种字段不管其,select-column set value会自动带出
							var columnName = sequenceColumnLi[i].Name;
							var field = rec.formFieldDict[columnName];
							if (rec.get(columnName)) {
								field.set("value", rec.get(columnName));
							} else {
								field.set("value", "");
							}
						}
					}
				}
			});
		},
		_isVirtualColumn: function(dataSetId, name) {
			var isVirtualColumn = false;
			var templateIterator = new TemplateIterator();
			var result = "";
			templateIterator.iterateAnyTemplateColumn(dataSetId, result, function IterateFunc(column, result) {
				if (column.Name == name) {
					if (column.XMLName.Local == "virtual-column") {
						isVirtualColumn = true;
					}
					return true;
				}
				return false;
			});
			
			return isVirtualColumn;
		},
		_isHidden: function(dataSetId, name) {
			var isHidden = false;
			var templateIterator = new TemplateIterator();
			var result = "";
			templateIterator.iterateAnyTemplateColumn(dataSetId, result, function IterateFunc(column, result) {
				if (column.Name == name) {
					if (column.Hideable == "true") {
						isHidden = true;
					}
					return true;
				}
				return false;
			});
			
			return isHidden;
		},
		_isSkip: function(dataSetId, name) {
			var self = this;
			if (self._isHidden(dataSetId, name)) {
				return true;
			}
			
			var isSkip = false;
			var templateIterator = new TemplateIterator();
			var result = "";
			templateIterator.iterateAnyTemplateColumnModel(result, function IterateFunc(columnModel, result) {
				if (columnModel.DataSetId == dataSetId) {
					if (columnModel.IdColumn.Name == name) {
						isSkip = true;
					} else if (columnModel.CheckboxColumn.Name == name) {
						isSkip = true;
					}
					return true;
				}
				return false;
			});
			
			return isSkip;
		},
		getRecords : function() {
			var self = this;
			var host = this.get('host');
			var result = [];
			host.get("data").each(function(rec, i) {
				var record = {};
				for ( var j = 0; j < rec.formFieldLi.length; j++) {
					var formField = rec.formFieldLi[j];
					record[formField.get("name")] = formField.get("value");
				}
				result.push(record);
			});
			return result;
		},
		_getColumnKey : function(e) {
			var quick_edit_re = /pformquickedit-key:([^\s]+)/;
			var m = quick_edit_re.exec(e.get('className'));
			return m[1];
		}
	});

	Y.namespace("Plugin");
	Y.Plugin.DataTablePFormQuickEdit = PFormQuickEdit;
}, '1.1.0', {
	"requires" : [ "datatable-base", "papersns-form" ]
});
