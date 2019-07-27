package main

import (
	"aviso"
	"flag"
	"fmt"
	"log"
	"time"
)

var av *aviso.Aviso
var method *string
var theme *string

func init() {
	method = flag.String("task", "scrape", "specify task")
	theme = flag.String("theme", "", "specify theme to find")
	flag.Parse()

	av = aviso.New("localhost", 5433, "postgres", "rlabssupersecret", "aviso")
	av.ConnectDB()
	//av.InitDB() // uncomment if news table in Postgres doesn't exist
}

func main() {
	switch *method {

	case "scrape":
		c := time.Tick(5 * time.Second)
		for {
			select {
			case <-c:
				go av.Start()
			}
		}
	case "find":
		if *theme == "" {
			log.Fatal(`
				Please specify theme:
  					aviso --task find --theme YOURTHEME
				`)
		}
		result, _ := av.FindByTheme(*theme)
		fmt.Println("Finded:")
		for _, queryresult := range result {
			fmt.Printf("%d. %s\nLink: %s\nSource: %s\nTime: %s\n", queryresult.Id, queryresult.Theme, queryresult.Link, queryresult.Site, queryresult.Time)
		}
	case "getall":
		//query from DB
		result, err := av.DB.QueryAll()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("From database:")
		for _, queryresult := range result {
			fmt.Printf("\n%d. %s\nLink: %s\nSource: %s\nTime: %s\n", queryresult.Id, queryresult.Theme, queryresult.Link, queryresult.Site, queryresult.Time)
		}
	}
}
