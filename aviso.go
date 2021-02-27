package aviso

import (
	"aviso/db"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type Fetcher interface {
	Fetch(url, theme string) (map[string]map[string]string, error)
}

type Aviso struct {
	DB     *db.DB
	wg     *sync.WaitGroup
	Config string
}

func New(config, host string, port int, user, password, dbname string) *Aviso {
	dbInstance := db.CreateDBConf(host, port, user, password, dbname)
	return &Aviso{DB: dbInstance, wg: &sync.WaitGroup{}, Config: config}
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
func (av *Aviso) Scrape(fetcher Fetcher, url string, theme string, ch chan map[string]map[string]string) {

	links, err := fetcher.Fetch(url, theme)
	if err != nil {
		log.Fatal(err)
	}

	ch <- links
	av.wg.Done()
}

func (av *Aviso) GetTargets() ([]string, []string, error) {
	data, err := ioutil.ReadFile(av.Config)
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

func (aviso *Aviso) Start(fetcher Fetcher) {

	// Reed url from yaml
	seedUrls, themes, err := aviso.GetTargets()
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
			go aviso.Scrape(fetcher, url, theme, chMap)
		}
		aviso.wg.Done()
	}

	go func() {
		for {
			select {
			case newNote := <-chMap:
				for source, mapLinks := range newNote {
					for k, v := range mapLinks {

						if aviso.isNew(v) {
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
			}
		}
	}()
	aviso.wg.Wait()
}

func (aviso *Aviso) EndlessScrape(fetcher Fetcher) {
	c := time.Tick(5 * time.Second)
	for {
		select {
		case <-c:
			go aviso.Start(fetcher)
		}
	}
}
