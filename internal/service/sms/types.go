package sms

import (
	"context"
	"fmt"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

type Service interface {
	Send(ctx context.Context, toNb string, body string, args ...string) error
}

// }}}
// {{{ Struct

// }}}
// {{{ Other structs

type SMSThirdPartyError struct {
	Code int
	Msg  string
}

// }}}
// {{{ Struct Methods

func (t *SMSThirdPartyError) Error() string {
	return fmt.Sprintf("sending sms error, code: %d, message: %s", t.Code, t.Msg)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
