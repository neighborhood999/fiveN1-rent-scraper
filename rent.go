package rent

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
	"github.com/vinta/pangu"
)

const (
	rootURL = "https://rent.591.com.tw/"
)

// HouseInfo is the representation rent house information.
type HouseInfo struct {
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
}

// Response is the representation http.Response.
type Response *http.Response

// HouseInfoCollection is the representation house information collection.
type HouseInfoCollection map[int][]*HouseInfo

// Document is the representation goquery.Document.
type Document struct {
	doc *goquery.Document
}

// FiveN1 is the representation page information.
type FiveN1 struct {
	records  int
	pages    int
	queryURL string
	rentList HouseInfoCollection
	wg       sync.WaitGroup
	rw       sync.RWMutex
	client   *http.Client
	cookie   *http.Cookie
}

// NewHouseInfo create a new `HouseInfo`.
func NewHouseInfo() *HouseInfo {
	return &HouseInfo{}
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
func NewFiveN1(url string) *FiveN1 {
	return &FiveN1{
		queryURL: url,
		rentList: make(map[int][]*HouseInfo),
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

func isBooleanNum(field string, n int) error {
	if !(n == 0 || n == 1) {
		return errors.New(field + " 請輸入 0 或是 1 的值！")
	}

	return nil
}

// GenerateURL is convert options to query parameters.
func GenerateURL(o Options) (string, error) {
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

// ConvertToJSON is convert rent house info collection to json.
func ConvertToJSON(list HouseInfoCollection) []byte {
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return b
}

// ExportJSON export json file.
func ExportJSON(b []byte) {
	f, err := os.Create("/tmp/rent.json")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(b)
	log.Println("✔️  Export Path: \x1b[42m/tmp/rent.json\x1b[0m")
}

func (d *Document) clone(res Response) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	d.doc = doc
}

// SetReqCookie set the region value.
func (f *FiveN1) SetReqCookie(region string) {
	f.cookie.Value = region
}

func (f *FiveN1) request(url string) Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.AddCookie(f.cookie)

	res, err := f.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (f *FiveN1) parseRentHouse(page int, doc *goquery.Document) {
	doc.Find("#content").Each(func(_ int, selector *goquery.Selection) {
		selector.Find(".listInfo.clearfix").Each(func(item int, listInfo *goquery.Selection) {
			houseInfo := NewHouseInfo()

			// Content Title
			title := listInfo.Find(".pull-left.infoContent > h3 > a[href]").Text()
			houseInfo.Title = stringReplacer(title)

			// Content URL
			var url string
			if href, ok := listInfo.Find(".pull-left.infoContent > h3 > a").Attr("href"); ok {
				url = stringReplacer(href)
			}
			houseInfo.URL = "https:" + url

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

				houseInfo.OptionType = trimTextSpace(splitDescription[0])
				houseInfo.Ping = trimTextSpace(splitDescription[1])
				houseInfo.RentType = trimTextSpace(splitDescription[2])
				houseInfo.Floor = trimTextSpace(splitDescription[3])

				// Rent House Address
				address := stringReplacer(infoContent.Find(".lightBox").Eq(1).Text())
				houseInfo.Address = address
			})

			// Rent Price
			listInfo.Find(".price").Each(func(_ int, price *goquery.Selection) {
				houseInfo.Price = stringReplacer(price.Text())
			})

			// New Rent House
			listInfo.Find(".newArticle").Each(func(_ int, n *goquery.Selection) {
				houseInfo.IsNew = true
			})

			// Add rent house into list
			f.rw.Lock()
			f.rentList[page+1] = append(f.rentList[page+1], houseInfo)
			f.rw.Unlock()
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

func (f *FiveN1) firstPage() {
	d := NewDocument()
	res := f.request(f.queryURL)
	d.clone(res)

	f.parseRecordsNum(d.doc) // Record pages number at first
	log.Println("------------------")
	log.Printf("| Total Page: \x1b[94;1m%d\x1b[0m |", f.pages)
	log.Println("------------------")
	f.parseRentHouse(0, d.doc)
}

func (f *FiveN1) worker(n int) {
	defer f.wg.Done()
	log.Printf("\x1b[30;1mStart worker at page number: %d\x1b[0m", n+1)

	d := NewDocument()
	r := strconv.Itoa(n * 30)
	res := f.request(f.queryURL + "&firstRow=" + r)
	d.clone(res)

	f.parseRentHouse(n, d.doc)
}

// Scrape will clone entire DOM for reuse.
// It will scrape first page at first, if specify page number > 1,
// it will start workers.
func (f *FiveN1) Scrape(page int) {
	f.firstPage()
	if page > f.pages {
		log.Fatal("\x1b[91;1mNo More Pages！\x1b[0m")
	}

	for i := 1; i < page; i++ {
		f.wg.Add(1)
		go f.worker(i)
	}

	f.wg.Wait()
}
