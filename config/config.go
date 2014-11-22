package config

import (
	"bufio"
	. "github.com/hongjinqiu/gometa/common"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	GOMETA_CONFIG_PATH = "GOMETA_CONFIG_PATH"
)

var rwmutex sync.RWMutex = sync.RWMutex{}

var configPath string = ""
var isLoad bool = false
var config map[string]string = map[string]string{}

func getConfigPath() string {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return configPath
}

func getIsLoad() bool {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return isLoad
}

func getConfigValue(key string) string {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return config[key]
}

func SetConfigPath(filePath string) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	configPath = filePath
}

func getConfigFileWithouLock() *os.File {
	if configPath == "" && os.Getenv(GOMETA_CONFIG_PATH) == "" {
		lookpath, err := exec.LookPath(os.Args[0])
		if err != nil {
			panic(err)
		}
		absPath, err := filepath.Abs(lookpath)
		if err != nil {
			panic(err)
		}
		lastIndex := strings.LastIndex(absPath, string(os.PathSeparator))
		execPath := absPath[0:lastIndex]
		configFile, err := os.Open(filepath.Join(execPath, "app.conf"))
		if err != nil {
			panic(err)
		}

		return configFile
	} else if configPath == "" && os.Getenv(GOMETA_CONFIG_PATH) != "" {
		configPath = os.Getenv(GOMETA_CONFIG_PATH)
		configFile, err := os.Open(configPath)
		if err != nil {
			panic(err)
		}

		return configFile
	} else {
		configFile, err := os.Open(configPath)
		if err != nil {
			panic(err)
		}

		return configFile
	}
}

func loadConfig() {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	if !isLoad || (isLoad && config["debug"] == "true") {
		configFile := getConfigFileWithouLock()
		defer configFile.Close()

		commonUtil := CommonUtil{}
		reader := bufio.NewReader(configFile)
		isDev := false
		isProd := false
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			if commonUtil.IsEmpty(string(line)) {
				continue
			}
			// 去掉备注部分,添加多行支持,最后赋值到data中
			keyValueLi := strings.Split(string(line), "=")
			if len(keyValueLi[0]) == 0 || keyValueLi[0][0:1] == "#" {
				continue
			}

			key := keyValueLi[0]
			if key == "[dev]" {
				isDev = true
				continue
			}
			if key == "[prod]" {
				isProd = true
				continue
			}
			value := ""
			if len(keyValueLi) > 1 {
				value = keyValueLi[1]
			}
			key = commonUtil.TrimString(key)
			value = commonUtil.TrimString(value)
			if len(value) > 0 && value[len(value)-1:] == `\` {
				value = value[0 : len(value)-1]
				for {
					nextLine, _, err := reader.ReadLine()
					if err == io.EOF {
						break
					}
					nextLineStr := commonUtil.TrimString(string(nextLine))
					if len(nextLineStr) > 0 && nextLineStr[len(nextLineStr)-1:] == `\` {
						value += nextLineStr[0 : len(nextLineStr)-1]
					} else {
						value += nextLineStr
						break
					}
				}
			}
			// value的环境变量解析
			value = getValueFromEnv(value)
			if !isDev && !isProd {
				config[key] = value
			} else if config["debug"] == "true" && isDev && !isProd { // 应用[dev]部分配置
				config[key] = value
			} else if config["debug"] == "false" && isDev && isProd { // 应用[prod]部分配置
				config[key] = value
			}
		}
		isLoad = true
	}
}

func getValueFromEnv(value string) string {
	regx := regexp.MustCompile(`\$[\da-zA-Z_.]*`)
	result := regx.FindAllString(value, -1)
	
	for _, item := range result {
		envname := item[1:]
		envvalue := os.Getenv(envname)
		sepLi := strings.Split(envvalue, string(os.PathListSeparator))
		value = strings.Replace(value, item, sepLi[0], -1)
	}
	
	regx = regexp.MustCompile(`%[\da-zA-Z_.]*%`)
	result = regx.FindAllString(value, -1)
	
	for _, item := range result {
		envname := item[1:len(item) - 1]
		envvalue := os.Getenv(envname)
		sepLi := strings.Split(envvalue, string(os.PathListSeparator))
		value = strings.Replace(value, item, sepLi[0], -1)
	}
	
	return value
}

func String(key string) string {
	if !getIsLoad() || (getIsLoad() && getConfigValue("debug") == "true") {
		loadConfig()
	}
	return getConfigValue(key)
}
