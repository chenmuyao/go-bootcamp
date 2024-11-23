package twilio

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	tw "github.com/twilio/twilio-go"
)

func TestService(t *testing.T) {
	// XXX: MUST set TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN in env
	client := tw.NewRestClient()
	service := NewService(client, "")

	err := service.Send(
		context.Background(),
		"",
		fmt.Sprintf("Verification code: %06d", rand.Intn(1000000)))
	// assert.NoError(t, err)
	assert.ErrorContains(t, err, "credentials")
}
