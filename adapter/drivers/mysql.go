package drivers

import (
	"errors"
	"fmt"
	gorose "github.com/gohouse/gorose/v2"
	"github.com/gohouse/gosms/adapter"
)

type MysqlDriver struct{}

func NewMysqlDriver() *MysqlDriver {
	return &MysqlDriver{}
}

func (*MysqlDriver) CreateTable(db *gorose.Engin) (err error) {
	sqlStr := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s%s (
  id int(11) NOT NULL AUTO_INCREMENT,
  code varchar(6) NOT NULL COMMENT '验证码',
  sms_status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态:默认0发送失败,1发送成功,2已核销',
  mobile_pre varchar(5) NOT NULL DEFAULT '86' COMMENT '国家代码',
  mobile varchar(12) NOT NULL COMMENT '手机号',
  ip varchar(15) NOT NULL DEFAULT '' COMMENT 'ip',
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信验证码';`, db.GetPrefix(), "sms")
	_, err = db.NewSession().Execute(sqlStr)

	return
}

// GenerateSms 生成短信验证码入库
func (*MysqlDriver) GenerateSms(db *gorose.Engin, sms *adapter.Sms) (pkid int64, err error)  {
	return db.NewOrm().InsertGetId(sms)
}

// UpdateSmsStatus 更新短信发送结果信息
func (*MysqlDriver) UpdateSmsSendResult(db *gorose.Engin, sms *adapter.Sms) (pkid int64, err error)  {
	return db.NewOrm().Where("id",sms.Id).ExtraCols("send_result").Update(sms)
}

// GetLatestSms 根据条件获取最新一条发送结果
func (*MysqlDriver) GetLatestSms(db *gorose.Engin, sms *adapter.Sms) (err error)  {
	var code = sms.Code
	var dba = db.NewOrm()

	err = dba.Table(sms).
		Where("mobile_pre", sms.MobilePre).
		Where("mobile", sms.Mobile).
		//Where("code", sms.Code).
		Where("sms_status", "<", 2).
		Order("id desc").
		Select()
	if err!=nil {
		return
	}
	if sms.Id == 0 {
		return errors.New("请先发送验证码")
	}
	if code != sms.Code {
		return errors.New("验证码有误")
	}
	return
}


func (*MysqlDriver) VerifySms(db *gorose.Engin, sms *adapter.Sms) (affected_rows int64, err error) {
	sms.SmsStatus = 2
	return db.NewOrm().Where("id",sms.Id).Update(sms)
}