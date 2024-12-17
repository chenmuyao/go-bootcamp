package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wantBody ginx.Result
	}{
		{
			name:   "send ok",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "12345",
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: ginx.CodeOK,
				Msg:  "Sent successfully",
			},
		},
		{
			name:   "no phone number",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
			},
			phone:    "",
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "empty phone number",
			},
		},
		{
			name: "sent too frequently",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				err := rdb.Set(ctx, key, "600123", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "600123", code)
			},
			phone:    "12345",
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "sent too many",
			},
		},
		{
			name: "sent too frequently",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				err := rdb.Set(ctx, key, "600123", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "600123", code)
			},
			phone:    "12345",
			wantCode: http.StatusInternalServerError,
			wantBody: ginx.InternalServerErrorResult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			req, err := http.NewRequest(
				http.MethodPost,
				"/user/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res ginx.Result
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
