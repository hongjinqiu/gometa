function ColumnSequenceService() {}

function SequenceStruct() {
	this.preColumnLi = [];// 前驱column li
	this.column = null;// 自身
	this.postColumnLi = [];// 后继column li
}

ColumnSequenceService.prototype._buildSequenceStructDict = function(columnLi) {
	var sequenceStructDict = {};
	var columnDict = {};
	for (var i = 0; i < columnLi.length; i++) {
		if (columnLi[i].Name != "") {
			columnDict[columnLi[i].Name] = columnLi[i];
		}
	}
	for (var i = 0; i < columnLi.length; i++) {
		if (columnLi[i].Name != "") {
			// 后继
			var sequenceStruct = new SequenceStruct();
			sequenceStruct.column = columnLi[i];
			sequenceStructDict[columnLi[i].Name] = sequenceStruct;
			
			if (columnLi[i].CRelationDS && columnLi[i].CRelationDS.CRelationItemLi) {
				for (var j = 0; j < columnLi[i].CRelationDS.CRelationItemLi.length; j++) {
					var relationItem = columnLi[i].CRelationDS.CRelationItemLi[j];
					if (relationItem.CCopyConfigLi) {
						for (var k = 0; k < relationItem.CCopyConfigLi.length; k++) {
							var copyColumnName = relationItem.CCopyConfigLi[k].CopyColumnName;
							var isIn = false;
							for (var l = 0; l < sequenceStruct.postColumnLi.length; l++) {
								if (sequenceStruct.postColumnLi[l].Name == copyColumnName) {
									isIn = true;
									break;
								}
							}
							if (!isIn) {
								if (columnDict[copyColumnName]) {
									sequenceStruct.postColumnLi.push(columnDict[copyColumnName]);
								}
							}
						}
					}
				}
			}
			
			// 前驱
			for (var j = 0; j < columnLi.length; j++) {
				if (j != i) {
					if (columnLi[j].CRelationDS && columnLi[j].CRelationDS.CRelationItemLi) {
						for (var k = 0; k < columnLi[j].CRelationDS.CRelationItemLi.length; k++) {
							var relationItem = columnLi[j].CRelationDS.CRelationItemLi[k];
							if (relationItem.CCopyConfigLi) {
								for (var l = 0; l < relationItem.CCopyConfigLi.length; l++) {
									var copyColumnName = relationItem.CCopyConfigLi[l].CopyColumnName;
									var isIn = false;
									for (var m = 0; m < sequenceStruct.preColumnLi.length; m++) {
										if (sequenceStruct.preColumnLi[m].Name == copyColumnName) {
											isIn = true;
											break;
										}
									}
									if (!isIn) {
										if (copyColumnName == columnLi[i].Name && columnDict[columnLi[j].Name]) {
											sequenceStruct.preColumnLi.push(columnDict[columnLi[j].Name]);
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return sequenceStructDict;
}

ColumnSequenceService.prototype.buildSequenceColumnLi = function(columnLi) {
	var self = this;
	var result = [];
	var sequenceStructDict = self._buildSequenceStructDict(columnLi);
	var sequenceStructDictClone = {};
	for (var name in sequenceStructDict) {
		sequenceStructDictClone[name] = sequenceStructDict[name];
	}
	var isEnd = false;
	while(!isEnd) {
		for (var name in sequenceStructDictClone) {
			if (sequenceStructDictClone[name].preColumnLi.length == 0) {
				result.push(sequenceStructDictClone[name].column);
			}
		}
		for (var i = 0; i < result.length; i++) {
			delete sequenceStructDictClone[result[i].Name];
		}
		for (var name in sequenceStructDictClone) {
			var preColumnLi = [];
			for (var i = 0; i < sequenceStructDictClone[name].preColumnLi.length; i++) {
				var isIn = false;
				for (var j = 0; j < result.length; j++) {
					if (sequenceStructDictClone[name].preColumnLi[i].Name == result[j].Name) {
						isIn = true;
						break;
					}
				}
				if (!isIn) {
					preColumnLi.push(sequenceStructDictClone[name].preColumnLi[i]);
				}
			}
			sequenceStructDictClone[name].preColumnLi = preColumnLi;
		}
		var dictLength = 0;
		for (var name in sequenceStructDictClone) {
			dictLength += 1;
			break;
		}
		if (dictLength == 0) {
			isEnd = true;
		}
	}
	return result;
}
