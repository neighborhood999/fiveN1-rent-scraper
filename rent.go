package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
	"github.com/vinta/pangu"
)

var (
	rootURL = "https://rent.591.com.tw/"
)

// RentHouseInfo is the representation rent house information.
type RentHouseInfo struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Address    string `json:"address"`
	RentType   string `json:"rentType"`
	OptionType string `json:"optionType"`
	Ping       string `json:"ping"` // a.k.a 坪數
	Floor      string `json:"floor"`
	Price      string `json:"price"`
	IsNew      bool   `json:"isNew"`
}

// Options is the representation query arugment.
type Options struct {
	Region    int    `url:"region"`    // 地區 - 預設：`1`
	Section   string `url:"section"`   // 鄉鎮 - 可選擇多個區域，例如：`section=7,4`
	Kind      int    `url:"kind"`      // 租屋類型 - `0`：不限、`1`：整層住家、`2`：獨立套房、`3`：分租套房、`4`：雅房、`5`：車位，`6`：其他
	RentPrice string `url:"rentprice"` // 租金 - `2`：5k - 10k、`3`：10k - 20k、`4`: 20k - 30k；或者可以輸入價格範圍，例如：`0,10000`
	Area      string `url:"area"`      // 坪數格式 - `10,20`（10 到 20 坪）
	Order     string `url:"order"`     // 貼文時間 - 預設：`posttime`
	OrderType string `url:"ordertype"` // 排序方式 - `desc` 或 `asc`
	Sex       int    `url:"sex"`       // 性別 - `0`：不限、`1`：男性、`2`：女性
	HasImg    int    `url:"hasimg"`    // 過濾是否有「房屋照片」 - `0`：否、`1`：是
	NotCover  int    `url:"not_cover"` // 過濾是否為「頂樓加蓋」 - `0`：否、`1`：是
	Role      int    `url:"role"`      // 過濾是否為「屋主刊登」 - `0`：否、`1`：是
	FirstRow  int    `url:"firstRow"`  // 頁數設定
}

// Response is the representation http.Response.
type Response *http.Response

// RentHouseInfoCollection is the representation house information collection.
type RentHouseInfoCollection map[int][]*RentHouseInfo

// Document is the representation goquery.Document.
type Document struct {
	doc *goquery.Document
}

// FiveN1 is the representation page information.
type FiveN1 struct {
	records  int
	pages    int
	queryURL string
	rentList RentHouseInfoCollection
	client   *http.Client
	cookie   *http.Cookie
}

// NewRentHouseInfo create a new `RentHouseInfo`.
func NewRentHouseInfo() *RentHouseInfo {
	return &RentHouseInfo{}
}

// NewOptions create a `Options` with default value.
func NewOptions() *Options {
	return &Options{
		Kind:      2,
		Region:    1,
		Section:   "0",
		RentPrice: "2",
		HasImg:    0,
		NotCover:  0,
		Role:      0,
		Order:     "posttime",
		OrderType: "desc",
	}
}

// NewFiveN1 create a `FiveN1` with default value.
func NewFiveN1() *FiveN1 {
	return &FiveN1{
		rentList: make(map[int][]*RentHouseInfo),
		client:   &http.Client{},
		cookie: &http.Cookie{
			Name:  "urlJumpIp",
			Value: "1",
		},
	}
}

// NewDocument create a `Document` with default value.
func NewDocument() *Document {
	return &Document{
		doc: &goquery.Document{},
	}
}

// SetReqCookie set the region value.
func (f *FiveN1) SetReqCookie(region string) {
	f.cookie.Value = region
}

func (d *Document) clone(res Response) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	d.doc = doc
}

func (f *FiveN1) request(url string) Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	f.queryURL = url
	req.AddCookie(f.cookie)

	res, err := f.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func isBooleanNum(field string, n int) error {
	if !(n == 0 || n == 1) {
		return errors.New(field + " 請輸入 0 或是 1 的值！")
	}
	return nil
}

func generateURL(o Options) (string, error) {
	if err := isBooleanNum("`HasImg`", o.HasImg); err != nil {
		return "", err
	}

	if err := isBooleanNum("`NotCover`", o.NotCover); err != nil {
		return "", err
	}

	if err := isBooleanNum("`Role`", o.Role); err != nil {
		return "", err
	}

	v, _ := query.Values(o)

	return rootURL + "?" + v.Encode(), nil
}

func stringReplacer(text string) string {
	replacer := strings.NewReplacer("\n", "", " ", "")

	return pangu.SpacingText(replacer.Replace(text))
}

func trimTextSpace(s string) string {
	return strings.Fields(s)[0]
}

func fillDescription(s []string) []string {
	s = append(s, s[2])
	s[2] = "沒有格局說明"

	return s
}

func convertToJSON(list RentHouseInfoCollection) []byte {
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func exportJSON(b []byte) {
	f, err := os.Create("/tmp/rent.json")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(b)
	fmt.Println("🎈 Done！Check out `/tmp/rent.json`.")
}

func (f *FiveN1) parseRentHouse(doc *goquery.Document) {
	// "https://rent.591.com.tw/?kind=2&region=1&rentprice=2&hasimg=1&not_cover=1&role=1&order=posttime&orderType=desc"
	doc.Find("#content").Each(func(_ int, selector *goquery.Selection) {
		selector.Find(".listInfo.clearfix").Each(func(item int, listInfo *goquery.Selection) {
			rentHouse := NewRentHouseInfo()

			// Content Title
			title := listInfo.Find(".pull-left.infoContent > h3 > a[href]").Text()
			rentHouse.Title = stringReplacer(title)

			// Content URL
			var url string
			if href, ok := listInfo.Find(".pull-left.infoContent > h3 > a").Attr("href"); ok {
				url = stringReplacer(href)
			}
			rentHouse.URL = "https:" + url

			listInfo.Find(".pull-left.infoContent").Each(func(_ int, infoContent *goquery.Selection) {
				// Rent House Description.
				description := stringReplacer(infoContent.Find(".lightBox").First().Text())

				splitDescription := strings.Split(description, "|")

				// Exchange
				if len(splitDescription) == 4 {
					tmp := splitDescription[2] // 坪數
					splitDescription[2] = splitDescription[1]
					splitDescription[1] = tmp
				}

				if len(splitDescription) < 4 {
					splitDescription = fillDescription(splitDescription)
				}

				rentHouse.OptionType = trimTextSpace(splitDescription[0])
				rentHouse.Ping = trimTextSpace(splitDescription[1])
				rentHouse.RentType = trimTextSpace(splitDescription[2])
				rentHouse.Floor = trimTextSpace(splitDescription[3])

				// Rent House Address
				address := stringReplacer(infoContent.Find(".lightBox").Eq(1).Text())
				rentHouse.Address = address
			})

			// Rent Price
			listInfo.Find(".price").Each(func(_ int, price *goquery.Selection) {
				rentHouse.Price = stringReplacer(price.Text())
			})

			// New Rent House
			listInfo.Find(".newArticle").Each(func(_ int, n *goquery.Selection) {
				rentHouse.IsNew = true
			})

			// Add rent house into list
			f.rentList[1] = append(f.rentList[1], rentHouse)
		})
	})
}

func (f *FiveN1) parseRecordsNum(doc *goquery.Document) {
	doc.Find(".page-limit > .pageBar > .TotalRecord > .R").Each(func(_ int, selector *goquery.Selection) {
		totalRecord, _ := strconv.Atoi(stringReplacer(selector.Text()))
		f.records = totalRecord
		f.pages = (totalRecord / 30) + 1
	})
}

// Scrape create a new Document and clone entire DOM for reuse,
// and parse rent houses information.
func (f *FiveN1) Scrape(url string) {
	d := NewDocument() // Initial document
	res := f.request(url)
	d.clone(res)

	f.parseRecordsNum(d.doc)
	f.parseRentHouse(d.doc)
	// exportJSON(convertToJSON(f.rentList))
}

func main() {
	o := NewOptions()
	url, err := generateURL(*o)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(url)
	f := NewFiveN1()
	f.Scrape(url)
	// fmt.Println(f)
}
