package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	mock "github.com/chenmuyao/go-bootcamp/internal/service/mocks"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	jwtmocks "github.com/chenmuyao/go-bootcamp/internal/web/jwt/mocks"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler)

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody ginx.Result
	}{
		{
			name: "signup ok",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				hdl := jwtmocks.NewMockHandler(ctrl)
				hdl.EXPECT().SetLoginToken(gomock.Any(), gomock.Any()).Return(nil)
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}, nil)
				return userSvc, nil, hdl
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password!123",
"confirm_password": "password!123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: ginx.CodeOK,
				Msg:  "signup success",
			},
		},
		{
			name: "wrong json format",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password!123",
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "not a valid email",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "notanemail.com",
"password": "password!123",
"confirm_password": "password!123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "not a valid email",
			},
		},
		{
			name: "not a valid password",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password",
"confirm_password": "password"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "not a valid password",
			},
		},
		{
			name: "passwords don't match",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password!123",
"confirm_password": "password!1234"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "2 passwords don't match",
			},
		},
		{
			name: "user exists",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{}, service.ErrDuplicatedUser)
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password!123",
"confirm_password": "password!123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "user exists",
			},
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{}, errors.New("db error"))
				return userSvc, nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/signup",
					bytes.NewReader([]byte(`{
"email": "ok@test.com",
"password": "password!123",
"confirm_password": "password!123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: ginx.InternalServerErrorResult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Prepare
			userSvc, codeSvc, jwtHdl := tc.mock(ctrl)
			l, _ := zap.NewDevelopment()
			hdl := NewUserHandler(
				logger.NewZapLogger(l),
				userSvc,
				codeSvc,
				jwtHdl,
			)
			ginx.InitCounter(prom.CounterOpts{
				Namespace: "my_company",
				Subsystem: "wetravel",
				Name:      "errcode",
				Help:      "Error code data",
				ConstLabels: prom.Labels{
					"instance_id": "instance",
				},
			})

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// Run Test
			server.ServeHTTP(recorder, req)

			// Check Results
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res ginx.Result
			_ = json.NewDecoder(recorder.Body).Decode(&res)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}

func TestLoginJWT(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler)

		reqBuilder func(*testing.T) *http.Request

		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "login ok",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, ijwt.Handler) {
				hdl := jwtmocks.NewMockHandler(ctrl)
				hdl.EXPECT().SetLoginToken(gomock.Any(), gomock.Any()).Return(nil)
				us := mock.NewMockUserService(ctrl)
				us.EXPECT().Login(gomock.Any(), "ok@test.com", "password123!").Return(domain.User{
					ID:       123,
					Email:    "ok@test.com",
					Password: "$2a$10$ak/.qMW4bKq3ksEXokuuquyXXNjHAv1t8wqwWvGVbje0rjyrZTqgy",
				}, nil)
				return us, nil, hdl
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("POST", "/user/login", bytes.NewBuffer([]byte(`{
"email": "ok@test.com",
"password": "password123!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Code: ginx.CodeOK,
				Msg:  "successful login",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			us, cs, jwtHdl := tc.mock(ctrl)
			l, _ := zap.NewDevelopment()
			hdl := NewUserHandler(logger.NewZapLogger(l), us, cs, jwtHdl)
			ginx.InitCounter(prom.CounterOpts{
				Namespace: "my_company",
				Subsystem: "wetravel",
				Name:      "errcode",
				Help:      "Error code data",
				ConstLabels: prom.Labels{
					"instance_id": "instance",
				},
			})

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res ginx.Result
			err := json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestUserEmailPattern(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{name: "no @", email: "1234", match: false},
		{name: "no @ suffix", email: "1234@", match: false},
		{name: "suffix not valid", email: "1234@1234", match: false},
		{name: "no username", email: "@1234.com", match: false},
		{name: "ok", email: "123@1234.com", match: true},
	}

	h := NewUserHandler(logger.NewNopLogger(), nil, nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRegex.MatchString(tc.email)
			assert.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
