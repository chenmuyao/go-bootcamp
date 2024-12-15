package gitea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/pkg/httpx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

// {{{ Consts

const (
	// url, clientId, redirectURI, state
	authURLPattern     = "https://%s/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s"
	userAPIPattern     = "https://%s/api/v1/user"
	accessTokenPattern = "https://%s/login/oauth/access_token"
	grantType          = "authorization_code"
)

// }}}
// {{{ Global Varirables

var redirectURI = url.PathEscape("http://localhost:8081/oauth2/gitea/callback")

// }}}
// {{{ Interface

type Service interface {
	AuthURL(ctx context.Context, state string) string
	VerifyCode(ctx context.Context, code string) (domain.GiteaInfo, error)
}

// }}}
// {{{ Struct

type service struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	l            logger.Logger
}

func NewService(baseURL string, clientID string, clientSecret string, l logger.Logger) Service {
	return &service{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   http.DefaultClient,
		l:            l,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (s *service) AuthURL(ctx context.Context, state string) string {
	s.l.Info("Gitea Auth", logger.Field{
		Key:   "state",
		Value: state,
	})
	return fmt.Sprintf(authURLPattern, s.baseURL, s.clientID, redirectURI, state)
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.GiteaInfo, error) {
	type Body struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		GrantType    string `json:"grant_type"`
	}
	body := Body{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		Code:         code,
		GrantType:    grantType,
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&body)
	if err != nil {
		s.l.Error("json encode error")
		return domain.GiteaInfo{}, err
	}

	accessTokenURL := fmt.Sprintf(accessTokenPattern, s.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accessTokenURL, &buf)
	if err != nil {
		s.l.Error("access token query construction error")
		return domain.GiteaInfo{}, err
	}

	req.Header.Add(httpx.Accept, httpx.ApplicationJSON)
	req.Header.Add(httpx.ContentType, httpx.ApplicationJSON)

	httpResp, err := s.httpClient.Do(req)
	if err != nil {
		s.l.Error("http client do error")
		return domain.GiteaInfo{}, err
	}
	if httpResp.StatusCode != http.StatusOK {
		return domain.GiteaInfo{}, errors.New("get access token error")
	}

	type Response struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		s.l.Error("resp decode error")
		return domain.GiteaInfo{}, err
	}

	// Get User Info
	apiURL := fmt.Sprintf(userAPIPattern, s.baseURL)
	apiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		s.l.Error("api request error")
		return domain.GiteaInfo{}, err
	}
	apiReq.Header.Add(httpx.Accept, httpx.ApplicationJSON)
	apiReq.Header.Add(httpx.Authorization, fmt.Sprintf("token %s", resp.AccessToken))
	apiReq.Header.Add(httpx.ContentType, httpx.ApplicationJSON)

	apiResp, err := s.httpClient.Do(apiReq)
	if err != nil {
		s.l.Error("api http error")
		return domain.GiteaInfo{}, err
	}
	if apiResp.StatusCode != http.StatusOK {
		return domain.GiteaInfo{}, errors.New("api query error")
	}

	var giteaInfo domain.GiteaInfo
	err = json.NewDecoder(apiResp.Body).Decode(&giteaInfo)
	if err != nil {
		s.l.Error("api decode error", logger.Field{Key: "apiResp", Value: apiResp})
		return domain.GiteaInfo{}, err
	}

	return giteaInfo, err
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
