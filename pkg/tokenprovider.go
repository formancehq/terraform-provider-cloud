package pkg

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg/otlp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2/clientcredentials"
)

//go:generate mockgen -typed -destination=tokenprovider_generated.go -package=pkg . TokenProviderImpl
type TokenProviderImpl interface {
	AccessToken(ctx context.Context) (*TokenInfo, error)
	RefreshToken(ctx context.Context) (*TokenInfo, error)
	OrganizationId(ctx context.Context) (string, error)
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

	scopes []string
	opts   []UrlOpts
}

type TokenProviderFactory func(transport http.RoundTripper, creds Creds, scopes []string, opts ...UrlOpts) TokenProviderImpl

type UrlOpts func(url.Values)

func WithResource(resource string) UrlOpts {
	return func(v url.Values) {
		v.Add("resource", resource)
	}
}

func NewTokenProvider(transport http.RoundTripper, creds Creds, scopes []string, opts ...UrlOpts) TokenProviderImpl {
	return TokenProvider{
		client: &http.Client{
			Transport: transport,
		},
		cloud:  &TokenInfo{},
		creds:  creds,
		scopes: scopes,
		opts:   opts,
	}
}

var (
	ScopeCloud = []string{
		"organization:CreateStack",
		"organization:ReadStack",
		"organization:UpdateStack",
		"organization:DeleteStack",
		"organization:UpgradeStack",
		"organization:ListStacks",

		"organization:ReadStackUser",
		"organization:UpdateStackUser",
		"organization:DeleteStackUser",

		"organization:ListStackModules",
		"organization:EnableStackModule",
		"organization:DisableStackModule",

		"organization:ListRegions",
		"organization:ReadRegion",

		"organization:ReadUser",
		"organization:CreateUser",
		"organization:UpdateUser",

		"organization:Read",
	}
	ScopeStack = []string{
		"organization:Read",
	}
)

func (p TokenProvider) AccessToken(ctx context.Context) (*TokenInfo, error) {
	ctx, span := otlp.Tracer.Start(ctx, "AccessToken")
	defer span.End()
	p.cloud.Lock()
	defer p.cloud.Unlock()

	if p.cloud.AccessToken != "" {
		return &TokenInfo{
			AccessToken:  p.cloud.AccessToken,
			RefreshToken: p.cloud.RefreshToken,
			Expiry:       p.cloud.Expiry,
		}, nil
	}

	logger := logging.FromContext(ctx).WithField("operation", "accesstoken")
	logger.Debugf("Getting access token for %s", p.creds.Endpoint())
	defer logger.Debugf("Getting access token done")

	rp, err := rp.NewRelyingPartyOIDC(ctx,
		p.creds.Endpoint(),
		p.creds.ClientId(),
		p.creds.ClientSecret(),
		"",
		p.scopes,
		rp.WithHTTPClient(p.client),
	)
	if err != nil {
		logger.Errorf("Unable to create OIDC client: %s", err.Error())
		return nil, err
	}

	config := &clientcredentials.Config{
		Scopes:         rp.OAuthConfig().Scopes,
		ClientID:       rp.OAuthConfig().ClientID,
		ClientSecret:   rp.OAuthConfig().ClientSecret,
		TokenURL:       rp.OAuthConfig().Endpoint.TokenURL,
		EndpointParams: make(url.Values),
	}

	for _, opt := range p.opts {
		opt(config.EndpointParams)
	}

	t, err := (config).Token(ctx)

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

func (p TokenProvider) OrganizationId(ctx context.Context) (string, error) {
	accessToken, err := p.AccessToken(ctx)
	if err != nil {
		return "", err
	}

	var claims jwt.MapClaims
	_, err = oidc.ParseToken(accessToken.AccessToken, &claims)
	if err != nil {
		return "", err
	}

	organizationId := claims["organization_id"].(string)
	return organizationId, nil

}

func (p TokenProvider) RefreshToken(ctx context.Context) (*TokenInfo, error) {
	ctx, span := otlp.Tracer.Start(ctx, "RefreshToken")
	defer span.End()
	logging.FromContext(ctx).Debugf("Getting refresh token for %s", p.creds.Endpoint())

	tokenInfo, err := p.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	if time.Now().Before(tokenInfo.Expiry) {
		return &TokenInfo{
			AccessToken:  tokenInfo.AccessToken,
			RefreshToken: tokenInfo.RefreshToken,
			Expiry:       tokenInfo.Expiry,
		}, nil
	}

	p.cloud.Lock()
	p.cloud.AccessToken = ""
	p.cloud.Unlock()
	tokenInfo, err = p.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	p.cloud.AccessToken = tokenInfo.AccessToken
	p.cloud.Expiry = tokenInfo.Expiry
	p.cloud.RefreshToken = tokenInfo.RefreshToken

	return &TokenInfo{
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
		Expiry:       tokenInfo.Expiry,
	}, nil

}
