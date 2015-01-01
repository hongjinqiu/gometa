function FormTemplateFactory() {
}

/**
 * 扩展listTemplate
 */
FormTemplateFactory.prototype.extendFormTemplate = function(modelExtraInfo) {
	var templateIterator = new TemplateIterator();
	if (modelExtraInfo.buttonConfig) {
		for (var key in modelExtraInfo.buttonConfig) {
			templateIterator.iterateAnyTemplateButton(result, function(toolbarOrColumnModel, button, result){
				if (button.Name == key) {
					if (!button.jsConfig) {
						button.jsConfig = {};
					}
					for (var item in modelExtraInfo.buttonConfig[key]) {
						button.jsConfig[item] = modelExtraInfo.buttonConfig[key][item];
					}
					return true;
				}
				return false;
			});
		}
	}
}
