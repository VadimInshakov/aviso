package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

var urlArray []string
var firstStart bool = true

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool, somethingNew chan bool) {

	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if isAnchor {

				// Extract the href value, if there is one
				ok, url := getHref(t)
				hasProto := strings.Index(url, "http") == 0

				if !ok {
					continue
				}

				if hasProto {
					urlArray = append(urlArray, url)

					// Make sure the url begines in http**
					isNew := func() bool {

						if len(urlArray) > 0 {
							for _, a := range urlArray {
								if a == url {
									return false
								}
							}
						}
						return true
					}()

					if (hasProto && isNew) || (hasProto && !isNew && firstStart) {
						somethingNew <- true
						ch <- url
					}
				}
			}

		}
	}
}

func getTargets() []string {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	type Targets struct {
		Urls []string `yaml:"urls"`
	}
	var t Targets
	yaml.Unmarshal([]byte(data), &t)
	return t.Urls
}

func start(somethingNew chan bool) {
	foundUrls := make(map[string]bool)
	// Reed url from yaml
	seedUrls := getTargets()

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished, somethingNew)
	}

	isNew := false
	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		case <-somethingNew:
			isNew = true
		}
	}

	if isNew || firstStart {
		var htmlData string = "<html><body>"
		fmt.Println("\nFound", len(foundUrls), "unique urls: \n")
		for url, _ := range foundUrls {
			fmt.Println(" - " + url)
			htmlPart := fmt.Sprintf("<li> <a href=%s /> %s </li>", url, url)
			htmlData = fmt.Sprintf("%s%s", htmlData, htmlPart)
		}
		htmlData = fmt.Sprintf("%s%s", htmlData, "<body/><html/>")
		ioutil.WriteFile("./index.html", []byte(htmlData), 0644)
	}

	firstStart = false
	close(chUrls)
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
