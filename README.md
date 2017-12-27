# fiveN1 Rent Scraper

![logo](./logo/fiveN1-rent-scraper-logo.png)
[![godoc](https://camo.githubusercontent.com/5771fd8cd24b1f8c34b82f152587dbce2294d9e1/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f746a2f6e6f64652d7072756e653f7374617475732e737667)](https://godoc.org/github.com/neighborhood999/fiveN1-rent-scraper)
[![Go Report Card](https://goreportcard.com/badge/github.com/neighborhood999/fiveN1-rent-scraper)](https://goreportcard.com/report/github.com/neighborhood999/fiveN1-rent-scraper)

> a.k.a 591 rent scraper.

Easy scraping [591](https://rent.591.com.tw/) rent information.

## Requirement

```sh
$ go get -u github.com/vinta/pangu
$ go get -u github.com/PuerkitoBio/goquery
$ go get -u github.com/google/go-querystring/query
```

## Installation

```sh
$ go get github.com/neighborhood999/fiveN1-rent-scraper
```

## Usage

Create default options and you can generate url:

```go
o := rent.NewOptions()
url, _ := rent.GenerateURL(o)
log.Println(url)
```

or you can setting options for your requirement.

## Example

```go
import (
	"log"

	"github.com/neighborhood999/fiveN1-rent-scraper"
)

func main() {
	o := rent.NewOptions()
	url, err := rent.GenerateURL(o)
	if err != nil {
		log.Fatalf("\x1b[91;1m%s\x1b[0m", err)
	}

	f := rent.NewFiveN1(url)
	f.Scrape(1)

	json := rent.ConvertToJSON(f.RentList)
	log.Println(string(json))
}
```

And output json:

```json
{
  "1": [
    {
      "title": "臨近北醫，精緻單人套房 (5 樓 B 室)",
      "url": "https://rent.591.com.tw/rent-detail-5926570.html",
      "address": "信義區 - 吳興街 336 巷",
      "rentType": "沒有格局說明",
      "optionType": "獨立套房",
      "ping": "5.5",
      "floor": "樓層：5/6",
      "price": "8,600 元 / 月",
      "isNew": true
    }
  ]
}
```

Index number is the representation **page number**, every index contain **30** items. 🏠

## Multiple Page Scrape

Default will scrape first page, if you want to scrape more pages, setting page amount in `Scrape` method:

```go
f := rent.NewFiveN1(url)
f.Scrape(5) // scrape page 1-5
```

If page amount `> 1`, it will start goroutine automatically for scrape to correspond to page number.

## LICENSE

MIT © [Peng Jie](https://github.com/neighborhood999)
