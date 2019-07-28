Minimalistic scraper that seek specified links on scpecified sites (config.yaml) and save it to DB, html and throw it to stdout. 

**What you need to start**
1) Specify sources and themes in config.yaml
2) Uncomment `av.InitDB()` in `/runner/main.go` for fisrt start initialization
3) start program from `runner` directory with `go run main.go`

You can specify flags:

`--task scrape` (start scraper for monitoring specified sites)

`--task find --theme X` (find specific news with X pattern)

`--task getall` (get all saved news)
