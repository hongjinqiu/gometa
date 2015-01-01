Y.LCheckboxField = Y.Base.create('l-checkbox-field', Y.RCheckboxField, [Y.WidgetChild], {
    initializer: function() {
    	Y.LCheckboxField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
