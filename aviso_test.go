package aviso

import (
	"aviso/mock"
	"log"
	"testing"
)

var av *Aviso

func TestStart(t *testing.T) {
	var myfetcher Fetcher = &mock.MockFetcher{Protocol: "https"}
	var err error
	av, err = New("./config.yaml", "aviso_test.db")

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

	go func() {
		for {
			select {
			case newNote := <-chMap:
				for _, mapLinks := range newNote {
					for k, v := range mapLinks {
						// write to stdout
						got := v
						want := "Путин акции простесты задержали"
						if got != want {
							t.Errorf("title mismatch:\ngot:%s\nwant:%s", got, want)
						}
						got = k
						want = "https://mock.mock/mock"
						if got != want {
							t.Errorf("link mismatch:\ngot:%s\nwant:%s", got, want)
						}
					}
				}
			}
		}
	}()
	av.wg.Wait()

}
