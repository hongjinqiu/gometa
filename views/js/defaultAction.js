function _recurionApplyCopyField(data, columnLi, columnName, columnValue) {
	var bo = new FormManager().getBo();
	data[columnName] = columnValue;
	var commonUtil = new CommonUtil();
	for (var i = 0; i < columnLi.length; i++) {
		if (columnLi[i].Name == columnName) {
			if (columnLi[i].XMLName.Local == "select-column") {
				if (columnLi[i].CRelationDS) {
					var relationItem = commonUtil.getCRelationItem(columnLi[i].CRelationDS, bo, data);
					if (relationItem.CCopyConfigLi) {
						for (var j = 0; j < relationItem.CCopyConfigLi.length; j++) {
							var copyValueField = relationItem.CCopyConfigLi[j].CopyValueField;
							var selectorDict = g_relationManager.getRelationBo(relationItem.CRelationConfig.SelectorName, columnValue);
							if (selectorDict) {
								var copyColumnValue = selectorDict[copyValueField];
								_recurionApplyCopyField(data, columnLi, relationItem.CCopyConfigLi[j].CopyColumnName, copyColumnValue);
							}
						}
					}
				}
			}
		}
	}
}

function selectRowBtnDefaultAction(dataSetId, toolbarOrColumnModel, button, inputValueLi) {
	var selectValueLi = [];
	if (button.CRelationDS && button.CRelationDS.CRelationItemLi) {
		var relationItem = button.CRelationDS.CRelationItemLi[0];
		var selectorName = relationItem.CRelationConfig.SelectorName;
		for (var i = 0; i < inputValueLi.length; i++) {
			var selectorDict = g_relationManager.getRelationBo(selectorName, inputValueLi[i]);
			selectValueLi.push(selectorDict);
		}
	}
	
	var formManager = new FormManager();
	var templateIterator = new TemplateIterator();
	var result = "";

	// use default action
	var dataLi = [];
	var columnResult = "";
	var columnLi = [];
	templateIterator.iterateAllTemplateColumn(dataSetId, columnResult, function IterateFunc(column, result) {
		columnLi.push(column);
	});
	for (var i = 0; i < selectValueLi.length; i++) {
		var data = formManager.getDataSetNewData(dataSetId);
		if (button.CRelationDS && button.CRelationDS.CRelationItemLi) {
			var relationItem = button.CRelationDS.CRelationItemLi[0];
			if (relationItem.CCopyConfigLi) {
				for (var j = 0; j < relationItem.CCopyConfigLi.length; j++) {
					var columnName = relationItem.CCopyConfigLi[j].CopyColumnName;
					var copyValueField = relationItem.CCopyConfigLi[j].CopyValueField;
					var columnValue = selectValueLi[i][copyValueField];
					_recurionApplyCopyField(data, columnLi, columnName, columnValue);
				}
			}
		}
		dataLi.push(data);
	}
	// 允许重复的判断,
	var gridDataLi = g_gridPanelDict["B"].dt.get("data").toJSON();
	var notAllowDuplicateColumn = [];
	var modelIterator = new ModelIterator();
	modelIterator.iterateAllField(g_dataSourceJson, result, function(fieldGroup, result){
		if (fieldGroup.getDataSetId() == dataSetId && fieldGroup.AllowDuplicate == "false") {
			notAllowDuplicateColumn.push(fieldGroup.Id);
		}
	});
	for (var i = 0; i < dataLi.length; i++) {
		var isIn = false;
		for (var j = 0; j < gridDataLi.length; j++) {
			var flag = notAllowDuplicateColumn.length > 0;
			for (var k = 0; k < notAllowDuplicateColumn.length; k++) {
				flag = flag && (dataLi[i][notAllowDuplicateColumn[k]] == gridDataLi[j][notAllowDuplicateColumn[k]]);
			}
			if (flag) {
				isIn = true;
				break
			}
		}
		if (!isIn) {
			gridDataLi.push(dataLi[i]);
		}
	}
	g_gridPanelDict["B"].dt.set("data", gridDataLi);
}
