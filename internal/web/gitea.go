package web

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
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
	ijwt.Handler
	svc             gitea.Service
	userSvc         service.UserService
	key             []byte
	stateCookieName string
}

func NewOAuth2GiteaHandler(
	svc gitea.Service,
	userSvc service.UserService,
	hdl ijwt.Handler,
) *OAuth2GiteaHandler {
	return &OAuth2GiteaHandler{
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
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2GiteaHandler) Auth2URL(ctx *gin.Context) {
	state := shortuuid.New()
	val := o.svc.AuthURL(ctx, state)

	err := o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: CodeOK,
		Data: val,
	})
}

func (o *OAuth2GiteaHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: CodeUserSide,
			Msg:  "OAuth authentication failed",
		})
		return
	}

	code := ctx.Query("code")

	giteaInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		slog.Error("wrong auth", "msg", err)
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: CodeUserSide,
			Msg:  "Wrong authentication code",
		})
		return
	}

	u, err := o.userSvc.FindOrCreateByGitea(ctx, giteaInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	ssid := uuid.New().String()
	refreshToken, err := o.GenerateRefreshToken(u.ID, ssid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	token, err := o.GenerateJWTToken(ctx, u.ID, ssid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	redirectURI := fmt.Sprintf(redirectPattern, token, refreshToken)
	ctx.Redirect(http.StatusPermanentRedirect, redirectURI)
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