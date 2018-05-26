# fiveN1 Rent Scraper

![logo](./logo/fiveN1-rent-scraper-logo.png)
[![godoc](https://camo.githubusercontent.com/5771fd8cd24b1f8c34b82f152587dbce2294d9e1/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f746a2f6e6f64652d7072756e653f7374617475732e737667)](https://godoc.org/github.com/neighborhood999/fiveN1-rent-scraper)
[![Build Status](https://img.shields.io/travis/neighborhood999/fiveN1-rent-scraper.svg?style=flat-square)](https://travis-ci.org/neighborhood999/fiveN1-rent-scraper)
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
	if err := f.Scrape(1); err != nil {
		log.Fatal(err)
	}

	json := rent.ConvertToJSON(f.RentList)
	log.Println(string(json))
}
```

And output json:

```json
{
  "1": [
    {
      "preview":
        "https://hp1.591.com.tw/house/active/2013/11/04/138356961541004002_765x517.water3.jpg",
      "title": "松山路套房 * 交通方便近捷運後山埤站 *",
      "url": "https://rent.591.com.tw/rent-detail-5857738.html",
      "address": "信義區 - 松山路 119 巷",
      "rentType": "沒有格局說明",
      "optionType": "獨立套房",
      "ping": "6",
      "floor": "樓層：3/4",
      "price": "9,000 元 / 月",
      "isNew": true
    }
  ]
}
```

Index number is the representation **page number**, every index contain **30** items. 🏠

## Options

Create and generate URL by options:

```go
o := rent.NewOptions()
url, err := rent.GenerateURL(o)
if err != nil {
	log.Fatal(err)
}

log.Println(url)
```

You may set more options for your requirement. Reference below `Options` struct comment:

```go
type Options struct {
	Region    int    `url:"region"`              // 地區 - 預設：`1`
	Section   string `url:"section,omitempty"`   // 鄉鎮 - 可選擇多個區域，例如：`section=7,4`
	Kind      int    `url:"kind"`                // 租屋類型 - `0`：不限、`1`：整層住家、`2`：獨立套房、`3`：分租套房、`4`：雅房、`8`：車位，`24`：其他
	RentPrice string `url:"rentprice,omitempty"` // 租金 - `2`：5k - 10k、`3`：10k - 20k、`4`: 20k - 30k；或者可以輸入價格範圍，例如：`0,10000`
	Area      string `url:"area,omitempty"`      // 坪數格式 - `10,20`（10 到 20 坪）
	Order     string `url:"order"`               // 貼文時間 - 預設使用刊登時間：`posttime`，或是使用價格排序：`money`
	OrderType string `url:"orderType"`           // 排序方式 - `desc` 或 `asc`
	Sex       int    `url:"sex,omitempty"`       // 性別 - `0`：不限、`1`：男性、`2`：女性
	HasImg    string `url:"hasimg,omitempty"`    // 過濾是否有「房屋照片」 - ``：空值（不限）、`1`：是
	NotCover  string `url:"not_cover,omitempty"` // 過濾是否為「頂樓加蓋」 - ``：空值（不限）、`1`：是
	Role      string `url:"role,omitempty"`      // 過濾是否為「屋主刊登」 - ``：空值（不限）、`1`：是
	Shape     string `url:"shape,omitempty"`     // 房屋類型 - `1`：公寓、`2`：電梯大樓、`3`：透天厝、`4`：別墅
	Floor     string `url:"floor,omitempty"`     // 樓層 - `0,0`：不限、`0,1`：一樓、`2,6`：二樓到六樓、`6,12`：六樓到十二樓、`12,`：十二樓以上
	Option    string `url:"option,omitempty"`    // 提供設備 - `tv`：電視、`cold`：冷氣、`icebox`：冰箱、`hotwater`：熱水器、`naturalgas`：天然瓦斯、`four`：第四台、`broadband`：網路、`washer`：洗衣機、`bed`：床、`wardrobe`：衣櫃、`sofa`：沙發。可選擇多個設備，例如：option=tv,cold
	Other     string `url:"other,omitempty"`     // 其他條件 - `cartplace`：有車位、`lift`：有電梯、`balcony_1`：有陽台、`cook`：可開伙、`pet`：可養寵物、`tragoods`：近捷運、`lease`：可短期租賃。可選擇多個條件，例如：other=cartplace,cook
	FirstRow  int    `url:"firstRow"`
}
```

## Multiple Page Scrape

If you want to get more rent pieces of information, setting page amount in `Scrape` method:

```go
f := rent.NewFiveN1(url)
f.Scrape(5) // Page 1-5
```

When amount `> 1`, it will start goroutine automatically and correspond to the page number to scrape.

## `urlJumpIp` Code List

| Area   | Code |
| :----- | :--- |
| 台北市 | 1    |
| 基隆市 | 2    |
| 新北市 | 3    |
| 新竹市 | 4    |
| 新竹縣 | 5    |
| 桃園市 | 6    |
| 苗栗縣 | 7    |
| 台中市 | 8    |
| 彰化縣 | 10   |
| 南投縣 | 11   |
| 嘉義市 | 12   |
| 嘉義縣 | 13   |
| 雲林縣 | 14   |
| 台南市 | 15   |
| 高雄市 | 17   |
| 屏東縣 | 19   |
| 宜蘭縣 | 21   |
| 台東縣 | 22   |
| 花蓮縣 | 23   |
| 澎湖縣 | 24   |
| 金門縣 | 25   |
| 連江縣 | 26   |

## Secion Code List

More section code list reference to: [Section Code List](./docs/section-code-list.md).

## LICENSE

MIT © [Peng Jie](https://github.com/neighborhood999)
