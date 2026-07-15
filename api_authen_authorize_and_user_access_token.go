package lark

import (
	"context"
	"fmt"
	"net/url"
)

// GenAuthorizeURL builds the current web authorization URL.
//
// This method only builds the URL. It does not launch a browser or handle callbacks.
func (r *AuthenService) GenAuthorizeURL(ctx context.Context, request *GenAuthorizeURLReq) string {
	_ = ctx

	if request == nil {
		request = &GenAuthorizeURLReq{}
	}

	clientID := request.ClientID
	if clientID == "" {
		clientID = r.cli.appID
	}

	responseType := request.ResponseType
	if responseType == "" {
		responseType = "code"
	}

	query := url.Values{}
	query.Set("client_id", clientID)
	query.Set("response_type", responseType)
	query.Set("redirect_uri", request.RedirectURI)

	if request.Scope != "" {
		query.Set("scope", request.Scope)
	}
	if request.State != "" {
		query.Set("state", request.State)
	}
	if request.CodeChallenge != "" {
		query.Set("code_challenge", request.CodeChallenge)
	}
	if request.CodeChallengeMethod != "" {
		query.Set("code_challenge_method", request.CodeChallengeMethod)
	}
	if request.Prompt != "" {
		query.Set("prompt", request.Prompt)
	}

	return fmt.Sprintf("%s/open-apis/authen/v1/authorize?%s", r.cli.oauthBaseURL(), query.Encode())
}

// CreateUserAccessToken exchanges an authorization code for a user_access_token.
func (r *AuthenService) CreateUserAccessToken(ctx context.Context, request *CreateUserAccessTokenReq, options ...MethodOptionFunc) (*CreateUserAccessTokenResp, *Response, error) {
	if r.cli.mock.mockAuthenCreateUserAccessToken != nil {
		r.cli.Log(ctx, LogLevelDebug, "[lark] Authen#CreateUserAccessToken mock enable")
		return r.cli.mock.mockAuthenCreateUserAccessToken(ctx, request, options...)
	}

	req := &RawRequestReq{
		Scope:              "Authen",
		API:                "CreateUserAccessToken",
		Method:             "POST",
		URL:                r.cli.openBaseURL + "/open-apis/authen/v1/access_token",
		Body:               request,
		MethodOption:       newMethodOption(options),
		NeedAppAccessToken: true,
	}
	resp := new(createUserAccessTokenResp)

	response, err := r.cli.RawRequest(ctx, req, resp)
	return resp.Data, response, err
}

// MockAuthenCreateUserAccessToken mock AuthenCreateUserAccessToken method
func (r *Mock) MockAuthenCreateUserAccessToken(f func(ctx context.Context, request *CreateUserAccessTokenReq, options ...MethodOptionFunc) (*CreateUserAccessTokenResp, *Response, error)) {
	r.mockAuthenCreateUserAccessToken = f
}

// UnMockAuthenCreateUserAccessToken un-mock AuthenCreateUserAccessToken method
func (r *Mock) UnMockAuthenCreateUserAccessToken() {
	r.mockAuthenCreateUserAccessToken = nil
}

// RefreshUserAccessToken refreshes a user_access_token.
func (r *AuthenService) RefreshUserAccessToken(ctx context.Context, request *RefreshUserAccessTokenReq, options ...MethodOptionFunc) (*RefreshUserAccessTokenResp, *Response, error) {
	if r.cli.mock.mockAuthenRefreshUserAccessToken != nil {
		r.cli.Log(ctx, LogLevelDebug, "[lark] Authen#RefreshUserAccessToken mock enable")
		return r.cli.mock.mockAuthenRefreshUserAccessToken(ctx, request, options...)
	}

	req := &RawRequestReq{
		Scope:              "Authen",
		API:                "RefreshUserAccessToken",
		Method:             "POST",
		URL:                r.cli.openBaseURL + "/open-apis/authen/v1/refresh_access_token",
		Body:               request,
		MethodOption:       newMethodOption(options),
		NeedAppAccessToken: true,
	}
	resp := new(refreshUserAccessTokenResp)

	response, err := r.cli.RawRequest(ctx, req, resp)
	return resp.Data, response, err
}

// MockAuthenRefreshUserAccessToken mock AuthenRefreshUserAccessToken method
func (r *Mock) MockAuthenRefreshUserAccessToken(f func(ctx context.Context, request *RefreshUserAccessTokenReq, options ...MethodOptionFunc) (*RefreshUserAccessTokenResp, *Response, error)) {
	r.mockAuthenRefreshUserAccessToken = f
}

// UnMockAuthenRefreshUserAccessToken un-mock AuthenRefreshUserAccessToken method
func (r *Mock) UnMockAuthenRefreshUserAccessToken() {
	r.mockAuthenRefreshUserAccessToken = nil
}

// GenAuthorizeURLReq ...
type GenAuthorizeURLReq struct {
	ClientID            string `json:"client_id,omitempty"`
	ResponseType        string `json:"response_type,omitempty"`
	RedirectURI         string `json:"redirect_uri,omitempty"`
	Scope               string `json:"scope,omitempty"`
	State               string `json:"state,omitempty"`
	CodeChallenge       string `json:"code_challenge,omitempty"`
	CodeChallengeMethod string `json:"code_challenge_method,omitempty"`
	Prompt              string `json:"prompt,omitempty"`
}

// CreateUserAccessTokenReq ...
type CreateUserAccessTokenReq struct {
	GrantType string `json:"grant_type,omitempty"`
	Code      string `json:"code,omitempty"`
}

// CreateUserAccessTokenResp ...
type CreateUserAccessTokenResp struct {
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

// RefreshUserAccessTokenReq ...
type RefreshUserAccessTokenReq struct {
	GrantType    string `json:"grant_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// RefreshUserAccessTokenResp ...
type RefreshUserAccessTokenResp struct {
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type createUserAccessTokenResp struct {
	Code  int64                      `json:"code,omitempty"`
	Msg   string                     `json:"msg,omitempty"`
	Data  *CreateUserAccessTokenResp `json:"data,omitempty"`
	Error *ErrorDetail               `json:"error,omitempty"`
}

type refreshUserAccessTokenResp struct {
	Code  int64                       `json:"code,omitempty"`
	Msg   string                      `json:"msg,omitempty"`
	Data  *RefreshUserAccessTokenResp `json:"data,omitempty"`
	Error *ErrorDetail                `json:"error,omitempty"`
}
