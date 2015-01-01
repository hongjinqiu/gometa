function ColumnDataSourceManager() {
}
/*
ColumnDataSourceManager.prototype.createVirtualColumn = function(columnModelName, columnModel, columnIndex) {
	var self = this;
	var yInst = self.yInst;
	var i = columnIndex;
	if (columnModel.ColumnLi[i].XMLName.Local == "virtual-column" && columnModel.ColumnLi[i].Hideable != "true") {
		var virtualColumn = columnModel.ColumnLi[i];
		return {
			key: columnModel.ColumnLi[i].Name,
			label: columnModel.ColumnLi[i].Text,
			allowHTML:  true, // to avoid HTML escaping
			formatter:      function(virtualColumn){
				return function(o){
					var htmlLi = [];
					var buttonBoLi = null;
					if (o.value) {
						buttonBoLi = o.value[virtualColumn.Buttons.XMLName.Local];
					}
					for (var j = 0; j < virtualColumn.Buttons.ButtonLi.length; j++) {
						var btnTemplate = null;
						if (virtualColumn.Buttons.ButtonLi[j].Mode == "fn") {
							btnTemplate = "<input type='button' value='{value}' onclick='doVirtualColumnBtnAction(\"{columnModelName}\", this, {handler})' class='{class}' />";
						} else if (virtualColumn.Buttons.ButtonLi[j].Mode == "url") {
							btnTemplate = "<input type='button' value='{value}' onclick='location.href=\"{href}\"' class='{class}' />";
						} else {
							btnTemplate = "<input type='button' value='{value}' onclick='window.open(\"{href}\")' class='{class}' />";
						}
						if (!buttonBoLi || buttonBoLi[j]["isShow"]) {
							var id = columnModel.IdColumn.Name;
							var isUsed = g_usedCheck && g_usedCheck[columnModel.DataSetId] && g_usedCheck[columnModel.DataSetId][o.data[id]];
							if (!(isUsed && virtualColumn.Buttons.ButtonLi[j].Name == "btn_delete")) {
								// handler进行值的预替换,
								var Y = yInst;
								var handler = virtualColumn.Buttons.ButtonLi[j].Handler;
								handler = Y.Lang.sub(handler, o.data);
								htmlLi.push(Y.Lang.sub(btnTemplate, {
									value: virtualColumn.Buttons.ButtonLi[j].Text,
									handler: handler,
									"class": virtualColumn.Buttons.ButtonLi[j].IconCls,
									href: handler,
									columnModelName: columnModelName
								}));
							}
						}
					}
					return htmlLi.join("");
				}
			}(virtualColumn)
		};
	}
	return null;
}*/

ColumnDataSourceManager.prototype.getColumns = function(columnModelName, columnModel, Y) {
	var self = this;
	self.yInst = Y;
	var columnManager = new ColumnManager();
	// 换掉createVirtualColumn,button需要做被用的判断,被用时,不显示删除按钮,
	//columnManager.createVirtualColumn = self.createVirtualColumn;
	var columns = columnManager.getColumns(columnModelName, columnModel, Y);
	if (g_dataSourceJson) {
		var modelIterator = new ModelIterator();
		var result = "";
		for (var i = 0; i < columns.length; i++) {
			columns[i].allowHTML = true;
			modelIterator.iterateAllField(g_dataSourceJson, result, function(fieldGroup, result){
				if (fieldGroup.getDataSetId() == columnModel.DataSetId && fieldGroup.Id == columns[i].key) {
					if (fieldGroup.AllowEmpty == "false") {
						columns[i].label = '<font style="color:red">*</font>' + columns[i].label;
					}
				}
			});
		}
	}
	return columns
}
