package log

import (
	"bufio"
	"fmt"
	"github.com/hongjinqiu/gometa/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	FLUSH_LOG_PERIOD = 1 * time.Second

	LevelTrace = iota
	LevelInfo
	LevelWarning
	LevelError
)

var stdout *log.Logger
var fileLog *log.Logger
var file *os.File
var bufWriter *bufio.Writer
var rwmutex sync.RWMutex = sync.RWMutex{}

var logChan chan string = make(chan string, 100)
var isLogToNewFileChan chan int = make(chan int)

var stopFlush chan int = make(chan int)
var startFlush chan int = make(chan int)

func init() {
	if isLogToStdout() {
		stdout = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}
	if isLogToFile() {
		logFilePath, logFileNameWithoutPath := getLogFileNameAndPath()
		writer, err := os.OpenFile(filepath.Join(logFilePath, logFileNameWithoutPath), os.O_RDWR|os.O_APPEND, os.ModeDevice|0666)
		if err != nil {
			panic(err)
		}

		setFile(writer)
		setBufWriter(bufio.NewWriter(writer))
		setFileLog(log.New(getBufWriter(), "", log.Ldate|log.Ltime))
		time.AfterFunc(FLUSH_LOG_PERIOD, flushLogPeriod)
		logToNewFilePeriod()

		logChan <- "avoid out of order"
		<-logChan
		go logDailyFile()
	}
}

func logDailyFile() {
	for {
		select {
		case info := <-logChan:
			getFileLog().Printf(info)
		case <-isLogToNewFileChan:
			logToNewFile()
		}
	}
}

func flushLogPeriod() {
	select {
	case <-stopFlush:
		<-startFlush
	default:
		err := getBufWriter().Flush()
		if err != nil {
			panic(err)
		}
	}
	time.AfterFunc(FLUSH_LOG_PERIOD, flushLogPeriod)
}

func logToNewFilePeriod() {
	now := time.Now()
	hour := 0
	minute := 0
	second := 0
	nextDate := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, time.Local)
	nextDate = nextDate.AddDate(0, 0, 1)
	time.AfterFunc(nextDate.Sub(now), func() {
		isLogToNewFileChan <- 1
		logToNewFilePeriod()
	})
}

func logToNewFile() {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	if isLogToFile() {
		logFilePath, logFileNameWithoutPath := getLogFileNameAndPath()
		logFilePathName := filepath.Join(logFilePath, logFileNameWithoutPath)
		currentTime := time.Now()
		yesterDay := currentTime.AddDate(0, 0, -1)
		yesterdayYmd := yesterDay.Format("20060102")
		currentYmd := currentTime.Format("20060102")
		if yesterdayYmd < currentYmd {
			stopFlush <- 1
			defer func() {
				startFlush <- 1
			}()

			err := bufWriter.Flush()
			if err != nil {
				panic(err)
			}

			err = file.Close()
			if err != nil {
				panic(err)
			}

			err = os.Rename(logFilePathName, logFilePathName+"."+yesterdayYmd)
			if err != nil {
				panic(err)
			}

			fi, err := os.Create(logFilePathName)
			if err != nil {
				panic(err)
			}
			file = fi
			bufWriter = bufio.NewWriter(file)
			fileLog = log.New(bufWriter, "", log.Ldate|log.Ltime)
		}
	}
}

func getLogFileNameAndPath() (string, string) {
	if isLogToFile() {
		logFileOutput := config.String("log.file.output")
		logFilePath := ""
		logFileNameWithoutPath := ""
		if strings.Contains(logFileOutput, string(os.PathSeparator)) {
			lastIndex := strings.LastIndex(logFileOutput, string(os.PathSeparator))
			logFilePath = logFileOutput[0:lastIndex]
			logFileNameWithoutPath = logFileOutput[lastIndex:]
		} else {
			lookpath, err := exec.LookPath(os.Args[0])
			if err != nil {
				panic(err)
			}
			absPath, err := filepath.Abs(lookpath)
			if err != nil {
				panic(err)
			}
			lastIndex := strings.LastIndex(absPath, string(os.PathSeparator))
			logFilePath = absPath[0:lastIndex]
			logFileNameWithoutPath = logFileOutput
		}
		if _, err := os.Stat(logFilePath); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(logFilePath, os.ModeDevice|0666)
			}
		}
		if _, err := os.Stat(filepath.Join(logFilePath, logFileNameWithoutPath)); err != nil {
			if os.IsNotExist(err) {
				tmpFile, tmpErr := os.Create(filepath.Join(logFilePath, logFileNameWithoutPath))
				if tmpErr != nil {
					panic(tmpErr)
				}
				tmpErr = tmpFile.Close()
				if tmpErr != nil {
					panic(tmpErr)
				}
			}
		}
		return logFilePath, logFileNameWithoutPath
	}
	return "", ""
}

func isLogToStdout() bool {
	logLogger := config.String("log.logger")
	logLoggerLower := strings.ToLower(logLogger)
	if strings.Contains(logLoggerLower, "stdout") {
		return true
	}
	return false
}

func isLogToFile() bool {
	logLogger := config.String("log.logger")
	logLoggerLower := strings.ToLower(logLogger)
	if strings.Contains(logLoggerLower, "file") {
		return true
	}
	return false
}

func Level() int {
	if strings.ToLower(config.String("log.level")) == "trace" {
		return LevelTrace
	}
	if strings.ToLower(config.String("log.level")) == "info" {
		return LevelInfo
	}
	if strings.ToLower(config.String("log.level")) == "warn" {
		return LevelWarning
	}
	if strings.ToLower(config.String("log.level")) == "error" {
		return LevelError
	}
	return LevelTrace
}

func Trace(v ...interface{}) {
	if Level() <= LevelTrace {
		if isLogToStdout() {
			stdout.Printf("%v %v\n", config.String("log.trace.prefix"), v)
		}
		if isLogToFile() {
			logChan <- fmt.Sprintf("%v %v\n", config.String("log.trace.prefix"), v)
		}
	}
}

func Info(v ...interface{}) {
	if Level() <= LevelInfo {
		if isLogToStdout() {
			stdout.Printf("%v %v\n", config.String("log.info.prefix"), v)
		}
		if isLogToFile() {
			logChan <- fmt.Sprintf("%v %v\n", config.String("log.info.prefix"), v)
		}
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	if Level() <= LevelWarning {
		if isLogToStdout() {
			stdout.Printf("%v %v\n", config.String("log.warn.prefix"), v)
		}
		if isLogToFile() {
			logChan <- fmt.Sprintf("%v %v\n", config.String("log.warn.prefix"), v)
		}
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	if Level() <= LevelError {
		if isLogToStdout() {
			stdout.Printf("%v %v\n", config.String("log.error.prefix"), v)
		}
		if isLogToFile() {
			logChan <- fmt.Sprintf("%v %v\n", config.String("log.error.prefix"), v)
		}
	}
}

func getFileLog() *log.Logger {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return fileLog
}

func setFileLog(log *log.Logger) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	fileLog = log
}

func getBufWriter() *bufio.Writer {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return bufWriter
}

func setBufWriter(writer *bufio.Writer) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	bufWriter = writer
}

func getFile() *os.File {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return file
}

func setFile(fi *os.File) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	file = fi
}
