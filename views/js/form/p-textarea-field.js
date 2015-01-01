Y.PTextareaField = Y.Base.create('p-textarea-field', Y.RTextareaField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell table-layout-textarea textarea_box',
	
    bindUI: function() {
    	Y.PTextareaField.superclass.bindUI.apply(this, arguments);
    	
    	var self = this;
    	new FormManager().applyEventBehavior(self, Y);
    },

    _validateReadonly: function(val) {
    	var self = this;
    	return new FormManager().validateReadonly(self, val, Y);
    },
    
    initializer: function() {
    	Y.PTextareaField.superclass.initializer.apply(this, arguments);
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
