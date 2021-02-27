package fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type LinksFetcher struct {
	Protocol string
}

func (f *LinksFetcher) Fetch(url string, theme string) (map[string]map[string]string, error) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	var links = make(map[string]map[string]string, 100)
	var submap = make(map[string]string)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), theme) {
			href, _ := s.Attr("href")
			// some links are relational (like /addr, not http://<host>:<port>/addr),
			// so we need concatenate missing part in this case)
			if !strings.Contains(href, fmt.Sprintf("%s:", f.Protocol)) {
				url = strings.ReplaceAll(url, "/?hl=ru", "")
				if len(href) > 0 {
					href = url + href[1:]
				}
			}
			submap[href] = s.Text()
		}
	})
	links[url] = submap

	return links, err
}
