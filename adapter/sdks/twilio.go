// documents: https://www.twilio.com/blog/2014/06/sending-sms-from-your-go-app.html
// author: kevin
package sdks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/gohouse/gosms/adapter"
	"strings"
)
// TwilioOptions 配置
type TwilioOptions struct {
	From         string
	Template     string
	AccessKeyId  string
	AccessSecret string
}
// TwilioSdk 主结构
type TwilioSdk struct {
	*TwilioOptions
}
// NewTwilioSdk 初始化TwilioSdk
func NewTwilioSdk(opts *TwilioOptions) *TwilioSdk {
	return &TwilioSdk{opts}
}
// SendSMS 实现发送接口
func (a *TwilioSdk) SendSMS(sms *adapter.Sms) (as adapter.SmsResponse) {
	// Set initial variables
	accountSid := a.AccessKeyId
	authToken := a.AccessSecret
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	//v.Set("To","+631234123412")
	v.Set("To", fmt.Sprintf("+%v%s", sms.MobilePre, sms.Mobile))
	v.Set("From", a.From)
	v.Set("Body", fmt.Sprintf(a.Template, sms.Code))
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		as = adapter.SmsResponse{}
		var data map[string]interface{}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(bodyBytes, &data)
		if err == nil {
			//fmt.Println(data["sid"])
			as.OrderNo = data["sid"].(string)
		}
		as.Result = string(bodyBytes)
		as.Error = err
		return as
	}

	b, _ := ioutil.ReadAll(resp.Body)
	return adapter.SmsResponse{
		Result:  string(b),
		OrderNo: "",
		Error:   errors.New(resp.Status),
	}
}

// CheckSMS 实现查询接口
func (a *TwilioSdk) CheckSMS(sms *adapter.Sms) (as adapter.SmsResponse) {
	return
}
