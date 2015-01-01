Y.LTextField = Y.Base.create('l-text-field', Y.RTextField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell inputWidth',
	
    bindUI: function() {
    	Y.LTextField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    },
    initializer: function() {
    	Y.LTextField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
