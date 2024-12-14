package auth

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type AuthSMSService struct {
	svc   sms.Service
	token string
	key   []byte
}

func NewAuthSMSService(svc sms.Service, key string, token string) *AuthSMSService {
	return &AuthSMSService{
		svc:   svc,
		key:   []byte(key),
		token: token,
	}
}

// }}}
// {{{ Other structs

type SMSClaims struct {
	jwt.RegisteredClaims
}

// }}}
// {{{ Struct Methods

func (a *AuthSMSService) Send(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(a.token, &claims, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return err
	}
	return a.svc.Send(ctx, toNb, body, args...)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
