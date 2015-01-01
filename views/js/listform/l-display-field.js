Y.LDisplayField = Y.Base.create('l-display-field', Y.RDisplayField, [Y.WidgetChild], {
    bindUI: function() {
    	Y.LDisplayField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    },
    initializer: function() {
    	Y.LDisplayField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
