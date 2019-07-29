package main

import (
	"aviso"
	"aviso/fetcher"
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
	init := flag.Bool("init", false, "init table (true) or not (false)")
	flag.Parse()

	av = aviso.New("../config.yaml", "localhost", 5433, "postgres", "dbsecret", "aviso")
	av.ConnectDB()
	if *init {
		av.InitDB()
	}
}

func main() {
	switch *method {

	case "scrape":
		var _ aviso.Fetcher = (*fetcher.LinksFetcher)(nil)
		var myfetcher *fetcher.LinksFetcher = &fetcher.LinksFetcher{Protocol: "https"}

		c := time.Tick(5 * time.Second)
		for {
			select {
			case <-c:
				go av.Start(myfetcher)
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
