package ioc

import (
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return &localsms.Service{}
}
