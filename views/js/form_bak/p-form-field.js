/**
 * @class PFormField
 * @extends Widget
 * @param config {Object} Configuration object
 * @constructor
 * @description A representation of an individual form field.
 */

Y.PFormField = Y.Base.create('p-form-field', Y.Widget, [Y.WidgetParent, Y.WidgetChild], {
    toString: function() {
        return this.name;
    },

    /**
     * @property PFormField.FIELD_TEMPLATE
     * @type String
     * @description Template used to render the field node
     */
    FIELD_TEMPLATE : '<input>',

    /**
     * @property PFormField.FIELD_CLASS
     * @type String
     * @description CSS class used to locate a placeholder for
     *     the field node and style it.
     */
    FIELD_CLASS : 'table-layout-cell',

    /**
     * @property PFormField.LABEL_TEMPLATE
     * @type String
     * @description Template used to draw a label node
     */
    LABEL_TEMPLATE : '<label></label>',

    /**
     * @property PFormField.LABEL_CLASS
     * @type String
     * @description CSS class used to locate a placeholder for
     *     the label node and style it.
     */
    LABEL_CLASS : 'table-layout-cell-title',

    /**
     * @property PFormField.HINT_TEMPLATE
     * @type String
     * @description Optionally a template used to draw a hint node. Derived
     *     classes can use it to provide additional information about the field
     */
    HINT_TEMPLATE : '',

    /**
     * @property PFormField.HINT_CLASS
     * @type String
     * @description CSS class used to locate a placeholder for
     *     the hint node and style it.
     */
    HINT_CLASS : 'hint',

    /**
     * @property PFormField.ERROR_TEMPLATE
     * @type String
     * @description Template used to draw an error node
     */
    ERROR_TEMPLATE : '<div></div>',

    /**
     * @property PFormField.ERROR_CLASS
     * @type String
     * @description CSS class used to locate a placeholder for
     *     the error node and style it.
     */
    ERROR_CLASS : 'x-form-invalid-field',

    /**
     * @property _labelNode
     * @protected
     * @type Object
     * @description The label node for this form field
     */
    _labelNode: null,

     /**
     * @property _hintNode
     * @protected
     * @type Object
     * @description The hint node with extra text describing the field
     */    
    _hintNode : null,

    /**
     * @property _fieldNode
     * @protected
     * @type Object
     * @description The form field itself
     */
    _fieldNode: null,

    /**
     * @property _errorNode
     * @protected
     * @type Object
     * @description If a validation error occurs, it will be displayed in this node
     */
    _errorNode: null,

    /**
     * @property _initialValue
     * @private
     * @type String
     * @description The initial value set on this field, reset will set the value to this
     */
    _initialValue: null,

    /**
     * @method _validateError
     * @protected
     * @param val {Mixed}
     * @description Validates the value passed to the error attribute
     * @return {Boolean}
     */
    _validateError: function(val) {
        if (Y.Lang.isString(val)) {
            return true;
        }
        if (val === null || typeof val == 'undefined') {
            return true;
        }

        return false;
    },

    /**
     * @method _validateValidator
     * @protected
     * @param val {Mixed}
     * @description Validates the input of the validator attribute
     * @return {Boolean}
     */
    _validateValidator: function(val) {
        if (Y.Lang.isFunction(val)) {
            return true;
        }
        return false;
    },
    
    _validateReadonly: function(val) {
    	if (!Y.Lang.isBoolean(val)) {
    		return false;
    	}
    	var self = this;
    	var modelIterator = new ModelIterator();
    	var result = "";
    	var validateResult = true;
    	modelIterator.iterateAllField(g_dataSourceJson, result, function(fieldGroup, result){
    		if (fieldGroup.Id == self.get("name") && fieldGroup.getDataSetId() == self.get("dataSetId")) {
    			if (fieldGroup.FixReadOnly == "true" && !val) {
    				validateResult = false;
    			}
    		}
    	});
    	return validateResult;
    },

    /**
     * @method _renderNode
     * @protected
     * @description Helper method to render new nodes, possibly replacing
     *     markup placeholders.
     */
    _renderNode : function (nodeTemplate, nodeClass, nodeBefore) {
        if (!nodeTemplate) {
            return null;
        }
        var contentBox = this.get('contentBox'),
            node = Y.Node.create(nodeTemplate),
            placeHolder = contentBox.one('.' + nodeClass);

        node.addClass(nodeClass);

        if (placeHolder) {
            placeHolder.replace(node);
        } else {
            if (nodeBefore) {
                contentBox.insertBefore(node, nodeBefore);
            } else {
                contentBox.appendChild(node);
            }
        }

        return node;
    },

    /**
     * @method _renderLabelNode
     * @protected
     * @description Draws the form field's label node into the contentBox
     */
    _renderLabelNode: function() {
        var contentBox = this.get('contentBox'),
        labelNode = contentBox.one('label');

        if (!labelNode || labelNode.get('for') != this.get('id')) {
            labelNode = this._renderNode(this.LABEL_TEMPLATE, this.LABEL_CLASS);
        }

        this._labelNode = labelNode;
    },

    /**
     * @method _renderHintNode
     * @protected
     * @description Draws the hint node into the contentBox. If a node is
     *     found in the contentBox with class HINT_CLASS, it will be
     *     considered a markup placeholder and replaced with the hint node.
     */
    _renderHintNode : function () {
        this._hintNode = this._renderNode(this.HINT_TEMPLATE,
                                          this.HINT_CLASS);
    },

    /**
     * @method _renderFieldNode
     * @protected
     * @description Draws the field node into the contentBox
     */
    _renderFieldNode: function() {
        var contentBox = this.get('contentBox'),
        field = contentBox.one('#' + this.get('id'));

        if (!field) {
            field = this._renderNode(this.FIELD_TEMPLATE, this.FIELD_CLASS);
        }

        this._fieldNode = field;
    },

    /**
     * @method _syncLabelNode
     * @protected
     * @description Syncs the the label node and this instances attributes
     */
    _syncLabelNode: function() {
        var label = this.get('label'),
            required = this.get('required'),
            requiredLabel = this.get('requiredLabel');
        if (this._labelNode) {
            this._labelNode.set("text", "");
            if (label) {
                this._labelNode.append("<span class='caption'>" + label + "</span>"); 
            }
            if (required && requiredLabel) {
                this._labelNode.append("<span class='separator'> </span>");
                this._labelNode.append("<span class='required'>" + requiredLabel + "</span>");
            }
            this._labelNode.setAttribute('for', this.get('id') + Y.PFormField.FIELD_ID_SUFFIX);
        }
    },

    /**
     * @method _syncHintNode
     * @protected
     * @description Syncs the hintNode
     */
    _syncHintNode : function () {
        if (this._hintNode) {
            this._hintNode.set("text", this.get("hint"));
        }
    },

    /**
     * @method _syncFieldNode
     * @protected
     * @description Syncs the fieldNode and this instances attributes
     */
    _syncFieldNode: function() {
        var nodeType = this.INPUT_TYPE || this.name.split('-')[1];
        if (!nodeType) {
            return;
        }

        this._fieldNode.setAttrs({
            name: this.get('name'),
            type: nodeType,
            id: this.get('id') + Y.PFormField.FIELD_ID_SUFFIX,
            value: this.get('value')
        });

        this._fieldNode.setAttribute('tabindex', Y.PFormField.tabIndex);
        Y.PFormField.tabIndex++;
    },

    /**
     * @method _syncError
     * @private
     * @description Displays any pre-defined error message
     */
    _syncError: function() {
        var err = this.get('error');
        if (err) {
            this._showError(err);
        }
    },

    _syncDisabled: function(e) {
        var dis = this.get('disabled');
        if (dis === true) {
            this._fieldNode.setAttribute('disabled', 'disabled');
        } else {
            this._fieldNode.removeAttribute('disabled');
        }
    },
    
    _syncReadonly: function(e) {
        var value = this.get('readonly');
        if (value === true) {
            this._fieldNode.setAttribute('readonly', 'readonly');
        } else {
            this._fieldNode.removeAttribute('readonly');
        }
    },

    /**
     * @method _checkRequired
     * @private
     * @description if the required attribute is set to true, returns whether or not a value has been set
     * @return {Boolean}
     */
    _checkRequired: function() {
        if (this.get('required') === true && this.get('value').length === 0) {
            return false;
        }
        return true;
    },

    /**
     * @method _showError
     * @param {String} errMsg
     * @private
     * @description Adds an error node with the supplied message
     */
    /*
    _showError: function(errMsg) {
        var errorNode = this._renderNode(this.ERROR_TEMPLATE, this.ERROR_CLASS, this._labelNode);

        errorNode.set("text", errMsg);
        this._errorNode = errorNode;
    },
    */
    _showError: function(errMsg) {
    	this._fieldNode.addClass(this.ERROR_CLASS);
    },

    /**
     * @method _clearError
     * @private
     * @description Removes the error node from this field
     */
    /*
    _clearError: function() {
        if (this._errorNode) {
            this._errorNode.remove();
            this._errorNode = null;
        }
    },
    */
    _clearError: function() {
    	this._fieldNode.removeClass(this.ERROR_CLASS);
    	this._clearErrorOnMouseOver();
    },
    
    _showErrorOnMouseOver: function() {
    	var err = this.get('error');
        if (err) {
        	var xy = this._fieldNode.getXY();
    		x = xy[0];
    		y = xy[1];
    		var fieldWidth = parseInt(this._fieldNode.getComputedStyle("width"));
    		var width = 200;
    		var height = 50;
//    		var height = parseInt(this._fieldNode.getComputedStyle("height"));
    		
    		var errorRenderId = this.get('id') + "_error";
        	var errorNode = Y.one("#" + errorRenderId);
    		if (!errorNode) {
    			var errorStyleLi = [];
    			errorStyleLi.push('position: absolute;');
    			errorStyleLi.push('z-index: 999;');
    			errorStyleLi.push('background-color: white;');
//    			errorStyleLi.push('opacity: 0.5;');
//    			errorStyleLi.push('filter:alpha(opacity=50);');
    			errorStyleLi.push('width: ' + width + 'px;');
    			errorStyleLi.push('height: ' + height + 'px;');
    			errorStyleLi.push('left: ' + (x + fieldWidth) + 'px;');
    			errorStyleLi.push('top: ' + y + 'px;');
    			errorStyleLi.push('display: none;');
    			errorStyleLi.push('border: 1px solid red;');


    			var htmlLi = [];
    			htmlLi.push('<div id="' + errorRenderId + '" style="' + errorStyleLi.join("") + '">');
    			htmlLi.push(err);
    			htmlLi.push('</div>');
    			Y.one("body").append(htmlLi.join(""));
    			errorNode = Y.one("#" + errorRenderId);
    		}
    		errorNode.setStyle("width", width + "px");
    		errorNode.setStyle("height", height + "px");
    		errorNode.setStyle("left", (x + fieldWidth) + "px");
    		errorNode.setStyle("top", y + "px");

    		errorNode.setStyle("display", "");
        }
    },
    
    _clearErrorOnMouseOver: function() {
    	var errorRenderId = this.get('id') + "_error";
    	var errorNode = Y.one("#" + errorRenderId);
    	if (errorNode) {
    		errorNode.remove();
    		errorNode = null;
        }
    },

    _enableInlineValidation: function() {
        this.after('valueChange', this.validateField, this);
    },

    _disableInlineValidation: function() {
        this.detach('valueChange', this.validateField, this);
    },

    /**
     * @method validateField
     * @description Runs the validation functions of this form field
     * @return {Boolean}
     */
    validateField: function(e) {
        var value = this.get('value'),
        validator = this.get('validator');

        this.set('error', null);

        if (e && e.src != 'ui') {
            return false;
        }

        if (!this._checkRequired()) {
            this.set('error', Y.PFormField.REQUIRED_ERROR_TEXT);
            return false;
        } else if (!value) {
            return true;
        }

        return validator.call(this, value, this);
    },

    resetFieldNode: function() {
        this.set('value', this._initialValue);
        this._fieldNode.set('value', this._initialValue);
        this.fire('nodeReset');
    },

    /**
     * @method clear
     * @description Clears the value AND the initial value of this field
     */
    clear: function() {
        this.set('value', '');
        this._fieldNode.set('value', '');
        this._initialValue = null;
        this.fire('clear');
    },

    initializer: function() {
    	var self = this;
        this.publish('blur');
        this.publish('change');
        this.publish('focus');
        this.publish('clear');
        this.publish('nodeReset');

        this._initialValue = this.get('value');
        if (g_dataSourceJson) {
        	var modelIterator = new ModelIterator();
        	var result = "";
        	modelIterator.iterateAllField(g_dataSourceJson, result, function(fieldGroup, result){
        		if (fieldGroup.Id == self.get("name") && fieldGroup.getDataSetId() == self.get("dataSetId")) {
        			if (fieldGroup.AllowEmpty != "true") {
        				self.set("required", true);
        			}
        			self.set("readonly", fieldGroup.FixReadOnly == "true");
        		}
        	});
        	var formManager = new FormManager();
        	self.set("validator", formManager.dsFormFieldValidator);
        }
    },

    destructor: function(config) {

    },

    renderUI: function() {
        this._renderLabelNode();
        this._renderFieldNode();
        this._renderHintNode();
    },

    bindUI: function() {
        this._fieldNode.on('change', Y.bind(function(e) {
            this.set('value', this._fieldNode.get('value'), {
                src: 'ui'
            });
        },
        this));

        this.on('valueChange', Y.bind(function(e) {
            if (e.src != 'ui') {
                this._fieldNode.set('value', e.newVal);
            }
        },
        this));
        
        this.after('readonlyChange', Y.bind(function(e) {
        	this._syncReadonly();
        },
        this));

        this._fieldNode.on('blur', Y.bind(function(e) {
            this.set('value', this._fieldNode.get('value'), {
                src: 'ui'
            });
        },
        this));

        this._fieldNode.on('focus', Y.bind(function(e) {
            this.fire('focus', e);
        },
        this));

        this.on('errorChange', Y.bind(function(e) {
            if (e.newVal) {
                this._showError(e.newVal);
            } else {
                this._clearError();
            }
        },
        this));
        
        this._fieldNode.on('mouseover', Y.bind(function(e) {
            this._showErrorOnMouseOver();
        },
        this));
        
        this._fieldNode.on('mouseout', Y.bind(function(e) {
            this._clearErrorOnMouseOver();
        },
        this));
        

        this.on('validateInlineChange', Y.bind(function(e) {
            if (e.newVal === true) {
                this._enableInlineValidation();
            } else {
                this._disableInlineValidation();
            }
        },
        this));

        this.after('disabledChange', Y.bind(function(e) {
            this._syncDisabled();
        },
        this));
        
        // 应用上js相关的操作,
        var self = this;
        var modelIterator = new ModelIterator();
    	var result = "";
    	modelIterator.iterateAllField(g_dataSourceJson, result, function(fieldGroup, result){
    		if (fieldGroup.Id == self.get("name") && fieldGroup.getDataSetId() == self.get("dataSetId")) {
    			if (fieldGroup.jsConfig && fieldGroup.jsConfig.listeners) {
    				for (var key in fieldGroup.jsConfig.listeners) {
    					self._fieldNode.on(key, function(key) {
    						return function(e) {
    							fieldGroup.jsConfig.listeners[key](e, self);
    						}
    					}(key));
    				}
    			}
    		}
    	});
    },

    syncUI: function() {
        this.get('boundingBox').removeAttribute('tabindex');
        this._syncLabelNode();
        this._syncHintNode();
        this._syncFieldNode();
        this._syncError();
        this._syncDisabled();
        this._syncReadonly();

        if (this.get('validateInline') === true) {
            this._enableInlineValidation();
        }
    }
    
},
{
    /**
     * @property PFormField.ATTRS
     * @type Object
     * @protected
     * @static
     */
    ATTRS: {
        /**
         * @attribute id
         * @type String
         * @default Either a user defined ID or a randomly generated by Y.guid()
         * @description A randomly generated ID that will be assigned to the field and used 
         * in the label's for attribute
         */
        id: {
            value: Y.guid(),
            validator: Y.Lang.isString,
            writeOnce: true
        },

        /**
         * @attribute name
         * @type String
         * @default ""
         * @writeOnce
         * @description The name attribute to use on the field
         */
        name: {
            validator: Y.Lang.isString,
            writeOnce: true
        },
        
        dataSetId: {
            validator: Y.Lang.isString,
            writeOnce: true
        },

        /**
         * @attribute value
         * @type String
         * @default ""
         * @description The current value of the form field
         */
        value: {
            value: '',
            //validator: Y.Lang.isString
            validator: function(val) {
            	return Y.Lang.isString(val) || Y.Lang.isNumber(val);
            },
            setter: function(val) {
            	if (val === undefined || val === null) {
            		val = "";
            	}
            	return val + "";
            }
        },

        /**
         * @attribute label
         * @type String
         * @default ""
         * @description Label of the form field
         */
        label: {
            value: '',
            validator: Y.Lang.isString
        },

        /**
         * @attribute hint
         * @type String
         * @default ""
         * @description Extra text explaining what the field is about.
         */
        hint : {
            value : '',
            validator : Y.Lang.isString
        },
        
        /**
         * @attribute validator
         * @type Function
         * @default "function () { return true; }"
         * @description Used to validate this field by the Form class
         */
        validator: {
            value: function(val) {
                return true;
            },
            validator: function(val) {
                return this._validateValidator(val);
            }
        },

        /**
         * @attribute error
         * @type String
         * @description An error message associated with this field. Setting this will
         *              cause validation to fail until a new value is entered
         */
        error: {
            value: false,
            validator: function(val) {
                return this._validateError(val);
            }
        },

        /**
         * @attribute required
         * @type Boolean
         * @default false
         * @description Set true if this field must be filled out when submitted
         */
        required: {
            value: false,
            validator: Y.Lang.isBoolean
        },

        /**
         * @attribute validateInline
         * @type Boolean
         * @default false
         * @description Set to true to validate this field whenever it's value is changed
         */
        validateInline: {
            value: false,
            validator: Y.Lang.isBoolean
        },

        /**
         * @attribute requiredLabel
         * @type String
         * @description Text to append to the labal caption for a required
         *     field, by default nothing will be appended.
         */
        requiredLabel : {
            value : '',
            validator : Y.Lang.isString
        },
        
        /**
         * 是否只读
         */
        readonly : {
            value : false,
            validator : function(val) {
                return this._validateReadonly(val);
            }
        }
    },

    /**
     * @property PFormField.tabIndex
     * @type Number
     * @description The current tab index of all Y.PFormField instances
     */
    tabIndex: 1,

    /**
     * @property PFormField.REQUIRED_ERROR_TEXT
     * @type String
     * @description Error text to display for a required field
     */
    REQUIRED_ERROR_TEXT: '不允许为空',

    /**
     * @property PFormField.FIELD_ID_SUFFIX
     * @type String
     */
    FIELD_ID_SUFFIX: '-field'
});
