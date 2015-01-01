function ColumnSelectorManager() {
}

ColumnSelectorManager.prototype.getColumns = function(columnModelName, columnModel, Y) {
	var self = this;
	var columnManager = new ColumnManager();
	return columnManager._getColumnsCommon(columnModelName, columnModel, Y, function(column){
		return column.UseIn == "selector";
	});
}
