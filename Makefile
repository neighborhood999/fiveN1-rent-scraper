GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean

all: test

test:
	$(GOTEST) -v

run:
	$(GORUN) ./_example/basic/main.go

race_detect:
	$(GORUN) -race ./_example/basic/main.go

install:
	$(GOGET) github.com/PuerkitoBio/goquery
	$(GOGET) github.com/google/go-querystring/query
	$(GOGET) github.com/vinta/pangu
	$(GOGET) github.com/stretchr/testify/assert

clean:
	$(GOCLEAN) -x -i ./...
