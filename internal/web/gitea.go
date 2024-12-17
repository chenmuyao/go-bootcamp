package web

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

// {{{ Consts

const (
	redirectPattern = "http://localhost:5173/oauth2success?token=%s&refresh_token=%s"
	cookieDomain    = "/oauth2/gitea/callback"
	cookieMaxAge    = 600
)

// }}}
// {{{ Global Varirables

var OAuthJWTKey = []byte("xQePmbb2TP9CUyFZkgOnV3JQdr22ZNBx")

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type OAuth2GiteaHandler struct {
	l logger.Logger
	ijwt.Handler
	svc             gitea.Service
	userSvc         service.UserService
	key             []byte
	stateCookieName string
}

func NewOAuth2GiteaHandler(
	l logger.Logger,
	svc gitea.Service,
	userSvc service.UserService,
	hdl ijwt.Handler,
) *OAuth2GiteaHandler {
	return &OAuth2GiteaHandler{
		l:               l,
		svc:             svc,
		userSvc:         userSvc,
		key:             OAuthJWTKey,
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}

// }}}
// {{{ Other structs

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

// }}}
// {{{ Struct Methods

func (o *OAuth2GiteaHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/gitea")
	g.GET("/authurl", ginx.WrapLog(o.l, o.Auth2URL))
	g.Any("/callback", ginx.WrapLog(o.l, o.Callback))
}

func (o *OAuth2GiteaHandler) Auth2URL(ctx *gin.Context) (ginx.Result, error) {
	state := shortuuid.New()
	val := o.svc.AuthURL(ctx, state)

	err := o.setStateCookie(ctx, state)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}

	return ginx.Result{
		Code: ginx.CodeOK,
		Data: val,
	}, nil
}

func (o *OAuth2GiteaHandler) Callback(ctx *gin.Context) (ginx.Result, error) {
	err := o.verifyState(ctx)
	if err != nil {
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "OAuth authentication failed",
		}, fmt.Errorf("OAuth2 state verification failed: %w", err)
	}

	code := ctx.Query("code")

	giteaInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		slog.Error("wrong auth", "msg", err)
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "Wrong authentication code",
		}, fmt.Errorf("OAuth2 code verification failed: %w", err)
	}

	u, err := o.userSvc.FindOrCreateByGitea(ctx, giteaInfo)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}

	ssid := uuid.New().String()
	refreshToken, err := o.GenerateRefreshToken(u.ID, ssid)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}

	token, err := o.GenerateJWTToken(ctx, u.ID, ssid)
	if err != nil {
		return ginx.InternalServerErrorResult, err
	}

	redirectURI := fmt.Sprintf(redirectPattern, token, refreshToken)
	ctx.Redirect(http.StatusPermanentRedirect, redirectURI)

	return ginx.Result{}, nil
}

func (o *OAuth2GiteaHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {
		slog.Error("token string generate error", "err", err)
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, cookieMaxAge, cookieDomain, "", false, true)
	return nil
}

func (o *OAuth2GiteaHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return err
	}

	var sc StateClaims

	_, err = jwt.ParseWithClaims(ck, &sc, func(*jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("%w, cookie invalid", err)
	}
	if sc.State != state {
		return errors.New("state is modified")
	}
	return nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
