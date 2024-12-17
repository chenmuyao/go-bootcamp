package web

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"text/template"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// {{{ Consts

const (
	emailRegexPattern    = `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
	codeSMSTemplate      = "Verification code for WeTravel: {{.Code}}\nExpires in 10 min.\n[WeTravel]"
)

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type UserHandler struct {
	l logger.Logger
	ijwt.Handler
	emailRegex    *regexp2.Regexp
	passwordRegex *regexp2.Regexp
	svc           service.UserService
	codeSvc       service.CodeService
}

func NewUserHandler(
	l logger.Logger,
	svc service.UserService,
	codeSvc service.CodeService,
	hdl ijwt.Handler,
) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegex: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		svc:           svc,
		codeSvc:       codeSvc,
		Handler:       hdl,
		l:             l,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	user := server.Group("/user/")
	user.POST("/signup", ginx.WrapBody(h.l, h.SignUp))
	// user.POST("/login", h.Login)
	user.POST("/login", ginx.WrapBody(h.l, h.LoginJWT))
	user.GET("/profile", ginx.WrapLog(h.l, h.Profile))
	user.GET("/profile/:id", ginx.WrapLog(h.l, h.Profile))
	user.POST("/edit", ginx.WrapBodyAndClaims(h.l, h.Edit))

	user.GET("/refresh_token", ginx.WrapLog(h.l, h.RefreshToken))

	// SMS code login
	user.POST("/login_sms/code/send", ginx.WrapBody(h.l, h.SendSMSLoginCode))
	user.POST("/login_sms", ginx.WrapBody(h.l, h.LoginSMS))

	user.POST("/logout", ginx.WrapLog(h.l, h.LogoutJWT))
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context, req UserSMSCodeReq) (ginx.Result, error) {
	if req.Phone == "" {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "empty phone number",
		}, errors.New("empty phone number")
	}
	tpl, err := template.New(bizLogin).Parse(codeSMSTemplate)
	if err != nil {
		return ginx.Result{
			Code: ginx.CodeServerSide,
			Msg:  "internal sever error",
		}, fmt.Errorf("cannot parse sms template: %w", err)
	}
	err = h.codeSvc.Send(ctx, bizLogin, req.Phone, tpl)
	switch err {
	case nil:
		return ginx.Result{
			Code: ginx.CodeOK,
			Msg:  "Sent successfully",
		}, nil
	case service.ErrCodeSendTooMany:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "sent too many",
		}, errors.New("sent too many")
	default:

		return ginx.InternalServerErrorResult, fmt.Errorf("sms send code error: %w", err)
	}
}

func (h *UserHandler) LoginSMS(ctx *gin.Context, req UserLoginSMSReq) (ginx.Result, error) {
	if req.Phone == "" || req.Code == "" {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "must have the phone number and the code",
		}, fmt.Errorf("invalid phone or code: %s/%s", req.Phone, req.Code)
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf("verify code internal error: %w", err)
	}
	if !ok {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "Wrong code, please retry",
		}, errors.New("wrong code")
	}

	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf(
			"find or create user internal error: %w",
			err,
		)
	}
	err = h.SetLoginToken(ctx, u.ID)
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf(
			"set login token internal error: %w",
			err,
		)
	}
	return ginx.Result{
		Code: ginx.CodeOK,
		Msg:  "successful login",
	}, nil
}

func (h *UserHandler) SignUp(ctx *gin.Context, req UserSignUpReq) (ginx.Result, error) {
	// Check request
	isEmail, err := h.emailRegex.MatchString(req.Email)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}
	if !isEmail {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "not a valid email",
		}, fmt.Errorf("invalid email")
	}

	validPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}
	if !validPassword {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "not a valid password",
		}, fmt.Errorf("invalid password")
	}

	if req.Password != req.ConfirmPassword {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "2 passwords don't match",
		}, fmt.Errorf("2 passwords don't match")
	}

	u, err := h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		err = h.SetLoginToken(ctx, u.ID)
		if err != nil {
			return ginx.InternalServerErrorResult, err
		}
		return ginx.Result{
			Code: ginx.CodeOK,
			Msg:  "signup success",
		}, nil
	case service.ErrDuplicatedUser:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "user exists",
		}, fmt.Errorf("user exists")
	default:
		return ginx.InternalServerErrorResult, err
	}
}

// func (h *UserHandler) LoginSession(ctx *gin.Context) {
// 	type Req struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
//
// 	var req Req
//
// 	if err := ctx.Bind(&req); err != nil {
// 		return
// 	}
//
// 	// NOTE: No need to check, because if it's not valid, we won't get
// 	// anything from the DB anyway.
//
// 	u, err := h.svc.Login(ctx, req.Email, req.Password)
// 	switch err {
// 	case nil:
// 		sess := sessions.Default(ctx)
// 		sess.Set("userID", u.ID)
// 		sess.Options(sessions.Options{
// 			MaxAge:   900, // 15min - expire time of the session (+ expire time of the userID entry in redis.)
// 			HttpOnly: true,
// 		})
// 		err = sess.Save()
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, Result{
// 			Code: CodeOK,
// 			Msg:  "successful login",
// 		})
// 	case service.ErrInvalidUserOrPassword:
// 		ctx.JSON(http.StatusBadRequest, Result{
// 			Code: CodeUserSide,
// 			Msg:  "wrong login or password",
// 		})
// 	default:
// 		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
// 	}
// }

func (h *UserHandler) LoginJWT(ctx *gin.Context, req UserLoginReq) (ginx.Result, error) {
	// NOTE: No need to check, because if it's not valid, we won't get
	// anything from the DB anyway.

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		err = h.SetLoginToken(ctx, u.ID)
		if err != nil {
			return ginx.InternalServerErrorResult, err
		}
		return ginx.Result{
			Code: ginx.CodeOK,
			Msg:  "successful login",
		}, nil
	case service.ErrInvalidUserOrPassword:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "wrong login or password",
		}, errors.New("wrong login or password")
	default:
		return ginx.InternalServerErrorResult, err
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) (ginx.Result, error) {
	var userID int64
	var err error
	id := ctx.Param("id")
	if id == "" {
		// get userID from session
		// userID, err = h.getUserIDFromSession(ctx)
		// if err != nil {
		// 	log.Println(err)
		// 	ctx.String(http.StatusInternalServerError, "system error")
		// 	return
		// }
		userID = ctx.MustGet("user").(ijwt.UserClaims).UID
	} else {
		if userID, err = strconv.ParseInt(id, 10, 64); err != nil {
			return ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  fmt.Sprintf("unknown userID: %s", id),
			}, fmt.Errorf("unknown userID %s: %w", id, err)
		}
	}

	u, err := h.svc.GetProfile(ctx, userID)
	switch err {
	case nil:
		break
	case service.ErrInvalidUserID:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "invalid userID",
		}, fmt.Errorf("unknown userID %d: %w", userID, err)
	default:
		log.Printf("failed to get user profile: %s\n", err.Error())
		return ginx.InternalServerErrorResult, fmt.Errorf(
			"failed to get user %d profile: %w",
			userID,
			err,
		)
	}

	resp := struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
		Profile  string `json:"profile"`
	}{
		Email:    u.Email,
		Phone:    u.Phone,
		Name:     u.Name,
		Birthday: u.Birthday.Format("2006-01-02"),
		Profile:  u.Profile,
	}

	return ginx.Result{
		Code: ginx.CodeOK,
		Data: resp,
	}, nil
}

func (h *UserHandler) Edit(
	ctx *gin.Context,
	req UserEditReq,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	userID := uc.UID

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}
	err = h.svc.EditProfile(ctx, &domain.User{
		ID:       userID,
		Name:     req.Name,
		Birthday: birthday,
		Profile:  req.Profile,
	})
	switch err {
	case nil:
		return ginx.Result{
			Code: ginx.CodeOK,
			Msg:  "user profile update success",
		}, nil
	case service.ErrInvalidUserID:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "unknown userID",
		}, fmt.Errorf("invalid user id %d: %w", userID, err)
	default:
		return ginx.InternalServerErrorResult, fmt.Errorf("failed to edit user profile: %w", err)
	}
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) (ginx.Result, error) {
	tokenStr := h.ExtractToken(ctx)

	var rc ijwt.RefreshClaims

	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(t *jwt.Token) (interface{}, error) {
		return ijwt.RefreshKey, nil
	})
	if err != nil {
		return ginx.UnauthorizedResult, err
	}
	if token == nil || !token.Valid {
		return ginx.UnauthorizedResult, err
	}

	err = h.CheckSession(ctx, rc.SSID)
	if err != nil {
		return ginx.UnauthorizedResult, err
	}

	err = h.SetLoginToken(ctx, rc.UID)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}

	return ginx.Result{
		Code: ginx.CodeOK,
		Msg:  "OK",
	}, nil
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) (ginx.Result, error) {
	err := h.ClearToken(ctx)
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf("failed to clear token: %w", err)
	}
	return ginx.Result{
		Code: ginx.CodeOK,
		Msg:  "Logout success",
	}, nil
}

// func (h *UserHandler) getUserIDFromSession(ctx *gin.Context) (int64, error) {
// 	sess := sessions.Default(ctx)
// 	userID, ok := sess.Get("userID").(int64)
// 	if !ok {
// 		return 0, errors.New("failed to get userID from session")
// 	}
// 	return userID, nil
// }

// func (h *UserHandler) LogoutSession(ctx *gin.Context) {
// 	sess := sessions.Default(ctx)
// 	sess.Options(sessions.Options{
// 		MaxAge: -1,
// 	})
// 	sess.Save()
// }

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
