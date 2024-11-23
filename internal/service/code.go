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

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

// NOTE: check -> do: racing condition
// Redis single thread, use lua for atomic operations

// NOTE: template example: "Verification code for webook: {{.Code}}\nExpires in 10 min.\n[webook]"
func (svc *CodeService) Send(
	ctx context.Context,
	biz string,
	phone string,
	tpl template.Template,
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

func (svc *CodeService) Verify(
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

func (svc *CodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (svc *CodeService) generateBody(code int) string {
	return fmt.Sprintf("Verification code for webook: %06d\nExpires in 10 min.\n[webook]", code)
}
