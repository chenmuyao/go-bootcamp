package gitea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/lithammer/shortuuid/v4"
)

type Service interface {
	AuthURL(ctx context.Context) string
	VerifyCode(ctx context.Context, code string) (domain.GiteaInfo, error)
}

var redirectURI = url.PathEscape("http://localhost:8081/oauth2/gitea/callback")

type service struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewService(baseURL string, clientID string, clientSecret string) Service {
	return &service{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   http.DefaultClient,
	}
}

func (s *service) AuthURL(ctx context.Context) string {
	// url, clientId, redirectURI, state
	authURLPattern := "https://%s/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s"

	state := shortuuid.New()

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
		GrantType:    "authorization_code",
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&body)
	if err != nil {
		slog.Error("json encode error")
		return domain.GiteaInfo{}, err
	}

	accessTokenURL := fmt.Sprintf("https://%s/login/oauth/access_token", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accessTokenURL, &buf)
	if err != nil {
		slog.Error("access token query construction error")
		return domain.GiteaInfo{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	httpResp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Error("http client do error")
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
		slog.Error("resp decode error")
		return domain.GiteaInfo{}, err
	}

	// Get User Info
	apiURL := fmt.Sprintf("https://%s/api/v1/user", s.baseURL)
	apiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		slog.Error("api request error")
		return domain.GiteaInfo{}, err
	}
	apiReq.Header.Add("accept", "application/json")
	apiReq.Header.Add("Authorization", fmt.Sprintf("token %s", resp.AccessToken))
	apiReq.Header.Add("Content-Type", "application/json")

	apiResp, err := s.httpClient.Do(apiReq)
	if err != nil {
		slog.Error("api http error")
		return domain.GiteaInfo{}, err
	}
	if apiResp.StatusCode != http.StatusOK {
		return domain.GiteaInfo{}, errors.New("api query error")
	}

	var giteaInfo domain.GiteaInfo
	err = json.NewDecoder(apiResp.Body).Decode(&giteaInfo)
	if err != nil {
		slog.Error("api decode error", "apiResp", apiResp)
		return domain.GiteaInfo{}, err
	}

	return giteaInfo, err
}
