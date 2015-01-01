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
					rec.allowHTML = true;
					rec.formatter = function(o) {
						var bodyHtmlLi = [];
						var btnTemplate = "<input type='button' value='{value}' onclick='doPluginVirtualColumnBtnAction(\"{columnModelName}\", this, {handler})' class='{class}' />";
						bodyHtmlLi.push(Y.Lang.sub(btnTemplate, {
							value: "复制",
							handler: "g_pluginCopyRow",
							"class": "test",
							columnModelName: host.dataSetId
						}));
						bodyHtmlLi.push(Y.Lang.sub(btnTemplate, {
							value: "删除",
							handler: "g_pluginRemoveSingleRow",
							"class": "test",
							columnModelName: host.dataSetId
						}));
//						console.log(Y.Lang.sub(btnTemplate, {
//							value: "删除",
//							handler: "g_pluginRemoveSingleRow",
//							"class": "test",
//							columnModelName: host.dataSetId
//						}));
//						bodyHtmlLi.push("<input type='button' value='复制' class='' onclick='showAlert(\"复制test\")'/>");
//						bodyHtmlLi.push("<input type='button' value='删除' class='' onclick='showAlert(\"删除test\")'/>");
						return bodyHtmlLi.join("");
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
				console.log("addRows");
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
			host.get('data').each(function(rec, recordIndex) {
				if (!host.getRecord(recordIndex).formFieldLi) {
					host.getRecord(recordIndex).formFieldLi = [];
					if (!host.getRecord(recordIndex).formFieldDict) {
						host.getRecord(recordIndex).formFieldDict = {};
					}
					//IdColumn,
					var iterateResult = "";
					templateIterator.iterateAnyTemplateColumnModel(iterateResult, function IterateFunc(columnModel, result) {
						if (columnModel.DataSetId == dataSetId) {
							//columnModel.IdColumn
							var field = formFieldFactory.getFormField(Y, columnModel.IdColumn.Name, dataSetId);
							field.set("value", rec.get(columnModel.IdColumn.Name));
							host.getRecord(recordIndex).formFieldLi.push(field);
							host.getRecord(recordIndex).formFieldDict[columnModel.IdColumn.Name] = field;
							return true;
						}
						return false;
					});
					//HiddenColumn,
					var iterateResult = "";
					templateIterator.iterateAllTemplateColumn(dataSetId, iterateResult, function IterateFunc(column, result) {
						if (column.Hideable == "true") {
							var field = formFieldFactory.getFormField(Y, column.Name, dataSetId);
							field.set("value", rec.get(column.Name));
							host.getRecord(recordIndex).formFieldLi.push(field);
							host.getRecord(recordIndex).formFieldDict[column.Name] = field;
						}
					});
					
					var list = rows.item(recordIndex).all('.pformquickedit-container');
					var field_count = list.size();
					for ( var j = 0; j < field_count; j++) {
						var fieldName = self._getColumnKey(list.item(j));
						
						var field = formFieldFactory.getFormField(Y, fieldName, dataSetId);
						
						field.render("#" + list.item(j).get("id"));
						if (rec.get(fieldName)) {
							field.set("value", rec.get(fieldName) + "");
						} else {
							field.set("value", "");
						}
						host.getRecord(recordIndex).formFieldLi.push(field);
						host.getRecord(recordIndex).formFieldDict[fieldName] = field;
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
				for ( var j = 0; j < host.getRecord(i).formFieldLi.length; j++) {
					var formField = host.getRecord(i).formFieldLi[j];
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
