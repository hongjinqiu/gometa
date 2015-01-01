Y.PNumberField = Y.Base.create('p-number-field', Y.RNumberField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell inputWidth',
	
	_getFieldDict: function() {
		var self = this;
		var dataSetId = self.get("dataSetId");
		var fieldDict = null;
		if (dataSetId == "A") {
			fieldDict = g_masterFormFieldDict;
		} else {
			if (g_gridPanelDict[dataSetId + "_addrow"]) {
				var record = g_gridPanelDict[dataSetId + "_addrow"].dt.getRecord(self._fieldNode);
				fieldDict = record.formFieldDict;
			}
		}
		return fieldDict;
	},
	
    bindUI: function() {
    	Y.PNumberField.superclass.bindUI.apply(this, arguments);
    	
    	var self = this;
    	new FormManager().applyEventBehavior(self, Y);
    },

    _validateReadonly: function(val) {
    	var self = this;
    	return new FormManager().validateReadonly(self, val, Y);
    },
    
    initializer: function() {
    	Y.PNumberField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new FormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	dataSetId: {
            validator: Y.Lang.isString,
            writeOnce: true
        }
    }
});
