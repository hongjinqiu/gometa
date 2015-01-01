var sysUserModel = {
	"A" : {
		"nick" : {
			listeners: {}
		},
		"jsConfig": {
			listeners: {}
		}
	},
	"B" : {
		"attachCount" : {
			displayField : "",// 可以为函数
			valueField : "",// 可以为函数
			selectorName : "",// 可以为函数
			selectionMode : "single",// single|multi
			selectFunc: function(datas, formObj){},// 单多选回调
			unSelectFunc: function(formObj){},// 单多选回调
			queryFunc: function(){},// 单多选回调
			listeners : {// 会用yui.on调用,
				focus: function(e){},
				blur: function(e){},
				tabchange: function(e){},//暂时没有实现
				select: function(e){},
				reset: function(e){},//暂时没有实现
				confirm: function(e){},//暂时没有实现
				change: function(e){},
				dblclick: function(e){},
				keydown: function(e){},
				click: function(e){}
			},
			formatter: function(o){},// 数据集字段函数,接受o作为参数,
			defaultValueExprForJs : function(bo, data) {},// 整个业务对象,单行数据
			calcValueExprForJs : function(bo, data) {}// 整个业务对象,单行数据,
//			validate: function(bo, data) {}// 业务的validate,覆盖不了字段上的validate方法,放到带数据集上处理,
		},
		afterNewData: function(dataSource, bo){},// defaultValueExprForJs->calcValueExprForJs->afterNewData->calcValueExprForJs
		listeners: {},
		beforeEdit: function(recordLi, record, recordIndex){},//数据集函数,表格控件函数,表格刚渲染完成时触发,参数:(recordLi, record, recordIndex),record.formFieldDict,可以取得到别的控件,recordLi.size(),可取长度
		validateEdit: function(jsonDataLi){},//数据集函数,表格控件函数,表格弹出框点击确定时触发,参数:(jsonDataLi),editor里面有record,点击确定时触发,
		edit: function(editor, e){},//数据集函数,表格控件函数,参数:(editor,e),未实现
		canceledit: function(editor, e){},//数据集函数,表格控件函数,参数:(editor,e),未实现
		celldblclick: function() {},// 暂不实现
//		validate: function(bo, data) {}// 整条记录的validate,未实现,请使用validateEdit方法,
	},
	validate: function(bo, masterMessageLi, detailMessageDict) {},// bo:整条业务对象,masterMessageLi:[],主数据集错误信息,detailMessageDict:{"B":[],"C":[]},分录数据集错误信息
	"buttonConfig": {
		"selectRowBtn": {
			selectFunc: function(datas){},// 单多选回调
			queryFunc: function(){}// 单多选回调
		}
	}
};

var listTemplateExtraInfo = {
	"ColumnModel" : {

	},
	"QueryParameter" : {
		"currencyTypeId" : {
//				selectFunc : function(datas, formObj) {
//				},// 单多选回调
//				unSelectFunc : function(formObj) {
//				},// 单多选回调
			queryFunc : function() {
				return {
					"code": "RMB",
					"name": "人"
				};
			},// 单多选回调
			listeners : {
				click: function(e, formObj){
					console.log("click");
				},
				valueChange: function(e, formObj) {
					console.log("value cc change outside");
				}
			}
		}
	}
};
