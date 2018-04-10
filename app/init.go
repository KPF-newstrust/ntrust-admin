package app

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"net/url"
	"encoding/json"
	
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo-contrib/session"
	"github.com/gorilla/sessions"
	"github.com/boj/redistore"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
	"github.com/patrickmn/go-cache"
)

const (
	SVK_AuthUserId = "SAUT_UID"
	SVK_AuthLevel = "SAUT_Level"
	SVK_AuthName = "SAUT_Name"

	MVK_AuthUserId = "_AuthUID"
	MVK_AuthLevel = "_AuthLevel"
	MVK_AuthName = "_AuthName"
	MVK_AuthPhoto = "_AuthPhoto"

	_FLK_Success = "flSuccess"
	_FLK_Error = "flError"


	
	COLLNAME_YEAR1NEWS = "news"

	COLLNAME_USER = "user"
	COLLNAME_NEWS = "news"
	COLLNAME_ASSTATS = "asStats"
	COLLNAME_NEWS_ENTITY = "news_entity"
	COLLNAME_SRCNEWS = "cm2016"
	COLLNAME_ANNOTATE = "annotate"
	COLLNAME_BUGREPORT = "bugs"
	COLLNAME_POSDIC = "posdic"
	COLLNAME_APILAB = "apilab"
	
	_MaxRedisConnections = 10	
)

type AppConfig struct {
	IsDebug bool
	UseMinified bool
	UseSSL bool
	ListenPort int
	GatewayPort int
	Hostname string
	DsnMongo string
	DsnRedis string
	DsnAmqp string
}

type AuthContext struct {
	echo.Context
	model map[string]interface{}
	sess *sessions.Session
	authUserId string
}

func (ctx *AuthContext) renderTemplate(filename string) error {
	return ctx.Render(http.StatusOK, filename, ctx.model)
}

func (ctx *AuthContext) renderError(err error) error {
	ctx.model["errorMsg"] = err.Error()
	return ctx.Render(http.StatusInternalServerError, "error.html", ctx.model)
}

func (ctx *AuthContext) flashSuccess(format string, args ...interface{}) {
	ctx.sess.AddFlash(fmt.Sprintf(format, args...), _FLK_Success)
}

func (ctx *AuthContext) flashError(format string, args ...interface{}) {
	ctx.sess.AddFlash(fmt.Sprintf(format, args...), _FLK_Error)
}

var (
	BaseHostname string

	mgoSession *mgo.Session
	mgoDB *mgo.Database
	mgoOldDB *mgo.Database
//	mgoFS *mgo.GridFS

	rdsPool *redis.Pool

	amqpConn *amqp.Connection
	amqpChn *amqp.Channel
	amQueue_Task amqp.Queue
	//amQueue_Notice amqp.Queue

	localCache *cache.Cache
	tzLocation *time.Location
)

func makeBaseHostname(config *AppConfig) {
	isHttps := ""
	if config.UseSSL {
		isHttps += "s"
	}

	addPort := ""
	if config.GatewayPort != 80 {
		addPort = fmt.Sprintf(":%d", config.GatewayPort)
	}

	BaseHostname = fmt.Sprintf("http%s://%s%s", isHttps, config.Hostname, addPort)
}

func _ensureSingleColumnUniqueIndex(collname string, fieldName string) {
	coll := mgoDB.With(mgoSession).C(collname)
	err := coll.EnsureIndex(mgo.Index{
		Key: []string{ fieldName },
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}
}

func ensureMongoIndices() {

	coll := mgoDB.With(mgoSession).C(COLLNAME_SRCNEWS)
	err := coll.EnsureIndex(mgo.Index{
		Key: []string{"insert_dt"},
		Unique: false,
		DropDups: false,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}

	coll = mgoDB.With(mgoSession).C(COLLNAME_ANNOTATE)
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{"newsId"},
		Unique: true,
		DropDups: false,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}

	coll = mgoDB.With(mgoSession).C(COLLNAME_BUGREPORT)
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{"newsId"},
		Unique: false,
		DropDups: false,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}

	coll = mgoOldDB.With(mgoSession).C(COLLNAME_YEAR1NEWS)
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{"date"},
		Unique: false,
		DropDups: false,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}

	coll = mgoDB.With(mgoSession).C(COLLNAME_POSDIC)
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{ "word", "pos" },
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}

	/*coll = mgoDB.With(mgoSession).C(COLLNAME_APITEST)
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{ "ts" },
		Unique: false,
		DropDups: false,
		Background: true,
		Sparse: true,
	})
	if err != nil {
		panic(err)
	}*/
	
}

func openMongoDB(config *AppConfig) bool {
	var err error

	fmt.Printf("MongoDB DSN=%s\n", config.DsnMongo)
	mgoSession, err = mgo.Dial(config.DsnMongo)
	if err != nil {
		log.Printf("MongoDB connect failed: %s\n", err.Error())
		return false
	}

	mgoSession.SetMode(mgo.Monotonic, true)
	mgoDB = mgoSession.DB("")
	mgoOldDB = mgoSession.DB("ntrust1")
	//mgoFS = mgoDB.GridFS("fs")

	ensureMongoIndices()
	return true
}

func ExportDic(config *AppConfig, filename string) bool {

	if !openMongoDB(config) {
		return false
	}

	exportPosDic(filename)

	Close()
	return true
}

////////////////////////////////////////////////////////////////////////////////

func openRabbitMQ(config *AppConfig) bool {
	var err error

	fmt.Printf("RabbitMQ DSN=%s\n", config.DsnAmqp)
	amqpConn, err = amqp.Dial(config.DsnAmqp)
	if err != nil {
		log.Printf("AMQP Dial failed: %v\n", err)
		return false
	}

	amqpChn, err = amqpConn.Channel()
	if err != nil {
		log.Printf("AMQP Channel open failed: %v\n", err)
		return false
	}

	amQueue_Task, err = amqpChn.QueueDeclare(
		"task",	// name
		true,	// durable
		false,	// delete when unused
		false,	// exclusive
		false,	// no-wait
		nil,	// arguments
	)
	if err != nil {
		log.Printf("QueueDeclare(task) failed: %v", err)
		return false
	}

	return true
}

func mqSend_Task(cmd string, ver int, data bson.M) error {
	data["cmd"] = cmd
	data["ver"] = ver
	jsonBody, err := json.Marshal(data)	
	if err == nil {
		err = amqpChn.Publish(
			"",     // exchange
			amQueue_Task.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body: jsonBody,
			})
	}

	return err
}

////////////////////////////////////////////////////////////////////////////////

func Close() {
	if amqpChn != nil {
		amqpChn.Close()
		amqpChn = nil
	}

	if amqpConn != nil {		
		amqpConn.Close()
		amqpConn = nil
	}

	if rdsPool != nil {
		rdsPool.Close()
		rdsPool = nil
	}

	if mgoSession != nil {
		mgoSession.Close()
		mgoSession = nil
	}
}

func Setup(ec *echo.Echo, config *AppConfig) bool {
	var err error
	
	tzLocation, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Printf("TimeZone loading failed: %v\n", err)
		return false
	}

	makeBaseHostname(config)

	// Middleware

	if !config.IsDebug {
		//ec.Use(middleware.Gzip())
		ec.Use(middleware.Logger())
		ec.Use(middleware.Recover())
	}

	if !openMongoDB(config) {
		return false
	}

	if !openRabbitMQ(config) {
		return false
	}
	
	// Redis
	// check https://github.com/soveran/redisurl/blob/master/redisurl.go
	fmt.Printf("Redis DSN=%s\n", config.DsnRedis)
	redisUrl, err := url.Parse(config.DsnRedis)
	rdsPool = redis.NewPool(func() (redis.Conn, error) {
		rdconn, err := redis.Dial("tcp", redisUrl.Host)
		if err != nil {
			return nil, err
		}

		return rdconn, err
	}, _MaxRedisConnections)
	// try connect redis
	{
		_rconn := rdsPool.Get()
		defer _rconn.Close()
		if _rconn.Err() != nil {
			log.Printf("Redis connect failed: %s\n", err.Error())
			return false
		}
	}

	store, err := redistore.NewRediStore(32, "tcp", redisUrl.Host, "", []byte("ntrust2nd"))
	if err != nil {
		log.Printf("NewRedisStore failed: %v\n", err)
		return false
	}

	ec.Use(session.Middleware(store))

	/*ec.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_csrf",
	}))*/

	ec.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sess, err := session.Get("session", ctx)
			if err != nil {
				return err
			}

			cc := &AuthContext{ctx,
				make(map[string]interface{}),
				sess,
				"",
			}

			authUserId := cc.sess.Values[SVK_AuthUserId]
			if authUserId != nil {
				cc.authUserId = authUserId.(string)
				cc.model[MVK_AuthUserId] = cc.authUserId
				cc.model[MVK_AuthLevel] = cc.sess.Values[SVK_AuthLevel]
				cc.model[MVK_AuthName] = cc.sess.Values[SVK_AuthName]
				cc.model[MVK_AuthPhoto] = "/public/img/noavatar.png"
			}
			
			_csrf := ctx.Get("csrf")
			if _csrf == nil {
				_csrf = "nil"
			}			
			cc.model["_csrf"] = _csrf
			
			cc.model["_flSuccess"] = cc.sess.Flashes(_FLK_Success)
			cc.model["_flError"] = cc.sess.Flashes(_FLK_Error)

			defer cc.sess.Save(ctx.Request(), ctx.Response())
			return h(cc)
		}
	})

	// Create a cache with a default expiration time of 5 minutes,
	// and which purges expired items every 30 seconds
	localCache = cache.New(5*time.Minute, 30*time.Second)

	setupAuthRoutes(ec)
	setupAdminRoutes(ec)
	setupApiRoutes(ec)
	setupClusterOpenRoutes(ec)

	ec.GET("/", func(_ctx echo.Context) error {
		ctx := _ctx.(*AuthContext)
		ctx.model["var1"] = 123;
		ctx.model["var2"] = "ABC";
		return ctx.Redirect(http.StatusFound, "/admin/cluster/list")
	})

	return true
}