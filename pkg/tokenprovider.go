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
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//go:generate mockgen -destination=tokenprovider_generated.go -package=pkg . TokenProviderImpl
type TokenProviderImpl interface {
	RunE(ctx context.Context) error
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

func NewTokenProvider(client *http.Client, creds Creds) TokenProvider {
	return TokenProvider{
		client: client,
		cloud:  &TokenInfo{},
		creds:  creds,
	}
}

func (p TokenProvider) RunE(ctx context.Context) error {
	logging.FromContext(ctx).Debugf("Running token provider for %s", p.creds.Endpoint())
	for {
		select {
		case <-ctx.Done():
			logging.FromContext(ctx).Debugf("Stopping token provider for %s", p.creds.Endpoint())
			return ctx.Err()
		case <-time.After(time.Until(p.cloud.Expiry)):
			if _, err := p.RefreshToken(ctx); err != nil {
				logging.FromContext(ctx).Errorf("Unable to refresh token: %s", err.Error())
				return err
			}
		}
	}
}

func (p TokenProvider) AccessToken(ctx context.Context) (*TokenInfo, error) {
	p.cloud.Lock()
	defer p.cloud.Unlock()

	logger := logging.FromContext(ctx).WithField("func", "AccessToken")
	logger.Debugf("Getting access token for %s", p.creds.Endpoint())
	defer logger.Debugf("Getting access token done")

	client := &http.Client{
		Transport: otlp.NewRoundTripper(http.DefaultTransport, true),
	}

	rp, err := rp.NewRelyingPartyOIDC(ctx, p.creds.Endpoint(), p.creds.ClientId(), p.creds.ClientSecret(), "", []string{
		"openid", "email", "offline_access", "supertoken",
	}, rp.WithHTTPClient(client))
	if err != nil {
		logger.Errorf("Unable to create OIDC client: %s", err.Error())
		return nil, err
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
		return nil, err
	}

	p.cloud.AccessToken = t.AccessToken
	p.cloud.Expiry = t.Expiry
	p.cloud.RefreshToken = t.RefreshToken

	return nil, nil
}

func (p TokenProvider) RefreshToken(ctx context.Context) (*TokenInfo, error) {
func (p TokenProvider) RefreshToken(ctx context.Context) (*TokenInfo, error) {
	logging.FromContext(ctx).Debugf("Getting refresh token for %s", p.creds.Endpoint())
	if p.cloud.AccessToken == "" {
		return p.AccessToken(ctx)
	}

	if time.Now().Before(p.cloud.Expiry) {
		return p.cloud, nil
	}

	p.cloud.Lock()
	defer p.cloud.Unlock()

	form := url.Values{
		"grant_type":    []string{string(oidc.GrantTypeRefreshToken)},
		"refresh_token": []string{p.cloud.RefreshToken},
		"client_id":     []string{p.creds.ClientId()},
		"client_secret": []string{p.creds.ClientSecret()},
	}

	discoveryConfiguration, err := client.Discover(ctx, p.creds.Endpoint(), http.DefaultClient)
	if err != nil {
		return nil, err
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, discoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
		return nil, err
	}
	req.SetBasicAuth(p.creds.ClientId(), p.creds.ClientSecret())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err := p.client.Do(req)
	if err != nil {
		return nil, err
		return nil, err
	}

	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			return nil, err
			return nil, err
		}
		return nil, errors.New(string(data))
		return nil, errors.New(string(data))
	}

	token := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&token); err != nil {
		return nil, err
		return nil, err
	}

	p.cloud.Lock()
	defer p.cloud.Unlock()
	p.cloud.AccessToken = token.AccessToken
	p.cloud.Expiry = token.Expiry
	p.cloud.RefreshToken = token.RefreshToken

	return p.cloud, nil
	return p.cloud, nil

}
