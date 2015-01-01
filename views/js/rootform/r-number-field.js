/**
 * @class RNumberField
 * @extends RTextField
 * @param config {Object} Configuration object
 * @constructor
 * @description A hidden field node
 */
Y.RNumberField = Y.Base.create('r-number-field', Y.RTextField, [Y.WidgetChild], {
    INPUT_TYPE: "text",
    
    bindUI: function() {
    	Y.RNumberField.superclass.bindUI.apply(this, arguments);
    	var self = this;
    	
    	this._fieldNode.detach('change');
    	this._fieldNode.detach('blur');
    	this.detach('valueChange');
    	
    	this._fieldNode.on('change', Y.bind(function(e) {
    		self.fieldNodeChange();
    	},
        this));
    	
    	this._fieldNode.on('blur', Y.bind(function(e) {
    		self.fieldNodeChange();
    	},
        this));

    	this.on('valueChange', Y.bind(function(e) {
            if (e.src != 'ui') {
    			var displayPattern = self.get("displayPattern");
    			if (displayPattern) {
    				var value = e.newVal + "";
    				value = this._trimZero(value);
    				var showValue = this.getShowValue(value);
    				if (showValue == "") {// 0显示空的折腾
    					this._fieldNode.set('value', showValue);
    				} else {
    					var displayValue = displayPattern(value);
    					
    					displayValue = this.getShowValue(displayValue);
    					this._fieldNode.set('value', displayValue);
    				}
    			} else {
    				var showValue = this.getShowValue(e.newVal);
                    this._fieldNode.set('value', showValue);
    			}
            }
        },
        this));
    },
    _trimZero: function(value){
    	value = value.replace(/^0+|0+$/gi, "");
    	if (value == ".") {
    		value = "0";
    	}
    	return value
    },
    fieldNodeChange: function() {
    	var self = this;
		var displayPattern = self.get("displayPattern");
		if (displayPattern) {
			var value = this.getValueFromShowValue(this._fieldNode.get('value'));
			value = value.replace(/[^\d.-]/gi, "");
			
			var showValue = this.getShowValue(value);
			if (showValue == "") {// 0显示空的折腾
				this._fieldNode.set('value', showValue);
			} else {
				var fieldValue = displayPattern(value);
				this._fieldNode.set('value', fieldValue);
			}
			
			this.set('value', value, {
				src: 'ui'
			});
		} else {
			var value = this.getValueFromShowValue(this._fieldNode.get('value'));
			this.set('value', value, {
				src: 'ui'
			});
		}
    }
},
{
    /**
	 * @property RNumberField.ATTRS
	 * @type Object
	 * @static
	 */
    ATTRS: {
    	displayPattern: {
            value: null,
            validator: Y.Lang.isFunction
        }
    }

});
