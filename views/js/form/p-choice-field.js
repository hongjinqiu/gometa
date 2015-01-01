Y.PChoiceField = Y.Base.create('p-choice-field', Y.RChoiceField, [Y.WidgetChild], {
    bindUI: function() {
    	Y.PChoiceField.superclass.bindUI.apply(this, arguments);
    	
    	var self = this;
    	new FormManager().applyEventBehavior(self, Y);
    },

    _validateReadonly: function(val) {
    	var self = this;
    	return new FormManager().validateReadonly(self, val, Y);
    },
    
    initializer: function() {
    	console.log("choices");
    	Y.PChoiceField.superclass.initializer.apply(this, arguments);
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
