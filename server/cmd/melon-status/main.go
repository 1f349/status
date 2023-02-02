package main

import (
	"code.mrmelon54.com/melon/status/server/notifier"
	"code.mrmelon54.com/melon/status/server/ping"
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/web"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
	"xorm.io/xorm"
)

func main() {
	log.Println("[Main] Starting up Melon Status")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Main] Loading database")
	dbEnv := os.Getenv("DB")
	var engine *xorm.Engine
	if strings.HasPrefix(dbEnv, "sqlite:") {
		engine, err = xorm.NewEngine("sqlite3", strings.TrimPrefix(dbEnv, "sqlite:"))
	} else if strings.HasPrefix(dbEnv, "mysql:") {
		engine, err = xorm.NewEngine("mysql", strings.TrimPrefix(dbEnv, "mysql:"))
	} else {
		log.Fatalln("[Main] Only mysql and sqlite are supported")
	}
	if err != nil {
		log.Fatalf("Unable to load database (\"%s\")\n", dbEnv)
	}
	check(engine.Sync(&structure.Service{}, &structure.Group{}, &structure.Hit{}, &structure.Failure{}))

	p := ping.New(engine, notifier.Init())
	p.Reload()
	p.Run()
	m := ping.NewMaintenance(engine)
	m.Run()

	s := web.New(engine, p).SetupWeb()
	err = s.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Printf("[Http] The HTTP server shutdown successfully\n")
		} else {
			log.Printf("[Http] Error trying to host the HTTP server: %s\n", err.Error())
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
