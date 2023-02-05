Minimalistic scraper that searches for specific links on the specific sites (config.yaml), save it to DB and throw it to stdout. 

**What you need to start**
1) specify sources and themes in config.yaml
2) build:
    
          go build -o aviso cmd/main.go
          
3) start:

          ./aviso --task server
  

You can specify flags:

`--task ui` (start scraper with UI)

`--task headless` (start scraper for headless monitoring specified sites)

`--task find --theme X` (find specific news with X pattern)

`--task getall` (get all saved news)