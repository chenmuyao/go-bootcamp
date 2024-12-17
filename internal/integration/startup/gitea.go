package startup

import (
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

func NewDummyGiteaHandler(
	userSvc service.UserService,
	hdl ijwt.Handler,
	l logger.Logger,
) *web.OAuth2GiteaHandler {
	dummySvc := gitea.NewService("dummy", "dummy", "dummy", l)
	return web.NewOAuth2GiteaHandler(l, dummySvc, userSvc, hdl)
}
