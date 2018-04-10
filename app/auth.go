package app

import (
	//"fmt"
	"time"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AdminUser struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"DT_RowId"`
	CreatedAt BJsonTime `bson:"createdAt" json:"createdAt"`
	LastLoginAt BJsonTime `bson:"lastLoginAt" json:"lastLoginAt"`

	LoginId string `bson:"loginId" json:"loginId"`
	Passwd bson.Binary `bson:"passwd" json:"-"`
	Level int `bson:"level" json:"level"`

	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Phone string `bson:"phone" json:"phone"`
	
	Roles []string `bson:"roles" json:"roles"`
}

func (ctx *AuthContext) isLoggedIn() bool {
	return (ctx.authUserId != "")
}

const supervisorLevelNumber = 9

func (ctx *AuthContext) isLevelSupervisor() bool {
	return (ctx.model[MVK_AuthLevel] == supervisorLevelNumber)
}

func doLogout(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	for key := range ctx.sess.Values {
		delete(ctx.sess.Values, key)
	}

	ctx.authUserId = ""
	delete(ctx.model, MVK_AuthUserId)
	delete(ctx.model, MVK_AuthLevel)
	delete(ctx.model, MVK_AuthName)

	return ctx.Redirect(http.StatusFound, "/login")
}

func showLoginPage(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	url := ctx.QueryParam("url")
	if url == "" {
		url = "/admin/cluster/list"
	}
	ctx.model["retUrl"] = url

	return ctx.renderTemplate("login.html")
}

func encodePasswd(raw string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return hashedPassword
}

func _initUser(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var user AdminUser
	user.CreatedAt = NewJsonNow()
	user.LoginId = "admin"
	user.Passwd.Data = encodePasswd("aa")
	user.Name = "관리자"
	user.Email = "admin@ntrust.com"
	user.Phone = "010-1234-5678"
	user.Level = supervisorLevelNumber;
	user.Roles = []string{"SUPERVISOR"}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)
	
	err := coll.Insert(user)
	if err != nil {
		panic(err)
	}

	return ctx.String(http.StatusOK, "Done")
}

func checkLogin(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	loginId := ctx.FormValue("login")
	passwdRaw := ctx.FormValue("passwd")
	remember := ctx.FormValue("remember")
	retUrl := ctx.FormValue("retUrl")

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoDB.With(sessionCopy).C(COLLNAME_USER)

	var errMsg string
	var userInfo AdminUser
	err := coll.Find(bson.M{"loginId":loginId}).Select(bson.M{"_id":1, "passwd":1, "name":1, "level":1}).One(&userInfo)
	if err != nil {
		if err == mgo.ErrNotFound {
			errMsg = "해당 유저를 찾을 수 없습니다."
		} else {
			panic(err)
		}
	} else if bcrypt.CompareHashAndPassword(userInfo.Passwd.Data, []byte(passwdRaw)) != nil {
		errMsg = "아이디 또는 비밀번호가 틀립니다."
	} else {
		if remember == "1" {
			// TODO
		}
		ctx.sess.Values[SVK_AuthUserId] = userInfo.Id.Hex()
		ctx.sess.Values[SVK_AuthLevel] = userInfo.Level
		ctx.sess.Values[SVK_AuthName] = userInfo.Name
		//fmt.Printf("Login3 %s, %s\n", ctx.sess.Get(SVK_AuthUserId), ctx.sess.Get(SVK_AuthName))		
		ctx.flashSuccess("로그인하였습니다.")
		err = ctx.sess.Save(ctx.Request(), ctx.Response())
		if err != nil {
			return ctx.showAdminMsg("Cookie save failed: " + err.Error())
		}

		coll.UpdateId(bson.ObjectId(userInfo.Id), bson.M{"$set":bson.M{"lastLoginAt":time.Now()}})
		LogInfo("AdminWeb",ctx, "%s(%s)가 로그인 하였습니다.", userInfo.Name, loginId)

		if retUrl != "" {
			return ctx.Redirect(http.StatusFound, retUrl)
		} else {
			return ctx.Redirect(http.StatusFound, "/admin/cluster/list")
		}
	}

	if errMsg == "" {
		errMsg = "알 수 없는 오류입니다."
	}	
	ctx.flashError(errMsg)

	return ctx.Redirect(http.StatusFound, "/login")
}

func cookieTest(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	/*var userInfo AdminUser
	userInfo.Id = "userInfo2.Id"
	userInfo.Name = "userInfo2.Name"
	ctx.sess.Set(SVK_AuthUserId, userInfo.Id.Hex())
	ctx.sess.Set(SVK_AuthName, userInfo.Name)
	fmt.Printf("Check1 %s, %s\n", ctx.sess.Get(SVK_AuthUserId), ctx.sess.Get(SVK_AuthName))*/
	ctx.flashSuccess("로그인하였습니까?")
	ctx.flashSuccess("로그인하였습니다.")
	ctx.flashError("로그인하였습니까?")
	ctx.flashError("로그인하였습니다.")
	return ctx.Redirect(http.StatusFound, "/admin/assess/list")
}

func setupAuthRoutes(ec *echo.Echo) {

	ec.GET("/ck3", cookieTest)

	ec.GET("/_initUser", _initUser)
	ec.GET("/logout", doLogout)
	ec.GET("/login", showLoginPage)
	ec.POST("/login_check", checkLogin)
}
