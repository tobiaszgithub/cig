package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/beeekind/go-authhttp"
	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
	"golang.org/x/oauth2/clientcredentials"
)

func GetClient(conf config.Configuration) *http.Client {

	ctx := context.Background()

	var client *http.Client

	if conf.Authorization.Type == "oauth" {
		oauthConf := clientcredentials.Config{
			ClientID:     conf.Authorization.ClientID,
			ClientSecret: conf.Authorization.ClientSecret,
			Scopes:       []string{},
			TokenURL:     conf.Authorization.TokenURL,
		}

		client = oauthConf.Client(ctx)
	}

	if conf.Authorization.Type == "basic" {
		client = authhttp.NewHTTPClient(authhttp.WithBasicAuth(conf.Authorization.Username,
			conf.Authorization.Password))
	}

	return client

}

func GetIntegrationPackages(conf config.Configuration) (*model.IPResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages"

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	rawRes, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer rawRes.Body.Close()

	var decodedRes model.IPResponse

	if err := json.NewDecoder(rawRes.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err

}
