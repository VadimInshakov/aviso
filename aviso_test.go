package aviso

import (
	"aviso/mock"
	"log"
	"reflect"
	"testing"
)

var av *Aviso

func TestStart(t *testing.T) {

	var myfetcher Fetcher = &mock.MockFetcher{Protocol: "https"}

	av = New("./config.yaml", "localhost", 5433, "postgres", "dbsecret", "aviso")
	av.ConnectDB()
	// av.InitDB() // uncomment if news table in Postgres doesn't exist

	// Reed url from yaml
	seedUrls, themes, err := av.GetTargets()
	if err != nil {
		log.Fatal(err)
	}

	// Channels
	chMap := make(chan map[string]map[string]string, 100)
	defer close(chMap)

	// Kick off the Scrape process (concurrently)
	for _, theme := range themes {
		av.wg.Add(1)
		for _, url := range seedUrls {
			av.wg.Add(1)
			go av.Scrape(myfetcher, url, theme, chMap)
		}
		av.wg.Done()
	}

	// Subscribe to both channels
	go func() {
		for {
			select {
			case newNote := <-chMap:
				for _, mapLinks := range newNote {
					for k, v := range mapLinks {
						// write to stdout
						got := reflect.TypeOf(v).Name()
						want := "string"
						if got != want {
							t.Errorf("title type mismatch:\ngot:%s\nwant:%s", got, want)
						}
						got = reflect.TypeOf(k).Name()
						want = "string"
						if got != want {
							t.Errorf("link type mismatch:\ngot:%s\nwant:%s", got, want)
						}
					}
				}
			}
		}
	}()
	av.wg.Wait()

}
