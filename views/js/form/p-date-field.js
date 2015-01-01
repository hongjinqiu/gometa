Y.PDateField = Y.Base.create('p-date-field', Y.RDateField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell trigger_input inputWidth Wdate',
	
    bindUI: function() {
    	Y.PDateField.superclass.bindUI.apply(this, arguments);
    	
    	var self = this;
    	new FormManager().applyEventBehavior(self, Y);
    },

    _validateReadonly: function(val) {
    	var self = this;
    	return new FormManager().validateReadonly(self, val, Y);
    },
    
    initializer: function() {
    	Y.PDateField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new FormManager().initializeAttr(self, Y);
    	
    	var dbPattern = "";
    	var displayPattern = "";
    	var templateIterator = new TemplateIterator();
		var result = "";
		templateIterator.iterateAnyTemplateColumn(self.get("dataSetId"), result, function(column, result) {
			if (column.Name == self.get("name")) {
				dbPattern = column.DbPattern;
				displayPattern = column.DisplayPattern;
				return true;
			}
			return false;
		});
		this.set("dbPattern", dbPattern);
		this.set("displayPattern", displayPattern);
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
