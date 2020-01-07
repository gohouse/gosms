package gosms

import (
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/gosms/adapter/drivers"
	"github.com/gohouse/gosms/adapter/sdks"
	"testing"
)

func TestNewGoSMS(t *testing.T) {
	engin,_ := gorose.Open()
	NewGoSMS(engin, drivers.NewMysqlDriver(), Sdk{CC_GLOBAL:sdks.NewTwilioSdk(&sdks.TwilioOptions{
		From:         "",
		Template:     "",
		AccessKeyId:  "",
		AccessSecret: "",
	})})
}
