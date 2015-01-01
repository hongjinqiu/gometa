/**
 * @class PNumberField
 * @extends PTextField
 * @param config {Object} Configuration object
 * @constructor
 * @description A hidden field node
 */
Y.PNumberField = Y.Base.create('p-number-field', Y.PTextField, [Y.WidgetChild], {
    INPUT_TYPE: "text"
},
{
    /**
	 * @property PNumberField.ATTRS
	 * @type Object
	 * @static
	 */
    ATTRS: {
    }

});
