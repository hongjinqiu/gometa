package session

import (
	"crypto/md5"
	"fmt"
	"github.com/hongjinqiu/gometa/common"
	"github.com/hongjinqiu/gometa/config"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
	"strings"
)

const (
	EXPIRES = "expires"
)

var sessionDict map[string]map[string]string = map[string]map[string]string{}
var rwmutex sync.RWMutex = sync.RWMutex{}

func init() {
	time.AfterFunc(time.Second*time.Duration(10), scanSession)
}

func scanSession() {
	expireSessionLi := []string{}
	now := time.Now()
	rwmutex.RLock()
	for key, value := range sessionDict {
		expires := value[EXPIRES]
		date, err := time.Parse("20060102150405", expires)
		if err != nil {
			panic(err)
		}
		if now.After(date) {
			expireSessionLi = append(expireSessionLi, key)
		}
	}
	rwmutex.RUnlock()

	if len(expireSessionLi) > 0 {
		rwmutex.Lock()
		defer rwmutex.Unlock()

		for _, item := range expireSessionLi {
			delete(sessionDict, item)
		}
	}
	time.AfterFunc(time.Second*time.Duration(10), scanSession)
}

func getSessionId() string {
	dateUtil := common.DateUtil{}
	dateStr := dateUtil.GetDateByFormat("20060102150405")
	float64Str := fmt.Sprint(rand.Float64())
	h := md5.New()
	io.WriteString(h, dateStr+float64Str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*
Cookie {
	Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string // Raw t
}
*/

func GetFromSession(w http.ResponseWriter, r *http.Request, name string) string {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	sessionName := config.String("SESSION_NAME")
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		if err != http.ErrNoCookie {
			panic(err)
		} else {
			// get cookie from response,
			cookieStr := w.Header().Get("Set-Cookie")
			cookieDict := getCookieDictFromStr(cookieStr)
			sessionId := cookieDict[sessionName]
			return sessionDict[sessionId][name]
//			return ""
		}
	} else {
		sessionId := cookie.Value
		return sessionDict[sessionId][name]
	}
}

/**
cookieStr的格式:
gometasessionid=67f663ae6d07dc255fdc650897140353; Path=/; Expires=Sun, 28 Dec 2014 08:15:14 UTC; Max-Age=14400
*/
func getCookieDictFromStr(cookieStr string) map[string]string {
	result := map[string]string{}
	li := strings.Split(cookieStr, ";")
	for _, item := range li {
		item = strings.TrimSpace(item)
		itemLi := strings.Split(item, "=")
		if len(itemLi) == 1 {
			result[itemLi[0]] = ""
		} else if len(itemLi) > 1 {
			result[itemLi[0]] = itemLi[1]
		}
	}
	return result
}

func SetToSession(w http.ResponseWriter, r *http.Request, name string, value string) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	sessionName := config.String("SESSION_NAME")
	sessionCookieAge := config.String("SESSION_COOKIE_AGE")
	sessionCookieAgeInt, err := strconv.Atoi(sessionCookieAge)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	expires := now.Add(time.Second * time.Duration(sessionCookieAgeInt))
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		if err != http.ErrNoCookie {
			panic(err)
		} else {
			cookie = &http.Cookie{
				Name:    sessionName,
				Value:   getSessionId(),
				Path:    "/",
				Expires: expires,
				MaxAge:  sessionCookieAgeInt,
			}
			http.SetCookie(w, cookie)
		}
	}
	sessionId := cookie.Value

	sessionValueDict := map[string]string{}
	if tmpSessionValueDict := sessionDict[sessionId]; tmpSessionValueDict != nil {
		sessionValueDict = tmpSessionValueDict
	}
	sessionValueDict[name] = value
	sessionValueDict[EXPIRES] = expires.Format("20060102150405")
	sessionDict[sessionId] = sessionValueDict

	cookie.Expires = expires
	cookie.MaxAge = sessionCookieAgeInt

	http.SetCookie(w, cookie)
}
