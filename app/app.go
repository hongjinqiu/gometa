package app

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	. "github.com/hongjinqiu/gometa/common"
	"github.com/hongjinqiu/gometa/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

var gzipRwlock sync.RWMutex = sync.RWMutex{}
var isRunTxnPeriod bool = false
var periodRwlock sync.RWMutex = sync.RWMutex{}

type App struct{}

type StringArraySort struct {
	objLi []string
}

func (o StringArraySort) Len() int {
	return len(o.objLi)
}

func (o StringArraySort) Less(i, j int) bool {
	return o.objLi[i] <= o.objLi[j]
}

func (o StringArraySort) Swap(i, j int) {
	o.objLi[i], o.objLi[j] = o.objLi[j], o.objLi[i]
}

func (self App) getFileNameConcatFromQuery(r *http.Request) string {
	queryLi := []string{}
	name := ""
	commonUtil := CommonUtil{}
	for k := range r.URL.Query() {
		//		name += k
		if !commonUtil.IsNumber(k) && k != "" {
			queryLi = append(queryLi, k)
		}
	}
	stringArraySort := StringArraySort{queryLi}
	sort.Sort(stringArraySort)
	name = strings.Join(stringArraySort.objLi, "")
	return name
}

func (self App) getComboFileContent(r *http.Request) string {
	jsPath := config.String("JS_PATH")
	content := ""
	commonUtil := CommonUtil{}
	for k := range r.URL.Query() {
		if !commonUtil.IsNumber(k) && k != "" {
			file, err := os.Open(path.Join(jsPath, k))
			defer file.Close()
			if err != nil {
				panic(err)
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				panic(err)
			}
			content += string(data) + "\n"
		}
	}
	return content
}

func (self App) isFileExist(name string) bool {
	h := md5.New()
	io.WriteString(h, name)
	gzipFileName := fmt.Sprintf("%x", h.Sum(nil))
	gzipPath := config.String("GZIP_PATH")
	if _, err := os.Stat(path.Join(gzipPath, gzipFileName)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return true
}

func (self App) getGzipContent(name string) []byte {
	h := md5.New()
	io.WriteString(h, name)
	gzipFileName := fmt.Sprintf("%x", h.Sum(nil))
	gzipPath := config.String("GZIP_PATH")
	filePath := path.Join(gzipPath, gzipFileName)

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (self App) gzipAndSave(name string, content string) []byte {
	gzipRwlock.Lock()
	defer gzipRwlock.Unlock()

	h := md5.New()
	io.WriteString(h, name)
	gzipFileName := fmt.Sprintf("%x", h.Sum(nil))
	gzipPath := config.String("GZIP_PATH")
	filePath := path.Join(gzipPath, gzipFileName)

	if !self.isFileExist(name) {
		data := bytes.Buffer{}
		w := gzip.NewWriter(&data)
		w.Write([]byte(content))
		w.Close()

		bytes := data.Bytes()
		err := ioutil.WriteFile(filePath, bytes, os.ModeDevice|0666)
		if err != nil {
			panic(err)
		}
		return bytes
	} else {
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		return bytes
	}
}

func (self App) Combo(w http.ResponseWriter, r *http.Request) {
	nameConcat := self.getFileNameConcatFromQuery(r)

	acceptEncoding := r.Header.Get("Accept-Encoding")
	if strings.Index(acceptEncoding, "gzip") > -1 {
		text := ""
		if config.String("debug") == "true" {
			content := self.getComboFileContent(r)
			data := bytes.Buffer{}
			gzipW := gzip.NewWriter(&data)
			gzipW.Write([]byte(content))
			gzipW.Close()
			text = data.String()
		} else {
			if self.isFileExist(nameConcat) {
				text = string(self.getGzipContent(nameConcat))
			} else {
				content := self.getComboFileContent(r)
				text = string(self.gzipAndSave(nameConcat, content))
			}
		}

		//		self.Response.Status = http.StatusOK
		if strings.Index(r.URL.RawQuery, ".css") <= -1 {
			w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
		} else {
			w.Header()["Content-Type"] = []string{"text/css;charset=UTF-8"}
		}
		w.Header()["Content-Encoding"] = []string{"gzip"}
		w.Write([]byte(text))
		//		end := time.Now()
		//		println("^^^^^^^^^^^^^^^^url is:" + url + " time spend is:" + fmt.Sprint((end.UnixNano() - start.UnixNano())))
		return
	}

	content := self.getComboFileContent(r)
//	self.Response.Status = http.StatusOK
	if strings.Index(r.URL.RawQuery, ".css") <= -1 {
		w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
	} else {
		w.Header()["Content-Type"] = []string{"text/css;charset=UTF-8"}
	}
	w.Write([]byte(content))
}

func (self App) getFormJsContent() string {
	jsPath := config.String("COMBO_VIEW_PATH")
	content := ""
	formJsLi := self.getFormJsLi()
	// 加入日期标记,gzip到目标文件时,有用,
	//	commonUtil := CommonUtil{}
	//	for k := range self.Params.Query {
	//		if commonUtil.IsNumber(k) && k != "" {
	//
	//		}
	//	}
	for _, k := range formJsLi {
		if strings.Index(k, ".js") == -1 && strings.Index(k, ".css") == -1 {
			panic("fileName is:" + k + ", expect ends with .js or .css")
		}
		isFileExist := false
		for _, filePath := range strings.Split(jsPath, ":") {
			if _, err := os.Stat(path.Join(filePath, k)); err != nil {
				if os.IsNotExist(err) {
					continue
				}
			}
			isFileExist = true
			file, err := os.Open(path.Join(filePath, k))
			defer file.Close()
			if err != nil {
				panic(err)
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				panic(err)
			}
			content += string(data) + "\n"
			break
		}
		if !isFileExist {
			panic(k + " is not exists")
		}
	}
	prefix := "YUI.add('papersns-form', function(Y) {\n"
	suffix := "}, '1.1.0' ,{requires:['node', 'widget-base', 'widget-htmlparser', 'io-form', 'widget-parent', 'widget-child', 'base-build', 'substitute', 'io-upload-iframe', 'collection', 'overlay', 'calendar', 'datatype-date']});\n"
	return prefix + content + suffix
}

func (self App) getFormJsLi() []string {
	formJsLi := []string{"js/rootform/r-form-field.js", "js/rootform/r-text-field.js", "js/rootform/r-hidden-field.js", "js/rootform/r-checkbox-field.js", "js/rootform/r-radio-field.js", "js/rootform/r-choice-field.js", "js/rootform/r-select-field.js", "js/rootform/r-trigger-field.js", "js/rootform/r-number-field.js", "js/rootform/r-display-field.js", "js/rootform/r-textarea-field.js", "js/rootform/r-date-field.js"}
	lFormJsLi := []string{"js/listform/lformcommon.js", "js/listform/l-form-field.js", "js/listform/l-text-field.js", "js/listform/l-hidden-field.js", "js/listform/l-checkbox-field.js", "js/listform/l-radio-field.js", "js/listform/l-choice-field.js", "js/listform/l-select-field.js", "js/listform/l-trigger-field.js", "js/listform/l-number-field.js", "js/listform/l-display-field.js", "js/listform/l-textarea-field.js", "js/listform/l-date-field.js"}
	pFormJsLi := []string{"js/form/p-form-field.js", "js/form/p-text-field.js", "js/form/p-hidden-field.js", "js/form/p-checkbox-field.js", "js/form/p-radio-field.js", "js/form/p-choice-field.js", "js/form/p-select-field.js", "js/form/p-trigger-field.js", "js/form/p-number-field.js", "js/form/p-display-field.js", "js/form/p-textarea-field.js", "js/form/p-date-field.js"}
	for _, k := range pFormJsLi {
		formJsLi = append(formJsLi, k)
	}
	for _, k := range lFormJsLi {
		formJsLi = append(formJsLi, k)
	}
	return formJsLi
}

func (self App) FormJS(w http.ResponseWriter, r *http.Request) {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if strings.Index(acceptEncoding, "gzip") > -1 {
		text := ""
		if config.String("debug") == "true" {
			content := self.getFormJsContent()
			data := bytes.Buffer{}
			gzipW := gzip.NewWriter(&data)
			gzipW.Write([]byte(content))
			gzipW.Close()
			text = data.String()
		} else {
			formJsNameLi := self.getFormJsLi()
			nameConcat := strings.Join(formJsNameLi, "")
			if self.isFileExist(nameConcat) {
				text = string(self.getGzipContent(nameConcat))
			} else {
				content := self.getFormJsContent()
				text = string(self.gzipAndSave(nameConcat, content))
			}
		}

//		self.Response.Status = http.StatusOK
		w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
		w.Header()["Content-Encoding"] = []string{"gzip"}

		//		end := time.Now()
		//		println("^^^^^^^^^^^^^^^^ formjs url time spend is:" + fmt.Sprint((end.UnixNano() - start.UnixNano())))
		w.Write([]byte(text))
	}

	content := self.getFormJsContent()
//	self.Response.Status = http.StatusOK
	w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
	w.Write([]byte(content))
}

func (self App) ComboView(w http.ResponseWriter, r *http.Request) {
	//	url := self.Request.URL.Path + "?" + self.Request.URL.RawQuery
	//	start := time.Now()

	nameConcat := self.getFileNameConcatFromQuery(r)

	acceptEncoding := r.Header.Get("Accept-Encoding")
	if strings.Index(acceptEncoding, "gzip") > -1 {
		text := ""
		if config.String("debug") == "true" {
			content := self.getComboViewFileContent(r)
			data := bytes.Buffer{}
			w := gzip.NewWriter(&data)
			w.Write([]byte(content))
			w.Close()
			text = data.String()
		} else {
			if self.isFileExist(nameConcat) {
				text = string(self.getGzipContent(nameConcat))
			} else {
				content := self.getComboViewFileContent(r)
				text = string(self.gzipAndSave(nameConcat, content))
			}
		}

//		self.Response.Status = http.StatusOK
		if strings.Index(r.URL.RawQuery, ".css") <= -1 {
			w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
		} else {
			w.Header()["Content-Type"] = []string{"text/css;charset=UTF-8"}
		}
		//		end := time.Now()
		//		println("^^^^^^^^^^^^^^^^ comboview url is:" + url + " time spend is:" + fmt.Sprint((end.UnixNano() - start.UnixNano())))
		w.Header()["Content-Encoding"] = []string{"gzip"}
		w.Write([]byte(text))
	}

	content := self.getComboViewFileContent(r)
//	self.Response.Status = http.StatusOK
	if strings.Index(r.URL.RawQuery, ".css") <= -1 {
		w.Header()["Content-Type"] = []string{"text/javascript;charset=UTF-8"}
	} else {
		w.Header()["Content-Type"] = []string{"text/css;charset=UTF-8"}
	}
	w.Write([]byte(content))
}

func (c App) getComboViewFileContent(r *http.Request) string {
	jsPath := config.String("COMBO_VIEW_PATH")
	content := ""
	commonUtil := CommonUtil{}
	for k := range r.URL.Query() {
		if !commonUtil.IsNumber(k) && k != "" {
			if strings.Index(k, ".js") == -1 && strings.Index(k, ".css") == -1 {
				panic("fileName is:" + k + ", expect ends with .js or .css")
			}
			isFileExist := false
			for _, filePath := range strings.Split(jsPath, ":") {
				if _, err := os.Stat(path.Join(filePath, k)); err != nil {
					if os.IsNotExist(err) {
						continue
					}
				}
				isFileExist = true
				file, err := os.Open(path.Join(filePath, k))
				defer file.Close()
				if err != nil {
					panic(err)
				}

				data, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}
				content += string(data) + "\n"
				break
			}
			if !isFileExist {
				panic(k + " is not exists")
			}
		}
	}
	return content
}
