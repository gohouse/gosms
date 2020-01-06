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

// 阿里云公共配置
var aliOpts = &sdks.AliyunOptions{
	SignName:     "123",
	TemplateCode: "gsfas",
	OutId:        "1",
	AccessKeyId:  "xxxxxxxxxxx",
	AccessSecret: "xxxxxxxxxxxxxx",
}

func main() {
	// 初始化短信,包括数据库,入库adapter,短信服务
	var gs = gosms.NewGoSMS(
		DB(),                                      // 数据库orm
		drivers.NewMysqlDriver(),                  // 短信入库和核销
		gosms.Sdk{gosms.CC_CN: sdks.NewAliyunSdk(aliOpts)}, // 短信服务商
	)

	// 短信验证码
	var code = "231532"
	// 入库的数据
	var sms = adapter.Sms{
		Code:      code,
		MobilePre: 12,
		Mobile:    "13212341234",
		Ip:        "8.8.8.8",
	}

	// 发送短信
	err := gs.SendSMS(&sms)
	fmt.Println(err)

	// 核销短信
	//err = gs.CheckSMS(&sms)
	//fmt.Println(err)
}

// 试用gorose, 初始化数据库
func DB() *gorose.Engin {
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
