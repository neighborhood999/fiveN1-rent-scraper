GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean

all: test

.PHONY: test
test:
	$(GOTEST) -v

.PHONY: run
run:
	$(GORUN) ./_example/basic/main.go

.PHONY: race_detect
race_detect:
	$(GORUN) -race ./_example/basic/main.go

.PHONY: install
install:
	$(GOGET) github.com/PuerkitoBio/goquery
	$(GOGET) github.com/google/go-querystring/query
	$(GOGET) github.com/vinta/pangu
	$(GOGET) github.com/stretchr/testify/assert

.PHONY: clean
clean:
	$(GOCLEAN) -x -i ./...
