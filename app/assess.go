package app

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type AssessListItem struct {
	NewsId string		`json:"newsId"`
	InsertDt string  	`json:"insertDt"`
	Category string		`json:"category"`
	Cluster string		`json:"cluster"`
	Title string        `json:"title"`
	ScoreAvg string		`json:"avg"`
	ScoreSum string		`json:"sco"`
	VanillaSum string	`json:"vsc"`
	JournalSum string	`json:"jsc"`
}

func (src *AssessStat) MakeListItem(applyWeight bool) AssessListItem {
	var jsco float64
	if applyWeight {
		jsco = src.WeightedJournalSum
	} else {
		jsco = src.JournalSum
	}
	return AssessListItem{
		src.NewsId,
		src.InsertDt.Format("2006-01-02 15:04"),
		src.Category,
		src.Cluster,
		src.Title,
		fmt.Sprintf("%.3f",	src.Average),
		fmt.Sprintf("%.3f",	src.ScoreSum),
		fmt.Sprintf("%.3f",	src.VanillaSum),
		fmt.Sprintf("%.3f",	jsco),
	}
}

type AssessStat struct {
	NewsId string `bson:"news_id"`
	Category string `bson:"category"`
	Cluster string `bson:"detailType"`

	InsertDt time.Time `bson:"insertDt"` 
	Title string `bson:"title"`
	ScoreSum float64 `bson:"score_totalSum"`
	JournalSum float64 `bson:"journal_totalSum"`
	
	Readability float64 `bson:"readability"`
	Clariry float64 `bson:"clariry"`
	Reality float64 `bson:"reality"`
	Usefulness float64 `bson:"usefulness"`
	Balance float64 `bson:"balance"`
	Variety float64 `bson:"variety"`
	Uniqueness float64 `bson:"uniqueness"`
	Importance float64 `bson:"importance"`
	Deep float64 `bson:"deep"`
	Inflammation float64 `bson:"inflammation"`
	Average float64 `bson:"average"`
	Sum float64 `bson:"sum"`
	Count float64 `bson:"count"`

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

	WeightedJournalSum float64 `bson:"weight"`
}

func jsonAssessList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}

	/*parms, _ := ctx.FormParams()
	for k,v := range parms {
		if !strings.HasPrefix(k, "columns") {
			fmt.Printf("%s => %s\n", k, v)
		}
	}*/

	// journal weight values
	applyWeight := ctx.FormValue("weight[apply]") == "1"
	weightRead,_ := strconv.Atoi(ctx.FormValue("weight[read]"))
	weightClear,_ := strconv.Atoi(ctx.FormValue("weight[clear]"))
	weightTruth,_ := strconv.Atoi(ctx.FormValue("weight[truth]"))
	weightUseful,_ := strconv.Atoi(ctx.FormValue("weight[useful]"))
	weightBalance,_ := strconv.Atoi(ctx.FormValue("weight[balance]"))
	weightVariety,_ := strconv.Atoi(ctx.FormValue("weight[variety]"))
	weightOriginal,_ := strconv.Atoi(ctx.FormValue("weight[original]"))
	weightImportant,_ := strconv.Atoi(ctx.FormValue("weight[important]"))
	weightDeep,_ := strconv.Atoi(ctx.FormValue("weight[deep]"))
	weightYellow,_ := strconv.Atoi(ctx.FormValue("weight[yellow]"))
	
    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterCategory := ctx.FormValue("columns[1][search][value]")
	filterCluster := ctx.FormValue("columns[2][search][value]")
	filterTitle := ctx.FormValue("columns[3][search][value]")
	filterNewsId := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)

	if filterNewsId != "" {
		findParam["newsId"] = bson.RegEx{filterNewsId, ""}
	} else {
		if filterCategory != "" {
			findParam["category"] = bson.RegEx{filterCategory, ""}
		}
		if filterCluster != "" {
			findParam["detailType"] = bson.RegEx{filterCluster, ""}
		}
		if filterTitle != "" {
			findParam["title"] = bson.RegEx{filterTitle, ""}
		}
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ASSTATS)

	json.Total, err = coll.Count()
	if err != nil {
		panic(err)
	}

	qry := coll.Find(&findParam)

	json.Filtered, err = qry.Count()
	if err != nil {
		panic(err)
	}

	selFields := bson.M{
		"insertDt":1,
		"category":1,
		"detailType":1,
		"title":1,
		"average":1,
		"score_totalSum":1,
		"vanilla_totalSum":1,
		"journal_totalSum":1,
		"news_id":1,
	}

	var results []AssessStat

	// 정렬
	if (applyWeight) {
		// 가중치 적용
		selFields["weight"] = bson.M{"$add": []bson.M {
			bson.M{ "$multiply": []interface{} { "$journal.readability", weightRead }},
			bson.M{ "$multiply": []interface{} { "$journal.transparency", weightClear }},
			bson.M{ "$multiply": []interface{} { "$journal.factuality", weightTruth }},
			bson.M{ "$multiply": []interface{} { "$journal.utility", weightUseful }},
			bson.M{ "$multiply": []interface{} { "$journal.fairness", weightBalance }},
			bson.M{ "$multiply": []interface{} { "$journal.diversity", weightVariety }},
			bson.M{ "$multiply": []interface{} { "$journal.originality", weightOriginal }},
			bson.M{ "$multiply": []interface{} { "$journal.importance", weightImportant }},
			bson.M{ "$multiply": []interface{} { "$journal.depth", weightDeep }},
			bson.M{ "$multiply": []interface{} { "$journal.sensationalism", (weightYellow*-1) }},
		}}

		pipeline := []bson.M{
			{ "$match": findParam },
			{ "$project": selFields },
			{ "$sort": bson.M{ "weight": -1 }},
			{ "$skip": skipRows },
			{ "$limit": fetchRows },
		}

		err = coll.Pipe(pipeline).All(&results)
		if err != nil {
			panic(err)
		}
	} else {
		// 컬럼 정렬	
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
				orderFieldName += "category"
			case 2:
				orderFieldName += "detailType"
			case 3:
				orderFieldName += "title"
			case 4:
				orderFieldName += "average"
			case 5:
				orderFieldName += "score_totalSum"
			case 6:
				orderFieldName += "vanilla_totalSum"
			case 7:
				orderFieldName += "journal_totalSum"
			default:
				orderFieldName += "title"
		}

		err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(selFields).All(&results)
		if err != nil {
			panic(err)
		}
	}

    json.Rows = make([]interface{}, len(results))
	for i,v := range results {
		json.Rows[i] = v.MakeListItem(applyWeight)
	}
	
	return ctx.JSON(http.StatusOK, json)
}

func showAssessList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3f"
	return ctx.renderTemplate("assess_list.html")
}

func showAssessDetail(_ctx echo.Context) error {
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
	collAss := mgoDB.With(sessionCopy).C(COLLNAME_ASSTATS)

	var data MdbNews
	err := coll.Find(where).One(&data)
	if err != nil {
		return ctx.showAdminMsg("Object not found")
	}

	if !strings.HasPrefix(data.Url, "http") {
		data.Url = "http://" + data.Url
	}

	ctx.model["item"] = &data

	var assess AssessStat
	err = collAss.Find(bson.M{"_id.news_id": data.NewsId}).One(&assess)
	if err != nil {
		fmt.Printf("AssessDetail, %v: %s\n", err, data.NewsId)
	}
	ctx.model["assess"] = &assess

	ctx.model["lmenu"] = "t3f"
	return ctx.renderTemplate("assess_view.html")
}

func setupAssessRoutes(grp *echo.Group) {
	grp.GET("/assess/list", showAssessList)
	grp.POST("/assess/list.json", jsonAssessList)
	grp.GET("/assess/view", showAssessDetail)
}
