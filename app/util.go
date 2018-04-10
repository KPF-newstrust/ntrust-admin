package app

import (
	"fmt"
	"log"
	"time"
	"net/http"

	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LEVEL_TRACE = 1
	LEVEL_DEBUG = 3
	LEVEL_INFO = 5
	LEVEL_WARNING = 6
	LEVEL_ERROR = 7
	LEVEL_FATAL = 9
)

const COLLNAME_EVENTLOG = "evtlog"

// https://gist.github.com/bsphere/8369aca6dde3e7b4392c
type BJsonTime time.Time

func NewJsonNow() BJsonTime {
	return BJsonTime(time.Now())
}

func (thiz *BJsonTime) SetBSON(raw bson.Raw) error {
	var tm time.Time

	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}

	*thiz = BJsonTime(tm)
	return nil
}

func (thiz BJsonTime) GetBSON() (interface{}, error) {
	tm := time.Time(thiz)
	if tm.IsZero() {
		return nil, nil
	}

	return tm, nil
}

func (thiz *BJsonTime) MarshalJSON() ([]byte, error) {
	tm := time.Time(*thiz)
	if tm.IsZero() {
		return []byte("null"), nil
	}

	str := fmt.Sprintf("\"%s\"", tm.Format("2006-01-02 15:04:05"))
	return []byte(str), nil
}

type appEvent struct {
	Id			bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Timestamp	BJsonTime `bson:"ts" json:"ts"`
	Level	   	int `bson:"lv" json:"lv"`
	IpAddr		string `bson:"ip" json:"ip"`
	Tag			string `bson:"tag" json:"tag"`
	Msg			string `bson:"msg" json:"msg"`
}

func saveAppEvent(evt *appEvent, ctx* AuthContext) {
	evt.Timestamp = NewJsonNow()
	if ctx == nil {
		evt.IpAddr = "internal"
	} else {
		evt.IpAddr = ctx.RealIP()
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()

	coll := mgoDB.With(sessionCopy).C(COLLNAME_EVENTLOG)
	err := coll.Insert(evt)
	if err != nil {
		log.Printf("saveAppEvent failed: %s\n", err.Error())
	}
	//return (err == nil)
}

func LogTrace(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_TRACE
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

func LogDebug(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_DEBUG
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

func LogInfo(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_INFO
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

func LogWarning(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_WARNING
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

func LogError(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_ERROR
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

func LogFatal(tag string, ctx *AuthContext, format string, args ...interface{}) {
	var evt appEvent
	evt.Level = LEVEL_FATAL
	evt.Tag = tag
	evt.Msg = fmt.Sprintf(format, args...)
	go saveAppEvent(&evt, ctx)
}

////////////////////////////////////////////////////////////////////////////////

type JqDataTable struct {
    Draw int `json:"draw"`
    Total int `json:"recordsTotal"`
    Filtered int `json:"recordsFiltered"`
    Rows []interface{} `json:"data"`
}

func (ret *JqDataTable) append(ptr interface{}) {
    ret.Rows = append(ret.Rows, ptr)
}

/*
type JqGridData struct {
    Page int `json:"page"`
    Total int `json:"total"`
    Records int `json:"records"`
    Rows []interface{} `json:"rows"`
}

func (ret *JqGridData) setPageParam(curPage, numRows int) {
    ret.Page = curPage
    ret.Total = (ret.Records + numRows -1) / numRows
}

func (ret *JqGridData) append(ptr interface{}) {
    ret.Rows = append(ret.Rows, ptr)
}
*/

////////////////////////////////////////////////////////////////////////////////

const (
	API_AUTH_FAILED = 1
	API_DATABASE_ERROR = 2

	API_INVALID_PARAM = 10
	API_INVALID_CONDITION = 11
	API_KEY_DUPLICATED = 12

	API_APILAB_FAILED = 20
	
	API_INTERNAL_SERVER_ERROR = 99
)

type apiResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (ctx *AuthContext) apiOk(msg string, data interface{}) error {
	return ctx.JSON(http.StatusOK, apiResult{Code: 0, Msg: msg, Data: data})
}

func apiCodeToMsg(errCode int) string {
	switch(errCode) {
	case API_AUTH_FAILED:
		return "사용자 인증에 실패하였습니다."
	case API_DATABASE_ERROR:
		return "데이터베이스 오류입니다."

	case API_INVALID_PARAM:
		return "입력 파라메터가 올바르지 않습니다."
	case API_INVALID_CONDITION:
		return "조건이 올바르지 않습니다."
	case API_KEY_DUPLICATED:
		return "중복된 데이터 입니다."

	case API_APILAB_FAILED:
		return "API 테스트 에러"

	case API_INTERNAL_SERVER_ERROR:
		return "서버 내부 에러입니다."
		
	default:
		return "알 수 없는 에러코드"
	}
}

func (ctx *AuthContext) apiError0(errCode int) error {		
	return ctx.JSON(http.StatusOK, apiResult{Code:errCode, Msg:apiCodeToMsg(errCode), Data:nil})
}

func (ctx *AuthContext) apiErrorMsg(errCode int, msg string) error {
	return ctx.JSON(http.StatusOK, apiResult{Code:errCode, Msg:msg, Data:nil})
}

func (ctx *AuthContext) apiErrorData(errCode int, data interface{}) error {
	return ctx.JSON(http.StatusOK, apiResult{Code:errCode, Msg:apiCodeToMsg(errCode), Data:data})
}
