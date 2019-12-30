package sdks

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gohouse/date"
	"github.com/gohouse/gosms/adapter"
)

// 请求参数
//名称	类型	是否必选	示例值	描述
//PhoneNumbers	String	是	15900000000
//接收短信的手机号码。
//
//格式：
//
//国内短信：11位手机号码，例如15951955195。
//国际/港澳台消息：国际区号+号码，例如85200000000。
//支持对多个手机号码发送短信，手机号码之间以英文逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
//
//说明 验证码类型短信，建议使用单独发送的方式。
//SignName	String	是	阿里云
//短信签名名称。请在控制台签名管理页面签名名称一列查看。
//
//说明 必须是已添加、并通过审核的短信签名。
//TemplateCode	String	是	SMS_153055065
//短信模板ID。请在控制台模板管理页面模板CODE一列查看。
//
//说明 必须是已添加、并通过审核的短信签名；且发送国际/港澳台消息时，请使用国际/港澳台短信模版。
//AccessKeyId	String	否	LTAIP00vvvvvvvvv
//主账号AccessKey的ID。
//
//Action	String	否	SendSms
//系统规定参数。取值：SendSms。
//
//OutId	String	否	abcdefgh
//外部流水扩展字段。
//
//SmsUpExtendCode	String	否	90999
//上行短信扩展码，无特殊需要此字段的用户请忽略此字段。
//
//TemplateParam	String	否	{"code":"1111"}
//短信模板变量对应的实际值，JSON格式。
//
//说明 如果JSON中需要带换行符，请参照标准的JSON协议处理。
type AliyunOptions struct {
	SignName string
	TemplateCode string
	OutId string
	AccessKeyId string
	AccessSecret string
}

type AliyunSdk struct {
	*AliyunOptions
	client *dysmsapi.Client
}

func NewAliyunSdk(opts *AliyunOptions) *AliyunSdk {
	cli,err := dysmsapi.NewClientWithAccessKey("cn-hangzhou",
		opts.AccessKeyId, opts.AccessSecret)
	if err!=nil {
		panic(err.Error())
	}
	return &AliyunSdk{opts,cli}
}

func (a *AliyunSdk) SendSMS(as *adapter.Sms) adapter.SmsResponse {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.SignName = a.SignName
	request.TemplateCode = a.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, as.Code)
	request.OutId = a.OutId

	request.PhoneNumbers = fmt.Sprintf("%v%s",as.MobilePre,as.Mobile)
	request.SmsUpExtendCode = as.Code

	response, err := a.client.SendSms(request)
	if err != nil {
		return adapter.SmsResponse{
			Result:response,
			Error:err}
	}
	if response.Code != "OK" {
		return adapter.SmsResponse{
			Result:response,
			Error:errors.New(response.Message)}
	}
	return adapter.SmsResponse{
		Result:  response,
		OrderNo: response.BizId,
		Error:   err,
	}
}

func (a *AliyunSdk) CheckSMS(sms *adapter.Sms) adapter.SmsResponse {
	request := dysmsapi.CreateQuerySendDetailsRequest()
	request.Scheme = "https"

	request.PhoneNumber = fmt.Sprintf("%v%s",sms.MobilePre,sms.Mobile)
	request.SendDate = date.NewDate(date.BindDateTime(sms.CreatedAt)).TodayDate()
	request.PageSize = requests.NewInteger(1)
	request.CurrentPage = requests.NewInteger(1)
	request.BizId = sms.OrderNo

	response, err := a.client.QuerySendDetails(request)

	return adapter.SmsResponse{
		Result:  response,
		OrderNo: "",
		Error:   err,
	}
}
