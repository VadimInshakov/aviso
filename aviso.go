package aviso

import (
	"aviso/config"
	"aviso/domain"
	"aviso/fetcher"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

type repo interface {
	Init() error
	Insert(theme string, link string, site string, t time.Time) error
	QueryAll() ([]domain.WebObj, error)
	GetByTheme(theme string) (*domain.WebObj, error)
	FindByTheme(theme string) ([]domain.WebObj, error)
}

type Aviso struct {
	repo    repo
	g       *errgroup.Group
	targets *domain.Targets
}

func New(configpath string, repo repo) (*Aviso, error) {
	// Reed url from yaml
	targets, err := config.GetTargets(configpath)
	if err != nil {
		return nil, err
	}
	a := &Aviso{repo: repo, g: &errgroup.Group{}, targets: targets}
	if err := a.repo.Init(); err != nil {
		return nil, errors.Wrap(err, "repo init error")
	}
	return a, nil
}

func (av *Aviso) isNew(theme string) bool {
	result, err := av.repo.GetByTheme(theme)
	if err != nil {
		return true
	}

	if result.Theme == theme {
		return false
	}

	return true
}

func (av *Aviso) FindSavedByTheme(theme string) ([]domain.WebObj, error) {
	result, err := av.repo.FindByTheme(theme)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Extract all links from a given webpage
func (av *Aviso) scrape(f fetcher.Fetcher, url string, theme string, ch chan map[string]map[string]string, selectors, excludeSelectors []string) error {
	links, err := f.Fetch(&fetcher.FetchQuery{Url: url, Theme: theme, Selectors: selectors, ExcludeSelectors: excludeSelectors})
	if err != nil {
		return err
	}

	ch <- links
	return nil
}

func (av *Aviso) scrapeThemes(fetcher fetcher.Fetcher) {
	chMap := make(chan map[string]map[string]string, 5500)

	// scrape sites concurrently
	for _, url := range av.targets.Urls {
		if len(av.targets.Themes) == 0 {
			av.g.Go(func() error {
				return av.scrape(fetcher, url.Link, "", chMap, url.Selectors, url.ExcludeSelectors)
			})
			continue
		}
		for _, theme := range av.targets.Themes {
			av.g.Go(func() error {
				return av.scrape(fetcher, url.Link, theme, chMap, url.Selectors, url.ExcludeSelectors)
			})
		}
	}

	av.g.Go(func() error {
		for newNote := range chMap {
			for source, mapLinks := range newNote {
				for k, v := range mapLinks {
					if av.isNew(v) {
						// write to stdout
						fmt.Printf("\n- %s\n%s", v, k)

						//write to database
						if err := av.repo.Insert(v, k, source, time.Now()); err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})

	if err := av.g.Wait(); err != nil {
		log.Fatal(err)
	}
	close(chMap)
}

func (av *Aviso) ScrapeLoop(fetcher fetcher.Fetcher) {
	c := time.Tick(5 * time.Second)
	for {
		select {
		case <-c:
			go av.scrapeThemes(fetcher)
		}
	}
}
