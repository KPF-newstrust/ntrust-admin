package app

import (
	"os"
	"io"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	POS_UNKNOWN = 0
	POS_NOUN    = "명사"
	POS_ADVERB  = "부사"
	POS_PREDICATE = "용언"
	POS_IRREGULAR = "불규칙활용"
	POS_JOSA    = "조사"
)

type PosCommon struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"DT_RowId"`
	Timestamp time.Time `bson:"ts" json:"-"`
	Word string `bson:"word" json:"word"`
	Pos string `bson:"pos" json:"pos"`
	Meta []string `bson:"meta,omitempty" json:"meta"`
	PosDetail string `bson:"dpos" json:"dpos"`
}

func (u *PosCommon) MarshalJSON() ([]byte, error) {
	type Alias PosCommon
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"ts"`
	}{
		Alias:    (*Alias)(u),
		Timestamp: u.Timestamp.Format("2006-01-02 15:04:05"),
	})
}

type PosJosa struct {
	PosCommon `bson:",inline"`
	Zong      int `bson:"zong" json:"zong"`
}

type PosIrrConj struct {
	PosCommon `bson:",inline"`
	Root string `bson:"root" json:"root"`
	Tail string `bson:"tail" json:"tail"`
}

type PosMergedAll struct {
	PosCommon `bson:",inline"`
	Zong int `bson:"zong" json:"zong"`
	Root string `bson:"root" json:"root"`
	Tail string `bson:"tail" json:"tail"`
}

type bulkInsertResult struct {
	PosCommon
	Error string
}

func jsonWordList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
    }

    orderDir := ctx.FormValue("order[0][dir]")
	orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterWord := ctx.FormValue("columns[1][search][value]")
	filterPos := ctx.FormValue("columns[2][search][value]")
	filterMeta := ctx.FormValue("columns[3][search][value]")	

	var jsonData JqDataTable
    jsonData.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	var findParam = make(bson.M)

	if filterWord != "" {
		findParam["word"] = bson.RegEx{filterWord, ""}
	}
	if filterPos != "" {
		findParam["pos"] = filterPos
	}
	if filterMeta != "" {
		metas := strings.Split(filterMeta,"\t")
		findParam["meta"] = bson.M{"$all":metas}
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_POSDIC)

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
		orderFieldName += "word"
	case 2:
		orderFieldName += "pos"
	default:
		orderFieldName += "ts"
	}

	var results []PosCommon
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(bson.M{
			"ts":1,
			"word":1,
			"pos":1,
			"dpos":1,
			"meta":1,
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

func showWordList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t5a"
	return ctx.renderTemplate("dict/list.html")
}


func (ctx *AuthContext) _readPosDic(data interface{}) error {
	posId := ctx.QueryParam("id")
	if posId == "" {
		return nil
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_POSDIC)

	return coll.FindId(bson.ObjectIdHex(posId)).One(data)
}

func showNoun(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	err := ctx._readPosDic(&pos)
	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.model["pos"] = &pos
	ctx.model["metaOpts"] = []string{ "시간","정치인","연예인","체육인","예술인","경제인","학자", "숫자","음식","교통수단","공산품" }
	ctx.model["lmenu"] = "t5b"
	return ctx.renderTemplate("dict/noun.html")
}

func showSolo(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	err := ctx._readPosDic(&pos)
	if err != nil {
		return ctx.showAdminError(err)
	}

	ctx.model["pos"] = &pos
	ctx.model["lmenu"] = "t5c"
	return ctx.renderTemplate("dict/solo.html")
}

func showJosa(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosJosa
	err := ctx._readPosDic(&pos)
	if err != nil {
		return ctx.showAdminError(err)
	}
	
	ctx.model["pos"] = &pos
	ctx.model["metaOpts"] = []string{ "주격","목적격","보격" }
	ctx.model["lmenu"] = "t5d"
	return ctx.renderTemplate("dict/josa.html")
}

func showPredicate(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	err := ctx._readPosDic(&pos)
	if err != nil {
		return ctx.showAdminError(err)
	}
	
	ctx.model["pos"] = &pos
	ctx.model["metaOpts"] = []string{ "행위","추상","색깔","맛" }
	ctx.model["lmenu"] = "t5e"
	return ctx.renderTemplate("dict/predicate.html")
}

func showIrregularConjugate(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosIrrConj
	err := ctx._readPosDic(&pos)
	if err != nil {
		return ctx.showAdminError(err)
	}
	
	ctx.model["pos"] = &pos
	ctx.model["lmenu"] = "t5f"
	return ctx.renderTemplate("dict/irregular.html")
}

func (ctx *AuthContext) _processPosPost(dicname string, pos *PosCommon, posAll interface{},
	fnInit func(), fnEdit func(bson.M)) error {

	cmd := ctx.FormValue("cmd")

	pos.Timestamp = time.Now()
	pos.Word = strings.TrimSpace(ctx.FormValue("word"))
	if pos.Word == "" && cmd != "del" {
		return ctx.showAdminMsg("어절이 입력되지 않았습니다.")
	}

	vars,_ := ctx.FormParams()
	for _, meta := range vars["meta"] {
		pos.Meta = append(pos.Meta, meta)
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_POSDIC)

	var err error
	if cmd == "new" {

		pos.Id = bson.NewObjectId()
		fnInit()
		err = coll.Insert(posAll)
		if err != nil {
			if mgo.IsDup(err) {
				ctx.flashError("이미 입력된 단어입니다.")
				return ctx.Redirect(http.StatusFound, "/admin/dict/"+dicname)
			}

			return ctx.showAdminError(err)
		}

		ctx.flashSuccess("신규 등록하였습니다: \"%s\" (%s)", pos.Word, pos.Pos)

	} else {

		posId := ctx.FormValue("id")
		if !bson.IsObjectIdHex(posId) {
			return ctx.showAdminMsg("ID가 올바르지 않습니다: " + posId)
		}

		pos.Id = bson.ObjectIdHex(posId)
		if cmd == "del" {
			err = coll.RemoveId(pos.Id)
			if err != nil {
				return ctx.showAdminError(err)
			}

			ctx.flashSuccess("삭제하였습니다.")
			return ctx.Redirect(http.StatusFound, "/admin/dict/"+dicname)
		}
		
		if cmd == "edit" {
			fnInit()

			sets := bson.M{"word":pos.Word, "ts":pos.Timestamp}
			upds := bson.M{"$set":sets}
			if len(pos.Meta) == 0 {
				upds["$unset"] = bson.M{"meta":""}
			} else {
				sets["meta"] = pos.Meta
			}

			fnEdit(sets)
			err = coll.Update(bson.M{"_id":pos.Id}, upds)
			if err != nil {
				return ctx.showAdminError(err)
			}

			ctx.flashSuccess("수정하였습니다: \"%s\"", pos.Word)
		}
	}

	return ctx.Redirect(http.StatusFound, fmt.Sprintf("/admin/dict/%s?id=%s", dicname, pos.Id.Hex()))
}

func postNoun(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	for n:=0; n<=10; n++ {
		neMeta := ctx.FormValue(fmt.Sprintf("ne%d", n))
		if neMeta != "" {
			pos.Meta = append(pos.Meta, neMeta)
		}
	}

	return ctx._processPosPost("noun", &pos, &pos, func() {
		pos.Pos = POS_NOUN
	}, func(_m bson.M) {
	})
}

func postSolo(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	return ctx._processPosPost("solo", &pos, &pos, func() {
		pos.Pos = ctx.FormValue("pos")
	}, func(_m bson.M) {
		_m["pos"] = pos.Pos
	})
}

func postJosa(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosJosa
	return ctx._processPosPost("josa", &pos.PosCommon, &pos, func() {
		pos.Pos = POS_JOSA
		pos.Zong, _ = strconv.Atoi(ctx.FormValue("zong"))
	}, func(_m bson.M) {
		_m["zong"] = pos.Zong
	})	
}

func postPredicate(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosCommon
	return ctx._processPosPost("predicate", &pos, &pos, func() {
		pos.Pos = POS_PREDICATE
		pos.PosDetail = ctx.FormValue("dpos")
	}, func(_m bson.M) {
		_m["dpos"] = pos.PosDetail
	})
}

func postIrregularConjugate(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var pos PosIrrConj
	return ctx._processPosPost("irregular", &pos.PosCommon, &pos, func() {
		pos.Pos = POS_IRREGULAR
		pos.Root = strings.TrimSpace(ctx.FormValue("root"))
		pos.Tail = strings.TrimSpace(ctx.FormValue("tail"))
	}, func(_m bson.M) {
		_m["root"] = pos.Root
		_m["tail"] = pos.Tail
	})
}

func showBulkForm(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	ctx.model["lmenu"] = "t5x"
	return ctx.renderTemplate("dict/bulk_form.html")
}

func postBulk(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	pos := ctx.FormValue("pos")
	var dpos string
	if pos == "동사" || pos == "형용사" {
		dpos = pos
		pos = "용언"
	}

	words := strings.FieldsFunc(ctx.FormValue("words"), func(r rune) bool {
		return (r == ' ' || r == '\t' || r == '\r' || r == '\n')
	})

	var metas []string
	vars,_ := ctx.FormParams()
	for _, meta := range vars["meta"] {
		metas = append(metas, meta)
	}

	if pos == POS_NOUN {
		neMeta := ctx.FormValue("nemeta")
		if neMeta != "" {
			metas = append(metas, neMeta)
		}
	}	

	if len(words) == 0 {
		ctx.flashError("입력할 단어가 없습니다.")
		return ctx.Redirect(http.StatusFound, "/admin/dict/bulk")
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_POSDIC)

	results := make([]bulkInsertResult, len(words))
	for i,word := range words {
		results[i].Word = word
		results[i].Pos = pos
		results[i].PosDetail = dpos

		if word == "" {
			results[i].Error = "빈 단어"
			continue
		}

		results[i].Timestamp = time.Now()
		
		var old PosCommon
		err := coll.Find(bson.M{"word":word, "pos":pos}).Select(bson.M{"meta":1}).One(&old)
		if err == nil {
			// duplicated word
			if len(metas) == 0 {
				results[i].Error = "이미 입력된 단어이나 병합할 메타 정보가 없습니다."
			} else {
				numOldMetas := len(old.Meta)
				if numOldMetas == 0 {
					old.Meta = metas
				} else {
					for _,newmeta := range metas {
						found := false
						for oi:=0; oi<numOldMetas; oi++ {
							if old.Meta[oi] == newmeta {
								found = true
								break
							}
						}
						if !found {
							old.Meta = append(old.Meta, newmeta)
						}
					}	
				}

				results[i].Meta = old.Meta
				err = coll.UpdateId(old.Id, bson.M{"$set":bson.M{"meta":old.Meta, "ts":results[i].Timestamp}})
				if err != nil {
					results[i].Error = err.Error()
				} else {
					results[i].Error = "성공: 메타 정보를 병합하였습니다."
				}
			}
		} else if err == mgo.ErrNotFound {
			// new word
			results[i].Id = bson.NewObjectId()
			if len(metas) > 0 {
				results[i].Meta = metas
			}

			err = coll.Insert(&results[i].PosCommon)
			if err != nil {
				if mgo.IsDup(err) {
					results[i].Error = "이미 입력된 정보입니다."
				} else {
					results[i].Error = err.Error()
				}
			}
		}

		//fmt.Printf("result: %s %s %v\n", results[i].Word, results[i].Pos, results[i].Meta)
		
	}

	ctx.model["results"] = results
	ctx.model["lmenu"] = "t5x"
	return ctx.renderTemplate("dict/bulk_result.html")
}

func exportPosDic(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	wcnt, err := writePosDic(file)
	if err != nil {
		return err
	}

	file.Sync()
	fmt.Printf("Dictionary source file exported %s (%d words)\n", filename, wcnt)
	
	return nil
}

func writePosDic(wr io.Writer) (int,error) {
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_POSDIC)

	var arr []PosMergedAll
	err := coll.Find(nil).Sort("word").All(&arr)
	if err != nil {
		return 0, err
	}

	wcnt := 0
	for _,v := range arr {
		line := fmt.Sprintf(`%s:{"pos":"%s"`, v.Word, v.Pos)
		
		if v.PosDetail != "" {
			line += fmt.Sprintf(`,"dpos":"%s"`, v.PosDetail)
		}

		if v.Pos == POS_JOSA {
			var zongCond string
			switch (v.Zong) {
				case 0:
					zongCond = "all"
				case 1:
					zongCond = "yes"
				case 2:
					zongCond = "no"
				default:
					zongCond = fmt.Sprintf("err(%d)", v.Zong)
			}
			line += fmt.Sprintf(`,"zong":"%s"`, zongCond)
		} else if v.Pos == POS_IRREGULAR {
			line += fmt.Sprintf(`,"root":"%s","tail":"%s"`, v.Root, v.Tail)
		}

		if len(v.Meta) > 0 {
			meta, _ := json.Marshal(v.Meta)
			line += fmt.Sprintf(`,"meta":%s`, string(meta))
		}

		_, err = wr.Write([]byte(fmt.Sprintf("%s}\n", line)))
		if err != nil {
			return wcnt, err
		}

		wcnt++
	}

	return wcnt, nil
}

func downloadPosdic(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	buf := new(bytes.Buffer)
	wcnt, err := writePosDic(buf)
	if err != nil {
		return ctx.showAdminError(err)
	}

	filename := fmt.Sprintf("posdic-%d.txt", wcnt)
	ctx.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	_, err = ctx.Response().Write(buf.Bytes())

	return err
}

func setupPosdicRoutes(grp *echo.Group) {
	grp.GET("/dict/list", showWordList)
	grp.POST("/dict/list.json", jsonWordList)

	grp.GET("/dict/noun", showNoun)
	grp.POST("/dict/noun", postNoun)

	grp.GET("/dict/solo", showSolo)
	grp.POST("/dict/solo", postSolo)

	grp.GET("/dict/josa", showJosa)
	grp.POST("/dict/josa", postJosa)

	grp.GET("/dict/predicate", showPredicate)
	grp.POST("/dict/predicate", postPredicate)

	grp.GET("/dict/irregular", showIrregularConjugate)
	grp.POST("/dict/irregular", postIrregularConjugate)
	
	grp.GET("/dict/bulk", showBulkForm)
	grp.POST("/dict/bulk", postBulk)

	grp.GET("/dict/_down", downloadPosdic)
}
