function syncSelection(Y, record) {
	// 是否添加
	var id = record["id"];
	var idInputItem = Y.one("#selectionResult .selectionItem input[value='" + id + "']");
	if (!idInputItem) {
		var tempLi = ['<div class="selectionItem">'];
		tempLi.push('<div class="left">{display}</div>');
		tempLi.push('<div class="right" onclick="removeSelection(this)"><input type="hidden" name="selectionId" value="{id}" /></div>');
		tempLi.push('</div>');
		
		var display = Y.Lang.sub(listTemplate.ColumnModel.SelectionTemplate, record);
		
		var tempContent = Y.Lang.sub(tempLi.join(""), {
			"display": display,
			"id": id
		});
		
		Y.one("#selectionResult").setHTML(Y.one("#selectionResult").getHTML() + tempContent);
	}
}

function syncSelectionWhenChangeCheckbox(Y, dataGrid, nodeLi) {
	nodeLi.each(function(node, index, nodeLi) {
		if (node.get("checked")) {
			var record = dataGrid.getRecord(node).toJSON();
			g_selectionManager.addSelectionBo(record);
			syncSelection(Y, record);
		} else {
			// 是否删除
			var id = dataGrid.getRecord(node).get("id");
			var idInputItem = Y.one("#selectionResult .selectionItem input[value='" + id + "']");
			if (idInputItem) {
				var selectionItemLi = Y.all("#selectionResult .selectionItem");
				selectionItemLi.each(function(selectionItem, selectionItemIndex, selectionItemLi){
					if (selectionItem.one("input[value='" + id + "']")) {
						Y.one("#selectionResult").removeChild(selectionItem);
					}
				});
			}
		}
	});
}

function removeSelection(elem) {
	executeGYUI(function(Y) {
		Y.one(elem).ancestor(".selectionItem").remove();
		syncCheckboxWhenChangeSelection(Y, dtInst.dt);
	});
}

function syncCheckboxWhenChangeSelection(Y, dataGrid) {
	var selectionInputValueLi = Y.all("#selectionResult .selectionItem input").get("value");
	var checkboxItemCssSelector = dtInst.getCheckboxCssSelector();
	var checkboxItemLi = yInst.all(checkboxItemCssSelector);
	checkboxItemLi.each(function(checkboxItem, index){
		var id = dataGrid.getRecord(checkboxItem).get("id");
		var isSelected = selectionInputValueLi.some(function(value){return value == id});
		if (isSelected) {
			if (!checkboxItem.get("checked")) {
				checkboxItem.set("checked", isSelected);
			}
		} else {
			if (checkboxItem.get("checked")) {
				checkboxItem.set("checked", isSelected);
			}
		}
	});
	var itemCheckLi = checkboxItemLi.get("checked");
	var checkboxAllCssSelector = dtInst.getCheckboxAllCssSelector();
	var checkAllNode = yInst.one(checkboxAllCssSelector);
	if (checkAllNode) {
		var isAllChecked = itemCheckLi.every(function(value){return value});
		if (isAllChecked) {
			if (!checkAllNode.get("checked")) {
				checkAllNode.set("checked", isAllChecked);
			}
		} else {
			if (checkAllNode.get("checked")) {
				checkAllNode.set("checked", isAllChecked);
			}
		}
	}
}

/**
 * form页面传id进来,后端将其放置到g_selectionBo中,选择器要回显这些内容放到选择区域内,并同步到选择框中
 */
function syncCallbackSelection() {
	if (g_selectionBo) {
		for (var key in g_selectionBo) {
			if (/^\d+$/g.test(key)) {
				syncSelection(dtInst.yInst, g_selectionBo[key]);
			}
		}
		syncCheckboxWhenChangeSelection(dtInst.yInst, dtInst.dt);
	}
}

function selectorMain(Y) {
	var id = listTemplate.SelectorId;
	if (!id) {
		id = listTemplate.Id;
	}
	var url = "/console/selectorschema?@name=" + id + "&format=json";
	createGridWithUrl(Y, url, {
		columnManager: new ColumnSelectorManager()
	});
//	YUI().use("node", "event", function(Y) {
//		Y.on("domready", function(e) {
	var dataGrid = dtInst.dt;
	var checkboxItemInnerCssSelector = dtInst.getCheckboxInnerCssSelector();
	var checkboxCssSelector = dtInst.getCheckboxCssSelector();
	dataGrid.delegate("click", function(e) {
		var nodeLi = yInst.all(checkboxCssSelector);
		syncSelectionWhenChangeCheckbox(Y, dataGrid, nodeLi);
	}, checkboxItemInnerCssSelector, dataGrid);
	
	var checkboxAllInnerCssSelector = dtInst.getCheckboxAllInnerCssSelector();
	dataGrid.delegate("click", function(e) {
		var nodeLi = yInst.all(checkboxCssSelector);
		syncSelectionWhenChangeCheckbox(Y, dataGrid, nodeLi);
	}, checkboxAllInnerCssSelector, dataGrid);
	
	Y.one("#confirmBtn").on("click", function(e){
//			syncCheckboxWhenChangeSelection(Y, dataGrid);
		if (parent && parent.g_relationManager) {
			var selectorId = listTemplate.SelectorId;
			if (!selectorId) {
				selectorId = listTemplate.Id;
			}
			var selectValueLi = Y.all("#selectionResult .selectionItem input").get("value");
			if (!selectValueLi || selectValueLi.length == 0) {
				//showAlert("请先选择");
				parent.s_selectFunc([]);
			} else {
				for (var i = 0; i < selectValueLi.length; i++) {
					if (g_selectionBo[selectValueLi[i]]) {
						parent.g_relationManager.addRelationBo(selectorId, g_selectionBo["url"], g_selectionBo[selectValueLi[i]]);
					}
				}
				parent.s_selectFunc(selectValueLi);
			}
			parent.s_closeDialog();
		} else {
			alert("找不到父窗口，无法赋值！");
		}
	});
	Y.one("#clearBtn").on("click", function(e){
		Y.one("#selectionResult").setHTML("");
		syncCheckboxWhenChangeSelection(Y, dtInst.dt);
	});
	// 取得父函数的queryFunc,并设置到页面上的hidden field里面,最后调用refresh,应用这些参数查询数据
	if (parent && parent.s_queryFunc) {
		var queryDict = parent.s_queryFunc();
		for (var key in queryDict) {
			if (g_masterFormFieldDict[key]) {
				g_masterFormFieldDict[key].set("value", queryDict[key]);
			}
		}
	}
	// 同步g_selectionBo到选择区域,
	syncCallbackSelection();
	
	//if (parent || location.href.indexOf("@entrance=true") > -1) {
	g_gridPanelDict["columnModel_1"].dt.refreshPaginator();
	//}
//		});
//	});
}


