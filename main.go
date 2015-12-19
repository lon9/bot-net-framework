package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/ChimeraCoder/anaconda"
	"strconv"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/binding"
)


func main(){

	// Parse options
	opts := parseFlags()

	// Prepare martini
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout: "layout",
		Extensions: []string{".html"},
		Funcs: []template.FuncMap{},
		Charset: "UTF-8",
		IndentJSON: true,
		IndentXML: true,
	}))
	m.Use(martini.Logger())
	store := sessions.NewCookieStore([]byte("secret"))
	m.Use(sessions.Sessions("session", store))

	// Init Database
	var db gorm.DB
	var err error
	switch opts.Database {
	case "mysql":
		db, err = gorm.Open("mysql", opts.DBOptions)
	case "postgres":
		db, err = gorm.Open("postgres", opts.DBOptions)
	case "sqlite":
		db, err = gorm.Open("sqlite", opts.DBOptions)
	}
	if err != nil{
		panic(err)
	}

	// Migrate database
	// Debug
	//db.DropTable(&Talk{}, &Bot{}, &Tweet{})
	db.AutoMigrate(&Talk{}, &Bot{}, &Tweet{})
	m.Map(db)

	// Init Twitter Api object
	anaconda.SetConsumerKey(opts.ConsumerKey)
	anaconda.SetConsumerSecret(opts.ConsumerSecret)

	// Index
	m.Get("/", Index)

	// Register bot handlers
	m.Group("/twitter", func(r martini.Router) {
		r.Get("/", TwitterLogin)
		r.Get("/callback", TwitterCallback)
	})

	// API handlers
	m.Group("/api", func(r martini.Router){
		m.Group("/bot", func(r martini.Router) {
			m.Get("", IndexBot)
			m.Get("/:id", GetBot)
			m.Post("", binding.Bind(Bot{}), CreateBot)
			m.Put("", binding.Bind(Bot{}), UpdateBot)
			m.Delete("/:id", DeleteBot)
		})
		m.Group("/talk", func(r martini.Router) {
			m.Get("", IndexTalk)
			m.Get("/:id", GetTalk)
			m.Post("", binding.Bind(Talk{}), CreateTalk)
			m.Put("", binding.Bind(Talk{}), UpdateTalk)
			m.Delete("/:id", DeleteTalk)
		})
		m.Group("/tweet", func(r martini.Router) {
			m.Get("", IndexTweet)
			m.Get("/:id", GetTweet)
			m.Post("", binding.Bind(Tweet{}), CreateTweet)
			m.Put("", binding.Bind(Tweet{}), UpdateTweet)
			m.Delete("/:id", DeleteTweet)
		})

		m.Get("/", StartTalk)
	})


	m.RunOnAddr(fmt.Sprintf(":%d", opts.Port))
}

func parseFlags() *Options{
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "bot-net"
	parser.Usage = "[OPTIONS]"
	_, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	if opts.Port == 0{
		if opts.Port, err = strconv.Atoi(os.Getenv("BN_PORT")); opts.Port == 0 || err != nil{
			fmt.Println("Port number must be set as option or environment value.")
			os.Exit(1)
		}
	}

	if opts.ConsumerKey == "" {
		if opts.ConsumerKey = os.Getenv("BN_CONSUMER_KEY"); opts.ConsumerKey == ""{
			fmt.Println("Consumer key must be set as option or environment value.")
			os.Exit(1)
		}
	}

	if opts.ConsumerSecret == ""{
		if opts.ConsumerSecret = os.Getenv("BN_CONSUMER_SECRET"); opts.ConsumerSecret == ""{
			fmt.Println("Consumer secret must be set as option or environment value.")
			os.Exit(1)
		}
	}

	if opts.Database == ""{
		if opts.Database = os.Getenv("BN_DATABASE"); opts.Database == "" {
			fmt.Println("Database: mysql, postgre, sqlite")
			os.Exit(1)
		}
	}

	if opts.DBOptions == "" {
		if opts.DBOptions = os.Getenv("BN_DB_OPTIONS"); opts.DBOptions == ""{
			fmt.Println("Datamase options must be set.")
			os.Exit(1)
		}
	}

	return &opts
}