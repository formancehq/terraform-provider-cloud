package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//go:generate mockgen -destination=tokenprovider_generated.go -package=pkg . TokenProviderImpl
type TokenProviderImpl interface {
	AccessToken(ctx context.Context) (*TokenInfo, error)
	RefreshToken(ctx context.Context) (*TokenInfo, error)
}

type TokenInfo struct {
	sync.Mutex

	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

var (
	_ TokenProviderImpl = &TokenProvider{}
)

type TokenProvider struct {
	client *http.Client

	creds Creds

	cloud *TokenInfo
}

func NewTokenProvider(transport http.RoundTripper, creds Creds) TokenProvider {
	return TokenProvider{
		client: &http.Client{
			Transport: transport,
		},
		cloud: &TokenInfo{},
		creds: creds,
	}
}

func (p TokenProvider) AccessToken(ctx context.Context) (*TokenInfo, error) {
	p.cloud.Lock()
	defer p.cloud.Unlock()

	logger := logging.FromContext(ctx).WithField("func", "AccessToken")
	logger.Debugf("Getting access token for %s", p.creds.Endpoint())
	defer logger.Debugf("Getting access token done")

	rp, err := rp.NewRelyingPartyOIDC(ctx, p.creds.Endpoint(), p.creds.ClientId(), p.creds.ClientSecret(), "", []string{
		"openid", "email", "offline_access", "supertoken",
	}, rp.WithHTTPClient(p.client))
	if err != nil {
		logger.Errorf("Unable to create OIDC client: %s", err.Error())
		return nil, err
	}

	t, err := (&clientcredentials.Config{
		Scopes:       rp.OAuthConfig().Scopes,
		ClientID:     rp.OAuthConfig().ClientID,
		ClientSecret: rp.OAuthConfig().ClientSecret,
		TokenURL:     rp.OAuthConfig().Endpoint.TokenURL,
	}).Token(ctx)

	if err != nil {
		logger.Errorf("Unable to get token: %s", err.Error())
		return nil, err
	}

	p.cloud.AccessToken = t.AccessToken
	p.cloud.Expiry = t.Expiry
	p.cloud.RefreshToken = t.RefreshToken

	return &TokenInfo{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
	}, nil
}

func (p TokenProvider) RefreshToken(ctx context.Context) (*TokenInfo, error) {
	logging.FromContext(ctx).Debugf("Getting refresh token for %s", p.creds.Endpoint())

	p.cloud.Lock()
	if p.cloud.AccessToken == "" {
		p.cloud.Unlock()
		return p.AccessToken(ctx)
	}

	defer p.cloud.Unlock()
	if time.Now().Before(p.cloud.Expiry) {
		return &TokenInfo{
			AccessToken:  p.cloud.AccessToken,
			RefreshToken: p.cloud.RefreshToken,
			Expiry:       p.cloud.Expiry,
		}, nil
	}

	form := url.Values{
		"grant_type":    []string{string(oidc.GrantTypeRefreshToken)},
		"refresh_token": []string{p.cloud.RefreshToken},
		"client_id":     []string{p.creds.ClientId()},
		"client_secret": []string{p.creds.ClientSecret()},
	}

	discoveryConfiguration, err := client.Discover(ctx, p.creds.Endpoint(), http.DefaultClient)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, discoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.creds.ClientId(), p.creds.ClientSecret())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(data))
	}

	token := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&token); err != nil {
		return nil, err
	}

	p.cloud.AccessToken = token.AccessToken
	p.cloud.Expiry = token.Expiry
	p.cloud.RefreshToken = token.RefreshToken

	return &TokenInfo{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}, nil

}
