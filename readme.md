Minimalistic scraper that searches for specific links on the specific sites (config.yaml), save it to DB and throw it to stdout. 

**What you need to start**
1) specify sources and themes in config.yaml
2) build:
    
          go build -o aviso runner/main.go
3) init db:

          ./aviso --init true
          
4) start:

          ./aviso --task server
  

You can specify flags:

`--task server` (start scraper with UI)

`--task scrape` (start scraper for monitoring specified sites without UI)

`--task find --theme X` (find specific news with X pattern)

`--task getall` (get all saved news)

`--init <true/false>` (init table in Postgres or not)