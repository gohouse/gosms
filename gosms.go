package gosms

import (
	"encoding/json"
	"errors"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/gosms/adapter"
)

// Sdk ...
type Sdk struct {
	// 中国的短信接口
	China adapter.SdkAdapter
	// 国际短信接口
	Global adapter.SdkAdapter
}

// Driver ...
type Driver adapter.DriverAdapter

// GoSMS ...
type GoSMS struct {
	engin  *gorose.Engin
	driver adapter.DriverAdapter
	sdk    Sdk
}

// NewGoSMS ...
func NewGoSMS(engin *gorose.Engin, driver Driver, sdk Sdk) *GoSMS {
	var s = &GoSMS{engin: engin, driver: driver, sdk: sdk}
	// 初始化sms表
	err := driver.CreateTable(engin)
	if err != nil {
		panic(err.Error())
	}
	return s
}

// SendSMS ...
func (s *GoSMS) SendSMS(sms *adapter.Sms) (err error) {
	// 生成短信发送记录
	lastInsertId, err := s.driver.GenerateSms(s.engin, sms)
	if err != nil {
		return
	}
	sms.Id = lastInsertId

	// 调用sdk发送短信
	var res adapter.SmsResponse
	if sms.MobilePre == 86 {
		res = s.sdk.China.SendSMS(sms)
	} else {
		res = s.sdk.Global.SendSMS(sms)
	}

	// 记录返回结果
	if res.Error == nil {
		sms.SmsStatus = 1
		sms.OrderNo = res.OrderNo
	}
	if res.Error != nil {
		sms.SendResult = res.Error.Error()
	}
	if res.Result != nil {
		js, _ := json.Marshal(res.Result)
		sms.SendResult = string(js)
	}
	//aff, err := s.engin.NewOrm().Where("id",lastInsertId).Update(sms)
	aff, err := s.driver.UpdateSmsSendResult(s.engin, sms)

	if err != nil || aff == 0 {
		//todo 只是记录原始返回失败, 但是验证通过,如有需要,可以做日志记录
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// CheckSMS ...
func (s *GoSMS) CheckSMS(sms *adapter.Sms) (err error) {
	// 检查是否已经发送了验证码
	err = s.driver.GetLatestSms(s.engin, sms)
	if err != nil {
		return
	}

	if sms.Id == 0 {
		return errors.New("请先发送验证码")
	}

	//// 向运营商核对验证码
	//err = s.sdk.CheckSMS(sms)
	//if err != nil {
	//	return
	//}

	// 核销验证码
	aff, err := s.driver.VerifySms(s.engin, sms)
	if err != nil || aff == 0 {
		//todo 只是核销失败, 但是验证通过,如有需要,可以做日志记录
	}

	return nil
}
