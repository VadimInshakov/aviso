Minimalistic scraper that searches for specific links on the scpecific sites (config.yaml), save it to DB and throw it to stdout. 

**What you need to start**
1) You need PostgreSQL installed with database named `aviso`
1) Specify sources and themes in config.yaml
2) start program from `runner` directory with `go run main.go`

You can specify flags:

`--task scrape` (start scraper for monitoring specified sites)

`--task find --theme X` (find specific news with X pattern)

`--task getall` (get all saved news)

`--init <true/false>` (init table in Postgres or not)