package script

import (
	"sync"
	"fmt"
//	"os"
//	"os/exec"
	"strings"
	"github.com/hongjinqiu/gometa/log"
//	"regexp"
//	"strconv"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"github.com/hongjinqiu/gometa/config"
	"github.com/robertkrimen/otto"
)

var pyrwlock sync.RWMutex = sync.RWMutex{}
var flag bool = false
var jsvm *otto.Otto = otto.New()

type ExpressionParser struct{}

func (o ExpressionParser) ParseGolang(bo map[string]interface{}, data map[string]interface{}, expression string) string {
	scriptManager := ScriptManager{}
	//Parse(classMethod string, param []interface{}) []interface{} {
	paramLi := []interface{}{bo, data}
	values := scriptManager.Parse(expression, paramLi)
	return fmt.Sprint(values[0])
}

/**
```go
func (self Otto) Call(source string, this interface{}, argumentList ...interface{}) (Value, error)
```
value, _ := vm.Call(`[ 1, 2, 3, undefined, 4 ].concat`, nil, 5, 6, 7, "abc")

```go
func (self Otto) Run(src interface{}) (Value, error)
```

*/
func (o ExpressionParser) parseExpression(recordJson, expression string) string {
	if recordJson == "" {
		recordJson = "{}"
	}
	jsFunc := `(function(){
		var data = {RECORD_JSON};
		return {EXPRESSION};
	})`
	jsFunc = strings.Replace(jsFunc, "{RECORD_JSON}", recordJson, -1)
	jsFunc = strings.Replace(jsFunc, "{EXPRESSION}", expression, -1)
	value, err := jsvm.Call(jsFunc, nil)
	if err != nil {
		log.Error("Parse(recordJson, expression string) bool")
		log.Error("recordJson:" + recordJson)
		log.Error("expression:" + expression)
		panic(err)
	}
	return value.String()
}

func (o ExpressionParser) Parse(recordJson, expression string) bool {
	if expression == "" {
		return true
	}
	return strings.ToLower(o.parseExpression(recordJson, expression)) == "true"
}

func (o ExpressionParser) Validate(text, expression string) bool {
	if text == "" || expression == "" {
		return true
	}
	
	return strings.ToLower(o.parseExpression(text, expression)) == "true"
}

func (o ExpressionParser) ParseString(recordJson, expression string) string {
	if expression == "" {
		return ""
	}
	return o.parseExpression(recordJson, expression)
}

func (o ExpressionParser) ParseModel(boJson, dataJson, expression string) string {
	if expression == "" {
		return ""
	}
	if boJson == "" {
		boJson = "{}"
	}
	if dataJson == "" {
		dataJson = "{}"
	}
	jsFunc := `(function(){
		var bo = {BO};
		var data = {RECORD_JSON};
		return {EXPRESSION};
	})`
	jsFunc = strings.Replace(jsFunc, "{BO}", boJson, -1)
	jsFunc = strings.Replace(jsFunc, "{RECORD_JSON}", dataJson, -1)
	jsFunc = strings.Replace(jsFunc, "{EXPRESSION}", expression, -1)
	value, err := jsvm.Call(jsFunc, nil)
	if err != nil {
		log.Error("ParseModel(boJson, dataJson, expression string) string")
		log.Error("boJson:" + boJson)
		log.Error("dataJson:" + dataJson)
		log.Error("expression:" + expression)
		panic(err)
	}
	return value.String()
}

//func (o ExpressionParser) ParseBeforeBuildQuery(classMethod string, paramMap map[string]string) map[string]string {
//	return paramMap
//}
//
//func (o ExpressionParser) ParseAfterBuildQuery(classMethod string, queryLi []map[string]interface{}) []map[string]interface{} {
//	return queryLi
//}

/*
func (o ExpressionParser) ParseAfterQueryData(classMethod string, items []interface{}) []interface{} {
	return items
}
*/
