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

var gs *gosms.GoSMS

func init() {
	var aliOpts = &sdks.AliyunOptions{
		SignName:     "123",
		TemplateCode: "gsfas",
		OutId:        "1",
		AccessKeyId:  "xxxxxxxxxxx",
		AccessSecret: "xxxxxxxxxxxxxx",
	}
	gs = gosms.NewGoSMS(DB(), drivers.NewMysqlDriver(), gosms.Sdk{China: sdks.NewAliyunSdk(aliOpts)})
}

func main() {
	var code = "231532"
	var sms = adapter.Sms{
		Code:      code,
		MobilePre: 12,
		Mobile:    "13212341234",
		Ip:        "8.8.8.8",
	}

	err := gs.SendSMS(&sms)
	fmt.Println(err)

	//err = gs.CheckSMS(&sms)
	//fmt.Println(err)
}

func DB() *gorose.Engin {
	engin, err := gorose.Open(&gorose.Config{
		Driver: "mysql",
		Dsn:    "root:root@tcp(10.10.35.204:3306)/novel?charset=utf8",
		Prefix: "nv_",
	})
	if err != nil {
		panic(err.Error())
	}
	return engin
}
