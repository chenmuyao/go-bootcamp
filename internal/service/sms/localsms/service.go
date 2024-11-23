package localsms

import (
	"context"
	"log/slog"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, toNb string, body string, args ...string) error {
	slog.Info("sms", "body", body)
	return nil
}
