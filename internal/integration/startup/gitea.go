package startup

import (
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
)

func NewDummyGiteaHandler(userSvc service.UserService, hdl ijwt.Handler) *web.OAuth2GiteaHandler {
	dummySvc := gitea.NewService("dummy", "dummy", "dummy")
	return web.NewOAuth2GiteaHandler(dummySvc, userSvc, hdl)
}
