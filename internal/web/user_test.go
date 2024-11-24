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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody Result
	}{
		{
			name: "signup ok",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}, nil)
				return userSvc, nil
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
			wantBody: Result{
				Code: CodeOK,
				Msg:  "signup success",
			},
		},
		{
			name: "wrong json format",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil
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
			wantBody: Result{
				Code: CodeUserSide,
				Msg:  "not a valid email",
			},
		},
		{
			name: "not a valid password",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil
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
			wantBody: Result{
				Code: CodeUserSide,
				Msg:  "not a valid password",
			},
		},
		{
			name: "passwords don't match",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				return userSvc, nil
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
			wantBody: Result{
				Code: CodeUserSide,
				Msg:  "2 passwords don't match",
			},
		},
		{
			name: "user exists",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{}, service.ErrDuplicatedUser)
				return userSvc, nil
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
			wantBody: Result{
				Code: CodeUserSide,
				Msg:  "user exists",
			},
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "ok@test.com",
					Password: "password!123",
				}).Return(domain.User{}, errors.New("db error"))
				return userSvc, nil
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
			wantBody: InternalServerErrorResult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Prepare
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// Run Test
			server.ServeHTTP(recorder, req)

			// Check Results
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Result
			_ = json.NewDecoder(recorder.Body).Decode(&res)
			assert.Equal(t, tc.wantBody, res)
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

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRegex.MatchString(tc.email)
			assert.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}
