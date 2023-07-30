package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
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
	query.Add("redirect_uri", u.cfgGoogle.RedirectUri)
	query.Add("client_id", u.cfgGoogle.ClientId)
	query.Add("access_type", "offline")
	query.Add("response_type", "code")
	query.Add("prompt", "consent")
	query.Add("scope", strings.Join([]string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
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
		return "", errors.New("cannot get token from google api, status code: " + response.Status)
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
	request, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	response, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("cannot get user from google api, status code: " + response.Status)
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
