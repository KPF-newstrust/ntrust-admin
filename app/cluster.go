package app

import (
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ClusterItem struct {
	Id         bson.ObjectId `bson:"_id"`
	NewsId     string        `bson:"newsId"`
	Title      string        `bson:"title"`
	Media      string        `bson:"mediaName"`
	Similarity float64       `bson:"clusterSimilarity"`
	Score      float64       `bson:"journal_totalSum"`

	Weight float64 `bson:"weight"`
}

type ClusterCard struct {
	Id     bson.ObjectId `bson:"_id"`
	NewsId string        `bson:"newsId"`
	Title  string        `bson:"title"`
	Media  string        `bson:"mediaName"`
	Score  float64       `bson:"journal_totalSum"`

	Items []ClusterItem
	Count int `bson:"clusterNewsCount`
}

func showClusterList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	pgNum, _ := strconv.Atoi(ctx.FormValue("pg"))
	const CARDS_PER_PAGE int = 10

	var findParam = bson.M{"clusterDelegate": true}

	title := strings.TrimSpace(ctx.FormValue("title"))
	if title != "" {
		findParam["title"] = bson.RegEx{title, ""}
	}

	cate := ctx.FormValue("cate")
	if cate == "전체" || cate == "" {
		cate = "전체"
	} else {
		findParam["categoryFinal"] = cate
	}

	date := ctx.FormValue("date")
	if date == "" {
		date = "2016-06-01"
	}
	dtBegin, err := time.ParseInLocation("2006-01-02", date, tzLocation)
	if err == nil {
		dtEnd := dtBegin.AddDate(0, 0, 1)
		findParam["insertDt"] = bson.M{"$gte": dtBegin, "$lt": dtEnd}
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	qry := coll.Find(&findParam)
	totalRows, err := qry.Count()
	totalPages := (totalRows + CARDS_PER_PAGE - 1) / CARDS_PER_PAGE
	if totalPages <= pgNum {
		pgNum = totalPages - 1
		if pgNum < 0 {
			pgNum = 0
		}
	}

	skipRows := pgNum * CARDS_PER_PAGE

	selFields := bson.M{
		"title":            1,
		"mediaName":        1,
		"newsId":           1,
		"journal_totalSum": 1,
	}

	var cards []ClusterCard
	err = qry.Sort("-clusterNewsCount").Skip(skipRows).Limit(CARDS_PER_PAGE).Select(selFields).All(&cards)
	if err != nil {
		panic(err)
	}

	ctx.model["cards"] = cards
	for i, v := range cards {
		qry = coll.Find(bson.M{"clusterNewsId": v.NewsId})
		cards[i].Count, err = qry.Count()
		qry.Sort("-journal_totalSum").Limit(5).Select(
			bson.M{"newsId": 1, "title": 1, "mediaName": 1, "journal_totalSum": 1}).All(&cards[i].Items)
	}

	ctx.model["page"] = pgNum
	ctx.model["total"] = totalPages

	ctx.model["date"] = date
	ctx.model["title"] = title
	ctx.model["selCate"] = cate
	ctx.model["cateOpts"] = []string{"전체", "정치", "경제", "사회", "국제", "문화 예술", "IT 과학", "교육", "스포츠", "연예", "라이프스타일", "사설·칼럼"}
	ctx.model["lmenu"] = "t3g"
	return ctx.renderTemplate("cluster_list.html")
}

func getDefaultWeight(ctx echo.Context, varname string) int {
	val, err := strconv.Atoi(ctx.FormValue(varname))
	if err != nil {
		val = 1
	}
	return val
}

func showClusterDetail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	clusterNewsId := ctx.Param("id")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS)

	var dele ClusterCard
	err := coll.Find(bson.M{"clusterDelegate": true, "newsId": clusterNewsId}).Select(bson.M{
		"title": 1, "newsId": 1, "mediaName": 1}).One(&dele)
	if err != nil {
		return ctx.showAdminMsg("대표 기사를 찾을 수 없습니다.")
	}

	ctx.model["dele"] = dele

	// journal weight values
	weightRead := getDefaultWeight(ctx, "read")
	weightClear := getDefaultWeight(ctx, "clear")
	weightTruth := getDefaultWeight(ctx, "truth")
	weightUseful := getDefaultWeight(ctx, "useful")
	weightBalance := getDefaultWeight(ctx, "balance")
	weightVariety := getDefaultWeight(ctx, "variety")
	weightOriginal := getDefaultWeight(ctx, "original")
	weightImportant := getDefaultWeight(ctx, "important")
	weightDeep := getDefaultWeight(ctx, "deep")
	weightYellow := getDefaultWeight(ctx, "yellow")

	ctx.model["read"] = weightRead
	ctx.model["clear"] = weightClear
	ctx.model["truth"] = weightTruth
	ctx.model["useful"] = weightUseful
	ctx.model["balance"] = weightBalance
	ctx.model["variety"] = weightVariety
	ctx.model["original"] = weightOriginal
	ctx.model["important"] = weightImportant
	ctx.model["deep"] = weightDeep
	ctx.model["yellow"] = weightYellow

	var findParam = bson.M{"clusterNewsId": clusterNewsId}

	selFields := bson.M{
		"title":             1,
		"mediaName":         1,
		"newsId":            1,
		"journal_totalSum":  1,
		"clusterSimilarity": 1,
	}

	selFields["weight"] = bson.M{"$add": []bson.M{
		bson.M{"$multiply": []interface{}{"$journal.readability", weightRead}},
		bson.M{"$multiply": []interface{}{"$journal.transparency", weightClear}},
		bson.M{"$multiply": []interface{}{"$journal.factuality", weightTruth}},
		bson.M{"$multiply": []interface{}{"$journal.utility", weightUseful}},
		bson.M{"$multiply": []interface{}{"$journal.fairness", weightBalance}},
		bson.M{"$multiply": []interface{}{"$journal.diversity", weightVariety}},
		bson.M{"$multiply": []interface{}{"$journal.originality", weightOriginal}},
		bson.M{"$multiply": []interface{}{"$journal.importance", weightImportant}},
		bson.M{"$multiply": []interface{}{"$journal.depth", weightDeep}},
		bson.M{"$multiply": []interface{}{"$journal.sensationalism", (weightYellow * -1)}},
	}}

	pipeline := []bson.M{
		{"$match": findParam},
		{"$project": selFields},
		{"$sort": bson.M{"weight": -1}},
	}

	var items []ClusterItem
	err = coll.Pipe(pipeline).All(&items)
	if err != nil {
		panic(err)
	}

	ctx.model["items"] = items
	ctx.model["lmenu"] = "t3g"
	return ctx.renderTemplate("cluster_view.html")
}

func setupClusterRoutes(grp *echo.Group) {
	grp.GET("/cluster/list", showClusterList)
	grp.POST("/cluster/list", showClusterList)
	grp.GET("/cluster/view/:id", showClusterDetail)
	grp.POST("/cluster/view/:id", showClusterDetail)
}
