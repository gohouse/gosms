package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/gosms"
	"github.com/gohouse/gosms/adapter"
	"github.com/gohouse/gosms/adapter/drivers"
	"github.com/gohouse/gosms/adapter/sdks"
)

var gs2 *gosms.GoSMS

func init() {
	var twilioOptions = &sdks.TwilioOptions{
		From:         "+12341234",
		Template:     "您的验证码是: %s",
		AccessKeyId:  "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AccessSecret: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}
	gs2 = gosms.NewGoSMS(DB2(), drivers.NewMysqlDriver(), gosms.Sdk{gosms.CC_GLOBAL: sdks.NewTwilioSdk(twilioOptions)})
}

func main() {
	var code = "231532"
	var sms = adapter.Sms{
		Code:      code,
		MobilePre: 12,
		Mobile:    "12341234",
		Ip:        "8.8.8.8",
	}

	err := gs2.SendSMS(&sms)
	fmt.Println(err)
	if err==nil {
		fmt.Println("发送成功:",code)
	}

	//err = gs.CheckSMS(&sms)
	//fmt.Println(err)
}

func DB2() *gorose.Engin {
	engin, err := gorose.Open(&gorose.Config{
		Driver: "mysql",
		Dsn:    "root:root@tcp(localhost:3306)/novel?charset=utf8",
		Prefix: "nv_",
	})
	if err != nil {
		panic(err.Error())
	}
	return engin
}
