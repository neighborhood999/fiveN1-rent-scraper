GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

run:
	$(GORUN) rent.go

raceDetect:
	$(GORUN) -race rent.go

deps:
	$(GOGET) github.com/PuerkitoBio/goquery
	$(GOGET) github.com/google/go-querystring/query
	$(GOGET) github.com/vinta/pangu
