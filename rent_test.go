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
