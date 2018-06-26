package app

import (
	//"fmt"
	"strconv"
	"time"
	"sort"
	"net/http"
	
	"github.com/labstack/echo"
//	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type NsrcListItem struct {
	Id bson.ObjectId	`json:"DT_RowId"`
	InsertDt string  	`json:"insertDt"`
	MediaId string   	`json:"mediaId"`
	Title string        `json:"title"`
	Byline string       `json:"byline"`
	CategoryOrig string `json:"categoryOrig"`
}

func MakeNsrcListItem(_src interface{}) NsrcListItem {
	src := _src.(bson.M)
	return NsrcListItem{
		src["_id"].(bson.ObjectId),
		src["insert_dt"].(time.Time).Format("2006-01-02 15:04:05"),
		src["media_id"].(string),
		src["title"].(string),
		src["writer_byline"].(string),
		src["subjectinfo"].(string),
	}
}

func jsonNsrcList(_ctx echo.Context) error {
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
	filterCategory := ctx.FormValue("columns[2][search][value]")
	filterByline := ctx.FormValue("columns[3][search][value]")
	filterTitle := ctx.FormValue("columns[4][search][value]")
	filterNewsId := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)

	if filterNewsId != "" {
		findParam["newsitem_id"] = bson.RegEx{filterNewsId, ""}
	} else {		
		ymd := ctx.FormValue("ymd")
		dtBegin, err := time.ParseInLocation("2006-01-02", ymd, tzLocation)
		if err == nil {
			dtEnd := dtBegin.AddDate(0,0,1)
			findParam["insert_dt"] = bson.M{ "$gte":dtBegin, "$lt":dtEnd }
			//fmt.Printf("DTR: %s ~ %s\n", dtBegin, dtEnd)
		}

		if filterMediaId != "" {
			findParam["media_id"] = bson.RegEx{filterMediaId, ""}
		}
		if filterCategory != "" {
			findParam["subjectinfo"] = bson.RegEx{filterCategory, ""}
		}
		if filterByline != "" {
			findParam["writer_byline"] = bson.RegEx{filterByline, ""}
		}
		if filterTitle != "" {
			findParam["title"] = bson.RegEx{filterTitle, ""}
		}
	}
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS_SRC)

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
			orderFieldName += "insert_dt"
		case 1:
			orderFieldName += "media_id"
		case 2:
			orderFieldName += "subjectinfo"
		case 3:
			orderFieldName += "writer_byline"
		case 4:
		default:
			orderFieldName += "title"
	}

	selFields := bson.M{
		"insert_dt":1,
		"media_id":1,
		"subjectinfo":1,
		"writer_byline":1,
		"title":1,
	}
	var results []interface{}
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(selFields).All(&results)
	if err != nil {
		panic(err)
	}

    json.Rows = make([]interface{}, len(results))
	for i,v := range results {
		json.Rows[i] = MakeNsrcListItem(v)
	}
	
	return ctx.JSON(http.StatusOK, json)
}

func showNsrcList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3b"
	ctx.model["targetDate"] = "2016-06-01"
	return ctx.renderTemplate("nsrc_list.html")
}

func showNsrcDetail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var where bson.D
	objId := ctx.QueryParam("obj")
	newsId := ctx.QueryParam("nws")
	if objId != "" {
		where = bson.D{{"_id", bson.ObjectIdHex(objId)}}
	} else if newsId != "" {
		where = bson.D{{"newsitem_id", newsId}}
	} else {
		return ctx.showAdminMsg("Invalid ID")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_NEWS_SRC)

	var data bson.M
	err := coll.Find(where).One(&data)
	if err != nil {
		return ctx.showAdminMsg("Object not found")
	}

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	data["_id"] = data["_id"].(bson.ObjectId).Hex()
	ctx.model["item"] = data
	ctx.model["keys"] = keys
	ctx.model["lmenu"] = "t3b"
	
	return ctx.renderTemplate("nsrc_view.html")
}


func setupNsrcRoutes(grp *echo.Group) {
	// 수집 DB
	grp.GET("/nsrc/list", showNsrcList)
	grp.POST("/nsrc/list.json", jsonNsrcList)
	grp.GET("/nsrc/view", showNsrcDetail)	
}