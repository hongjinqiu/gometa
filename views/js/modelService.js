function ModelIterator() {
}

function IterateFunc(fieldGroup, data, rowIndex, result) {
}

ModelIterator.prototype.iterateAllFieldBo = function(dataSource, bo, result, iterateFunc) {
	var self = this;
	self.iterateDataBo(dataSource, bo, result, function(fieldGroupLi, data, rowIndex, result){
		for (var i = 0; i < fieldGroupLi.length; i++) {
			iterateFunc(fieldGroupLi[i], data, rowIndex, result);
		}
	})
}

function IterateFieldFunc(fieldGroup, result){}

ModelIterator.prototype.iterateAllField = function(dataSource, result, iterateFunc) {
	var self = this;
	var fieldGroupLi = self._getDataSetFieldGroupLi(dataSource.MasterData.FixField, dataSource.MasterData.BizField)
	for (var i = 0; i < fieldGroupLi.length; i++) {
		iterateFunc(fieldGroupLi[i], result);
	}
	if (dataSource.DetailDataLi) {
		for (var i = 0; i < dataSource.DetailDataLi.length; i++) {
			var fieldGroupLi = self._getDataSetFieldGroupLi(dataSource.DetailDataLi[i].FixField, dataSource.DetailDataLi[i].BizField);
			for (var j = 0; j < fieldGroupLi.length; j++) {
				iterateFunc(fieldGroupLi[j], result);
			}
		}
	}
}


ModelIterator.prototype.getFixFieldLi = function(fixField) {
	var fixFieldLi = [];
	fixFieldLi.push(fixField.PrimaryKey);
	fixFieldLi.push(fixField.CreateBy);
	fixFieldLi.push(fixField.CreateTime);
	fixFieldLi.push(fixField.CreateUnit);
	fixFieldLi.push(fixField.ModifyBy);
	fixFieldLi.push(fixField.ModifyTime);
	fixFieldLi.push(fixField.ModifyUnit);
	fixFieldLi.push(fixField.BillStatus);
	fixFieldLi.push(fixField.AttachCount);
	fixFieldLi.push(fixField.Remark);
	return fixFieldLi;
}

ModelIterator.prototype._getDataSetFieldGroupLi = function(fixField, bizField) {
	var self = this;
	var fieldGroupLi = self.getFixFieldLi(fixField);
	for (var i = 0; i < bizField.FieldLi.length; i++) {
		fieldGroupLi.push(bizField.FieldLi[i]);
	}
	return fieldGroupLi;
}

function IterateDataFunc(fieldGroupLi, data, rowIndex, result) {}

ModelIterator.prototype.iterateDataBo = function(dataSource, bo, result, iterateFunc) {
	var self = this;
	self._iterateMasterDataBo(dataSource, bo, result, iterateFunc);
	self._iterateDetailDataBo(dataSource, bo, result, iterateFunc)
}

ModelIterator.prototype._iterateMasterDataBo = function(dataSource, bo, result, iterateFunc) {
	var self = this;
	var data = bo["A"];
	var fieldGroupLi = self._getDataSetFieldGroupLi(dataSource.MasterData.FixField, dataSource.MasterData.BizField)
	var rowIndex = 0;
	iterateFunc(fieldGroupLi, data, rowIndex, result)
}

ModelIterator.prototype._iterateDetailDataBo = function(dataSource, bo, result, iterateFunc) {
	var self = this;
	if (dataSource.DetailDataLi) {
		for (var i = 0; i < dataSource.DetailDataLi.length; i++) {
			var item = dataSource.DetailDataLi[i];
			var fieldGroupLi = self._getDataSetFieldGroupLi(item.FixField, item.BizField);
			var dataLi = bo[item.Id];
			for (var j = 0; j < dataLi.length; j++) {
				var data = dataLi[j];
				iterateFunc(fieldGroupLi, data, j, result);
			}
		}
	}
}

function IterateFunc(dataSet, result){}

ModelIterator.prototype.iterateAllDataSet = function(dataSource, result, iterateFunc) {
	var self = this;
	iterateFunc(dataSource.MasterData, result);
	if (dataSource.DetailDataLi) {
		for (var i = 0; i < dataSource.DetailDataLi.length; i++) {
			iterateFunc(dataSource.DetailDataLi[i], result);
		}
	}
}

//function IterateFunc(btnName, btnValue, result){}
//
//ModelIterator.prototype.iterateAllButton = function(dataSource, result, iterateFunc) {
//	var self = this;
//	if (dataSource.buttonConfig) {
//		for (var key in dataSource.buttonConfig) {
//			iterateFunc(key, dataSource.buttonConfig[key], result);
//		}
//	}
//}



