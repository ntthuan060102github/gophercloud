package tokens

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// Scope allows a created token to be limited to a specific domain or project.
type Scope struct {
	ProjectID   string
	ProjectName string
	DomainID    string
	DomainName  string
	System      bool
	TrustID     string
}

// AuthOptionsBuilder provides the ability for extensions to add additional
// parameters to AuthOptions. Extensions must satisfy all required methods.
type AuthOptionsBuilder interface {
	// ToTokenV3CreateMap assembles the Create request body, returning an error
	// if parameters are missing or inconsistent.
	ToTokenV3CreateMap(map[string]any) (map[string]any, error)
	ToTokenV3HeadersMap(map[string]any) (map[string]string, error)
	ToTokenV3ScopeMap() (map[string]any, error)
	CanReauth() bool
}

// AuthOptions represents options for authenticating a user.
type AuthOptions struct {
	// IdentityEndpoint specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed
	// by all of the identity services, it will often be populated by a
	// provider-level function.
	IdentityEndpoint string `json:"-"`

	// Username is required if using Identity V2 API. Consult with your provider's
	// control panel to discover your account's username. In Identity V3, either
	// UserID or a combination of Username and DomainID or DomainName are needed.
	Username string `json:"username,omitempty"`
	UserID   string `json:"id,omitempty"`

	Password string `json:"password,omitempty"`

	// Passcode is used in TOTP authentication method
	Passcode string `json:"passcode,omitempty"`

	// At most one of DomainID and DomainName must be provided if using Username
	// with Identity V3. Otherwise, either are optional.
	DomainID   string `json:"-"`
	DomainName string `json:"name,omitempty"`

	// AllowReauth should be set to true if you grant permission for Gophercloud
	// to cache your credentials in memory, and to allow Gophercloud to attempt
	// to re-authenticate automatically if/when your token expires.  If you set
	// it to false, it will not cache these settings, but re-authentication will
	// not be possible.  This setting defaults to false.
	AllowReauth bool `json:"-"`

	// TokenID allows users to authenticate (possibly as another user) with an
	// authentication token ID.
	TokenID string `json:"-"`

	// Authentication through Application Credentials requires supplying name, project and secret
	// For project we can use TenantID
	ApplicationCredentialID     string `json:"-"`
	ApplicationCredentialName   string `json:"-"`
	ApplicationCredentialSecret string `json:"-"`

	Scope Scope `json:"-"`
}

// ToTokenV3CreateMap builds a request body from AuthOptions.
func (opts *AuthOptions) ToTokenV3CreateMap(scope map[string]any) (map[string]any, error) {
	gophercloudAuthOpts := gophercloud.AuthOptions{
		Username:                    opts.Username,
		UserID:                      opts.UserID,
		Password:                    opts.Password,
		Passcode:                    opts.Passcode,
		DomainID:                    opts.DomainID,
		DomainName:                  opts.DomainName,
		AllowReauth:                 opts.AllowReauth,
		TokenID:                     opts.TokenID,
		ApplicationCredentialID:     opts.ApplicationCredentialID,
		ApplicationCredentialName:   opts.ApplicationCredentialName,
		ApplicationCredentialSecret: opts.ApplicationCredentialSecret,
	}

	return gophercloudAuthOpts.ToTokenV3CreateMap(scope)
}

// ToTokenV3ScopeMap builds a scope request body from AuthOptions.
func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]any, error) {
	scope := gophercloud.AuthScope(opts.Scope)

	gophercloudAuthOpts := gophercloud.AuthOptions{
		Scope:      &scope,
		DomainID:   opts.DomainID,
		DomainName: opts.DomainName,
	}

	return gophercloudAuthOpts.ToTokenV3ScopeMap()
}

func (opts *AuthOptions) CanReauth() bool {
	if opts.Passcode != "" {
		// cannot reauth using TOTP passcode
		return false
	}

	return opts.AllowReauth
}

// ToTokenV3HeadersMap allows AuthOptions to satisfy the AuthOptionsBuilder
// interface in the v3 tokens package.
func (opts *AuthOptions) ToTokenV3HeadersMap(map[string]any) (map[string]string, error) {
	return nil, nil
}

func subjectTokenHeaders(subjectToken string) map[string]string {
	return map[string]string{
		"X-Subject-Token": subjectToken,
	}
}

// Create authenticates and either generates a new token, or changes the Scope
// of an existing token.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts AuthOptionsBuilder) (r CreateResult) {
	scope, err := opts.ToTokenV3ScopeMap()
	if err != nil {
		r.Err = err
		return
	}

	b, err := opts.ToTokenV3CreateMap(scope)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Post(ctx, tokenURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OmitHeaders: []string{"X-Auth-Token"},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get validates and retrieves information about another token.
func Get(ctx context.Context, c *gophercloud.ServiceClient, token string) (r GetResult) {
	resp, err := c.Get(ctx, tokenURL(c), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(token),
		OkCodes:     []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Validate determines if a specified token is valid or not.
func Validate(ctx context.Context, c *gophercloud.ServiceClient, token string) (bool, error) {
	resp, err := c.Head(ctx, tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(token),
		OkCodes:     []int{200, 204, 404},
	})
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200 || resp.StatusCode == 204, nil
}

// Revoke immediately makes specified token invalid.
func Revoke(ctx context.Context, c *gophercloud.ServiceClient, token string) (r RevokeResult) {
	resp, err := c.Delete(ctx, tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(token),
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
