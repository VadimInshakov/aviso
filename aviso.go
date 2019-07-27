package aviso

import (
	"aviso/db"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

type Aviso struct {
	DB *db.DB
	wg *sync.WaitGroup
}

func New(host string, port int, user, password, dbname string) *Aviso {
	dbInstance := db.CreateDBConf(host, port, user, password, dbname)
	return &Aviso{DB: dbInstance, wg: &sync.WaitGroup{}}
}

func (aviso *Aviso) InitDB() {
	// create table
	err := aviso.DB.Init()
	if err != nil {
		log.Fatalln("db connection failed:", err.Error())
	}
	log.Println("DB initialized")
}

func (aviso *Aviso) ConnectDB() {
	//create db config and init connection
	err := aviso.DB.Connect()
	if err != nil {
		log.Fatalln("DB connection failed:", err.Error())
	}
	log.Println("Connected to Postgres successfully")
}

// Extract all http** links from a given webpage
func (av *Aviso) Scrape(url string, theme string, ch chan map[string]map[string]string) {

	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
	}
	var links = make(map[string]map[string]string, 100)
	var submap = make(map[string]string)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), theme) {
			href, _ := s.Attr("href")
			if !strings.Contains(href, "https:") {
				url = strings.ReplaceAll(url, "/?hl=ru", "")
				href = url + href[1:]
			}
			submap[href] = s.Text()
		}
	})
	links[url] = submap

	ch <- links
	av.wg.Done()
}

func GetTargets() ([]string, []string, error) {
	data, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		log.Println("Reading file error: ")
		return nil, nil, err
	}
	type Targets struct {
		Urls   []string `yaml:"urls"`
		Themes []string `yaml:"themes"`
	}
	var t Targets
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Println("Unmarshalling error: ")
		return nil, nil, err
	}

	return t.Urls, t.Themes, nil
}

func (aviso *Aviso) isNew(theme string) bool {
	result, err := aviso.DB.GetByTheme(theme)
	if err != nil {
		log.Println(err)
		return true
	}

	if result.Theme == theme {
		return false
	}

	return true
}

func (aviso *Aviso) FindByTheme(theme string) ([]db.QueryResult, error) {
	result, err := aviso.DB.FindByTheme(theme)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (aviso *Aviso) Start() {

	// Reed url from yaml
	seedUrls, themes, err := GetTargets()
	if err != nil {
		log.Fatal(err)
	}

	// Channels
	chMap := make(chan map[string]map[string]string, 100)
	defer close(chMap)

	// Kick off the Scrape process (concurrently)
	for _, theme := range themes {
		aviso.wg.Add(1)
		for _, url := range seedUrls {
			aviso.wg.Add(1)
			go aviso.Scrape(url, theme, chMap)
		}
		aviso.wg.Done()
	}

	// Subscribe to both channels
	go func() {
		for {
			select {
			case newNote := <-chMap:

				//write to file
				var htmlData = "<html><body>"
				for source, mapLinks := range newNote {
					for k, v := range mapLinks {

						if aviso.isNew(v) {
							htmlPart := fmt.Sprintf("<li> <a href=%s /> %s </li>", k, v)
							htmlData = fmt.Sprintf("%s%s", htmlData, htmlPart)

							// write to stdout
							fmt.Printf("\n- %s\n%s", v, k)

							//write to database
							err = aviso.DB.Insert(v, k, source, time.Now())
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}
				htmlData = fmt.Sprintf("%s%s", htmlData, "<body/><html/>")
				ioutil.WriteFile("./index.html", []byte(htmlData), 0644)

			}
		}
	}()
	aviso.wg.Wait()
}
