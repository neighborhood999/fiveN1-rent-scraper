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
	Preview    string `json:"preview"`
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

// Options is the representation query argument.
type Options struct {
	Region      int    `url:"region"`                // 地區 - 預設：`1`
	Section     string `url:"section,omitempty"`     // 鄉鎮 - 可選擇多個區域，例如：`section=7,4`
	Kind        int    `url:"kind"`                  // 租屋類型 - `0`：不限、`1`：整層住家、`2`：獨立套房、`3`：分租套房、`4`：雅房、`8`：車位，`24`：其他
	RentPrice   string `url:"rentprice,omitempty"`   // 租金 - `2`：5k - 10k、`3`：10k - 20k、`4`: 20k - 30k；或者可以輸入價格範圍，例如：`0,10000`
	Area        string `url:"area,omitempty"`        // 坪數格式 - `10,20`（10 到 20 坪）
	Order       string `url:"order"`                 // 貼文時間 - 預設使用刊登時間：`posttime`，或是使用價格排序：`money`
	OrderType   string `url:"orderType"`             // 排序方式 - `desc` 或 `asc`
	Sex         int    `url:"sex,omitempty"`         // 性別 - `0`：不限、`1`：男性、`2`：女性
	HasImg      string `url:"hasimg,omitempty"`      // 過濾是否有「房屋照片」 - ``：空值（不限）、`1`：是
	NotCover    string `url:"not_cover,omitempty"`   // 過濾是否為「頂樓加蓋」 - ``：空值（不限）、`1`：是
	Role        string `url:"role,omitempty"`        // 過濾是否為「屋主刊登」 - ``：空值（不限）、`1`：是
	Shape       string `url:"shape,omitempty"`       // 房屋類型 - `1`：公寓、`2`：電梯大樓、`3`：透天厝、`4`：別墅
	Pattern     string `url:"pattern,omitempty"`     // 格局單選 - `0`：不限、`1`：一房、`2``：兩房、`3`：三房、`4`：四房、`5`：五房以上
	PatternMore string `url:"patternMore,omitempty"` // 格局多選 - 參考「格局單選」，可以選多種格局，例如：`1,2,3,4,5`
	Floor       string `url:"floor,omitempty"`       // 樓層 - `0,0`：不限、`0,1`：一樓、`2,6`：二樓到六樓、`6,12`：六樓到十二樓、`12,`：十二樓以上
	Option      string `url:"option,omitempty"`      // 提供設備 - `tv`：電視、`cold`：冷氣、`icebox`：冰箱、`hotwater`：熱水器、`naturalgas`：天然瓦斯、`four`：第四台、`broadband`：網路、`washer`：洗衣機、`bed`：床、`wardrobe`：衣櫃、`sofa`：沙發。可選擇多個設備，例如：option=tv,cold
	Other       string `url:"other,omitempty"`       // 其他條件 - `cartplace`：有車位、`lift`：有電梯、`balcony_1`：有陽台、`cook`：可開伙、`pet`：可養寵物、`tragoods`：近捷運、`lease`：可短期租賃。可選擇多個條件，例如：other=cartplace,cook
	FirstRow    int    `url:"firstRow"`
}

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
	RentList HouseInfoCollection
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
		Order:     "posttime",
		OrderType: "desc",
		FirstRow:  0,
	}
}

// NewFiveN1 create a `FiveN1` with default value.
func NewFiveN1(url string) *FiveN1 {
	return &FiveN1{
		queryURL: url,
		RentList: make(map[int][]*HouseInfo),
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

// GenerateURL is convert options to query parameters.
func GenerateURL(o *Options) (string, error) {
	v, error := query.Values(o)

	return rootURL + "?" + v.Encode(), error
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
	defer f.Close()

	f.Write(b)
	log.Println("✔️  Export Path: \x1b[42m/tmp/rent.json\x1b[0m")
}

func (d *Document) clone(res *http.Response) {
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

func (f *FiveN1) request(url string) *http.Response {
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

			if crop, ok := listInfo.Find(".pull-left.imageBox > img").Attr("data-original"); ok {
				preview := strings.Replace(crop, "210x158.crop.jpg", "765x517.water3.jpg", 1)
				houseInfo.Preview = preview
			}

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
			f.RentList[page+1] = append(f.RentList[page+1], houseInfo)
			f.rw.Unlock()
		})
	})
}

func (f *FiveN1) parseRecordsNum(doc *goquery.Document) {
	doc.Find(".pull-left.hasData > i").Each(func(_ int, selector *goquery.Selection) {
		recordString := stringReplacer(selector.Text())
		replaceComman := strings.Replace(recordString, ",", "", -1)
		totalRecord, _ := strconv.Atoi(replaceComman)
		f.records = totalRecord
		f.pages = (totalRecord / 30) + 1
	})
}

func (f *FiveN1) firstPage() {
	d := NewDocument()
	res := f.request(f.queryURL)
	d.clone(res)

	f.parseRecordsNum(d.doc) // Record pages number at first
	log.Println("---------------------")
	log.Printf("| Query URL:    \x1b[94;1m%v\x1b[0m |", f.queryURL)
	log.Printf("| Total Record: \x1b[94;1m%d\x1b[0m |", f.records)
	log.Printf("| Total Page:   \x1b[94;1m%d\x1b[0m  |", f.pages)
	log.Println("---------------------")
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
func (f *FiveN1) Scrape(page int) error {
	f.firstPage()
	if page > f.pages {
		return errors.New("\x1b[91;1mNo More Pages！\x1b[0m")
	}

	for i := 1; i < page; i++ {
		f.wg.Add(1)
		go f.worker(i)
	}

	f.wg.Wait()

	return nil
}
