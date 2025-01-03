package fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

type LinksFetcher struct {
	protocol string
	client   *http.Client
}

func NewLinksFetcher(protocol string) *LinksFetcher {
	lf := &LinksFetcher{protocol: protocol}
	lf.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          6000,
			MaxIdleConnsPerHost:   2,
			MaxConnsPerHost:       2,
			IdleConnTimeout:       1 * time.Minute,
			ResponseHeaderTimeout: 10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	return lf
}

type FetchQuery struct {
	Url              string
	Theme            string
	Selectors        []string
	ExcludeSelectors []string
}

func (f *LinksFetcher) Fetch(q *FetchQuery) (map[string]map[string]string, error) {
	resp, err := f.client.Get(q.Url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links = make(map[string]map[string]string, 100)
	var submap = make(map[string]string)

	for _, s := range q.Selectors {
		doc.Find(s).Find("a").Each(func(i int, s *goquery.Selection) {
			if q.Theme != "" {
				if !strings.Contains(s.Text(), q.Theme) {
					return
				}
			}
			if class, ok := s.Attr("class"); ok {
				for _, e := range q.ExcludeSelectors {
					if class == e {
						return
					}
				}
			}

			href, _ := s.Attr("href")

			if strings.Contains(href, "comment") {
				return
			}

			// some links are relational (like /addr, not http://<host>:<port>/addr),
			// so we need concatenate missing part in this case)
			if !strings.Contains(href, fmt.Sprintf("%s:", f.protocol)) {
				q.Url = strings.ReplaceAll(q.Url, "/?hl=ru", "")
				if len(href) > 0 {
					href = q.Url + href
				}
			}

			text := strings.TrimSpace(s.Text())
			if len(text) < 45 {
				return
			}

			submap[href] = text
		})
	}

	links[q.Url] = submap

	return links, err
}
