/**
 * @class RCheckboxField
 * @extends RFormField
 * @param config {Object} Configuration object
 * @constructor
 * @description A checkbox field node
 */

Y.RCheckboxField = Y.Base.create('r-checkbox-field', Y.RFormField, [Y.WidgetChild], {
    _syncChecked : function () {
    	this._fieldNode.set('checked', this.get('checked'));
    	this._fieldNode.set('value', this.get('value'));
    },

    initializer : function () {
        Y.RCheckboxField.superclass.initializer.apply(this, arguments);
    },

    renderUI : function () {
        this._renderFieldNode();
        this._renderLabelNode();
    },

    syncUI : function () {
        Y.RCheckboxField.superclass.syncUI.apply(this, arguments);
        this._syncChecked();
    },

    bindUI :function () {
        Y.RCheckboxField.superclass.bindUI.apply(this, arguments);
        this.after('checkedChange', Y.bind(function(e) {
            if (e.src != 'ui') {
                this._fieldNode.set('checked', e.newVal);
            }
        }, this));

        this._fieldNode.after('change', Y.bind(function (e) {
            this.set('checked', e.currentTarget.get('checked'), {src : 'ui'});
        }, this));
    },
    
    getShowValue: function(value) {
    	return value;
    }
}, {
    ATTRS : {
        'checked' : {
            value : false,
            validator : Y.Lang.isBoolean
        },
        value: {
            value: '',
            validator: function(val) {
                return Y.Lang.isString(val) || Y.Lang.isNumber(val);
            }
        }
    }
});
