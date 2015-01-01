function ListTemplateFactory() {
}

/**
 * 扩展listTemplate
 */
ListTemplateFactory.prototype.extendListTemplate = function(listTemplate, modelExtraInfo) {
	var listTemplateIterator = new ListTemplateIterator();
	var result = "";
	var queryParameterConfig = listTemplateExtraInfo["QueryParameter"];
	listTemplateIterator.iterateAllTemplateQueryParameter(result, function(queryParameter, result){
		if (queryParameterConfig[queryParameter.Name]) {
			for (var key in queryParameterConfig[queryParameter.Name]) {
				if (!queryParameter.jsConfig) {
					queryParameter.jsConfig = {};
				}
				queryParameter.jsConfig[key] = queryParameterConfig[queryParameter.Name][key];
			}
		}
	});
}
