Y.LSelectField = Y.Base.create('l-select-field', Y.RSelectField, [Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell selectWidth',
	
	initializer : function () {
		Y.LSelectField.superclass.initializer.apply(this, arguments);
		var self = this;
		
		new LFormManager().initializeAttr(self, Y);
		
		var choiceFieldManager = new ChoiceFieldManager();
		this.set("choices", choiceFieldManager.getChoices(self.get("name")));
    },
    
    bindUI: function() {
    	Y.LSelectField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	this.after('valueChange', Y.bind(function(e) {
    		
    	},
        this));
    	new LFormManager().applyEventBehavior(self);
    }
},
{

    ATTRS: {
    	
    }
});
