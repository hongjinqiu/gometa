package app

import (
	. "github.com/hongjinqiu/gometa/common"
	. "github.com/hongjinqiu/gometa/component"
	. "github.com/hongjinqiu/gometa/error"
	"github.com/hongjinqiu/gometa/global"
	. "github.com/hongjinqiu/gometa/lock"
	"github.com/hongjinqiu/gometa/mongo"
	. "github.com/hongjinqiu/gometa/mongo"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var rwlock sync.RWMutex = sync.RWMutex{}
var unitLockDict map[int]sync.RWMutex = map[int]sync.RWMutex{}

type StepService struct{}

func (o StepService) getUnitLock(int) {

}

func (o StepService) Run(sysUser map[string]interface{}) {
	sysUserManster := sysUser["A"].(map[string]interface{})
	sysUnitId := CommonUtil{}.GetIntFromMap(sysUserManster, "createUnit")
	// 加锁
	lockService := LockService{}
	unitLock := lockService.GetUnitLock(fmt.Sprint(sysUnitId))
	(*unitLock).Lock()
	defer (*unitLock).Unlock()

	sessionId := global.GetSessionId()
	defer global.CloseSession(sessionId)
	defer global.RollbackTxn(sessionId)

	txnId := global.GetTxnId(sessionId)
	_, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	stepLi := []map[string]interface{}{}
	err := db.C("SysStep").Find(map[string]interface{}{
		"A.sysUnitId": sysUnitId,
	}).Sort("A.type").All(&stepLi)
	if err != nil {
		panic(err)
	}

	for _, item := range stepLi {
		master := item["A"].(map[string]interface{})
		item["A"] = master

		status := fmt.Sprint(master["status"])
		if status == "1" { // 未开始
			stepType := fmt.Sprint(master["type"])
			if stepType == "3" { //初始化供应商类别
				o.InitProviderType(sessionId, sysUser)
			} else if stepType == "5" { //初始化币别
				o.InitCurrencyType(sessionId, sysUser)
			} else if stepType == "6" { //初始化银行资料
				o.InitBank(sessionId, sysUser)
			} else if stepType == "7" { //初始化计量单位
				o.InitMeasureUnit(sessionId, sysUser)
			} else if stepType == "9" { //初始化客户类别
				o.InitCustomerType(sessionId, sysUser)
			} else if stepType == "12" { //初始化税率类别
				o.InitTaxType(sessionId, sysUser)
			} else if stepType == "14" { //初始化收入费用类别
				o.InitIncomeType(sessionId, sysUser)
			} else if stepType == "15" { //初始化收入费用项目
				o.InitIncomeItem(sessionId, sysUser)
			} else if stepType == "16" { //初始化会计期
				o.InitAccountingPeriod(sessionId, sysUser)
			} else if stepType == "18" { //初始化收款单类型参数
				o.InitBillReceiveTypeParameter(sessionId, sysUser)
			} else if stepType == "19" { //初始化付款单类型参数
				o.InitBillPaymentTypeParameter(sessionId, sysUser)
			} else if stepType == "20" { //初始化系统参数
				o.InitSystemParameter(sessionId, sysUser)
			}
		}
		master["status"] = 2
		_, updateResult := txnManager.Update(txnId, "SysStep", item)
		if !updateResult {
			panic(BusinessError{Message: "更新SysStep" + fmt.Sprint(master["name"]) + "失败"})
		}
	}

	global.CommitTxn(sessionId)
}

/**
初始化供应商类别
*/
func (o StepService) InitProviderType(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "SUP10",
			"name": "厂商",
		},
		map[string]interface{}{
			"code": "SUP20",
			"name": "非厂商",
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "ProviderType", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "providerTypeId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "ProviderType", data)
		}
	}
}

/**
初始化币别
*/
func (o StepService) InitCurrencyType(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code":             "RMB",
			"name":             "人民币",
			"currencyTypeSign": "",
			"roundingWay":      2, // 四舍五入
			"amtDecimals":      3, // 3代表2位小数
			"upDecimals":       3, // 3代表2位小数
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "CurrencyType", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "currencyTypeId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "CurrencyType", data)
		}
	}
}

/**
初始化银行资料
*/
func (o StepService) InitBank(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "BPOSZT", "name": "POS在途资金银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95599", "name": "中国农业银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95595", "name": "光大银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95588", "name": "中国工商银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95580", "name": "邮政储蓄", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95577", "name": "华夏银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95568", "name": "中国民生银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95566", "name": "中国银行", "bankShort": "中行", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95561", "name": "兴业银行", "bankShort": "兴业", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95559", "name": "交通银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95558", "name": "中信银行", "bankShort": "", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95555", "name": "招商银行", "bankShort": "招行", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
		map[string]interface{}{
			"code": "B95533", "name": "中国建设银行", "bankShort": "建行", "linkman": "", "bankUrl": "", "linkPhone": "", "cusPhone": "", "complainPhone": "", "bankAddress": "",
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "Bank", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "bankId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "Bank", data)
		}
	}
}

/**
初始化计量单位
*/
func (o StepService) InitMeasureUnit(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "L", "name": "辆", "type": 1, "unitDecimals": 1, "roundingWay": 1,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "MeasureUnit", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "measureUnitId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "MeasureUnit", data)
		}
	}
}

/**
初始化客户类别
*/
func (o StepService) InitCustomerType(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "CUS10", "name": "个人客户",
		},
		map[string]interface{}{
			"code": "CUS20", "name": "法人客户",
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "CustomerType", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "customerTypeId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "CustomerType", data)
		}
	}
}

/**
初始化税率类别
*/
func (o StepService) InitTaxType(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "VAT7", "name": "7%增值税扣除税率", "taxRate": 7, "isDeductTax": 2, "isDeduct": 2,
		},
		map[string]interface{}{
			"code": "VAT5", "name": "5%增值税", "taxRate": 5, "isDeductTax": 1, "isDeduct": 2,
		},
		map[string]interface{}{
			"code": "VAT17", "name": "17%增值税", "taxRate": 17, "isDeductTax": 1, "isDeduct": 2,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "TaxType", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "taxTypeId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "TaxType", data)
		}
	}
}

/**
初始化收入费用类别
*/
func (o StepService) InitIncomeType(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "IN10", "name": "收入",
		},
		map[string]interface{}{
			"code": "EX10", "name": "费用",
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "IncomeType", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "incomeTypeId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "IncomeType", data)
		}
	}
}

/**
初始化收入费用项目
*/
func (o StepService) InitIncomeItem(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"code": "SHYYSR", "name": "售后营业收入", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "LSZLSR", "name": "租赁收入（零售）", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "SCZC", "name": "非常规市场推广支持", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "FPSR", "name": "废品收入", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "HYF", "name": "会员费收入", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "GZZJ", "name": "固定资产折旧", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "ZXF", "name": "装修成本摊销", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "JBGZ", "name": "基本工资总额", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "SBF", "name": "公司承担社保", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "CDZJ", "name": "场地租金", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "SCF", "name": "市场投入", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "GZCZSY", "name": "固定资产处置损益", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "ZJCB", "name": "资金成本", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "PXF", "name": "职员培训费用", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "BGF", "name": "办公费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "POS", "name": "POS手续费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "YHSXF", "name": "其他银行手续费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "JTF", "name": "公务交通费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "CLF", "name": "差旅费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "SDF", "name": "水电费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "WLF", "name": "电信网络等费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "XLF", "name": "设备及车修理费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "FLF", "name": "福利费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "WYF", "name": "物业费与公修金", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "BAF", "name": "保安费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "ZLCB", "name": "租赁成本", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "AQFXJ", "name": "安全风险基金", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "JZZC", "name": "捐赠支出", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "YWZDF", "name": "业务招待费", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "YZF_SH", "name": "运杂费（售后专用）", "incomeTypeId": 2,
		},
		map[string]interface{}{
			"code": "YWTGF", "name": "业务推广费", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "SJ", "name": "税金", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
		map[string]interface{}{
			"code": "WXZR", "name": "维修折让", "incomeTypeId": 1, "taxTypeId": 0, "measureUnitId": 0,
		},
	}
	taxType1Query := map[string]interface{}{
		"A.code":       "IN10",
		"A.createUnit": sysUserMaster["createUnit"],
	}
	taxType1, found := qb.FindByMapWithSession(session, "IncomeType", taxType1Query)
	if !found {
		queryByte, err := json.MarshalIndent(&taxType1Query, "", "\t")
		if err != nil {
			panic(err)
		}
		panic(BusinessError{Message: "未找到收入费用类别，查询条件为：" + string(queryByte)})
	}
	taxType1Id := taxType1["id"]

	taxType2Query := map[string]interface{}{
		"A.code":       "EX10",
		"A.createUnit": sysUserMaster["createUnit"],
	}
	taxType2, found := qb.FindByMapWithSession(session, "IncomeType", taxType2Query)
	if !found {
		queryByte, err := json.MarshalIndent(&taxType2Query, "", "\t")
		if err != nil {
			panic(err)
		}
		panic(BusinessError{Message: "未找到收入费用类别，查询条件为：" + string(queryByte)})
	}
	taxType2Id := taxType2["id"]

	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "IncomeItem", map[string]interface{}{
			"A.code":       initData["code"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "incomeItemId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				if k == "incomeTypeId" && fmt.Sprint(v) == "1" {
					aData[k] = taxType1Id
				} else if k == "incomeTypeId" && fmt.Sprint(v) == "2" {
					aData[k] = taxType2Id
				} else {
					aData[k] = v
				}
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "IncomeItem", data)
		}
	}
}

/**
初始化会计期
*/
func (o StepService) InitAccountingPeriod(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"accountingYear": 2014, "numAccountingPeriod": 12,
		},
		map[string]interface{}{
			"accountingYear": 2015, "numAccountingPeriod": 12,
		},
		map[string]interface{}{
			"accountingYear": 2016, "numAccountingPeriod": 12,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "AccountingPeriod", map[string]interface{}{
			"A.accountingYear":       initData["accountingYear"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "accountingPeriodId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			detailDataLi := []interface{}{}
			year := aData["accountingYear"].(int)
			for i := 0; i < aData["numAccountingPeriod"].(int); i++ {
				id := mongo.GetSequenceNo(db, "accountingPeriodId")
				data := map[string]interface{}{
					"createBy":    sysUserMaster["id"],
					"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
					"createUnit":  sysUserMaster["createUnit"],
					"modifyBy":    0,
					"modifyTime":  0,
					"modifyUnit":  0,
					"billStatus":  0,
					"attachCount": 0,
					"remark":      "",
				}
				data["id"] = id
				data["sequenceNo"] = i + 1
				numStr := fmt.Sprint(i + 1)
				if i+1 < 10 {
					numStr = "0" + numStr
				}
				startDateStr := fmt.Sprint(year) + numStr + "01"
				startDate, err := strconv.Atoi(startDateStr)
				if err != nil {
					panic(err)
				}
				data["startDate"] = startDate
				startTime, err := time.Parse("20060102", startDateStr)
				if err != nil {
					panic(err)
				}
				nextMonthTime := startTime.AddDate(0, 1, -1)
				data["endDate"], err = strconv.Atoi(nextMonthTime.Format("20060102"))
				if err != nil {
					panic(err)
				}
				detailDataLi = append(detailDataLi, data)
			}
			data["B"] = detailDataLi
			txnManager.Insert(txnId, "AccountingPeriod", data)
		}
	}
}

/**
初始化收款单类型参数
*/
func (o StepService) InitBillReceiveTypeParameter(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"billTypeId": 1,
			"property": 2,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "BillTypeParameter", map[string]interface{}{
			"A.billTypeId":       initData["billTypeId"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "billTypeParameterId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "BillTypeParameter", data)
		}
	}
}

/**
初始化付款单类型参数
*/
func (o StepService) InitBillPaymentTypeParameter(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"billTypeId": 2,
			"property": 2,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "BillTypeParameter", map[string]interface{}{
			"A.billTypeId":       initData["billTypeId"],
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "billTypeParameterId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "BillTypeParameter", data)
		}
	}
}

/**
初始化系统参数
*/
func (o StepService) InitSystemParameter(sessionId int, sysUser map[string]interface{}) {
	sysUserMaster := sysUser["A"].(map[string]interface{})
	txnId := global.GetTxnId(sessionId)
	session, db := global.GetConnection(sessionId)
	txnManager := TxnManager{db}
	qb := QuerySupport{}
	initDataLi := []map[string]interface{}{
		map[string]interface{}{
			"percentDecimals": 3,
			"percentRoundingWay": 2,
			"thousandDecimals": 2,
//			"currencyTypeId": xxxx,
			"costDecimals": 3,
//			"taxTypeId": xxxx,
		},
	}
	for _, initData := range initDataLi {
		_, found := qb.FindByMapWithSession(session, "SystemParameter", map[string]interface{}{
			"A.createUnit": sysUserMaster["createUnit"],
		})
		if !found {
			id := mongo.GetSequenceNo(db, "systemParameterId")
			aData := map[string]interface{}{
				"id":          id,
				"createBy":    sysUserMaster["id"],
				"createTime":  DateUtil{}.GetCurrentYyyyMMddHHmmss(),
				"createUnit":  sysUserMaster["createUnit"],
				"modifyBy":    0,
				"modifyTime":  0,
				"modifyUnit":  0,
				"billStatus":  0,
				"attachCount": 0,
				"remark":      "",
			}
			for k, v := range initData {
				aData[k] = v
			}
			// 币别
			currencyType, found := qb.FindByMapWithSession(session, "CurrencyType", map[string]interface{}{
				"A.createUnit": sysUserMaster["createUnit"],
			})
			if found {
				aData["currencyTypeId"] = currencyType["id"]
			}
			taxType, found := qb.FindByMapWithSession(session, "TaxType", map[string]interface{}{
				"A.createUnit": sysUserMaster["createUnit"],
			})
			if found {
				aData["taxTypeId"] = taxType["id"]
			}
			data := map[string]interface{}{
				"_id": id,
				"id":  id,
				"A":   aData,
			}
			txnManager.Insert(txnId, "SystemParameter", data)
		}
	}
}
