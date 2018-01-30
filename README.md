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

## Options

Create default options then you can generate url:

```go
o := rent.NewOptions()
url, _ := rent.GenerateURL(o)
log.Println(url)
```

or you can setting more detail for your requirement, please reference below:

```go
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
```

## Multiple Page Scrape

Default will scrape first page, if you want to scrape more pages, setting page amount in `Scrape` method:

```go
f := rent.NewFiveN1(url)
f.Scrape(5) // scrape page 1-5
```

If page amount `> 1`, it will start goroutine automatically for scrape to correspond to page number.

### `urlJumpIp` Code List

| Area  | Code |
| :---- | :--- |
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

### Secion Code List

### 台北市

| 名稱   | 代碼 |
| :---- | :--- |
| 中正區 | 1    |
| 大同區 | 2    |
| 中山區 | 3    |
| 松山區 | 4    |
| 大安區 | 5    |
| 萬華區 | 6    |
| 信義區 | 7    |
| 士林區 | 8    |
| 北投區 | 9    |
| 內湖區 | 10   |
| 南港區 | 11   |
| 文山區 | 12   |

<details>
  <summary>More</summary>

<h3>基隆市</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
	    <tr>
	        <td>仁愛區</td>
	        <td>13</td>
	    </tr>
			<tr>
	        <td>信義區</td>
	        <td>14</td>
	    </tr>
			<tr>
	        <td>中正區</td>
	        <td>15</td>
	    </tr>
			<tr>
	        <td>中山區</td>
	        <td>16</td>
	    </tr>
			<tr>
	        <td>安樂區</td>
	        <td>17</td>
	    </tr>
			<tr>
	        <td>暖暖區</td>
	        <td>18</td>
	    </tr>
			<tr>
	        <td>七堵區</td>
	        <td>19</td>
	    </tr>
		</tbody>
  </table>

<h3>新北市</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
	    <tr>
	        <td>萬里區</td>
	        <td>20</td>
	    </tr>
			<tr>
	        <td>金山區</td>
	        <td>21</td>
	    </tr>
			<tr>
	        <td>板橋區</td>
	        <td>26</td>
	    </tr>
			<tr>
	        <td>汐止區</td>
	        <td>27</td>
	    </tr>
			<tr>
	        <td>深坑區</td>
	        <td>28</td>
	    </tr>
			<tr>
	        <td>石碇區</td>
	        <td>29</td>
	    </tr>
			<tr>
	        <td>瑞芳區</td>
	        <td>30</td>
	    </tr>
			<tr>
	        <td>平溪區</td>
	        <td>31</td>
	    </tr>
			<tr>
	        <td>雙溪區</td>
	        <td>32</td>
	    </tr>
			<tr>
	        <td>貢寮區</td>
	        <td>33</td>
	    </tr>
			<tr>
	        <td>新店區</td>
	        <td>34</td>
	    </tr>
			<tr>
	        <td>坪林區</td>
	        <td>35</td>
	    </tr>
			<tr>
	        <td>烏來區</td>
	        <td>36</td>
	    </tr>
			<tr>
	        <td>永和區</td>
	        <td>37</td>
	    </tr>
			<tr>
	        <td>中和區</td>
	        <td>38</td>
	    </tr>
			<tr>
	        <td>土城區</td>
	        <td>39</td>
	    </tr>
			<tr>
	        <td>三峽區</td>
	        <td>40</td>
	    </tr>
			<tr>
	        <td>樹林區</td>
	        <td>41</td>
	    </tr>
			<tr>
	        <td>鶯歌區</td>
	        <td>42</td>
	    </tr>
			<tr>
	        <td>三重區</td>
	        <td>43</td>
	    </tr>
			<tr>
	        <td>新莊區</td>
	        <td>44</td>
	    </tr>
			<tr>
	        <td>泰山區</td>
	        <td>45</td>
	    </tr>
			<tr>
	        <td>林口區</td>
	        <td>46</td>
	    </tr>
			<tr>
	        <td>蘆洲區</td>
	        <td>47</td>
	    </tr>
			<tr>
	        <td>五股區</td>
	        <td>48</td>
	    </tr>
			<tr>
	        <td>八里區</td>
	        <td>49</td>
	    </tr>
			<tr>
	        <td>淡水區</td>
	        <td>50</td>
	    </tr>
			<tr>
	        <td>三芝區</td>
	        <td>51</td>
	    </tr>
			<tr>
	        <td>石門區</td>
	        <td>52</td>
	    </tr>
		</tbody>
  </table>

<h3>新竹市</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
			<tr>
				<td>香山區</td>
				<td>370</td>
			</tr>
	    <tr>
	        <td>東區</td>
	        <td>371</td>
	    </tr>
			<tr>
	        <td>北區</td>
	        <td>372</td>
	    </tr>
		</tbody>
  </table>

<h3>新竹縣</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
			<tr>
				<td>竹北市</td>
				<td>54</td>
			</tr>
	    <tr>
	        <td>湖口鄉</td>
	        <td>55</td>
	    </tr>
			<tr>
	        <td>新豐鄉</td>
	        <td>56</td>
	    </tr>
			<tr>
	        <td>新埔鎮</td>
	        <td>57</td>
	    </tr>
			<tr>
	        <td>關西鎮</td>
	        <td>58</td>
	    </tr>
			<tr>
	        <td>芎林鄉</td>
	        <td>59</td>
	    </tr>
			<tr>
	        <td>寶山鄉</td>
	        <td>60</td>
	    </tr>
			<tr>
	        <td>竹東鎮</td>
	        <td>61</td>
	    </tr>
			<tr>
	        <td>五峰鄉</td>
	        <td>62</td>
	    </tr>
			<tr>
	        <td>橫山鄉</td>
	        <td>63</td>
	    </tr>
			<tr>
	        <td>尖石鄉</td>
	        <td>64</td>
	    </tr>
			<tr>
	        <td>北埔鄉</td>
	        <td>65</td>
	    </tr>
			<tr>
	        <td>峨嵋鄉</td>
	        <td>66</td>
	    </tr>
		</tbody>
  </table>

<h3>桃園市</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
			<tr>
				<td>中壢區</td>
				<td>67</td>
			</tr>
	    <tr>
	        <td>平鎮區</td>
	        <td>68</td>
	    </tr>
			<tr>
	        <td>龍潭區</td>
	        <td>69</td>
	    </tr>
			<tr>
	        <td>楊梅區</td>
	        <td>70</td>
	    </tr>
			<tr>
	        <td>新屋區</td>
	        <td>71</td>
	    </tr>
			<tr>
	        <td>觀音區</td>
	        <td>72</td>
	    </tr>
			<tr>
	        <td>桃園區</td>
	        <td>73</td>
	    </tr>
			<tr>
	        <td>龜山區</td>
	        <td>74</td>
	    </tr>
			<tr>
	        <td>八德區</td>
	        <td>75</td>
	    </tr>
			<tr>
	        <td>大溪區</td>
	        <td>76</td>
	    </tr>
			<tr>
	        <td>復興區</td>
	        <td>77</td>
	    </tr>
			<tr>
	        <td>大園區</td>
	        <td>78</td>
	    </tr>
			<tr>
	        <td>蘆竹區</td>
	        <td>79</td>
	    </tr>
		</tbody>
  </table>

<h3>苗栗縣</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
			<tr>
				<td>竹南鎮</td>
				<td>80</td>
			</tr>
	    <tr>
	        <td>頭份市</td>
	        <td>81</td>
	    </tr>
			<tr>
	        <td>三灣鄉</td>
	        <td>82</td>
	    </tr>
			<tr>
	        <td>南庄鄉</td>
	        <td>83</td>
	    </tr>
			<tr>
	        <td>獅潭鄉</td>
	        <td>84</td>
	    </tr>
			<tr>
	        <td>後龍鎮</td>
	        <td>85</td>
	    </tr>
			<tr>
	        <td>通霄鎮</td>
	        <td>86</td>
	    </tr>
			<tr>
	        <td>苑裡鎮</td>
	        <td>87</td>
	    </tr>
			<tr>
	        <td>苗栗市</td>
	        <td>88</td>
	    </tr>
			<tr>
	        <td>造橋鄉</td>
	        <td>89</td>
	    </tr>
			<tr>
	        <td>頭屋鄉</td>
	        <td>90</td>
	    </tr>
			<tr>
	        <td>公館鄉</td>
	        <td>91</td>
	    </tr>
			<tr>
	        <td>大湖鄉</td>
	        <td>92</td>
	    </tr>
			<tr>
	        <td>泰安鄉</td>
	        <td>93</td>
	    </tr>
			<tr>
	        <td>銅鑼鄉</td>
	        <td>94</td>
	    </tr>
			<tr>
	        <td>三義鄉</td>
	        <td>95</td>
	    </tr>
			<tr>
	        <td>西湖鄉</td>
	        <td>96</td>
	    </tr>
			<tr>
	        <td>卓蘭鎮</td>
	        <td>97</td>
	    </tr>
		</tbody>
  </table>

<h3>台中市</h3>
	<table>
		<thead>
    	<th>名稱</th>
    	<th>代碼</th>
		</thead>
		<tbody>
			<tr>
				<td>中區</td>
				<td>98</td>
			</tr>
	    <tr>
	        <td>東區</td>
	        <td>99</td>
	    </tr>
			<tr>
	        <td>南區</td>
	        <td>100</td>
	    </tr>
			<tr>
	        <td>西區</td>
	        <td>101</td>
	    </tr>
			<tr>
	        <td>北區</td>
	        <td>102</td>
	    </tr>
			<tr>
	        <td>北屯區</td>
	        <td>103</td>
	    </tr>
			<tr>
	        <td>西屯區</td>
	        <td>104</td>
	    </tr>
			<tr>
	        <td>南屯區</td>
	        <td>105</td>
	    </tr>
			<tr>
	        <td>太平區</td>
	        <td>106</td>
	    </tr>
			<tr>
	        <td>大里區</td>
	        <td>107</td>
	    </tr>
			<tr>
	        <td>霧峰區</td>
	        <td>108</td>
	    </tr>
			<tr>
	        <td>烏日區</td>
	        <td>109</td>
	    </tr>
			<tr>
	        <td>豐原區</td>
	        <td>110</td>
	    </tr>
			<tr>
	        <td>后里區</td>
	        <td>111</td>
	    </tr>
			<tr>
	        <td>石岡區</td>
	        <td>112</td>
	    </tr>
			<tr>
	        <td>東勢區</td>
	        <td>113</td>
	    </tr>
			<tr>
	        <td>和平區</td>
	        <td>114</td>
	    </tr>
			<tr>
	        <td>新社區</td>
	        <td>115</td>
	    </tr>
			<tr>
	        <td>潭子區</td>
	        <td>116</td>
	    </tr>
			<tr>
	        <td>大雅區</td>
	        <td>117</td>
	    </tr>
			<tr>
	        <td>神岡區</td>
	        <td>118</td>
	    </tr>
			<tr>
	        <td>大肚區</td>
	        <td>119</td>
	    </tr>
			<tr>
	        <td>沙鹿區</td>
	        <td>120</td>
	    </tr>
			<tr>
	        <td>龍井區</td>
	        <td>121</td>
	    </tr>
			<tr>
	        <td>梧棲區</td>
	        <td>122</td>
	    </tr>
			<tr>
	        <td>清水區</td>
	        <td>123</td>
	    </tr>
			<tr>
	        <td>大甲區</td>
	        <td>124</td>
	    </tr>
			<tr>
	        <td>外埔區</td>
	        <td>125</td>
	    </tr>
			<tr>
	        <td>大安區</td>
	        <td>126</td>
	    </tr>
		</tbody>
  </table>
</details>

## LICENSE

MIT © [Peng Jie](https://github.com/neighborhood999)
