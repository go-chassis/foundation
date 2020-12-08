package httpclient_test

import (
	"context"
	"os"

	"github.com/go-chassis/foundation/httpclient"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpDo(t *testing.T) {
	os.Setenv("HTTP_DEBUG", "1")
	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.Requests)
	uc.Client = htc

	htServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("some thing"))
	}))

	resp, err := uc.Get(context.Background(), htServer.URL, nil)
	assert.NotNil(t, resp)
	t.Log(resp)
	assert.NoError(t, err)
}

func TestHttpDoHeadersNil(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.Requests)
	uc.Client = htc

	resp, err := uc.Do(context.Background(), "GET", "https://fakeURL", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}

func TestHttpDoURLInvalid(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.Requests)
	uc.Client = htc

	resp, err := uc.Do(context.Background(), "abc", "url", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}
func TestGetURLClient(t *testing.T) {

	tduration := time.Second * 2

	var uc = new(httpclient.Options)
	uc.Compressed = true
	uc.SSLEnabled = true
	uc.HandshakeTimeout = tduration
	uc.ResponseHeaderTimeout = tduration
	uc.RequestTimeout = tduration

	_, err := httpclient.New(uc)

	assert.NoError(t, err)

}

func TestGetURLClientURLClientOptionNil(t *testing.T) {

	var uc1 *httpclient.Options

	_, err := httpclient.New(uc1)

	assert.NoError(t, err)

}

func TestGetURLClientSSLEnabledFalse(t *testing.T) {

	tduration := time.Second * 2

	var uc2 = new(httpclient.Options)
	uc2.Compressed = true
	uc2.SSLEnabled = false
	uc2.HandshakeTimeout = tduration
	uc2.ResponseHeaderTimeout = tduration

	_, err := httpclient.New(uc2)

	assert.NoError(t, err)

}
