Y.LHiddenField = Y.Base.create('l-hidden-field', Y.RHiddenField, [Y.WidgetChild], {
    bindUI: function() {
    	Y.LHiddenField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    },
    initializer: function() {
    	Y.LHiddenField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
