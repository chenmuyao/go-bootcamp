package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"text/template"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
)

// TODO: Implement Email and Voice Code service

// {{{ Consts

// }}}
// {{{ Global Varirables

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
)

// }}}
// {{{ Interface

//go:generate mockgen -source=./code.go -package=svcmocks -destination=./mocks/code.mock.go
type CodeService interface {
	Send(ctx context.Context, biz string, phone string, tpl *template.Template) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// }}}
// {{{ Struct

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  sms,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

// NOTE: check -> do: racing condition
// Redis single thread, use lua for atomic operations

// NOTE: template example: "Verification code for WeTravel: {{.Code}}\nExpires in 10 min.\n[WeTravel]"
func (svc *codeService) Send(
	ctx context.Context,
	biz string,
	phone string,
	tpl *template.Template,
) error {
	code := svc.generateCode()

	// NOTE: In old school DDD theory, we should implement the Set and verify
	// logic here at the service level. But it's more complicated to solve
	// racing condition. We must introduce a distributed lock.
	// However, the problem of putting that logic into infra level is that
	// we have to reimplement the logic for other implementations like
	// memcached.
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	var buff bytes.Buffer

	val := struct {
		Code string
	}{
		Code: code,
	}
	if err := tpl.Execute(&buff, val); err != nil {
		// template error
		return err
	}
	return svc.sms.Send(ctx, phone, buff.String())
}

func (svc *codeService) Verify(
	ctx context.Context,
	biz string,
	phone string,
	inputCode string,
) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// shielding this error from the outside
		return false, nil
	}
	return ok, err
}

func (svc *codeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
