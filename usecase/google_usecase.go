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
)

type googleUsecase struct {
	cfgGoogle  domain.ConfigGoogle
	httpClient *http.Client
}

func NewGoogleUsecase(cfgGoogle domain.ConfigGoogle) domain.GoogleUsecase {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &googleUsecase{
		cfgGoogle:  cfgGoogle,
		httpClient: httpClient,
	}
}

func (u *googleUsecase) GetOAuthUrl() string {
	query := url.Values{}
	query.Add("client_id", u.cfgGoogle.ClientId)
	query.Add("redirect_uri", u.cfgGoogle.RedirectUri)
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
		"client_id":     u.cfgGoogle.ClientId,
		"client_secret": u.cfgGoogle.ClientSecret,
		"redirect_uri":  u.cfgGoogle.RedirectUri,
		"grant_type":    "authorization_code",
	})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	response, err := u.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("cannot get token from Google API,  status: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result domain.GoogleTokenResponse
	if err = json.Unmarshal(data, &result); err != nil {
		return "", err
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
		return nil, err
	}

	response, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("cannot get user from Google API, code: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result domain.GoogleUserResponse
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}
