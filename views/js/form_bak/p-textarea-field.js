/**
 * @class PTextareaField
 * @extends PFormField
 * @param config {Object} Configuration object
 * @constructor
 * @description A hidden field node
 */
Y.PTextareaField = Y.Base.create('p-textarea-field', Y.PFormField, [Y.WidgetChild], {

    FIELD_TEMPLATE : '<textarea></textarea>'

});
