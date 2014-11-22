package filter

import (
	"encoding/json"
	. "github.com/hongjinqiu/gometa/error"
	"net/http"
	"reflect"
	"github.com/hongjinqiu/gometa/config"
	"github.com/hongjinqiu/gometa/log"
	"runtime/debug"
)

type FilterFunc func(w http.ResponseWriter, r *http.Request, li []FilterFunc)

func BusinessPanicFilter(w http.ResponseWriter, r *http.Request, li []FilterFunc) {
	defer func() {
		if x := recover(); x != nil {
			if reflect.TypeOf(x).Name() == "BusinessError" {
				err := x.(BusinessError)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				jsonData, jsonErr := json.Marshal(&map[string]interface{}{
					"success": false,
					"code":    err.Code,
					"message": err.Error(),
				})
				if jsonErr != nil {
					panic(jsonErr)
				}
				w.Write(jsonData)
			} else {
				if config.String("debug") != "true" {
					log.Error(x, "\n", string(debug.Stack()))
				}
				panic(x)
			}
		}
	}()
	li[0](w, r, li[1:])
}

func UTF8HtmlFilter(w http.ResponseWriter, r *http.Request, li []FilterFunc) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	li[0](w, r, li[1:])
}
