Y.LChoiceField = Y.Base.create('l-choice-field', Y.RChoiceField, [Y.WidgetChild], {
	initializer : function () {
		Y.LChoiceField.superclass.initializer.apply(this, arguments);
		var self = this;
		
		new LFormManager().initializeAttr(self, Y);
		
		var choiceFieldManager = new ChoiceFieldManager();
		this.set("choices", choiceFieldManager.getChoices(self.get("name")));
    },
    
    bindUI: function() {
    	Y.LChoiceField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	new LFormManager().applyEventBehavior(self);
    }
},
{

    ATTRS: {
    	
    }
});
