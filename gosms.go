package gosms

import (
	"errors"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/gosms/adapter"
	"github.com/gohouse/t"
)

// Sdk 运营商sdk
type Sdk map[CallingCode]adapter.SdkAdapter
//type Sdk struct {
//	// 中国的短信接口
//	China adapter.SdkAdapter
//	// 国际短信接口
//	Global adapter.SdkAdapter
//}

// Driver 保存数据库的适配器,晚上自动入库与核销
type Driver adapter.DriverAdapter

// GoSMS sms主体
type GoSMS struct {
	engin  *gorose.Engin
	driver adapter.DriverAdapter
	sdk    Sdk
}

// NewGoSMS 初始化sms
func NewGoSMS(engin *gorose.Engin, driver Driver, sdk Sdk) *GoSMS {
	if engin==nil || driver == nil || sdk == nil {
		panic("对想为空")
	}
	var s = &GoSMS{engin: engin, driver: driver, sdk: sdk}
	// 初始化sms表
	err := driver.CreateTable(engin)
	if err != nil {
		panic(err.Error())
	}
	return s
}

// SendSMS 发送短信
func (s *GoSMS) SendSMS(sms *adapter.Sms) (err error) {
	// 生成短信发送记录,使用数据库适配器保存本地数据库
	lastInsertId, err := s.driver.GenerateSms(s.engin, sms)
	if err != nil {
		return
	}
	sms.Id = lastInsertId

	// 调用sdk发送短信
	var res adapter.SmsResponse
	mpre := CallingCode(sms.MobilePre)
	if sdkTmp,ok := s.sdk[mpre]; ok {
		res = sdkTmp.SendSMS(sms)
	} else {
		sdkTmp = s.sdk[CC_GLOBAL]
		if sdkTmp == nil {
			return errors.New("短信sdk未就绪")
		}
		res = sdkTmp.SendSMS(sms)
	}
	//if sms.MobilePre == 86 {	// 调用国内的短信服务
	//	res = s.sdk.China.SendSMS(sms)
	//} else {	// 国际服务
	//	res = s.sdk.Global.SendSMS(sms)
	//}

	// 记录返回结果
	if res.Error == nil {
		sms.SmsStatus = 1
		sms.OrderNo = res.OrderNo
	}
	sms.SendResult = ""
	if res.Error != nil {
		sms.SendResult = res.Error.Error()
	} else if res.Result != nil {
		sms.SendResult = t.New(res.Result).String()
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

// CheckSMS 核销短信
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
