package route

import (
	"github.com/hongjinqiu/gometa/config"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type HttpHandleFilterFunc func(http.ResponseWriter, *http.Request, []HttpHandleFilterFunc)

var FilterLi []HttpHandleFilterFunc = []HttpHandleFilterFunc{}

func init() {
	//	FilterLi = append(FilterLi, func(w http.ResponseWriter, r *http.Request, li []HttpHandleFilterFunc){
	//		println("^^^^ before filter 0")
	//		li[0](w, r, li[1:])
	//		println("^^^^ after filter 0")
	//	})
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, filterHandleFunc(handler))
}

func filterHandleFunc(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	if len(FilterLi) == 0 {
		return handler
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			tmpFilterLi := FilterLi
			tmpFilterLi = append(tmpFilterLi, func(w http.ResponseWriter, r *http.Request, li []HttpHandleFilterFunc) {
				handler(w, r)
			})
			tmpFilterLi[0](w, r, tmpFilterLi[1:])
		}
	}
}

func RegisteStaticFilePath() {
	staticPath := "/public/"
	HandleFunc(staticPath, func(w http.ResponseWriter, r *http.Request) {
		urlStaticPathPart := r.URL.Path[len(staticPath):]
		if os.PathSeparator != '/' {
			urlStaticPathPart = strings.Replace(urlStaticPathPart, "/", string(os.PathSeparator), -1)
		}
		// 删除外部设置的header,让go自己侦测content-type
		delete(w.Header(), "Content-Type")
		file := filepath.Join(config.String("gmeta.static"), urlStaticPathPart)
		http.ServeFile(w, r, file)
	})
}

func RegisteReflectController(constrollersDict []reflect.Type) {
	for _, item := range constrollersDict {
		structNameLi := []string{item.Name()}
		if strings.ToLower(item.Name()) != item.Name() {
			structNameLi = append(structNameLi, strings.ToLower(item.Name()))
		}
		for i := 0; i < item.NumMethod(); i++ {
			methodNameLi := []string{}
			firstCharacter := item.Method(i).Name[0]
			if 'A' <= firstCharacter && firstCharacter <= 'Z' {
				methodNameLi = append(methodNameLi, item.Method(i).Name)
				methodNameLi = append(methodNameLi, strings.ToLower(item.Method(i).Name))
			}

			for _, structName := range structNameLi {
				for _, methodName := range methodNameLi {
					for _, subfix := range []string{"", "/"} {
						HandleFunc("/"+structName+"/"+methodName+subfix, func(reflectType reflect.Type, index int) func(http.ResponseWriter, *http.Request) {
							return func(w http.ResponseWriter, r *http.Request) {
								inst := reflect.New(reflectType).Elem().Interface()
								instValue := reflect.ValueOf(inst)
								method := instValue.Method(index)
								in := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
								method.Call(in)
							}
						}(item, i))
					}
				}
			}

		}
	}
}
