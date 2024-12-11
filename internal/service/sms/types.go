package sms

import (
	"context"
	"fmt"
)

type SMSThirdPartyError struct {
	Code int
	Msg  string
}

func (t *SMSThirdPartyError) Error() string {
	return fmt.Sprintf("sending sms error, code: %d, message: %s", t.Code, t.Msg)
}

type Service interface {
	Send(ctx context.Context, toNb string, body string, args ...string) error
}
