package httpclient_test

import (
	"github.com/go-chassis/foundation/httpclient"
	_ "github.com/go-chassis/go-chassis/security/plugins/aes"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpDo(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.URLClient)
	uc.Client = htc

	htServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resp, err := uc.HTTPDo(http.MethodGet, htServer.URL, nil, nil)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
}

func TestHttpDoHeadersNil(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.URLClient)
	uc.Client = htc

	resp, err := uc.HTTPDo("GET", "https://fakeURL", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}

func TestHttpDoURLInvalid(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.URLClient)
	uc.Client = htc

	resp, err := uc.HTTPDo("abc", "url", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}
func TestGetURLClient(t *testing.T) {

	tduration := time.Second * 2

	var uc = new(httpclient.URLClientOption)
	uc.Compressed = true
	uc.SSLEnabled = true
	uc.HandshakeTimeout = tduration
	uc.ResponseHeaderTimeout = tduration

	c, err := httpclient.GetURLClient(uc)
	expectedc := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   tduration,
				ResponseHeaderTimeout: tduration,
				DisableCompression:    false,
			},
		},
	}

	assert.Equal(t, expectedc.Client, c.Client)
	assert.NoError(t, err)

}

func TestGetURLClientURLClientOptionNil(t *testing.T) {

	option := httpclient.DefaultURLClientOption
	expectedclient := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   option.HandshakeTimeout,
				ResponseHeaderTimeout: option.ResponseHeaderTimeout,
				DisableCompression:    !option.Compressed,
			},
		},
		TLS: option.TLSConfig,
	}

	var uc1 *httpclient.URLClientOption

	c1, err := httpclient.GetURLClient(uc1)

	assert.Equal(t, expectedclient.Client, c1.Client)
	assert.NoError(t, err)

}

func TestGetURLClientSSLEnabledFalse(t *testing.T) {

	tduration := time.Second * 2

	expectedc := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   tduration,
				ResponseHeaderTimeout: tduration,
				DisableCompression:    false,
			},
		},
	}

	var uc2 = new(httpclient.URLClientOption)
	uc2.Compressed = true
	uc2.SSLEnabled = false
	uc2.HandshakeTimeout = tduration
	uc2.ResponseHeaderTimeout = tduration

	c2, err := httpclient.GetURLClient(uc2)

	assert.Equal(t, expectedc.Client, c2.Client)
	assert.NoError(t, err)

}
