package adapter

// SmsResponse 发送短信后的返回结构体
type SmsResponse struct {
	Result interface{}
	OrderNo string
	Error error
}
