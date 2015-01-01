Y.PSelectField = Y.Base.create('p-select-field', Y.RSelectField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell selectWidth',
	
    bindUI: function() {
    	Y.PSelectField.superclass.bindUI.apply(this, arguments);
    	
    	var self = this;
    	new FormManager().applyEventBehavior(self, Y);
    },

    _validateReadonly: function(val) {
    	var self = this;
    	return new FormManager().validateReadonly(self, val, Y);
    },
    
    initializer: function() {
    	Y.PSelectField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	var formManager = new FormManager();
    	formManager.initializeAttr(self, Y);
    	formManager.setChoices(self);
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
