package rent

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	noMorePageErrMessage = "\x1b[91;1mNo More Pages！\x1b[0m"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	html, _ := ioutil.ReadFile("_fixtures/index.html")
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func TestGenerateURL(t *testing.T) {
	var err error
	o := NewOptions()

	_, err = GenerateURL(o)
	assert.EqualValues(t, nil, err)
}

func TestSetReqCookie(t *testing.T) {
	o := NewOptions()
	url, _ := GenerateURL(o)
	f := NewFiveN1(url)
	f.SetReqCookie("2")

	assert.Equal(t, "urlJumpIp=2", f.cookie.String())
}

func TestConvertToJSON(t *testing.T) {
	hm := make(map[int][]*HouseInfo)
	h := NewHouseInfo()

	h.Title = "Test Rent House Inforation"
	h.URL = "https://github.com/neighborhood999/fiveN1-rent-scraper"
	h.Address = "Taipei, Taiwan"
	h.RentType = "沒有格局說明"
	h.OptionType = "獨立套房"
	h.Ping = "9"
	h.Floor = "樓層：7/9"
	h.Price = "10,000 元 / 月"
	h.IsNew = true

	s := make([]*HouseInfo, 1)
	s[0] = h
	hm[0] = s

	b := ConvertToJSON(hm)

	var expectedType []byte
	assert.IsType(t, expectedType, b)
}

func TestExportJSON(t *testing.T) {
	text := []byte("Hello, World")
	ExportJSON(text)
}

func TestScrape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(callbackHandler))
	defer server.Close()
	mockURL := server.URL

	expectedPages := 12
	expectedRecords := 333

	f := NewFiveN1(mockURL)
	f.Scrape(1)

	assert.Equal(t, expectedPages, f.TotalPages)
	assert.Equal(t, expectedRecords, f.records)

	err := f.Scrape(13)
	assert.EqualError(t, err, noMorePageErrMessage)
}

func TestConcurrentScrape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(callbackHandler))
	defer server.Close()
	mockURL := server.URL + "/?"

	f := NewFiveN1(mockURL)
	f.Scrape(2)

	assert.IsType(t, HouseInfoCollection{}, f.RentList)
	assert.Equal(t, 2, len(f.RentList))
}

func TestGetTotalPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(callbackHandler))
	defer server.Close()
	mockURL := server.URL

	f := NewFiveN1(mockURL)
	assert.Equal(t, 0, f.TotalPages)
	assert.Equal(t, false, f.isExecuted)

	f.GetTotalPage()
	assert.Equal(t, 12, f.TotalPages)
	assert.Equal(t, true, f.isExecuted)
}
