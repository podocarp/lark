package lark

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type messageSendHTTPClient struct {
	request *http.Request
}

func (c *messageSendHTTPClient) Do(_ context.Context, request *http.Request) (*http.Response, error) {
	c.request = request
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"code":0,"data":{"message_id":"om_message"}}`)),
		Header:     make(http.Header),
	}, nil
}

func TestSendRawMessageAsUserUsesUserAccessToken(t *testing.T) {
	httpClient := &messageSendHTTPClient{}
	client := New(WithHttpClient(httpClient))

	response, _, err := client.Message.SendRawMessageAsUser(context.Background(), &SendRawMessageReq{
		ReceiveIDType: IDTypeChatID,
		ReceiveID:     "oc_message",
		MsgType:       MsgTypeText,
		Content:       `{"text":"hello"}`,
	}, WithUserAccessToken("user-access-token"))

	require.NoError(t, err)
	require.Equal(t, "om_message", response.MessageID)
	require.Equal(t, "Bearer user-access-token", httpClient.request.Header.Get(httpRequestHeaderAuthorization))
	require.Equal(t, "chat_id", httpClient.request.URL.Query().Get("receive_id_type"))
}

func TestSendRawMessageAsUserRequiresUserAccessToken(t *testing.T) {
	client := New()

	_, _, err := client.Message.SendRawMessageAsUser(context.Background(), &SendRawMessageReq{})

	require.EqualError(t, err, "user access token is required")
}

func TestSendRawMessageUsesTenantAccessToken(t *testing.T) {
	httpClient := &messageSendHTTPClient{}
	client := New(WithHttpClient(httpClient))
	client.Mock().MockGetTenantAccessToken(func(context.Context) (*TokenExpire, *Response, error) {
		return &TokenExpire{Token: "tenant-access-token"}, nil, nil
	})
	defer client.Mock().UnMockGetTenantAccessToken()

	response, _, err := client.Message.SendRawMessage(context.Background(), &SendRawMessageReq{
		ReceiveIDType: IDTypeChatID,
		ReceiveID:     "oc_message",
		MsgType:       MsgTypeText,
		Content:       `{"text":"hello"}`,
	})

	require.NoError(t, err)
	require.Equal(t, "om_message", response.MessageID)
	require.Equal(t, "Bearer tenant-access-token", httpClient.request.Header.Get(httpRequestHeaderAuthorization))
}
