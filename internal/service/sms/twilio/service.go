package twilio

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	tw "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Service struct {
	client *tw.RestClient
	fromNb string
}

func NewService(client *tw.RestClient, fromNb string) *Service {
	return &Service{
		client: client,
		fromNb: fromNb,
	}
}

func (s *Service) Send(ctx context.Context, toNb string, body string, args ...string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(toNb).SetFrom(s.fromNb).SetBody(body)

	res, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	if res.Status != nil && (*res.Status == "failed" || *res.Status == "undelivered") {
		if res.ErrorCode != nil && res.ErrorMessage != nil {
			return &sms.SMSThirdPartyError{Code: *res.ErrorCode, Msg: *res.ErrorMessage}
		}
	}
	return nil
}
