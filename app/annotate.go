package app

import (
	"fmt"
	//"strings"
	"net/http"
	"strconv"
	"time"
	"sort"
	"encoding/json"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type annotItem struct {
	Line int	`bson:"line" json:"line"`
	Seq int		`bson:"seq" json:"seq"`
	Text string	`bson:"text" json:"text"`
	Meta string	`bson:"meta" json:"meta"`
}

func (tp *annotItem) UnmarshalJSON(data []byte) error {
    var v []interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        fmt.Printf("Error whilde decoding %v\n", err)
        return err
	}
	
    tp.Line = int(v[0].(float64))
	tp.Seq = int(v[1].(float64))
	tp.Text = v[2].(string)
	tp.Meta = v[3].(string)
    return nil
}

type MdbAnnotate struct {
	Id bson.ObjectId	`bson:"_id,omitempty" json:"DT_RowId"`
	NewsId string       `bson:"newsId" json:"newsId"`
	InsertDt time.Time  `bson:"insertDt" json:"insertDt"`
	UpdateDt time.Time  `bson:"updateDt" json:"updateDt"`
	Category string		`bson:"category" json:"category"`
	Title string        `bson:"title" json:"title"`
	Content string  	`bson:"content" json:"content"`

	MediaId string   	`bson:"mediaId" json:"mediaId"`
	MediaName string	`bson:"mediaName" json:"mediaName"`
	Url string          `bson:"url" json:"url"`

	RealInfo []annotItem	`bson:"a_realInfo" json:"real_info"`
	RealQuote []annotItem	`bson:"a_realQuote" json:"real_quote"`
	RealIndir []annotItem	`bson:"a_realIndir" json:"real_indir"`
	AnonInfo []annotItem	`bson:"a_anonInfo" json:"anon_info"`
	AnonQuote []annotItem	`bson:"a_anonQuote" json:"anon_quote"`
	AnonIndir []annotItem	`bson:"a_anonIndir" json:"anon_indir"`
	ManJob []annotItem		`bson:"a_manJob" json:"man_job"`
	ManPub []annotItem		`bson:"a_manPub" json:"man_pub"`
	ManNews []annotItem		`bson:"a_manNews" json:"man_news"`
	Institute []annotItem	`bson:"a_institute" json:"institute"`
	Place []annotItem		`bson:"a_place" json:"place"`
	Number []annotItem		`bson:"a_number" json:"number"`
	Disgust []annotItem		`bson:"a_disgust" json:"disgust"`
}

type AnnotListItem struct {
	Id bson.ObjectId	`json:"DT_RowId"`
	UpdateDt string		`bson:"updateDt" json:"updateDt"`
	Category string		`bson:"category" json:"category"`
	Title string        `bson:"title" json:"title"`
	MediaName string        `bson:"mediaName" json:"mediaName"`
}

func (src *MdbAnnotate) MakeListItem() AnnotListItem {
	return AnnotListItem{
		src.Id,
		src.UpdateDt.Format("2006-01-02 15:04:05"),
		src.Category,
		src.Title,
		src.MediaName,
	}
}


func jsonAnnotateList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}

    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterMediaName := ctx.FormValue("columns[1][search][value]")
	filterCategory := ctx.FormValue("columns[2][search][value]")
	filterTitle := ctx.FormValue("columns[3][search][value]")
	filterNewsId := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var findParam = make(bson.M)
	if filterMediaName != "" {
		findParam["mediaName"] = bson.RegEx{filterMediaName, ""}
	}
	if filterCategory != "" {
		findParam["category"] = bson.RegEx{filterCategory, ""}
	}
	if filterTitle != "" {
		findParam["title"] = bson.RegEx{filterTitle, ""}
	}
	if filterNewsId != "" {
		findParam["newsId"] = bson.RegEx{filterNewsId, ""}
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)

	var err error
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
			orderFieldName += "updateDt"
		case 1:
			orderFieldName += "category"
		case 2:
		default:
			orderFieldName += "title"
	}

	var results []MdbAnnotate
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).All(&results)
	if err != nil {
		return ctx.showAdminError(err)
	}

    json.Rows = make([]interface{}, len(results))
	for i,v := range results {
		json.Rows[i] = v.MakeListItem()
	}
	
	return ctx.JSON(http.StatusOK, json)
}

func showAnnotateList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3d"
	return ctx.renderTemplate("annot_list.html")
}

func showAnnotateDetail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var where bson.M
	objId := ctx.QueryParam("id")
	newsId := ctx.QueryParam("nws")
	if objId != "" {
		where = bson.M{"_id": bson.ObjectIdHex(objId)}
	} else if newsId != "" {
		where = bson.M{"newsId": newsId}
	} else {
		return ctx.showAdminMsg("Invalid ID")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)
	
	var doc MdbAnnotate
	err := coll.Find(where).One(&doc)
	if err != nil {
		return ctx.showAdminError(err)
	}

	coll = mgoDB.With(sessionCopy).C(COLLNAME_NEWS)
	var docMore MdbNews
	err = coll.Find(bson.M{"newsId":doc.NewsId}).One(&docMore)
	if err != nil {
		return ctx.showAdminError(err)
	}

	doc.MediaId = docMore.MediaId
	doc.MediaName = docMore.MediaName
	doc.Url = docMore.Url

	ctx.model["item"] = doc
	ctx.model["lmenu"] = "t3d"
	
	return ctx.renderTemplate("annot_view.html")
}

var mapJsonToBsonKeys = map[string]string {
	"real_info": "a_realInfo",
	"real_quote": "a_realQuote",
	"real_indir": "a_realIndir",
	"anon_info": "a_anonInfo",
	"anon_quote": "a_anonQuote",
	"anon_indir": "a_anonIndir",
	"man_job": "a_manJob",
	"man_pub": "a_manPub",
	"man_news": "a_manNews",
	"institute": "a_institute",
	"place": "a_place",
	"number": "a_number",
	"disgust": "a_disgust",
}

func saveAnnotateItem(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")
	jsonTxt := ctx.FormValue("json")

	var dat map[string][]annotItem
	if err := json.Unmarshal([]byte(jsonTxt), &dat); err != nil {
        return ctx.showAdminError(err)
    }

	var chgs = make(bson.M)
	var drops = make(bson.M)

	for jk,bk := range mapJsonToBsonKeys {
		if arr, ok := dat[jk]; ok {
			if len(arr) > 1 {
				sort.Slice(arr, func(a, b int) bool {
					return (arr[a].Line * 1000 + arr[a].Seq) < (arr[b].Line * 1000 + arr[b].Seq)
				})
			}
			chgs[bk] = arr
		} else {
			drops[bk] = nil
		}
	}

	chgs["updateDt"] = time.Now()

	var updParam = make(bson.M)
	updParam["$set"] = chgs
	if len(drops) > 0 {
		updParam["$unset"] = drops
	}
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)
	
	err := coll.UpdateId(bson.ObjectIdHex(objId), updParam)
	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.flashSuccess("저장되었습니다.");
	return ctx.Redirect(http.StatusFound, "view?id="+objId)
}

func deleteAnnotateItem(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)

	err := coll.RemoveId(bson.ObjectIdHex(objId))
	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.flashSuccess("삭제되었습니다.");
	return ctx.Redirect(http.StatusFound, "list")
}

func updateAnnotateCategory(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	objId := ctx.FormValue("id")
	cate := ctx.FormValue("cate")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_ANNOTATE)

	var data MdbAnnotate
	coll.FindId(bson.ObjectIdHex(objId)).One(&data)
	if data.Category == cate {
		ctx.flashError("기존 카테고리 값과 같아서 수정하지 않았습니다.")
	} else {
		change := bson.M{ "$set": bson.M{"category": cate} }
		err := coll.UpdateId(bson.ObjectIdHex(objId), change)
		if err != nil {
			return ctx.showAdminError(err)
		}

		ctx.flashSuccess("카테고리를 변경하였습니다.")
	}
	
	return ctx.Redirect(http.StatusFound, "view?id="+objId)
}

func setupAnnotateRoutes(grp *echo.Group) {
	grp.GET("/annotate/list", showAnnotateList)
	grp.POST("/annotate/list.json", jsonAnnotateList)
	grp.GET("/annotate/view", showAnnotateDetail)
	grp.POST("/annotate/save", saveAnnotateItem)
	grp.POST("/annotate/delete", deleteAnnotateItem)
	grp.POST("/annotate/editCate", updateAnnotateCategory)
}