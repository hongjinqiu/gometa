package interceptor

import (
	. "github.com/hongjinqiu/gometa/common"
	"github.com/hongjinqiu/gometa/global"
	"fmt"
	"log"
	"strconv"
)

type AccountInOutItemInterceptor struct{}

func (o AccountInOutItemInterceptor) BeforeBuildQuery(sessionId int, paramMap map[string]string) map[string]string {
	queryMode := paramMap["queryMode"]
	if queryMode == "1" { // 按日期查询
		billDateBegin := paramMap["billDateBegin"]
		if billDateBegin != "" && billDateBegin != "0" {
			paramMap["qBillDateBegin"] = billDateBegin
		}
		billDateEnd := paramMap["billDateEnd"]
		if billDateEnd != "" && billDateEnd != "0" {
			paramMap["qBillDateEnd"] = billDateEnd
		}
	} else if queryMode == "2" { // 按期间查询
		billDateBegin := 19700101
		billDateEnd := 99991231

		commonUtil := CommonUtil{}
		//		minAccountingPeriod, _ := o.GetMinAccountingPeriod(sessionId)
		//		minAccountingPeriodData := minAccountingPeriod["A"].(map[string]interface{})
		//		minAccountingPeriodYear := commonUtil.GetIntFromMap(minAccountingPeriodData, "accountingYear")

		maxAccountingPeriod, _ := o.GetMaxAccountingPeriod(sessionId)
		maxAccountingPeriodData := maxAccountingPeriod["A"].(map[string]interface{})
		maxAccountingPeriodYear := commonUtil.GetIntFromMap(maxAccountingPeriodData, "accountingYear")

		accountingYearStart := commonUtil.GetIntFromString(paramMap["accountingYearStart"])
		accountingPeriodStart := commonUtil.GetIntFromString(paramMap["accountingPeriodStart"])
		accountingYearEnd := commonUtil.GetIntFromString(paramMap["accountingYearEnd"])
		accountingPeriodEnd := commonUtil.GetIntFromString(paramMap["accountingPeriodEnd"])

		if accountingYearStart == 0 {
			billDateBegin = 19700101
		} else if accountingYearStart > maxAccountingPeriodYear {
			billDateBegin = 99991231
		} else { // accountingYearStart >= minAccountingPeriodYear && accountingYearStart <= maxAccountingPeriodYear
			accountingPeriod, _ := o.GetAccountingPeriod(sessionId, accountingYearStart)
			bDataSetLi := accountingPeriod["B"].([]interface{})
			//			firstSequenceNo := commonUtil.GetIntFromMap(bDataSetLi[0].(map[string]interface{}), "sequenceNo")
			lastSequenceNo := commonUtil.GetIntFromMap(bDataSetLi[len(bDataSetLi)-1].(map[string]interface{}), "sequenceNo")
			if accountingPeriodStart == 0 {
				billDateBegin = commonUtil.GetIntFromMap(bDataSetLi[0].(map[string]interface{}), "startDate")
			} else if accountingPeriodStart > lastSequenceNo {
				endDate := commonUtil.GetIntFromMap(bDataSetLi[len(bDataSetLi)-1].(map[string]interface{}), "endDate")
				dateUtil := DateUtil{}
				billDateBegin = dateUtil.GetNextDate(endDate)
			} else {
				for _, item := range bDataSetLi {
					line := item.(map[string]interface{})
					sequenceNo := commonUtil.GetIntFromMap(line, "sequenceNo")
					if sequenceNo == accountingPeriodStart {
						billDateBegin = commonUtil.GetIntFromMap(line, "startDate")
						break
					}
				}
			}
		}

		if accountingYearEnd == 0 {
			billDateEnd = 99991231
		} else if accountingYearEnd > maxAccountingPeriodYear {
			billDateEnd = 99991231
		} else { // accountingYearEnd >= minAccountingPeriodYear && accountingYearEnd <= maxAccountingPeriodYear
			accountingPeriod, _ := o.GetAccountingPeriod(sessionId, accountingYearEnd)
			bDataSetLi := accountingPeriod["B"].([]interface{})
			//			firstSequenceNo := commonUtil.GetIntFromMap(bDataSetLi[0].(map[string]interface{}), "sequenceNo")
			lastSequenceNo := commonUtil.GetIntFromMap(bDataSetLi[len(bDataSetLi)-1].(map[string]interface{}), "sequenceNo")
			if accountingPeriodEnd == 0 {
				billDateEnd = commonUtil.GetIntFromMap(bDataSetLi[len(bDataSetLi)-1].(map[string]interface{}), "endDate")
			} else if accountingPeriodEnd > lastSequenceNo {
				billDateEnd = commonUtil.GetIntFromMap(bDataSetLi[len(bDataSetLi)-1].(map[string]interface{}), "endDate")
			} else {
				for _, item := range bDataSetLi {
					line := item.(map[string]interface{})
					sequenceNo := commonUtil.GetIntFromMap(line, "sequenceNo")
					if sequenceNo == accountingPeriodEnd {
						billDateEnd = commonUtil.GetIntFromMap(line, "endDate")
						break
					}
				}
			}
		}

		paramMap["qBillDateBegin"] = fmt.Sprint(billDateBegin)
		paramMap["qBillDateEnd"] = fmt.Sprint(billDateEnd)
	}
	return paramMap
}

func (o AccountInOutItemInterceptor) GetMinAccountingPeriod(sessionId int) (map[string]interface{}, bool) {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	session, _ := global.GetConnection(sessionId)
	collectionName := "AccountingPeriod"
	interceptorCommon := InterceptorCommon{}
	queryMap := map[string]interface{}{
		"A.createUnit": interceptorCommon.GetCreateUnitByUserId(session, userId),
	}
	pageNo := 1
	pageSize := 1
	orderBy := "A.accountingYear"
	queryResult := interceptorCommon.IndexWithSession(session, collectionName, queryMap, pageNo, pageSize, orderBy)
	items := queryResult["items"].([]interface{})
	if len(items) > 0 {
		return items[0].(map[string]interface{}), true
	}
	return nil, false
}

func (o AccountInOutItemInterceptor) GetMaxAccountingPeriod(sessionId int) (map[string]interface{}, bool) {
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	session, _ := global.GetConnection(sessionId)
	collectionName := "AccountingPeriod"
	interceptorCommon := InterceptorCommon{}
	queryMap := map[string]interface{}{
		"A.createUnit": interceptorCommon.GetCreateUnitByUserId(session, userId),
	}
	pageNo := 1
	pageSize := 1
	orderBy := "-A.accountingYear"
	queryResult := interceptorCommon.IndexWithSession(session, collectionName, queryMap, pageNo, pageSize, orderBy)
	items := queryResult["items"].([]interface{})
	if len(items) > 0 {
		return items[0].(map[string]interface{}), true
	}
	return nil, false
}

func (o AccountInOutItemInterceptor) GetAccountingPeriod(sessionId int, year int) (map[string]interface{}, bool) {
	session, _ := global.GetConnection(sessionId)
	collectionName := "AccountingPeriod"

	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}

	queryMap := map[string]interface{}{
		"A.accountingYear": year,
		"A.createUnit":     InterceptorCommon{}.GetCreateUnitByUserId(session, userId),
	}
	accountingPeriod, found := InterceptorCommon{}.FindByMapWithSession(session, collectionName, queryMap)
	if !found {
		//		panic(BusinessError{Message: "会计年度:" + fmt.Sprint(year) + ",会计期序号:" + fmt.Sprint(sequenceNo) + "未找到对应会计期"})
		log.Println("会计年度:" + fmt.Sprint(year) + "未找到对应会计期")
		return nil, false
	}
	return accountingPeriod, true
}

func (o AccountInOutItemInterceptor) GetAccountingPeriodFirstStartEndDate(sessionId int, year int) (int, int) {
	session, _ := global.GetConnection(sessionId)
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	collectionName := "AccountingPeriod"
	queryMap := map[string]interface{}{
		"A.accountingYear": year,
		"A.createUnit":     InterceptorCommon{}.GetCreateUnitByUserId(session, userId),
	}
	accountingPeriod, found := InterceptorCommon{}.FindByMapWithSession(session, collectionName, queryMap)
	if !found {
		//		panic(BusinessError{Message: "会计年度:" + fmt.Sprint(year) + ",会计期序号:" + fmt.Sprint(sequenceNo) + "未找到对应会计期"})
		log.Println("会计年度:" + fmt.Sprint(year) + "未找到对应会计期")
		return 0, 0
	}
	var startDate int
	var endDate int
	bDataSetLi := accountingPeriod["B"].([]interface{})
	commonUtil := CommonUtil{}
	for _, item := range bDataSetLi {
		line := item.(map[string]interface{})
		startDate = commonUtil.GetIntFromMap(line, "startDate")
		endDate = commonUtil.GetIntFromMap(line, "endDate")
		break
	}
	return startDate, endDate
}

func (o AccountInOutItemInterceptor) GetAccountingPeriodLastStartEndDate(sessionId int, year int) (int, int) {
	session, _ := global.GetConnection(sessionId)
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	collectionName := "AccountingPeriod"
	queryMap := map[string]interface{}{
		"A.accountingYear": year,
		"A.createUnit":     InterceptorCommon{}.GetCreateUnitByUserId(session, userId),
	}
	accountingPeriod, found := InterceptorCommon{}.FindByMapWithSession(session, collectionName, queryMap)
	if !found {
		//		panic(BusinessError{Message: "会计年度:" + fmt.Sprint(year) + ",会计期序号:" + fmt.Sprint(sequenceNo) + "未找到对应会计期"})
		log.Println("会计年度:" + fmt.Sprint(year) + "未找到对应会计期")
		return 0, 0
	}
	var startDate int
	var endDate int
	bDataSetLi := accountingPeriod["B"].([]interface{})
	commonUtil := CommonUtil{}
	if len(bDataSetLi) > 0 {
		line := bDataSetLi[len(bDataSetLi)-1].(map[string]interface{})
		startDate = commonUtil.GetIntFromMap(line, "startDate")
		endDate = commonUtil.GetIntFromMap(line, "endDate")
	}
	return startDate, endDate
}

func (o AccountInOutItemInterceptor) GetAccountingPeriodStartEndDate(sessionId int, year int, sequenceNo int) (int, int) {
	session, _ := global.GetConnection(sessionId)
	userId, err := strconv.Atoi(fmt.Sprint(global.GetGlobalAttr(sessionId, "userId")))
	if err != nil {
		panic(err)
	}
	collectionName := "AccountingPeriod"
	queryMap := map[string]interface{}{
		"A.accountingYear": year,
		"B.sequenceNo":     sequenceNo,
		"A.createUnit":     InterceptorCommon{}.GetCreateUnitByUserId(session, userId),
	}
	accountingPeriod, found := InterceptorCommon{}.FindByMapWithSession(session, collectionName, queryMap)
	if !found {
		//		panic(BusinessError{Message: "会计年度:" + fmt.Sprint(year) + ",会计期序号:" + fmt.Sprint(sequenceNo) + "未找到对应会计期"})
		log.Println("会计年度:" + fmt.Sprint(year) + ",会计期序号:" + fmt.Sprint(sequenceNo) + "未找到对应会计期")
		return 0, 0
	}
	var startDate int
	var endDate int
	bDataSetLi := accountingPeriod["B"].([]interface{})
	commonUtil := CommonUtil{}
	for _, item := range bDataSetLi {
		line := item.(map[string]interface{})
		if fmt.Sprint(line["sequenceNo"]) == fmt.Sprint(sequenceNo) {
			startDate = commonUtil.GetIntFromMap(line, "startDate")
			endDate = commonUtil.GetIntFromMap(line, "endDate")
			break
		}
	}
	return startDate, endDate
}
