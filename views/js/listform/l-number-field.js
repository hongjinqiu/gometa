Y.LNumberField = Y.Base.create('l-number-field', Y.RNumberField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell inputWidth',
	
    bindUI: function() {
    	Y.LNumberField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    },
    initializer: function() {
    	Y.LNumberField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
