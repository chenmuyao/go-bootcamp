package sms

import "context"

type Service interface {
	Send(ctx context.Context, toNb string, body string, args ...string) error
}
