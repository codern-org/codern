package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
)

type googleUsecase struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewGoogleUsecase(cfg *config.Config) domain.GoogleUsecase {
	return &googleUsecase{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (u *googleUsecase) GetOAuthUrl() string {
	query := url.Values{}
	query.Add("client_id", u.cfg.Google.ClientId)
	query.Add("redirect_uri", u.cfg.Google.RedirectUri)
	query.Add("response_type", "code")
	query.Add("prompt", "consent")
	query.Add("scope", strings.Join([]string{
		"openid",
		"email",
		"profile",
	}, " "))
	return "https://accounts.google.com/o/oauth2/v2/auth?" + query.Encode()
}

func (u *googleUsecase) GetToken(code string) (string, error) {
	body, err := json.Marshal(&map[string]string{
		"code":          code,
		"client_id":     u.cfg.Google.ClientId,
		"client_secret": u.cfg.Google.ClientSecret,
		"redirect_uri":  u.cfg.Google.RedirectUri,
		"grant_type":    "authorization_code",
	})
	if err != nil {
		return "", errs.New(errs.ErrGoogleAuth, "cannot construct auth payload", err)
	}

	request, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewReader(body))
	if err != nil {
		return "", errs.New(errs.ErrGoogleAuth, "cannot request to google api", err)
	}

	response, err := u.httpClient.Do(request)
	if err != nil {
		return "", errs.New(errs.ErrGoogleAuth, "cannot request to google api", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errs.New(errs.ErrGoogleAuth, "cannot get token from google api with status %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errs.New(errs.ErrGoogleAuth, "cannot read response from google api", err)
	}

	var result domain.GoogleTokenResponse
	if err = json.Unmarshal(data, &result); err != nil {
		return "", errs.New(errs.ErrGoogleAuth, "cannot read response from google api", err)
	}
	return result.AccessToken, nil
}

func (u *googleUsecase) GetUser(accessToken string) (*domain.GoogleUserResponse, error) {
	query := url.Values{}
	query.Add("alt", "json")
	query.Add("access_token", accessToken)
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?%s", query.Encode())

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errs.New(errs.ErrGoogleAuth, "cannot request to google api", err)
	}

	response, err := u.httpClient.Do(request)
	if err != nil {
		return nil, errs.New(errs.ErrGoogleAuth, "cannot request to google api", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errs.New(errs.ErrGoogleAuth, "cannot get user from google api with status %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.New(errs.ErrGoogleAuth, "cannot read response from google api", err)
	}

	var result domain.GoogleUserResponse
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, errs.New(errs.ErrGoogleAuth, "cannot read response from google api", err)
	}
	return &result, nil
}
