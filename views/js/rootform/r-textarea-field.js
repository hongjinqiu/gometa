/**
 * @class RTextareaField
 * @extends RFormField
 * @param config {Object} Configuration object
 * @constructor
 * @description A hidden field node
 */
Y.RTextareaField = Y.Base.create('r-textarea-field', Y.RFormField, [Y.WidgetChild], {

    FIELD_TEMPLATE : '<textarea></textarea>',
    
    FIELD_CLASS : 'table-layout-cell table-layout-textarea'

});
