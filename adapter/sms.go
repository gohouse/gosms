package adapter

type Sms struct {
	Id         int64  `gorose:"id" json:"id"`
	Code       string `gorose:"code" json:"code"`               // 验证码
	SmsStatus  int64  `gorose:"sms_status" json:"sms_status"`   // 状态:默认0发送失败,1发送成功,2已核销
	MobilePre  int64 `gorose:"mobile_pre" json:"mobile_pre"`   // 国家代码
	Mobile     string `gorose:"mobile" json:"mobile"`           // 手机号
	Ip         string `gorose:"ip" json:"ip"`                   // ip
	OrderNo    string `gorose:"order_no" json:"order_no"`       // 唯一编号
	SendResult string `gorose:"send_result" json:"send_result"` // 原始返回的结果
	CreatedAt  string `gorose:"created_at" json:"created_at"`
}

func (Sms) TableName() string {
	return "sms"
}
