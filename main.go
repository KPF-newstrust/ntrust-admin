package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"math"
	"strings"
	
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/CloudyKit/jet"
	//"gopkg.in/mgo.v2/bson"

	"ntrust_admin/app"
)

const APP_VERSION = "Dev-20171124a"

type templateImpl struct {
	jetSet *jet.Set
	jetVars jet.VarMap
}

func castTime(arg interface{}) time.Time {
	t1, ok := arg.(app.BJsonTime)
	if ok {
		return time.Time(t1)
	}
	t2, ok := arg.(*app.BJsonTime)
	if ok {
		return time.Time(*t2)
	}
	t3, ok := arg.(time.Time)
	if ok {
		return t3
	}
	t4, ok := arg.(*time.Time)
	if ok {
		return *t4
	}

	log.Panicf("not time? [%v]", arg)
	panic("Not a time object")
}

func initTemplateImpl(config *app.AppConfig) *templateImpl {
	t := new(templateImpl)

	if config.UseMinified {
		t.jetSet = jet.NewHTMLSet("./min.templates")
		log.Println("Uses minified HTML templates.")
	} else {
		t.jetSet = jet.NewHTMLSet("./templates")
	}
	
	t.jetSet.SetDevelopmentMode(config.IsDebug)

	t.jetVars = make(jet.VarMap)
	t.jetSet.AddGlobal("_version", APP_VERSION)
	t.jetSet.AddGlobal("_jsrev", "1")

	t.jetSet.AddGlobal("formatTillSec", func(mayTime interface{}) string {
		if mayTime == nil {
			return ""
		}

		tm := castTime(mayTime)
		return tm.Format("2006-01-02 15:04:05")
	})

	t.jetSet.AddGlobal("formatTillSecKor", func(mayTime interface{}) string {
		if mayTime == nil {
			return ""
		}

		tm := castTime(mayTime)
		return tm.Format("2006년 01월 02일 15시 04분 05초")
	})

	t.jetSet.AddGlobal("isZero", func(dt app.BJsonTime) bool {
		return time.Time(dt).IsZero()
	})

	t.jetSet.AddGlobal("linebreaksbr", func(text string) string {
		return strings.Replace(text, "\n", "<br>", -1)
	})

	t.jetSet.AddGlobal("linebreaksCRBR", func(_text interface{}) string {
		if text, ok := _text.(string); ok {
			text = strings.Replace(text, "<", "&lt;", -1)
			return strings.Replace(text, "\n", "<i>\u23CE</i><br />", -1)			
		}
		return "(ERROR: no string)"
	})

	t.jetSet.AddGlobal("noneIfNil", func(_text interface{}) string {
		if text, ok := _text.(string); ok {
			return text
		}
		if num, ok := _text.(int); ok {
			return fmt.Sprintf("%d", num)
		}
		if _text == nil {
			return ""
		}
		return fmt.Sprintf("%v", _text)
	})

	t.jetSet.AddGlobal("join", func(arr []string, sp string) string {
		return strings.Join(arr, sp)
	})

	t.jetSet.AddGlobal("hasElement", func(arr []string, elem string) bool {
		for _,a := range arr {
			if elem == a {
				return true
			}
		}
		return false
	})

	t.jetSet.AddGlobal("toFixed", func(val float64, dot int) string {		
		shift := 10.0
		for d:=1;d<dot; d++ {
			shift *= 10
		}
		val = math.Floor(val * shift);
		return fmt.Sprintf("%v", val / shift)
	})

	return t
}

func (t *templateImpl) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	jt, err := t.jetSet.GetTemplate(name)
	if err != nil {    	
		return err // template could not be loaded
	}

	//return jt.Execute(w, t.jetVars, data)
	
	// dev
	err = jt.Execute(w, t.jetVars, data)
	if err != nil {
		if !ctx.Echo().Debug {
			app.LogError("TemplateError",nil, "템플릿(%s) 렌더링중 에러 발생: %s", name, err.Error())
		}

		log.Panic(err)
	}

	return err
}

func getConfigInt(v *viper.Viper, key string, defValue int) int {
	if v.IsSet(key) {
		return v.GetInt(key)
	}

	return defValue
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".") 
	err := viper.ReadInConfig()
	if err != nil {
    	log.Printf("Fatal error config file: %s \n", err)
		return
	}
	
	var config app.AppConfig
	var exportDicFilename string

	// check command line options
	for idx:=1; idx<len(os.Args); idx++ {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "--dev") || strings.HasPrefix(arg, "--debug") {
			config.IsDebug = true
			continue
		}

		if strings.HasPrefix(arg, "--min") {
			config.UseMinified = true
			continue
		}

		if arg == "--exportdic" {
			idx++
			if idx < len(os.Args) {
				exportDicFilename = os.Args[idx]
				continue
			} else {
				fmt.Println("Output dictionary filename not specified.")
				return
			}
		}
	}

	// TODO: read default section values

	var subconf *viper.Viper;
	if config.IsDebug {
		log.Println("Debug mode")
		subconf = viper.Sub("dev")
	} else {
		log.Println("Production mode")
		subconf = viper.Sub("prod")
	}

	config.UseSSL = subconf.GetBool("use_ssl")
	config.Hostname = subconf.GetString("hostname")
	config.ListenPort = subconf.GetInt("listen_port")
	config.GatewayPort = subconf.GetInt("gateway_port")
	config.DsnMongo = subconf.GetString("dsn_mongo")
	config.DsnRedis = subconf.GetString("dsn_redis")
	config.DsnAmqp = subconf.GetString("dsn_amqp")

	if exportDicFilename != "" {
		app.ExportDic(&config, exportDicFilename)
		return
	}

	if config.ListenPort == 0 {
		log.Printf("Listen port is not specified")
		return
	}


	ec := echo.New()
	ec.Debug = config.IsDebug

	ec.Renderer = initTemplateImpl(&config)

	if subconf.GetBool("serve_static") {
		if config.UseMinified {
			ec.Static("/public", "min.static")
		} else {
			ec.Static("/public", "static")
		}		
	}

	ec.Static("/robots.txt", "static/robots.txt")
	
	if !app.Setup(ec, &config) {
		log.Println("App server setup failed.")
		return
	}

	//http.ListenAndServe(":8080", context.ClearHandler(e))
	ec.Logger.Fatal(ec.Start(fmt.Sprintf(":%d", config.ListenPort)))

	app.Close()
}
