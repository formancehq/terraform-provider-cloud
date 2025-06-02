package stack

import (
	"github.com/formancehq/terraform-provider/internal"
	"github.com/formancehq/terraform-provider/internal/membership"
)

type TokenProvider struct {
	*membership.TokenProvider
	stack *internal.TokenInfo
}

func NewTokenProvider(cloudprovider *membership.TokenProvider) TokenProvider {
	return TokenProvider{
		TokenProvider: cloudprovider,
		stack:         &internal.TokenInfo{},
	}
}

type Stack struct {
	Id             string `json:"stack_id"`
	OrganizationId string `json:"organization_id"`
	Url            string `json:"url"`
}

// func (p TokenProvider) StackSecurityToken(ctx context.Context, stack Stack) error {
// 	p.cloud.Lock()
// 	defer p.cloud.Unlock()

// 	form := url.Values{
// 		"grant_type":         []string{string(oidc.GrantTypeTokenExchange)},
// 		"audience":           []string{fmt.Sprintf("stack://%s/%s", stack.OrganizationId, stack.Id)},
// 		"subject_token":      []string{p.cloud.AccessToken},
// 		"subject_token_type": []string{"urn:ietf:params:oauth:token-type:access_token"},
// 	}

// 	membershipDiscoveryConfiguration, err := client.Discover(ctx, p.Endpoint(), p.client)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, membershipDiscoveryConfiguration.TokenEndpoint,
// 		bytes.NewBufferString(form.Encode()))
// 	if err != nil {
// 		return err
// 	}
// 	req.SetBasicAuth(p.ClientId(), p.ClientSecret())
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	ret, err := p.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	if ret.StatusCode != http.StatusOK {
// 		data, err := io.ReadAll(ret.Body)
// 		if err != nil {
// 			panic(err)
// 		}
// 		return errors.New(string(data))
// 	}

// 	securityToken := oauth2.Token{}
// 	if err := json.NewDecoder(ret.Body).Decode(&securityToken); err != nil {
// 		return err
// 	}

// 	// apiUrl := p.ApiUrl(stack, "auth")
// 	// form = url.Values{
// 	// 	"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
// 	// 	"assertion":  []string{securityToken.AccessToken},
// 	// 	"scope":      []string{"openid email"},
// 	// }

// 	// stackDiscoveryConfiguration, err := client.Discover(ctx, apiUrl.String(), p.client)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// req, err = http.NewRequestWithContext(ctx, http.MethodPost, stackDiscoveryConfiguration.TokenEndpoint,
// 	// 	bytes.NewBufferString(form.Encode()))
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	// ret, err = p.client.Do(req)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// if ret.StatusCode != http.StatusOK {
// 	// 	data, err := io.ReadAll(ret.Body)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	return errors.New(string(data))
// 	// }

// 	// stackToken := oauth2.Token{}
// 	// if err := json.NewDecoder(ret.Body).Decode(&stackToken); err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }
