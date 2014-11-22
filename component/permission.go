package component

import (
	"github.com/hongjinqiu/gometa/global"
	"strconv"
	"fmt"
)

type PermissionSupport struct{}

func (o PermissionSupport) GetPermissionQueryDict(sessionId int, security Security) map[string]interface{} {
	if security.ByUnit == "true" {
		userIdI := global.GetGlobalAttr(sessionId, "userId")
		userId, err := strconv.Atoi(fmt.Sprint(userIdI))
		if err != nil {
			panic(err)
		}
		querySupport := QuerySupport{}
		session, _ := global.GetConnection(sessionId)
		collectionName := "SysUser"
		sysUser, found := querySupport.FindByMapWithSession(session, collectionName, map[string]interface{}{
			"_id": userId,
		})
		if found {
			sysUserMaster := sysUser["A"].(map[string]interface{})
			return map[string]interface{}{
				"A.createUnit": sysUserMaster["createUnit"],
			}
		}
		return map[string]interface{}{
			"_id": -1,
		}
	}
	if security.ByAdmin == "true" {
		adminUserId := global.GetGlobalAttr(sessionId, "adminUserId")
		if adminUserId != nil && fmt.Sprint(adminUserId) != "" {
			return map[string]interface{}{}
		}
		return map[string]interface{}{
			"_id": -1,
		}
	}
	return map[string]interface{}{}
}
