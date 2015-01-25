package app

import (
	. "github.com/hongjinqiu/gometa/common"
	. "github.com/hongjinqiu/gometa/component"
	. "github.com/hongjinqiu/gometa/error"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/mongo"
	"github.com/hongjinqiu/gometa/mongo"
)

type LoginService struct{}

func (o LoginService) saveOrUpdateLastSessionData(sessionId int, resStruct map[string]interface{}, sysUnitId int, sysUserId int) {
	session, db := global.GetConnection(sessionId)
	qb := QuerySupport{}
	lastSessionData, found := qb.FindByMapWithSession(session, "LastSessionData", map[string]interface{}{
		"A.sysUserId": sysUserId,
		"A.sysUnitId": sysUnitId,
	})
	txnManager := TxnManager{db}
	if !found {
		id := mongo.GetSequenceNo(db, "lastSessionDataId")
		txnId := global.GetTxnId(sessionId)
		lastSessionData := map[string]interface{}{
			"_id": id,
			"id":  id,
			"A": map[string]interface{}{
				"id":          id,
				"sysUserId":   sysUserId,
				"sysUnitId":   sysUnitId,
				"resStruct":   resStruct,
				"createBy":    sysUserId,
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUnitId,
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"attachCount": 0,
				"remark":      "",
			},
		}
		txnManager.Insert(txnId, "LastSessionData", lastSessionData)
	} else {
		txnId := global.GetTxnId(sessionId)
		lastSessionDataMaster := lastSessionData["A"].(map[string]interface{})
		lastSessionData["A"] = lastSessionDataMaster

		lastSessionDataMaster["modifyBy"] = sysUserId
		lastSessionDataMaster["modifyTime"] = DateUtil{}.GetCurrentYyyyMMddHHmmss()
		lastSessionDataMaster["modifyBy"] = sysUnitId

		_, updateResult := txnManager.Update(txnId, "LastSessionData", lastSessionData)
		if !updateResult {
			panic(BusinessError{Message: "更新LastSessionData失败"})
		}
	}
}

/**
@param url /asdfas/zasdfasdf/?param1=value1&param2=value2
*/
func (o LoginService) DealLogin(sessionId int, url string) (resStruct map[string]interface{}, userId int, isStep bool) {
	isStep = false

	return nil, 0, isStep
}

func (o LoginService) GetStepTypeLi(appKey string) []int {
	return []int{3, 5, 6, 7, 9, 12, 14, 15, 16, 18, 19, 20}
}
