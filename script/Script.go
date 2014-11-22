package script

import (
	"sync"
	"reflect"
	"strings"
)

var rwlock sync.RWMutex = sync.RWMutex{}
var scriptDict map[string]reflect.Type = map[string]reflect.Type{}

func init() {
	rwlock.Lock()
	defer rwlock.Unlock()
	//scriptDict[reflect.TypeOf(SysUser{}).Name()] = reflect.TypeOf(SysUser{})
}

func GetScriptDict() map[string]reflect.Type {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return scriptDict
}

type ScriptManager struct{}

func (o ScriptManager) Parse(classMethod string, param []interface{}) []interface{} {
	exprContent := classMethod
	scriptStruct := strings.Split(exprContent, ".")[0]
	scriptStructMethod := strings.Split(exprContent, ".")[1]
	scriptType := GetScriptDict()[scriptStruct]
	inst := reflect.New(scriptType).Elem().Interface()
	instValue := reflect.ValueOf(inst)
	in := []reflect.Value{}
	for i, _ := range param {
		in = append(in, reflect.ValueOf(param[i]))
	}
	callValues := instValue.MethodByName(scriptStructMethod).Call(in)
	result := []interface{}{}
	for i, _ := range callValues {
		result = append(result, callValues[i].Interface())
	}
	return result
}
