package ioc

import (
	"os"

	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
)

func InitGiteaService() gitea.Service {
	baseURL := "git.vinchent.xyz"
	clientID, ok := os.LookupEnv("GITEA_CLIENT_ID")
	if !ok {
		panic("Gitea client id not found")
	}
	clientSecret, ok := os.LookupEnv("GITEA_CLIENT_SECRET")
	if !ok {
		panic("Gitea client secret not found")
	}
	return gitea.NewService(baseURL, clientID, clientSecret)
}
