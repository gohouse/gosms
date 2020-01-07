# gosms
golang通用发送短信库,包含aliyun短信可发国内,twilio可发国际,同时可以自由添加更多短信服务商,自动入库和验证

## 安装
```shell script
go get github.com/gohouse/gosms
```

## 特点
- 支持国内和国际短信  
- 自动维护数据库短信入库和核销  
- 可以自己添加更多短信服务商支持  

## 默认支持
- 短信服务商,默认支持 aliyun(国内), twilio(国际)  
- 数据库适配器: 默认提供了mysql的适配器, 可以自行添加,自动维护数据库表  

> 说明: 之所以不用 aliyun 发送国际短信,主要是因为太贵了,贵了将近1倍

> 如果想添加更多数据库支持,只需要实现接口 `datapter.DriverAdapter` 即可  
> 如果想添加更多短信服务商支持,只需要实现接口 `datapter.SdkAdapter` 即可  
>
## 使用
- 使用 aliyun 示例
```go
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
		DB(),	// 数据库orm
		drivers.NewMysqlDriver(),	// 短信入库和核销
		// 配置sdk, gosms.CC_CN 为中国
		// gosms.CC_Global 为全球, 另外可以指定人员国家,具体参考 gosms.CallingCode
		gosms.Sdk{gosms.CC_CN: sdks.NewAliyunSdk(aliOpts)},	// 短信服务商
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

```