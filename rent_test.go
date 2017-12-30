package rent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hasImgErrMessage   = "`HasImg` 請輸入 0 或是 1 的值！"
	notCoverErrMessage = "`NotCover` 請輸入 0 或是 1 的值！"
	roleErrMessage     = "`Role` 請輸入 0 或是 1 的值！"
)

func TestGenerateURL(t *testing.T) {
	var err error
	o := NewOptions()

	_, err = GenerateURL(o)
	assert.EqualValues(t, nil, err)

	o.HasImg = 2
	_, err = GenerateURL(o)
	assert.EqualError(t, err, hasImgErrMessage)

	o.HasImg = 1
	o.NotCover = 2
	_, err = GenerateURL(o)
	assert.EqualError(t, err, notCoverErrMessage)

	o.HasImg = 1
	o.NotCover = 1
	o.Role = 2
	_, err = GenerateURL(o)
	assert.EqualError(t, err, roleErrMessage)
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
