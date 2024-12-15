package ioc

import (
	"log/slog"

	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

func InitGiteaService(l logger.Logger) gitea.Service {
	baseURL := config.Cfg.OAuth2.BaseURL
	clientID := config.Cfg.OAuth2.ClientID
	if clientID == "" {
		slog.Error("Gitea client id not found")
	}
	clientSecret := config.Cfg.OAuth2.ClientSecret
	if clientSecret == "" {
		slog.Error("Gitea client secret not found")
	}
	return gitea.NewService(baseURL, clientID, clientSecret, l)
}
