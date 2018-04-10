package app

import (
	//"fmt"
	"strconv"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func showUserList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t9a"
	return ctx.renderTemplate("user_list.html")
}

func jsonUserList(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}
	
    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterSearch := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)

	if filterSearch != "" {
		findParam["name"] = bson.RegEx{filterSearch, ""}
	}
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

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
			orderFieldName += "createdAt"
		case 1:
			orderFieldName += "email"
		case 2:
			orderFieldName += "name"
		case 3:
			orderFieldName += "lastLoginAt"
		case 4:
		default:
			orderFieldName += "level"
	}

	selFields := bson.M{
		"createdAt":1,
		"loginId":1,
		"name":1,
		"email":1,
		"phone":1,
		"lastLoginAt":1,
		"level":1,
		"roles":1,
	}
	
	var results []AdminUser
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(selFields).All(&results)
	if err != nil {
		panic(err)
	}

    json.Rows = make([]interface{}, len(results))
	for i := range results {
		json.Rows[i] = &results[i]
	}
/*
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).Select(selFields).All(&json.Rows)
	if err != nil {
		panic(err)
	}
*/
	return ctx.JSON(http.StatusOK, json)
}

func addNewUser(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	loginId := ctx.FormValue("loginid")
	passwd := ctx.FormValue("passwd")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

	cntDupId, _ := coll.Find(bson.M{"loginId":loginId}).Count()
	if cntDupId > 0 {
		ctx.model["lmenu"] = "t9a"
		return ctx.showAdminMsg("이미 등록된 로그인 아이디 입니다.")
	}
/*
	cntDupEmail, _ := coll.Find(bson.M{"email":email}).Count()
	if cntDupEmail > 0 {
		ctx.model["lmenu"] = "t9a"
		return ctx.showAdminMsg("이미 등록된 이메일 주소 입니다.")
	}
*/	
	var user AdminUser
	user.Id = bson.NewObjectId()
	user.CreatedAt = NewJsonNow()
	user.LoginId = loginId
	user.Passwd.Data = encodePasswd(passwd)
	user.Name = "(이름 미등록)"
	user.Email = "(이메일 미등록)"
	user.Phone = "(전화번호 미등록)"
	user.Level = 0;
	
	err := coll.Insert(user)
	if err != nil {
		return ctx.showAdminError(err)
	}

	adminName := ctx.model[MVK_AuthName]
	LogInfo("AdminWeb",ctx, "관리자 %s가 신규 사용자 %s를 생성하였습니다.", adminName, user.LoginId)

	return ctx.Redirect(http.StatusFound, "/admin/user/view?id="+user.Id.Hex())
}

func showUserDetail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	objId := ctx.FormValue("id")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

	var data AdminUser
	err := coll.FindId(bson.ObjectIdHex(objId)).One(&data)
	if err != nil {
		ctx.flashError("관리자 정보를 찾을 수 없습니다.")
		return ctx.Redirect(http.StatusFound, "/admin/user/list")
	}

	ctx.model["lmenu"] = "t9a"
	ctx.model["item"] = data

	return ctx.renderTemplate("user_view.html")
}

func changeUserPasswd(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")
	rawPasswd := ctx.FormValue("passwd")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

	var data AdminUser
	objIdHex := bson.ObjectIdHex(objId)
	err := coll.FindId(objIdHex).Select(bson.M{"name":1}).One(&data)
	if err != nil {
		return ctx.showAdminError(err)
	}

	data.Passwd.Data = encodePasswd(rawPasswd)

	err = coll.UpdateId(objIdHex, bson.M{"$set":bson.M{"passwd":data.Passwd}})
	if err != nil {
		return ctx.showAdminError(err)
	}
	
	adminName := ctx.model[MVK_AuthName]
	LogInfo("AdminWeb", ctx, "%s님이 %s의 비밀번호를 변경하였습니다.", adminName, data.Name)

	ctx.flashSuccess("비밀번호를 변경하였습니다.")
	return ctx.Redirect(http.StatusFound, "/admin/user/view?id="+objId)
}

func changeUserInfo(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")
	name := ctx.FormValue("name")
	email := ctx.FormValue("email")
	phone := ctx.FormValue("phone")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

	var data AdminUser
	objIdHex := bson.ObjectIdHex(objId)
	err := coll.FindId(objIdHex).Select(bson.M{"name":1, "email":1, "phone":1}).One(&data)
	if err != nil {
		return ctx.showAdminError(err)
	}

	//nameChanged := false
	var changes = make(bson.M)
	if data.Name != name {
		//nameChanged = true
		changes["name"] = name
		ctx.flashSuccess("이름을 변경하였습니다.")
	}

	if data.Email != email {
		changes["email"] = email
		ctx.flashSuccess("이메일 주소를 변경하였습니다.")
	}

	if data.Phone != phone {
		changes["phone"] = phone
		ctx.flashSuccess("전화번호를 변경하였습니다.")
	}

	if len(changes) == 0 {
		ctx.flashError("변경사항이 없어서 업데이트하지 않았습니다.")
	} else {
		err := coll.UpdateId(objIdHex, bson.M{"$set":changes})
		if err != nil {
			return ctx.showAdminError(err)
		}
	}

	adminName := ctx.model[MVK_AuthName]
	LogInfo("AdminWeb", ctx, "%s님이 %s의 기본 정보를 변경하였습니다.", adminName, name)

	return ctx.Redirect(http.StatusFound, "/admin/user/view?id="+objId)
}

func changeUserPermmision(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	objId := ctx.FormValue("id")

	adminName := ctx.model[MVK_AuthName]
	LogInfo("AdminWeb",ctx, "관리자 %s가 누구의 권한 정보를 변경하였습니다.", adminName)
	
	ctx.flashSuccess("권한 정보를 변경하였습니다.")
	return ctx.Redirect(http.StatusFound, "/admin/user/view?id="+objId)
}

func setupStaffRoutes(grp *echo.Group) {
	grp.GET("/user/list", showUserList)
	grp.POST("/user/list.json", jsonUserList)
grp.GET("/user/list.json", jsonUserList)
	grp.POST("/user/newadmin", addNewUser)
	grp.GET("/user/view", showUserDetail)
	grp.POST("/user/chgPW", changeUserPasswd)
	grp.POST("/user/chgInfo", changeUserInfo)
	grp.POST("/user/chgPerm", changeUserPermmision)
}