package main

import (
	"aviso"
	"aviso/db/sqlite"
	"aviso/fetcher"
	"aviso/rest"
	"flag"
	"fmt"
	"log"
)

func main() {
	var av *aviso.Aviso
	var method *string
	var theme *string

	method = flag.String("task", "scrape", "specify task")
	theme = flag.String("theme", "", "specify theme to find")
	flag.Parse()

	database, err := sqlite.New("aviso.db")
	if err != nil {
		panic(err)
	}
	av, err = aviso.New("./config/config.yaml", database)
	if err != nil {
		panic(err)
	}

	lf := fetcher.NewLinksFetcher("https")
	switch *method {
	case "headless":
		av.ScrapeLoop(lf)
	case "ui":
		go rest.Run(database, "0.0.0.0", "8000")
		av.ScrapeLoop(lf)

		// get from database
	case "find":
		if *theme == "" {
			log.Fatal(`
				Please specify theme:
  					aviso --task find --theme YOURTHEME
				`)
		}
		result, _ := av.FindSavedByTheme(*theme)
		fmt.Println("Finded:")
		for _, queryresult := range result {
			fmt.Printf("%s\nLink: %s\nSource: %s\nTime: %s\n", queryresult.Theme, queryresult.Link, queryresult.Site, queryresult.Time)
		}
	case "getall":
		result, err := database.QueryAll()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("From database:")
		for _, queryresult := range result {
			fmt.Printf("\n%s\nLink: %s\nSource: %s\nTime: %s\n", queryresult.Theme, queryresult.Link, queryresult.Site, queryresult.Time)
		}
	}
}
