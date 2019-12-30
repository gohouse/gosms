package adapter

import (
	"github.com/gohouse/gorose/v2"
)

type DriverAdapter interface {
	// 创建对应的发短信表
	CreateTable(db *gorose.Engin) (err error)
	// 生成短信记录入库
	GenerateSms(db *gorose.Engin, sms *Sms) (pkid int64, err error)
	// 更新短信发送结果信息
	UpdateSmsSendResult(db *gorose.Engin, sms *Sms) (pkid int64, err error)
	// 根据条件获取最新一条发送结果
	GetLatestSms(db *gorose.Engin, sms *Sms) (err error)
	// 核销短息
	VerifySms(db *gorose.Engin, sms *Sms) (affected_rows int64, err error)
}

type SdkAdapter interface {
	SendSMS(*Sms) SmsResponse
	CheckSMS(*Sms) SmsResponse
}
