package app

import (
	"fmt"
	"strings"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type BugItem struct {
	Id bson.ObjectId	`bson:"_id,omitempty" json:"id"`
	Reporter bson.ObjectId	`bson:"reporter" json:"reporter"`
	CreatedAt time.Time	`bson:"createdAt" json:"createdAt"`
	RepliedAt *time.Time	`bson:"repliedAt" json:"repliedAt"`
	NewsId string   	`bson:"newsId" json:"newsId"`
	Category string   	`bson:"category" json:"category"`
	Message string      `bson:"message" json:"message"`
	Reply string      	`bson:"reply" json:"reply"`
	Confirmed bool		`bson:"confirmed" json:"confirmed"`
}

type bugListItem struct {
	Id string			`json:"DT_RowId"`
	CreatedAt string	`json:"createdAt"`
	NewsId string   	`json:"newsId"`
	Media string		`json:"media"`
	Title string      	`json:"title"`
	Category string   	`json:"category"`
	Message string      `json:"message"`
	Reply string      	`json:"reply"`
	Confirmed bool		`json:"confirmed"`	
}





type wordPosItem struct {
	Word string `bson:"word" json:"word"`
	Pos string `bson:"pos" json:"pos"`
}

type quotedSentenceItem struct {
	Sentence string `bson:"sentence" json:"sentence"`
	Length int `bson:"length" json:"length"`
}

type writerByline struct {
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Post string `bson:"post" json:"-"`
}

type scoreInfo struct {
	TotalSum float64 `bson:"totalSum" json:"totalSum"`
	Average float64 `bson:"average" json:"average"`
	Byline float64 `bson:"byline" json:"byline"`
	ContentLength float64 `bson:"contentLength" json:"contentLength"`
	QuoteCount float64 `bson:"quoteCount" json:"quoteCount"`
	TitleLength float64 `bson:"titleLength" json:"titleLength"`
	TitlePuncCount float64 `bson:"titlePuncCount" json:"titlePuncCount"`
	NumberCount float64 `bson:"numberCount" json:"numberCount"`	
	ImageCount float64 `bson:"imageCount" json:"imageCount"`
	AvgSentenceLength float64 `bson:"avgSentenceLength" json:"avgSentenceLength"`
	TitleAdverbCount float64 `bson:"titleAdverbCount" json:"titleAdverbCount"`
	AvgAdverbCountPerSentence float64 `bson:"avgAdverbCountPerSentence" json:"avgAdverbCountPerSentence"`
	QuotePercent float64 `bson:"quotePercent" json:"quotePercent"`
	AnonPredicateCount float64 `bson:"anonPredicateCount" json:"anonPredicateCount"`
	ForeignWordCount float64 `bson:"foreignWordCount" json:"foreignWordCount"`
	InformantRealCount float64 `bson:"informantRealCount" json:"informantRealCount"`
	QuoteRatioRealAnon float64 `bson:"quoteRatioRealAnon" json:"quoteRatioRealAnon"`
}

type MdbNews struct {
	Id bson.ObjectId	`bson:"_id,omitempty" json:"DT_RowId"`
	InsertDt time.Time  `bson:"insertDt" json:"insertDt"`
	NewsId string       `bson:"newsId" json:"newsId"`
	MediaId string   	`bson:"mediaId" json:"mediaId"`
	MediaName string	`bson:"mediaName" json:"mediaName"`
	Title string        `bson:"title" json:"title"`
	Content string  	`bson:"content" json:"content"`
	PubDate string     	`bson:"pubDate" json:"pubDate"`
	Byline string       `bson:"byline" json:"byline"`	// DELETE soon
	BylineWriter writerByline `bson:"byline_writer" json:"bylineWriter"`
	Bylines []writerByline    `bson:"bylines" json:"bylines"`
	Url string          `bson:"url" json:"url"`
	
	CategoryOrig string `bson:"categoryOrig" json:"categoryOrig"`
	CategoryXls string  `bson:"categoryXls" json:"categoryXls"`
	CategoryMan string  `bson:"categoryMan" json:"categoryMan"`
	CategoryCalc string	`bson:"categoryCalc" json:"categoryCalc"`
	CategoryFinal string	`bson:"categoryFinal" json:"categoryFinal"`

	ImageCount int			`bson:"image_count"`
	ContentLength int		`bson:"content_length" json:"content_length"`
	ContentNumberCount int	`bson:"content_numNumber"`
	ContentAvgSentenceLength float64	`bson:"content_avgSentenceLength"`
	ContentAvgAdverbsPerSentence float64	`bson:"content_avgAdverbsPerSentence"`
	ContentQuotePercent float64	`bson:"content_quotePercent"`
	ContentAnonPredicates []string	`bson:"content_anonPredicates"`
	ContentForeignWords []string	`bson:"content_foreignWords"`	// TODO
	//InformantReal int	`bson:"informantReal"`
	//QuoteRatioRealAnon int	`bson:"quoteRatioRealAnon"`
	
	TitleLength int			`bson:"title_length" json:"title_length"`
	TitleAdverbs []string	`bson:"title_adverbs"`
	
	TitleNumExclamation int `bson:"title_numExclamation" json:"title_numExclamation"`
	TitleNumQuestion int	`bson:"title_numQuestion" json:"title_numQuestion"`
	TitleNumSingleQuote int `bson:"title_numSingleQuote" json:"title_numSingleQuote"`
	TitleNumDoubleQuote int `bson:"title_numDoubleQuote" json:"title_numDoubleQuote"`
	TitleHasShock int		`bson:"title_hasShock" json:"title_hasShock"`
	TitleHasExclusive int	`bson:"title_hasExclusive" json:"title_hasExclusive"`
	TitleHasBreaking int	`bson:"title_hasBreaking" json:"title_hasBreaking"`
	TitleHasPlan int		`bson:"title_hasPlan" json:"title_hasPlan"`
	QuotedSentences []quotedSentenceItem `bson:"quotes" json:"quotes"`	
	Score scoreInfo			`bson:"score" json:"score"`

	NerPersons []string		`bson:"ner_PS" json:"ner_PS"`
	NerOrganizations []string	`bson:"ner_OG" json:"ner_OG"`
	NerLocations []string	`bson:"ner_LC" json:"ner_LC"`

	InformantReal []string		`bson:"informant_real" json:"informant_real"`
	InformantAnno []string		`bson:"informant_anno" json:"informant_anno"`

	ScoreSum float64 `bson:"score_totalSum"`
	JournalSum float64 `bson:"journal_totalSum"`
	Journal struct {
		Readability float64 `bson:"readability"`
		Transparency float64 `bson:"transparency"`
		Factuality float64 `bson:"factuality"`
		Utility float64 `bson:"utility"`
		Fairness float64 `bson:"fairness"`
		Diversity float64 `bson:"diversity"`
		Originality float64 `bson:"originality"`
		Importance float64 `bson:"importance"`
		Depth float64 `bson:"depth"`
		Sensationalism float64 `bson:"sensationalism"`
	} `bson:"journal"`

	VanillaSum float64 `bson:"vanilla_totalSum"`	
	Vanilla struct {
		Readability float64 `bson:"readability"`
		Transparency float64 `bson:"transparency"`
		Factuality float64 `bson:"factuality"`
		Utility float64 `bson:"utility"`
		Fairness float64 `bson:"fairness"`
		Diversity float64 `bson:"diversity"`
		Originality float64 `bson:"originality"`
		Importance float64 `bson:"importance"`
		Depth float64 `bson:"depth"`
		Sensationalism float64 `bson:"sensationalism"`
	} `bson:"vanilla"`
	
	ClusterDelegate bool `bson:"clusterDelegate"`
	ClusterNewsId string `bson:"clusterNewsId"`
	ClusterId string `bson:"clusterId"`
	ClusterSimilarity float64 `bson:"clusterSimilarity"`

	WeightedJournalSum float64 `bson:"weight"`
}

type MdbNewsEntity struct {
	MecabTime float32 `bson:"mecab_time" json:"mecab_time"`
	MecabTags []wordPosItem `bson:"mecab_postag" json:"mecab_postag"`
	MecabPersons []string `bson:"mecab_PS" json:"mecab_PS"`
	MecabOrganizations []string `bson:"mecab_OG" json:"mecab_OG"`
	MecabLocations []string `bson:"mecab_LC" json:"mecab_LC"`
	MecabPlans []string `bson:"mecab_PL" json:"mecab_PL"`
	MecabProducts []string `bson:"mecab_PR" json:"mecab_PR"`
	MecabEvents []string `bson:"mecab_EV" json:"mecab_EV"` 

	HannanumTime float32 `bson:"hannanum_time" json:"hannanum_time"`
	HannanumTags []wordPosItem `bson:"hannanum_postag" json:"hannanum_postag"`
	HannanumPersons []string `bson:"hannanum_PS" json:"hannanum_PS"`
	HannanumOrganizations []string `bson:"hannanum_OG" json:"hannanum_OG"`
	HannanumLocations []string `bson:"hannanum_LC" json:"hannanum_LC"`
	HannanumPlans []string `bson:"hannanum_PL" json:"hannanum_PL"`
	HannanumProducts []string `bson:"hannanum_PR" json:"hannanum_PR"`
	HannanumEvents []string `bson:"hannanum_EV" json:"hannanum_EV"` 
	
	KkmaTime float32 `bson:"kkma_time" json:"kkma_time"`
	KkmaTags []wordPosItem `bson:"kkma_postag" json:"kkma_postag"`
	KkmaPersons []string `bson:"kkma_PS" json:"kkma_PS"`
	KkmaOrganizations []string `bson:"kkma_OG" json:"kkma_OG"`
	KkmaLocations []string `bson:"kkma_LC" json:"kkma_LC"`
	KkmaPlans []string `bson:"kkma_PL" json:"kkma_PL"`
	KkmaProducts []string `bson:"kkma_PR" json:"kkma_PR"`
	KkmaEvents []string `bson:"kkma_EV" json:"kkma_EV"`

	TwitterTime float32 `bson:"twitter_time" json:"twitter_time"`
	TwitterTags []wordPosItem `bson:"twitter_postag" json:"twitter_postag"`
	TwitterPersons []string `bson:"twitter_PS" json:"twitter_PS"`
	TwitterOrganizations []string `bson:"twitter_OG" json:"twitter_OG"`
	TwitterLocations []string `bson:"twitter_LC" json:"twitter_LC"`
	TwitterPlans []string `bson:"twitter_PL" json:"twitter_PL"`
	TwitterProducts []string `bson:"twitter_PR" json:"twitter_PR"`
	TwitterEvents []string `bson:"twitter_EV" json:"twitter_EV"` 
}

type NewsListItem struct {
	Id bson.ObjectId	`json:"DT_RowId"`
	NewsId string		`json:"newsId"`
	InsertDt string  	`json:"insertDt"`
	MediaId string   	`json:"mediaId"`
	MediaName string	`json:"mediaName"`
	Title string        `json:"title"`
	Byline string       `json:"byline"`
	CategoryXls string	`json:"categoryXls"`
	CategoryMan string	`json:"categoryMan"`
	CategoryCalc string	`json:"categoryCalc"`
}

func (src *MdbNews) MakeListItem() NewsListItem {
	var byline string
	if len(src.Bylines) > 0 {
		byline = fmt.Sprintf("%s %s", src.Bylines[0].Name, src.Bylines[0].Email)
	}
	
	return NewsListItem{
		src.Id,
		src.NewsId,
		src.InsertDt.Format("2006-01-02 15:04:05"),
		src.MediaId,
		src.MediaName,
		src.Title,
		byline,
		src.CategoryXls,
		src.CategoryMan,
		src.CategoryCalc,
	}
}

func jsonNewsList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}
	
	/*parms, _ := ctx.FormParams()
	for k,v := range parms {
		if strings.HasPrefix(k, "columns") {
			fmt.Printf("%s => %s\n", k, v)
		}
	}*/

    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterMediaId := ctx.FormValue("columns[1][search][value]")
	filterMediaName := ctx.FormValue("columns[2][search][value]")
	filterCategory1 := ctx.FormValue("columns[3][search][value]")
	filterCategory2 := ctx.FormValue("columns[4][search][value]")
	filterCategory3 := ctx.FormValue("columns[5][search][value]")
	filterByline := ctx.FormValue("columns[6][search][value]")
	filterTitle := ctx.FormValue("columns[7][search][value]")
	filterNewsId := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)

	if filterNewsId != "" {
		findParam["newsId"] = bson.RegEx{filterNewsId, ""}
	} else {
		ymd := ctx.FormValue("ymd")
		dtBegin, err := time.ParseInLocation("2006-01-02", ymd, tzLocation)
		if err == nil {
			dtEnd := dtBegin.AddDate(0,0,1)
			findParam["insertDt"] = bson.M{ "$gte":dtBegin, "$lt":dtEnd }
		}

		if filterMediaId != "" {
			findParam["mediaId"] = bson.RegEx{filterMediaId, ""}
		}
		if filterMediaName != "" {
			findParam["mediaName"] = bson.RegEx{filterMediaName, ""}
		}
		if filterCategory1 != "" {
			findParam["categoryXls"] = bson.RegEx{filterCategory1, ""}
		}
		if filterCategory2 != "" {
			findParam["categoryMan"] = bson.RegEx{filterCategory2, ""}
		}
		if filterCategory3 != "" {
			findParam["categoryCalc"] = bson.RegEx{filterCategory3, ""}
		}
		if filterTitle != "" {
			findParam["title"] = bson.RegEx{filterTitle, ""}
		}

		if strings.HasPrefix(filterByline, "~!") {
			notIncl := filterByline[2:]
			if notIncl == "" {
				findParam["bylines"] = nil	
			} else {
				findParam["bylines"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{
					bson.M{"name": bson.M{"$exists":true}},
					bson.M{"name": bson.M{"$not":bson.RegEx{notIncl, ""}}},
				}}}
			}			
		} else if filterByline != "" {
			findParam["bylines"] = bson.M{"$elemMatch": bson.M{"$or": []bson.M{
				bson.M{"name":bson.RegEx{filterByline, ""}},
				bson.M{"email":bson.RegEx{filterByline, ""}},
			}}}
		}
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	json.Total, err = coll.Count()
	if err != nil {
		panic(err)
	}

	qry := coll.Find(&findParam)

	json.Filtered, err = qry.Count()
	if err != nil {
		panic(err)
	}

	// 정렬
	var orderFieldName string
	if orderDir == "desc" {
		orderFieldName = "-"
	} else {
		orderFieldName = ""
	}

	switch (orderColumn) {
		case 0:
			orderFieldName += "insertDt"
		case 1:
			orderFieldName += "mediaId"
		case 2:
			orderFieldName += "mediaName"
		case 3:
			orderFieldName += "categoryXls"
		case 4:
			orderFieldName += "categoryMan"
		case 5:
			orderFieldName += "categoryCalc"
		case 6:
			orderFieldName += "byline"
		case 7:
		default:
			orderFieldName += "title"
	}

	selFields := bson.M{
		"insertDt":1,
		"newsId":1,
		"mediaId":1,
		"mediaName":1,
		"categoryXls":1,
		"categoryMan":1,
		"categoryCalc":1,
		"bylines":1,
		"title":1,
	}
	var results []MdbNews
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(selFields).All(&results)
	if err != nil {
		panic(err)
	}

    json.Rows = make([]interface{}, len(results))
	for i,v := range results {
		json.Rows[i] = v.MakeListItem()
	}
	
	return ctx.JSON(http.StatusOK, json)
}

func showNewsList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3c"
	ctx.model["targetDate"] = "2016-06-01"
	return ctx.renderTemplate("news_list.html")
}

func showNewsDetail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var where bson.D
	objId := ctx.QueryParam("obj")
	newsId := ctx.QueryParam("nws")
	if objId != "" {
		where = bson.D{{"_id", bson.ObjectIdHex(objId)}}
	} else if newsId != "" {
		where = bson.D{{"newsId", newsId}}
	} else {
		return ctx.showAdminMsg("Invalid ID")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	var data MdbNews
	err := coll.Find(where).One(&data)
	if err != nil {
		return ctx.showAdminMsg("Object not found")
	}
	
	// 강원도민일보는 http:// 가 안붙어있더라..
	if !strings.HasPrefix(data.Url, "http") {
		data.Url = "http://" + data.Url
	}

	ctx.model["item"] = &data

	coll = mgoDB.With(sessionCopy).C(COLLNAME_NEWS_SRC)
	var nsrcData bson.M
	err = coll.Find(bson.M{"newsitem_id":data.NewsId}).Select(bson.M{"news_content":1}).One(&nsrcData)
	if err != nil {
		if err != mgo.ErrNotFound {
			return ctx.showAdminMsg("수집DB 항목 없음: " + data.NewsId)
		}
	} else {
		ctx.model["src_content"] = nsrcData["news_content"].(string)
	}

	coll = mgoDB.With(sessionCopy).C(COLLNAME_BUGREPORT)
	var bugs []BugItem
	err = coll.Find(bson.M{"newsId":data.NewsId}).All(&bugs)
	ctx.model["bugs"] = bugs

	hasNewReply := false
	for _,bug := range bugs {
		if bug.RepliedAt != nil && bug.Confirmed == false {
			hasNewReply = true
			break
		}
	}
	ctx.model["newReply"] = hasNewReply

	ctx.model["lmenu"] = "t3c"
	return ctx.renderTemplate("news_view.html")
}

func showNewsDetailEntity(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	newsId := ctx.QueryParam("nws")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	
	coll_news := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)
	var dataNews MdbNews
	err := coll_news.Find(bson.M{"newsId":newsId}).Select(bson.M{"ner_PS":1, "ner_OG":1, "ner_LC":1}).One(&dataNews)
	if err != nil {
		return ctx.showAdminMsg("Object not found")
	} else {
		ctx.model["item"] = &dataNews
	}

	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS_ENTITY)
	var dataEntity MdbNewsEntity
	err = coll.Find(bson.M{"newsId":newsId}).One(&dataEntity)
	if err == nil {
		ctx.model["entity"] = &dataEntity
	}

	// 수동 정답셋
	coll = mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)
	var dataAnnot MdbAnnotate
	err = coll.Find(bson.M{"newsId":newsId}).One(&dataAnnot)
	ctx.model["annot"] = &dataAnnot
	
	// 1차년도 DB
	coll = mgoOldDB.With(sessionCopy).C(COLLNAME_YEAR1NEWS)
	var dataYear1 MdbYear1News
	coll.Find(bson.M{"news_id":newsId}).Select(bson.M{
			"TMS_NE_LOCATION":1,
			"TMS_NE_ORGANIZATION":1,
			"TMS_NE_PERSON":1,
			"TMS_NE_STREAM":1,
			"TMS_RAW_STREAM":1,
			"TMS_SIMILARITY":1,
		}).One(&dataYear1)
	ctx.model["year1"] = &dataYear1

	return ctx.renderTemplate("news_view_entity.html")
}

func showNewsDetailCluster(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	newsId := ctx.QueryParam("nws")
	ctx.model["thisNewsId"] = newsId

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	
	coll_news := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)
	var dataNews MdbNews
	err := coll_news.Find(bson.M{"newsId":newsId}).
		Select(bson.M{"clusterDelegate":1, "clusterNewsId":1, "clusterId":1, "clusterSimilarity":1}).One(&dataNews)
	if err != nil {
		return ctx.showAdminMsg("Object not found")
	}
	ctx.model["item"] = dataNews

	var clusterNewsId string
	if dataNews.ClusterDelegate {
		clusterNewsId = newsId
	} else if dataNews.ClusterNewsId != "" {
		clusterNewsId = dataNews.ClusterNewsId

		var dataDele MdbNews
		err = coll_news.Find(bson.M{"newsId":clusterNewsId}).Select(bson.M{"title":1, "newsId":1 }).One(&dataDele)
		if err != nil {
			return ctx.showAdminMsg("Dele Object not found")
		}
		ctx.model["dele"] = dataDele
	} else {
		clusterNewsId = ""
	}
	
	if clusterNewsId != "" {
		var dataClusters []MdbNews
		err = coll_news.Find(bson.M{"clusterNewsId":clusterNewsId}).Sort("-clusterSimilarity").
			Select(bson.M{"title":1, "newsId":1, "clusterSimilarity":1}).All(&dataClusters)
		if err != nil {
			return ctx.showAdminMsg("Cluster object not found")
		}

		ctx.model["clusters"] = &dataClusters
	}

	return ctx.renderTemplate("news_view_cluster.html")
}

func updateNewsManualCategory(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	objId := ctx.FormValue("id")
	cate := ctx.FormValue("cate")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	var data MdbNews
	coll.FindId(bson.ObjectIdHex(objId)).One(&data)
	if data.CategoryMan == cate {
		ctx.flashError("기존 카테고리 값과 같아서 수정하지 않았습니다.")
	} else {
		change := bson.M{ "$set": bson.M{"categoryMan": cate} }
		err := coll.UpdateId(bson.ObjectIdHex(objId), change)
		if err != nil {
			return ctx.showAdminError(err)
		}

		ctx.flashSuccess("수동 카테고리를 변경하였습니다.")
	}
	
	return ctx.Redirect(http.StatusFound, "view?obj="+objId)
}

func copyToAnnotateDB(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	var src MdbNews
	err := coll.FindId(bson.ObjectIdHex(objId)).One(&src)
	if err != nil {
		return ctx.showAdminError(err)
	}

	coll = mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)
	cnt, err := coll.Find(bson.M{"newsId":src.NewsId}).Count()
	if cnt > 0 {
		ctx.flashSuccess("이미 복사된 기사입니다.")
		return ctx.Redirect(http.StatusFound, "/admin/annotate/view?nws="+src.NewsId)
	}

	var dst MdbAnnotate
	dst.Id = bson.NewObjectId()
	dst.NewsId = src.NewsId
	dst.InsertDt = time.Now()
	dst.UpdateDt = time.Now()
	dst.Category = src.CategoryFinal
	dst.Title = src.Title
	dst.Content = src.Content
	dst.MediaId = src.MediaId
	dst.MediaName = src.MediaName

	err = coll.Insert(dst)
	if err != nil {
		return ctx.showAdminError(err)
	}

	return ctx.Redirect(http.StatusFound, "/admin/annotate/view?id="+dst.Id.Hex())
}

func reportBug(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	newsId := ctx.FormValue("newsId")
	bugId := ctx.FormValue("id")
	cmd := ctx.FormValue("cmd")
	
	cate := ctx.FormValue("cate")
	msg := ctx.FormValue("msg")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_BUGREPORT)

	var err error
	var successMsg string
	switch cmd {
	case "new":
		var item BugItem
		item.CreatedAt = time.Now()
		item.Reporter = bson.ObjectIdHex(ctx.authUserId)
		item.NewsId = newsId
		item.Message = msg
		item.Category = cate

		err = coll.Insert(&item)
		successMsg = "새로운 문제점을 보고하였습니다."

	case "edit":
		err = coll.UpdateId(bson.ObjectIdHex(bugId), bson.M{"$set":bson.M{
			"category":cate,
			"message":msg,
		}})
		successMsg = "문제점 내용을 수정하였습니다."

	case "del":
		err = coll.RemoveId(bson.ObjectIdHex(bugId))
		successMsg = "문제점 신고 내역을 삭제하였습니다."

	case "confirm":
		err = coll.UpdateId(bson.ObjectIdHex(bugId), bson.M{"$set":bson.M{"confirmed":true}})
		successMsg = "문제점 신고 내역을 확인하였습니다."
		
	default:
		return ctx.showAdminMsg("Invalid command: " + cmd)
	}

	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.flashSuccess(successMsg)
	return ctx.Redirect(http.StatusFound, "/admin/news/view?nws="+newsId)
}

func replyBug(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	newsId := ctx.FormValue("newsId")
	bugId := ctx.FormValue("id")
	reply := ctx.FormValue("reply")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_BUGREPORT)
	err := coll.UpdateId(bson.ObjectIdHex(bugId), bson.M{"$set":bson.M{
		"repliedAt": time.Now(),
		"reply": reply,
	}})

	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.flashSuccess("응답 메세지를 수정하였습니다.")
	return ctx.Redirect(http.StatusFound, "/admin/news/view?nws="+newsId)
}

func showBugsList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3e"
	return ctx.renderTemplate("bugs_list.html")
}

func jsonBugsList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}
	
    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterCategory := ctx.FormValue("columns[2][search][value]")
	filterMessage := ctx.FormValue("columns[3][search][value]")
	filterReply := ctx.FormValue("columns[4][search][value]")
	filterConfirm := ctx.FormValue("columns[5][search][value]")

	var jsonData JqDataTable
    jsonData.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)


	if filterCategory != "" {
		findParam["category"] = bson.RegEx{filterCategory, ""}
	}
	if filterMessage != "" {
		findParam["message"] = bson.RegEx{filterMessage, ""}
	}
	if filterReply != "" {
		findParam["reply"] = bson.RegEx{filterReply, ""}
	}
	if filterConfirm != "" {
		findParam["confirmed"] = (filterConfirm == "Y")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_BUGREPORT)

	jsonData.Total, err = coll.Count()
	if err != nil {
		panic(err)
	}

	qry := coll.Find(&findParam)

	jsonData.Filtered, err = qry.Count()
	if err != nil {
		panic(err)
	}

	// 정렬
	var orderFieldName string
	if orderDir == "desc" {
		orderFieldName = "-"
	} else {
		orderFieldName = ""
	}

	switch (orderColumn) {
		case 0:
			orderFieldName += "createdAt"
		case 1:
			orderFieldName += "panic"
		case 2:
			orderFieldName += "category"
		case 3:
			orderFieldName += "message"
		case 4:
			orderFieldName += "reply"
		case 5:
			orderFieldName += "confirmed"
		default:
			orderFieldName += "shit"
	}

	var results []BugItem
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).All(&results)
	if err != nil {
		panic(err)
	}

	coll = mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

    jsonData.Rows = make([]interface{}, len(results))
	for i,v := range results {
		var item bugListItem
		item.Id = v.Id.Hex()
		item.CreatedAt = v.CreatedAt.Format("2006-01-02 15:04:05")
		item.NewsId = v.NewsId
		item.Category = v.Category
		item.Message = v.Message
		item.Reply = v.Reply
		item.Confirmed = v.Confirmed

		var newsInfo MdbNews
		err = coll.Find(bson.M{"newsId":item.NewsId}).Select(bson.M{"title":1, "mediaName":1, "_id":0}).One(&newsInfo)
		if err == nil {
			item.Title = newsInfo.Title
			item.Media = newsInfo.MediaName
		} else {
			item.Title = fmt.Sprintf("에러(%v)", err)
		}
		
		jsonData.Rows[i] = &item
	}
	
	return ctx.JSON(http.StatusOK, jsonData)
}

func setupNewsRoutes(grp *echo.Group) {
	setupNsrcRoutes(grp)
	setupAssessRoutes(grp)
	setupClusterRoutes(grp)

	// 처리 DB -- 군더더기 제거, 문장분리 등
	grp.GET("/news/list", showNewsList)
	grp.POST("/news/list.json", jsonNewsList)
	grp.GET("/news/view", showNewsDetail)
	grp.GET("/news/entity", showNewsDetailEntity)
	grp.GET("/news/cluster", showNewsDetailCluster)
	grp.POST("/news/editCate", updateNewsManualCategory)
	grp.POST("/news/copyToAnnot", copyToAnnotateDB)

	grp.POST("/bugs/report", reportBug)
	grp.POST("/bugs/reply", replyBug)
	grp.GET("/bugs/list", showBugsList)
	grp.POST("/bugs/list.json", jsonBugsList)
	
}
