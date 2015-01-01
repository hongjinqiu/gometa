/**
 * @class RDateField
 * @extends RFormField
 * @param config {Object} Configuration object
 * @constructor
 * @description A form field which allows one or multiple values from a 
 * selection of choices
 */
Y.RDateField = Y.Base.create('r-date-field', Y.RFormField, [Y.WidgetParent, Y.WidgetChild], {
	FIELD_CLASS : 'table-layout-cell trigger_input Wdate',
	INPUT_TYPE: "text",
	DATE_TEMPLATE: '<a></a>',
//	DATE_CLASS: 'trigger_date',
//	_dateNode: null,
	_olNode: null,
	_overlay: null,
	_calendar: null,
	
//	_renderDateNode : function () {
//        this._dateNode = this._renderNode(this.DATE_TEMPLATE, this.DATE_CLASS);
//    },
	
	renderUI: function() {
		Y.RDateField.superclass.renderUI.apply(this, arguments);
		
//		this._renderDateNode();
		this._olNode = Y.Node.create('<div></div>');
		this.get('contentBox').appendChild(this._olNode);
    },
    
    _destroyOverlay: function() {
    	if (this._calendar) {
    		this._calendar.destroy();
    		this._calendar = null;
    	}
    	if (this._overlay) {
    		this._overlay.destroy();
    		this._overlay = null;
    	}
    },
    
    _hideOverlay: function() {
    	if (this._calendar) {
    		this._calendar.hide();
    	}
    	if (this._overlay) {
    		this._overlay.hide();
    	}
    },
    
    _createAndShowOverlay: function() {
    	var self = this;
    	if (this._overlay) {
    		if (this._calendar) {
    			this._calendar.show();
    		}
    		this._overlay.show();
    		return;
    	}
    	if (this.get("readonly")) {
    		return;
    	}
    	//this._hideOverlay();
    	
    	this._overlay = new Y.Overlay({
    		bodyContent:"<div></div>",
    		visible : false,
    		width : '250px',
    		zIndex     : (++panelZIndex)
    		//,xy : [500, 100]
    	});
    	
    	var node = this._fieldNode;
    	var xy = node.getXY();
    	var x = xy[0];
    	var y = xy[1];
    	//var width = parseInt(node.getComputedStyle("width"));
    	var height = parseInt(node.getComputedStyle("height"), 10);
    	this._overlay.setAttrs({
    		x: x,
    		y: y + height + 4
    	});
    	this._overlay.render(this._olNode);
    	
    	var selectDate = self._getSelectDate();
    	this._calendar = new Y.Calendar({
    		width:'250px',
    		showPrevMonth: true,
    		showNextMonth: true,
    		date: selectDate || new Date()}).render(this._overlay.get("bodyContent").getDOMNodes()[0]);
    	if (selectDate) {
    		this._calendar.selectDates(selectDate);
    	}
    	
    	
    	var dtdate = Y.DataType.Date;
		this._calendar.on("selectionChange", function (ev) {
			var dbPattern = self.get("dbPattern");
			
			var columnManager = new ColumnManager();
			var displayPattern = columnManager.convertDate2DisplayPattern(dbPattern);
			var newDate = ev.newSelection[0];
			self.set("value", dtdate.format(newDate, {
				format: displayPattern
			}));
			self._hideOverlay();
		});
		this._overlay.show();
    },
    
    _getSelectDate: function() {
    	var self = this;
    	if (this.get("value") && this.get("value").length > 1) {
    		var dbPattern = self.get("dbPattern");
			var value = self.get("value") + "";
			var index = dbPattern.indexOf("yyyy");
			var yyyy = value.substr(index, 4);
			index = dbPattern.indexOf("MM");
			var MM = value.substr(index, 2);
			index = dbPattern.indexOf("dd");
			var dd = value.substr(index, 2);
			
			var selectDate = new Date();
			selectDate.setFullYear(parseInt(yyyy, 10));
			selectDate.setMonth(parseInt(MM, 10) - 1);
			selectDate.setDate(parseInt(dd, 10));
			return selectDate;
    	}
    	return null;
    },
    
    _isClickInBoundingBox: function(e) {
    	if (this._calendar) {
    		var isInCalendar = e.target.ancestor("#" + this._calendar.get("boundingBox").get("id"));
    		if (isInCalendar != null) {
    			return true;
    		}
    		var isInFieldBox = e.target.ancestor("#" + this.get("boundingBox").get("id"));
    		if (isInFieldBox != null) {
    			return true;
    		}
    	}
    	return false;
    },
    
    bindUI: function() {
    	var self = this;
    	Y.RDateField.superclass.bindUI.apply(this, arguments);
    	
    	this._fieldNode.detach('change');
    	this._fieldNode.detach('blur');
    	
    	this.detach('valueChange');
    	
    	this._fieldNode.on('change', Y.bind(function(e) {
    		var value = this._fieldNode.get('value');
    		value = value.replace(/[-\/\s:]/gi, "");
    		value = this.getValueFromShowValue(value);
            this.set('value', value, {
                src: 'ui'
            });
        },
        this));
    	
    	this._fieldNode.on('blur', Y.bind(function(e) {
    		var value = this._fieldNode.get('value');
    		value = value.replace(/[-\/\s:]/gi, "");
    		value = this.getValueFromShowValue(value);
            this.set('value', value, {
                src: 'ui'
            });
        },
        this));
    	
    	this.on('valueChange', Y.bind(function(e) {
            if (e.src != 'ui') {
    			var dbPattern = self.get("dbPattern");
    			var displayPattern = self.get("displayPattern");
    			var value = e.newVal + "";
    			var index = 0;
    			index = dbPattern.indexOf("yyyy");
    			var yyyy = value.substr(index, 4);
    			index = dbPattern.indexOf("MM");
    			var MM = value.substr(index, 2);
    			index = dbPattern.indexOf("dd");
    			var dd = value.substr(index, 2);
    			index = dbPattern.indexOf("HH");
    			var HH = value.substr(index, 2);
    			index = dbPattern.indexOf("mm");
    			var mm = value.substr(index, 2);
    			index = dbPattern.indexOf("ss");
    			var ss = value.substr(index, 2);
    			
    			var displayValue = displayPattern;
    			displayValue = displayValue.replace("yyyy", yyyy);
    			displayValue = displayValue.replace("MM", MM);
    			displayValue = displayValue.replace("dd", dd);
    			displayValue = displayValue.replace("HH", HH);
    			displayValue = displayValue.replace("mm", mm);
    			displayValue = displayValue.replace("ss", ss);
    			displayValue = displayValue.replace(/[a-z]/gi, "");
    			displayValue = displayValue.replace(/^[-\/\s:]*|[-\/\s:]*$/gi, "");
    			displayValue = this.getShowValue(displayValue);
                this._fieldNode.set('value', displayValue);
            }
        },
        this));
    	

    	this._fieldNode.on('focus', Y.bind(function () {
    		if (!this.get("readonly")) {
    			this._createAndShowOverlay();
    		} else {
    			this._hideOverlay();
    		}
		}, this));
    	
//    	this._dateNode.on("click", Y.bind(function(e) {
    		/*if (!this._overlay) {
    			this._createAndShowOverlay();
    		} else {
    			this._hideOverlay();
    		}*/
//    		this._createAndShowOverlay();
//    	}, this));
    	
    	Y.one("document").on("click", Y.bind(function(e){
    		if (!this._isClickInBoundingBox(e)) {
    			this._hideOverlay();
    		}
    	}, this));
    },
    
    getDisplayValue: function(e) {
    	return this._fieldNode.get('value');
    },
    
    _syncReadonly: function(e) {
    	Y.RDateField.superclass._syncReadonly.apply(this, arguments);
    	
    	var value = this.get('readonly');
        if (value === true) {
//        	this._dateNode.setStyle("display", "none");
        	this._hideOverlay();
        	this._fieldNode.addClass('readonly');
        } else {
//        	this._dateNode.setStyle("display", "");
        	this._fieldNode.removeClass('readonly');
        }
    },
    
//    _syncDateNode: function() {
//    	if (this._dateNode) {
//    		this._dateNode.setAttrs({
//    			href: "javascript:void(0);",
//    			title: "日期选择",
//    			id: Y.guid() + Y.RFormField.FIELD_ID_SUFFIX
//    		});
//    	}
//    },
    
    syncUI: function() {
    	Y.RDateField.superclass.syncUI.apply(this, arguments);
    	
//    	this._syncDateNode();
    },
    
    initializer: function() {
    	Y.RDateField.superclass.initializer.apply(this, arguments);
    	
    }

},
{
    ATTRS: {
    	dbPattern: {
            value: '',
            validator: Y.Lang.isString
        },
        displayPattern: {
            value: '',
            validator: Y.Lang.isString
        }
    }
});
