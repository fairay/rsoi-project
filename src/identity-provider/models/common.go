package models

import (
	"encoding/json"
	"fmt"
	"identity-provider/objects"
	"identity-provider/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AuthM struct {
	client *http.Client
}

func NewAuthM(client *http.Client) *AuthM {
	return &AuthM{client: client}
}

type Models struct {
	Auth *AuthM
}

func InitModels() *Models {
	models := new(Models)
	client := &http.Client{}

	models.Auth = NewAuthM(client)
	return models
}

func (model *AuthM) Auth(username string, password string) (*objects.AuthResponse, error) {
	authRequest := url.Values{}
	authRequest.Set("scope", "openid")
	authRequest.Set("grant_type", "password")
	authRequest.Set("username", username)
	authRequest.Set("password", password)
	authRequest.Set("client_id", utils.Config.OktaClientId)
	authRequest.Set("client_secret", utils.Config.OktaClientSecret)
	encodedData := authRequest.Encode()

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/default/v1/token", utils.Config.OktaEndpoint),
		strings.NewReader(encodedData),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(authRequest.Encode())))

	resp, err := model.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data := &objects.AuthResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth failed, code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
