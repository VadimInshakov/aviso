package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"gopkg.in/yaml.v2"
)

var urlArray []string
var firstStart bool = true

// Extract all http** links from a given webpage
func crawl(url string, theme string, ch chan map[string]string, chFinished chan bool, somethingNew chan bool) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
	}
	var links map[string]string = make(map[string]string)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == theme {
			href, _ := s.Attr("href")
			links[href] = s.Text()
		}
	})

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	ch <- links

}

func getTargets() ([]string, []string) {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	type Targets struct {
		Urls   []string `yaml:"urls"`
		Themes []string `yaml:"themes"`
	}
	var t Targets
	yaml.Unmarshal([]byte(data), &t)
	return t.Urls, t.Themes
}

func start(somethingNew chan bool) {

	// Reed url from yaml
	seedUrls, themes := getTargets()

	// Channels
	chMap := make(chan map[string]string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, theme := range themes {
		for _, url := range seedUrls {
			go crawl(url, theme, chMap, chFinished, somethingNew)
		}
	}

	isNew := false
	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case mapLinks := <-chMap:
			if isNew || firstStart {
				var htmlData string = "<html><body>"
				for k, v := range mapLinks {
					fmt.Println(k, v)
					htmlPart := fmt.Sprintf("<li> <a href=%s /> %s </li>", k, v)
					htmlData = fmt.Sprintf("%s%s", htmlData, htmlPart)
				}
				htmlData = fmt.Sprintf("%s%s", htmlData, "<body/><html/>")
				ioutil.WriteFile("./index.html", []byte(htmlData), 0644)
			}
		case <-chFinished:
			c++
		case <-somethingNew:
			isNew = true
		}
	}

	firstStart = false
	close(chMap)
}

func main() {
	c := time.Tick(2 * time.Second)
	somethingNew := make(chan bool)
	for {
		select {
		case <-c:
			go start(somethingNew)
		}
	}
}
