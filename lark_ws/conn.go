package lark_ws

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/chyroc/lark"
)

func (c *Client) connect(ctx context.Context) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return
	}

	endpoint, err := c.getEndpoint(ctx)
	if err != nil {
		return
	}
	if err := c.saveURL(endpoint.URL); err != nil {
		return err
	}
	c.saveClientConfig(endpoint.ClientConfig)

	// Dial the (possibly host-rewritten) saved URL. If a connection host
	// override was configured, saveURL swapped in the internal host so we do
	// not hit the public msg-frontier endpoint.
	conn, resp, err := c.wsDialer.Dial(c.connURL.String(), nil)
	if err != nil && resp == nil {
		return err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return parseWebsocketErr("Callback", "ConnWebsocket", resp)
	}
	c.conn = conn

	c.logInfo(ctx, "connected to %s", c.connURL)

	go c.receiveMessageLoop(ctx)

	return
}

func (c *Client) disconnect(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return
	}

	_ = c.conn.Close()
	c.conn = nil
	c.connURL = nil
	c.connID = ""
	c.serviceID = ""
	c.logInfo(ctx, "disconnected to %s", c.connURL)
}

const (
	codeForbidden       = 403
	codeAuthFailed      = 514
	codeExceedConnLimit = 1000040350
)

func isRetryErr(err error) bool {
	var e *lark.Error
	if errors.As(err, &e) {
		return e.Code != codeForbidden && e.Code != codeExceedConnLimit
	}
	return true
}

func parseWebsocketErr(scope, funcName string, resp *http.Response) error {
	code, _ := strconv.ParseInt(resp.Header.Get("Handshake-Status"), 10, 64)
	msg := resp.Header.Get("Handshake-Msg")
	switch code {
	case codeAuthFailed:
		authCode, _ := strconv.ParseInt(resp.Header.Get("Handshake-Autherrcode"), 10, 64)
		if authCode != 0 {
			return lark.NewError(scope, funcName, authCode, msg)
		}
		return lark.NewError(scope, funcName, code, msg)
	default:
		return lark.NewError(scope, funcName, code, msg)
	}
}
