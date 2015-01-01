function TemplateIterator() {}

TemplateIterator.prototype._iterateTemplateColumn = function(dataSetId, result, isContinue, iterateFunc) {
	var listTemplateIterator = new ListTemplateIterator();
	for (var j = 0; j < g_formTemplateJsonData.FormElemLi.length; j++) {
		var formElem = g_formTemplateJsonData.FormElemLi[j];
		if (formElem.XMLName.Local == "column-model") {
			if (formElem.ColumnModel.DataSetId == dataSetId) {
				if (formElem.ColumnModel.ColumnLi) {
					var columnLi = [];
					listTemplateIterator.recursionGetColumnItem(formElem.ColumnModel, columnLi);
					for (var k = 0; k < columnLi.length; k++) {
						var column = columnLi[k];
						var iterateResult = iterateFunc(column, result);
						if (!isContinue && iterateResult) {
							return;
						}
					}
				}
			}
		}
	}
}

function IterateFunc(column, result) {
}

TemplateIterator.prototype.iterateAllTemplateColumn = function(dataSetId, result, iterateFunc) {
	var self = this;
	var isContinue = true;
	self._iterateTemplateColumn(dataSetId, result, isContinue, iterateFunc);
}

function IterateFunc(column, result) {
}

TemplateIterator.prototype.iterateAnyTemplateColumn = function(dataSetId, result, iterateFunc) {
	var self = this;
	var isContinue = false;
	self._iterateTemplateColumn(dataSetId, result, isContinue, iterateFunc);
}

TemplateIterator.prototype._iterateTemplateColumnModel = function(result, isContinue, iterateFunc) {
	for (var j = 0; j < g_formTemplateJsonData.FormElemLi.length; j++) {
		var formElem = g_formTemplateJsonData.FormElemLi[j];
		if (formElem.XMLName.Local == "column-model") {
			var iterateResult = iterateFunc(formElem.ColumnModel, result);
			if (!isContinue && iterateResult) {
				return;
			}
		}
	}
}

function IterateFunc(columnModel, result) {
}

TemplateIterator.prototype.iterateAllTemplateColumnModel = function(result, iterateFunc) {
	var self = this;
	var isContinue = true;
	self._iterateTemplateColumnModel(result, isContinue, iterateFunc);
}

function IterateFunc(columnModel, result) {
}

TemplateIterator.prototype.iterateAnyTemplateColumnModel = function(result, iterateFunc) {
	var self = this;
	var isContinue = false;
	self._iterateTemplateColumnModel(result, isContinue, iterateFunc);
}

/**
 * 按钮的分布,toolbar/button,columnModel/toolbar,columnModel/editor-toolbar,columnModel/virtual-column/buttons/button
 */
TemplateIterator.prototype._iterateTemplateButton = function(result, isContinue, iterateFunc) {
	for (var j = 0; j < g_formTemplateJsonData.FormElemLi.length; j++) {
		var formElem = g_formTemplateJsonData.FormElemLi[j];
		if (formElem.XMLName.Local == "toolbar") {
			if (formElem.Toolbar && formElem.Toolbar.ButtonLi) {
				for (var k = 0; k < formElem.Toolbar.ButtonLi.length; k++) {
					var button = formElem.Toolbar.ButtonLi[k];
					var iterateResult = iterateFunc(formElem.Toolbar, button, result);
					if (!isContinue && iterateResult) {
						return;
					}
				}
			}
		} else if (formElem.XMLName.Local == "column-model") {
			if (formElem.ColumnModel.Toolbar && formElem.ColumnModel.Toolbar.ButtonLi) {
				for (var k = 0; k < formElem.ColumnModel.Toolbar.ButtonLi.length; k++) {
					var button = formElem.ColumnModel.Toolbar.ButtonLi[k];
					var iterateResult = iterateFunc(formElem.ColumnModel, button, result);
					if (!isContinue && iterateResult) {
						return;
					}
				}
			}
			if (formElem.ColumnModel.EditorToolbar && formElem.ColumnModel.EditorToolbar.ButtonLi) {
				for (var k = 0; k < formElem.ColumnModel.EditorToolbar.ButtonLi.length; k++) {
					var button = formElem.ColumnModel.EditorToolbar.ButtonLi[k];
					var iterateResult = iterateFunc(formElem.ColumnModel, button, result);
					if (!isContinue && iterateResult) {
						return;
					}
				}
			}
			if (formElem.ColumnModel.ColumnLi) {
				for (var k = 0; k < formElem.ColumnModel.ColumnLi.length; k++) {
					var column = formElem.ColumnModel.ColumnLi[k];
					if (column.XMLName.Local == "virtual-column") {
						if (column.Buttons && column.Buttons.ButtonLi) {
							for (var l = 0; l < column.Buttons.ButtonLi.length; l++) {
								var button = column.Buttons.ButtonLi[l];
								var iterateResult = iterateFunc(formElem.ColumnModel, button, result);
								if (!isContinue && iterateResult) {
									return;
								}
							}
						}
					}
				}
			}
		}
	}
}

function IterateFunc(toolbarOrColumnModel, button, result) {
}

TemplateIterator.prototype.iterateAllTemplateButton = function(result, iterateFunc) {
	var self = this;
	var isContinue = true;
	self._iterateTemplateButton(result, isContinue, iterateFunc);
}

function IterateFunc(toolbarOrColumnModel, button, result) {
}

TemplateIterator.prototype.iterateAnyTemplateButton = function(result, iterateFunc) {
	var self = this;
	var isContinue = false;
	self._iterateTemplateButton(result, isContinue, iterateFunc);
}
