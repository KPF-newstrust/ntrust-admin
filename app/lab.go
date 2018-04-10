package app

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
	"strconv"
	"net/http"
	"encoding/json"
	
	"github.com/labstack/echo"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LabCommon struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"DT_RowId"`
	Type string `bson:"type" json:"type"`
	Title string `bson:"title" json:"title"`
	Content string `bson:"content" json:"-"`
	RequestedAt time.Time `bson:"requestedAt" json:"-"`
	UserId string `bson:"userId"`

	MediaType string `bson:"mediaType,omitempty" json:"-"`
	Category string `bson:"category,omitempty" json:"-"`
	//Byline string `bson:"byline,omitempty" json:"-"`

	CompletedAt *time.Time `bson:"completedAt,omitempty" json:"-"`
	Result string `bson:"result,omitempty" json:"-"`
}

func (u *LabCommon) MarshalJSON() ([]byte, error) {
	var state string
	if u.CompletedAt == nil {
		state = "대기중"
	} else {
		state = u.CompletedAt.Format("02일 15:04:05에 완료됨")
	}

	type Alias LabCommon
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"ts"`
		State string `json:"state"`
	}{
		Alias:    (*Alias)(u),
		Timestamp: u.RequestedAt.Format("2006-01-02 15:04:05"),
		State: state,
	})
}

type LabIntegrate struct {
	LabCommon `bson:",inline"`

	// postag
	MecabTime float32 `bson:"mecab_time"`
	MecabTags []wordPosItem `bson:"mecab_postag"`
	MecabPersons []string `bson:"mecab_PS"`
	MecabOrganizations []string `bson:"mecab_OG"`
	MecabLocations []string `bson:"mecab_LC"`
	MecabPlans []string `bson:"mecab_PL"`
	MecabProducts []string `bson:"mecab_PR"`
	MecabEvents []string `bson:"mecab_EV"`

	HannanumTime float32 `bson:"hannanum_time"`
	HannanumTags []wordPosItem `bson:"hannanum_postag"`
	HannanumPersons []string `bson:"hannanum_PS"`
	HannanumOrganizations []string `bson:"hannanum_OG"`
	HannanumLocations []string `bson:"hannanum_LC"`
	HannanumPlans []string `bson:"hannanum_PL"`
	HannanumProducts []string `bson:"hannanum_PR"`
	HannanumEvents []string `bson:"hannanum_EV"`
	
	KkmaTime float32 `bson:"kkma_time"`
	KkmaTags []wordPosItem `bson:"kkma_postag"`
	KkmaPersons []string `bson:"kkma_PS"`
	KkmaOrganizations []string `bson:"kkma_OG"`
	KkmaLocations []string `bson:"kkma_LC"`
	KkmaPlans []string `bson:"kkma_PL"`
	KkmaProducts []string `bson:"kkma_PR"`
	KkmaEvents []string `bson:"kkma_EV"`

	TwitterTime float32 `bson:"twitter_time"`
	TwitterTags []wordPosItem `bson:"twitter_postag"`
	TwitterPersons []string `bson:"twitter_PS"`
	TwitterOrganizations []string `bson:"twitter_OG"`
	TwitterLocations []string `bson:"twitter_LC"`
	TwitterPlans []string `bson:"twitter_PL"`
	TwitterProducts []string `bson:"twitter_PR"`
	TwitterEvents []string `bson:"twitter_EV"`

	// metric
	LineCount int
	ContentLength int		`bson:"content_length"`
	TitleLength int			`bson:"title_length"`
	TitleNumExclamation int `bson:"title_numExclamation"`
	TitleNumQuestion int	`bson:"title_numQuestion"`
	TitleNumSingleQuote int `bson:"title_numSingleQuote"`
	TitleNumDoubleQuote int `bson:"title_numDoubleQuote"`
	TitleHasShock int		`bson:"title_hasShock"`
	TitleHasExclusive int	`bson:"title_hasExclusive"`
	TitleHasBreaking int	`bson:"title_hasBreaking"`
	TitleHasPlan int		`bson:"title_hasPlan"`
	TitleAdverbs []string	`bson:"title_adverbs"`
	
	ContentAvgSentenceLength float64	`bson:"content_avgSentenceLength"`
	ContentAvgAdverbsPerSentence float64	`bson:"content_avgAdverbsPerSentence"`
	ContentQuotePercent float64	`bson:"content_quotePercent"`
	ContentAnonPredicates []string	`bson:"content_anonPredicates"`
	ContentForeignWords []string	`bson:"content_foreignWords"`	// TODO
	ContentNumberCount int	`bson:"content_numNumber"`

	QuotedSentences []quotedSentenceItem `bson:"quotes"`
	Bylines []writerByline `bson:"bylines"`

	Score scoreInfo			`bson:"score" json:"score"`

	Journal struct {
		Readability float64 `bson:"readability" json:"readability"`
		Transparency float64 `bson:"transparency" json:"transparency"`
		Factuality float64 `bson:"factuality" json:"factuality"`
		Utility float64 `bson:"utility" json:"utility"`
		Fairness float64 `bson:"fairness" json:"fairness"`
		Diversity float64 `bson:"diversity" json:"diversity"`
		Originality float64 `bson:"originality" json:"originality"`
		Importance float64 `bson:"importance" json:"importance"`
		Depth float64 `bson:"depth" json:"depth"`
		Sensationalism float64 `bson:"sensationalism" json:"sensationalism"`
	} `bson:"journal"`

	Vanilla struct {
		Readability float64 `bson:"readability" json:"readability"`
		Transparency float64 `bson:"transparency" json:"transparency"`
		Factuality float64 `bson:"factuality" json:"factuality"`
		Utility float64 `bson:"utility" json:"utility"`
		Fairness float64 `bson:"fairness" json:"fairness"`
		Diversity float64 `bson:"diversity" json:"diversity"`
		Originality float64 `bson:"originality" json:"originality"`
		Importance float64 `bson:"importance" json:"importance"`
		Depth float64 `bson:"depth" json:"depth"`
		Sensationalism float64 `bson:"sensationalism" json:"sensationalism"`
	} `bson:"vanilla"`

	JournalSum float64 `bson:"journal_totalSum"`
	VanillaSum float64 `bson:"vanilla_totalSum"`
}

func jsonLabList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	skipRows, _ := strconv.Atoi(ctx.FormValue("start"))
	fetchRows, _ := strconv.Atoi(ctx.FormValue("length"))
	if fetchRows <= 0 {
		fetchRows = 10
	}

	orderDir := ctx.FormValue("order[0][dir]")
	orderColumn, _ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterType := ctx.FormValue("columns[1][search][value]")
	filterContent := ctx.FormValue("columns[2][search][value]")
	filterCompleted := ctx.FormValue("columns[3][search][value]")

	var jsonData JqDataTable
	jsonData.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	var findParam = make(bson.M)

	if filterType != "" {
		findParam["type"] = filterType
	}
	if filterContent != "" {
		findParam["content"] = bson.RegEx{filterContent, ""}
	}
	if filterCompleted == "Y" {
		findParam["completedAt"] = bson.M{"$exists":true}
	} else if filterCompleted != "" {
		findParam["completedAt"] = bson.M{"$exists":false}
	}
	if ctx.model[MVK_AuthLevel].(int) < 5 { //어드민 이하는 자신이 요청한 이력만 조회
		findParam["userId"] = ctx.authUserId
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

	var err error
	jsonData.Total, err = coll.Count()
	if err != nil {
		panic(err)
	}

	qry := coll.Find(&findParam)

	jsonData.Filtered, err = qry.Count()
	if err != nil {
		panic(err)
	}

	var orderFieldName string
	if orderDir == "desc" {
		orderFieldName = "-"
	} else {
		orderFieldName = ""
	}

	switch (orderColumn) {
	case 1:
		orderFieldName += "type"
	case 2:
		orderFieldName += "title"
	case 3:
		orderFieldName += "completedAt"
	default:
		orderFieldName += "requestedAt"
	}

	var results []LabCommon
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(bson.M{
		"requestedAt":1,
		"type":1,
		"title":1,
		"completedAt":1,
	}).All(&results)
	if err != nil {
		panic(err)
	}

	jsonData.Rows = make([]interface{}, len(results))
	for i := range results {
		jsonData.Rows[i] = &results[i]
	}

	return ctx.JSON(http.StatusOK, jsonData)
}

func showLabIntro(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t8z"
	return ctx.renderTemplate("lab/intro.html")
}

func showLabList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t8a"
	return ctx.renderTemplate("lab/list.html")
}

func postLabNew(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var data LabCommon
	data.Type = ctx.FormValue("type")
	
	data.Title = strings.TrimSpace(ctx.FormValue("title"))
	data.Content = strings.TrimSpace(ctx.FormValue("content"))
	data.Content = strings.Replace(data.Content, "⏎","\n",-1)
	if data.Content == "" {
		return ctx.apiErrorMsg(API_INVALID_PARAM, "본문이 입력되지 않았습니다.")
	}

	if data.Title == "" {
		contentRunes := []rune(data.Content)
		lenTitle := len(contentRunes)
		if lenTitle > 25 {
			lenTitle = 25
		}
		data.Title = string(contentRunes[:lenTitle])
	}

	if data.Type == "trust" || data.Type == "integrate" {
		data.MediaType = ctx.FormValue("mtype")
		data.Category = ctx.FormValue("category")
	}

	data.RequestedAt = time.Now()
	data.Id = bson.NewObjectId()
	data.UserId = ctx.authUserId

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

	err := coll.Insert(&data)
	if err != nil {
		return ctx.apiErrorMsg(API_DATABASE_ERROR, err.Error())
	}

	payload := bson.M{"id": data.Id.Hex()}
	err = mqSend_Task("lab_"+data.Type, 1, payload)
	if err != nil {
		panic(err)
	}

	return ctx.apiOk("OK", payload)
}

func postLabQuery(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	cmd := ctx.FormValue("cmd")
	objId := ctx.FormValue("id")
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

	var data LabCommon
	err := coll.FindId(bson.ObjectIdHex(objId)).Select(bson.M{"type":1, "completedAt":1}).One(&data)
	if err == nil {
		if cmd == "wait" {
			// API 작업이 완료되었는지 검사
			return ctx.apiOk("OK", bson.M{"done": (data.CompletedAt != nil) })
		} else if cmd == "resend" {
			// RabbitMQ 재전송 
			payload := bson.M{"id": objId}
			err = mqSend_Task("lab_"+data.Type, 1, payload)
			if err != nil {
				return ctx.apiErrorMsg(API_APILAB_FAILED, fmt.Sprintf("API 요청 에러: %v", err))
			}

			return ctx.apiOk("OK", payload)
		} else if cmd == "delete" {
			err = coll.RemoveId(data.Id)
			if err != nil {
				return ctx.apiErrorMsg(API_DATABASE_ERROR, err.Error())
			}

			return ctx.apiOk("OK", nil)
		}

		return ctx.apiErrorMsg(API_APILAB_FAILED, "Invalid command: " + cmd)
	} else if err == mgo.ErrNotFound {
		return ctx.apiErrorMsg(API_APILAB_FAILED, "테스트 ID가 올바르지 않습니다.")
	} else {
		return ctx.apiErrorMsg(API_APILAB_FAILED, fmt.Sprintf("DB에러: %v", err))
	}
}

func (ctx *AuthContext) _showLabData(lmenu, templ string, dataptr interface{}, postproc func()) error {
	labId := ctx.QueryParam("id")
	if labId != "" {
		sessionCopy := mgoSession.Copy()
		defer sessionCopy.Close()
		coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

		err := coll.FindId(bson.ObjectIdHex(labId)).One(dataptr)
		if err != nil {
			return ctx.showAdminError(err)
		}

		if postproc != nil {
			postproc()
		}		
	}
	
	ctx.model["lmenu"] = lmenu
	ctx.model["data"] = dataptr
	return ctx.renderTemplate("lab/"+templ+".html")
}

func showLabSplit(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabCommon
	return ctx._showLabData("t8b", "split", &data, nil)
}

func showLabPostag(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabIntegrate
	return ctx._showLabData("t8c", "postag", &data, nil)
}

func showLabSanitize(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabCommon
	return ctx._showLabData("t8d", "sanitize", &data, nil)
}

func showLabMetric(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabIntegrate
	return ctx._showLabData("t8e", "metric", &data, func() {
		data.LineCount = len(strings.Split(data.Result, "\n"))
	})
}

func showLabTrust(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabIntegrate
	return ctx._showLabData("t8f", "trust", &data, nil)
}

func showLabIntegrate(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	var data LabIntegrate
	return ctx._showLabData("t8g", "integrate", &data, nil)
}

func setupLabRoutes(grp *echo.Group) {
	grp.GET("/lab/intro", showLabIntro)

	grp.GET("/lab/list", showLabList)
	grp.POST("/lab/list.json", jsonLabList)

	grp.POST("/lab/new.json", postLabNew)
	grp.POST("/lab/query.json", postLabQuery)

	grp.GET("/lab/split", showLabSplit)
	grp.GET("/lab/postag", showLabPostag)
	grp.GET("/lab/sanitize", showLabSanitize)
	grp.GET("/lab/metric", showLabMetric)
	grp.GET("/lab/trust", showLabTrust)
	grp.GET("/lab/integrate", showLabIntegrate)
}

var _validCategories = map[string]bool{
	"정치":true,
	"경제":true,
	"사회":true,
	"국제":true,
	"문화 예술":true,
	"IT 과학":true,
	"교육":true,
	"스포츠":true,
	"연예":true,
	"라이프스타일":true,
	"사설·칼럼":true,
}

func apiPostLabNew(_ctx echo.Context, dtype string) error {
	ctx := _ctx.(*AuthContext)

	var data LabCommon
	data.Type = dtype
	data.Title = strings.TrimSpace(ctx.FormValue("title"))
	data.Content = strings.TrimSpace(ctx.FormValue("content"))
	
	data.Content = strings.Replace(data.Content, "⏎","\n",-1)
	
	if data.Content == "" {
		return ctx.apiErrorMsg(API_INVALID_PARAM, "본문이 입력되지 않았습니다.")
	}

	if data.Title == "" {
		contentRunes := []rune(data.Content)
		lenTitle := len(contentRunes)
		if lenTitle > 25 {
			lenTitle = 25
		}
		data.Title = string(contentRunes[:lenTitle])
	}
	
	if data.Type == "trust" || data.Type == "integrate" {
		data.MediaType = ctx.FormValue("mediaType")
		if data.MediaType != "신문" && data.MediaType != "방송" {
			return ctx.apiErrorMsg(API_INVALID_PARAM, "mediaType 값이 올바르지 않습니다: " + data.MediaType)
		}

		data.Category = ctx.FormValue("category")
		if _, ok := _validCategories[data.Category]; !ok {
			return ctx.apiErrorMsg(API_INVALID_PARAM, "category 값이 올바르지 않습니다: " + data.Category)
		}
	}

	data.RequestedAt = time.Now()
	data.Id = bson.NewObjectId()

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

	err := coll.Insert(&data)
	if err != nil {
		return ctx.apiErrorMsg(API_DATABASE_ERROR, err.Error())
	}

	payload := bson.M{"id": data.Id.Hex()}
	err = mqSend_Task("lab_"+data.Type, 1, payload)
	if err != nil {
		panic(err)
	}

	return ctx.apiOk("OK", map[string]string{"id":data.Id.Hex()})
}

func apiGetLabResult(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.Param("objId")
	if !bson.IsObjectIdHex(objId) {
		return ctx.apiErrorMsg(API_INVALID_PARAM, "Invalid Object ID")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_APILAB)

	var data LabIntegrate
	err := coll.FindId(bson.ObjectIdHex(objId)).One(&data)
	if err != nil {
		return ctx.apiErrorMsg(API_DATABASE_ERROR, err.Error())
	}

	if data.CompletedAt == nil {
		return ctx.apiErrorMsg(API_APILAB_FAILED, "API Request not completed yet.")
	}

	result := bson.M{"requestedAt":data.RequestedAt}	
	result["completedAt"] = *data.CompletedAt

	result["id"] = data.Id.Hex()
	result["cmd"] = data.Type
	
	if data.Type == "split" || data.Type == "sanitize" {
		result["result"] = data.Result
	} else if data.Type == "postag" {
		result["result"] = data.MecabTags
	} else if data.Type == "metric" {
		result["title"] = data.Title
		result["content"] = data.Content
		result["result"] = bson.M{
			"bylines": data.Bylines,
			"quotes": data.QuotedSentences,
			"title_length": data.TitleLength,
			"title_numQuestion": data.TitleNumQuestion,
			"title_numExclamation": data.TitleNumExclamation,
			"title_numSingleQuote": data.TitleNumSingleQuote,
			"title_numDoubleQuote": data.TitleNumDoubleQuote,
			"title_hasExclusive": data.TitleHasExclusive,
			"title_hasBreaking": data.TitleHasBreaking,
			"title_hasPlan": data.TitleHasPlan,
			"content_length": data.ContentLength,
			"content_avgAdverbsPerSentence": data.ContentAvgAdverbsPerSentence,
			"content_avgSentenceLength": data.ContentAvgSentenceLength,
			"content_numberCount": data.ContentNumberCount,
			"content_quotePercent": data.ContentQuotePercent,
		}
	} else if data.Type == "trust" {
		result["title"] = data.Title
		result["content"] = data.Content
		result["result"] = bson.M{
			"scores": data.Score,
			"journal": data.Journal,
		}
	} else if data.Type == "integrate" {
		result["title"] = data.Title
		result["content"] = data.Content
		result["result"] = bson.M{
			"bylines": data.Bylines,
			"quotes": data.QuotedSentences,
			"title_length": data.TitleLength,
			"title_numQuestion": data.TitleNumQuestion,
			"title_numExclamation": data.TitleNumExclamation,
			"title_numSingleQuote": data.TitleNumSingleQuote,
			"title_numDoubleQuote": data.TitleNumDoubleQuote,
			"title_hasExclusive": data.TitleHasExclusive,
			"title_hasBreaking": data.TitleHasBreaking,
			"title_hasPlan": data.TitleHasPlan,
			"content_length": data.ContentLength,
			"content_avgAdverbsPerSentence": data.ContentAvgAdverbsPerSentence,
			"content_avgSentenceLength": data.ContentAvgSentenceLength,
			"content_numberCount": data.ContentNumberCount,
			"content_quotePercent": data.ContentQuotePercent,

			"scores": data.Score,
			"journal": data.Journal,
		}
	} else {
		return ctx.apiErrorMsg(API_APILAB_FAILED, "Unsupported api type: " + data.Type)
	}
	
	return ctx.apiOk("OK", result)
}

func setupApiRoutes(ec *echo.Echo) {
	ec.POST("/api/split", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "split")
	})
	ec.POST("/api/postag", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "postag")
	})
	ec.POST("/api/sanitize", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "sanitize")
	})
	ec.POST("/api/metric", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "metric")
	})
	ec.POST("/api/trust", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "trust")
	})
	ec.POST("/api/integrate", func(ctx echo.Context) error {
		return apiPostLabNew(ctx, "integrate")
	})
	ec.GET("/api/result/:objId", func(ctx echo.Context) error {
		return apiGetLabResult(ctx)
	})
}
