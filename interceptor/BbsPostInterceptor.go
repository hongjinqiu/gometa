package interceptor

import (
	"github.com/hongjinqiu/gometa/global"
	"fmt"
	"strconv"
	. "github.com/hongjinqiu/gometa/common"
)

type BbsPostInterceptor struct{}

func (o BbsPostInterceptor) AfterQueryData(sessionId int, dataSetId string, items []interface{}) []interface{}  {
	session, _ := global.GetConnection(sessionId)
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	collectionName := "BbsPostRead"
	interceptorCommon := InterceptorCommon{}
	commonUtil := CommonUtil{}
	for i, item := range items {
		data := item.(map[string]interface{})
		dataA := data["A"].(map[string]interface{})
		data["A"] = dataA
		items[i] = data
		
		bbsPostReadQuery := map[string]interface{}{
			"A.bbsPostId": data["id"],
			"A.readBy": userId,
		}
		bbsPostRead, found := interceptorCommon.FindByMapWithSession(session, collectionName, bbsPostReadQuery)
		if !found {
			data["bbsPostReadType"] = 1// 未读
		} else {
			bbsPostReadA := bbsPostRead["A"].(map[string]interface{})
			lastReadTime := commonUtil.GetFloat64FromMap(bbsPostReadA, "lastReadTime")
			lastReplyTime := commonUtil.GetFloat64FromMap(dataA, "lastReplyTime")
			if lastReadTime < lastReplyTime {
				dataA["bbsPostReadType"] = 1// 未读
			} else {
				dataA["bbsPostReadType"] = 2// 已读
			}
		}
	}
	return items
}
