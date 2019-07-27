Minimalistic scraper that seek specified links on scpecified sites (config.yaml) and save it to DB, html and throw it to stdout. 

**What you need to start**
1) Specify sources and themes in config.yaml
2) Uncomment `av.InitDB()` in `/runner/main.go` for fisrt start initialization
3) start program from `runner` directory with `go run main.go`

You can specify flags:

`--method scrape` (start scraper for monitoring specified sites)

`--method find --theme X` (find specific news with X pattern)

`--method getall` (get all saved news)
