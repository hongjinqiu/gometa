function SelectionManager() {
}

/**
 * 选中后,添加到g_selectionBo里面去
 */
SelectionManager.prototype.addSelectionBo = function(obj) {
	if (!g_selectionBo) {
		g_selectionBo = {};
	}
	g_selectionBo[obj["id"]] = obj;
}

SelectionManager.prototype.getSelectionBo = function(id) {
	if (!id || id === "" || id === 0 || id === "0") {
		return null;
	}
	if (g_selectionBo[id]) {
		return g_selectionBo[id];
	}
	return null;
}
