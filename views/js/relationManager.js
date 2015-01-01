function RelationManager() {
}

RelationManager.prototype.mergeRelationBo = function(relationBo) {
	for (var key in relationBo) {
		if (!g_relationBo[key]) {
			g_relationBo[key] = relationBo[key];
		} else {
			var relationBoItem = g_relationBo[key];
			for (var subKey in relationBo[key]) {
				if (!relationBoItem[subKey]) {
					relationBoItem[subKey] = relationBo[key][subKey];
				}
			}
		}
	}
}

/**
 * 开窗选回后,加入到g_relationBo中
 */
RelationManager.prototype.addRelationBo = function(selectorId, url, obj) {
	if (!g_relationBo[selectorId]) {
		g_relationBo[selectorId] = {};
	}
	g_relationBo[selectorId][obj["id"]] = obj;
	if (url) {
		g_relationBo[selectorId]["url"] = url;
	}
}

RelationManager.prototype.getRelationBo = function(selectorId, id) {
	if (!selectorId || selectorId === "" || selectorId === 0 || selectorId === "0") {
		return null;
	}
	if (!id || id === "" || id === 0 || id === "0") {
		return null;
	}
	if (g_relationBo[selectorId] && g_relationBo[selectorId][id]) {
		return g_relationBo[selectorId][id];
	}
	var self = this;
	var result = null;
	ajaxRequest({
		url : "/console/relation?selectorId=" + selectorId + "&id=" + id + "&date=" + new Date(),
		method: "GET",
		callback : function(o) {
			result = o["result"];
			if (result) {
				var url = o["url"];
				self.addRelationBo(selectorId, url, result);
			}
		}
	});
	return result;
}
