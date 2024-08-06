package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/barmi/go-admin-themes/adminlte"
	_ "github.com/barmi/go-admin/adapter/fasthttp"
	_ "github.com/barmi/go-admin/modules/db/drivers/mysql"

	"github.com/barmi/go-admin/engine"
	"github.com/barmi/go-admin/examples/datamodel"
	"github.com/barmi/go-admin/modules/config"
	"github.com/barmi/go-admin/modules/language"
	"github.com/barmi/go-admin/plugins/example"
	"github.com/barmi/go-admin/template"
	"github.com/barmi/go-admin/template/chartjs"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()

	eng := engine.Default()

	cfg := config.Config{
		Env: config.EnvLocal,
		Databases: config.DatabaseList{
			"default": {
				Host:            "127.0.0.1",
				Port:            "3306",
				User:            "root",
				Pwd:             "root",
				Name:            "godmin",
				MaxIdleConns:    50,
				MaxOpenConns:    150,
				ConnMaxLifetime: time.Hour,
				Driver:          config.DriverMysql,
			},
		},
		UrlPrefix: "admin",
		IndexUrl:  "/",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Debug:    true,
		Language: language.CN,
	}

	template.AddComp(chartjs.NewChart())

	// customize a plugin

	examplePlugin := example.NewExample()

	// load from golang.Plugin
	//
	// examplePlugin := plugins.LoadFromPlugin("../datamodel/example.so")

	// customize the login page
	// example: https://github.com/GoAdminGroup/demo.go-admin.cn/blob/master/main.go#L39
	//
	// template.AddComp("login", datamodel.LoginPage)

	// load config from json file
	//
	// eng.AddConfigFromJSON("../datamodel/config.json")

	if err := eng.AddConfig(&cfg).
		AddGenerators(datamodel.Generators).
		AddDisplayFilterXssJsFilter().
		// add generator, first parameter is the url prefix of table when visit.
		// example:
		//
		// "user" => http://localhost:9033/admin/info/user
		//
		AddGenerator("user", datamodel.GetUserTable).
		AddPlugins(examplePlugin).
		Use(router); err != nil {
		panic(err)
	}

	router.ServeFiles("/uploads/*filepath", "./uploads")

	eng.HTML("GET", "/admin", datamodel.GetContent)

	go func() {
		_ = fasthttp.ListenAndServe(":8897", router.Handler)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
}
