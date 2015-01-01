Y.LRadioField = Y.Base.create('l-radio-field', Y.RRadioField, [Y.WidgetChild], {
    initializer: function() {
    	Y.LRadioField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
