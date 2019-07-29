package mock

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type MockFetcher struct {
	Protocol string
}

func (f *MockFetcher) Fetch(url string, theme string) (map[string]map[string]string, error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(Mockdocument))

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
				href = url + href[1:]
			}
			submap[href] = s.Text()
		}
	})
	links["https://mock.mock/"] = submap

	return links, err
}
