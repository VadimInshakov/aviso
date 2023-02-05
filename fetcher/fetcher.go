package fetcher

type Fetcher interface {
	Fetch(query *FetchQuery) (map[string]map[string]string, error)
}
