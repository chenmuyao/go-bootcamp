package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	"github.com/gin-gonic/gin"
)

type OAuth2GiteaHandler struct {
	jwtHandler
	svc     gitea.Service
	userSvc service.UserService
}

func NewOAuth2GiteaHandler(svc gitea.Service, userSvc service.UserService) *OAuth2GiteaHandler {
	return &OAuth2GiteaHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (o *OAuth2GiteaHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/gitea")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2GiteaHandler) Auth2URL(ctx *gin.Context) {
	val := o.svc.AuthURL(ctx)

	ctx.JSON(http.StatusOK, Result{
		Code: CodeOK,
		Data: val,
	})
}

func (o *OAuth2GiteaHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	// state := ctx.Query("state")

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

	token, err := o.generateJWTToken(ctx, u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResult)
		return
	}

	redirectURI := fmt.Sprintf(
		"http://localhost:5173/oauth2success?token=%s", // TODO: should not hard code
		token,
	)
	ctx.Redirect(http.StatusPermanentRedirect, redirectURI)
}
