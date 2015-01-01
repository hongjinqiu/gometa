Y.LTextareaField = Y.Base.create('l-textarea-field', Y.RTextareaField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell table-layout-textarea textarea_box',
	
    bindUI: function() {
    	Y.LTextareaField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    },
    initializer: function() {
    	Y.LTextareaField.superclass.initializer.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().initializeAttr(self, Y);
    }
},
{

    ATTRS: {
    	
    }
});
