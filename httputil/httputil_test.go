package httputil_test

import (
	"bytes"
	"encoding/json"
	"github.com/go-chassis/foundation/httputil"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetURI(t *testing.T) {
	req := &http.Request{}

	t.Run("set wrong url to request", func(t *testing.T) {
		httputil.SetURI(req, "127.0.0.1:8080")
		assert.Nil(t, req.URL)
	})

	t.Run("set right url to request", func(t *testing.T) {
		httputil.SetURI(req, "http://127.0.0.1:8080")
		assert.NotNil(t, req.URL)
		assert.Equal(t, req.URL.Host, "127.0.0.1:8080")
		assert.Equal(t, req.URL.Scheme, "http")
	})

}
func TestSetBody(t *testing.T) {
	req := &http.Request{}
	t.Run("set data of body to request", func(t *testing.T) {
		httputil.SetBody(req, nil)
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)
		assert.Zero(t, len(body))

		data := map[string]string{
			"Test1": "test1",
			"Test2": "test2",
		}
		b, err := json.Marshal(data)
		assert.Nil(t, err)
		httputil.SetBody(req, b)
		body, err = ioutil.ReadAll(req.Body)
		assert.NotNil(t, body)
		assert.NotZero(t, len(body))
		assert.Equal(t, b, body)
	})
}
func TestSetGetCookie(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	t.Run("set data to req cookie", func(t *testing.T) {
		ck := "cookie_key"
		cv := "cookie_value"
		httputil.SetCookie(req, ck, cv)

		v, err := req.Cookie(ck)
		assert.Nil(t, err)
		assert.Equal(t, v.Value, cv)

		gv := httputil.GetCookie(req, ck)
		assert.NotEmpty(t, gv)
		assert.Equal(t, gv, cv)
	})

}

func TestSetContentType(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	t.Run("set value to req.header of ContentType", func(t *testing.T) {
		ct := "application/json"
		httputil.SetContentType(req, ct)
		cv := httputil.GetContentType(req)
		assert.Equal(t, cv, ct)
	})
}

func TestRespBody(t *testing.T) {
	resp := &http.Response{
		Body: nil,
	}
	t.Run("resp or body is nil , did not reply any data ", func(t *testing.T) {
		b := httputil.ReadBody(nil)
		assert.Nil(t, b)
		b = httputil.ReadBody(resp)
		assert.Nil(t, b)
	})
	bodies := []byte("test read resp bodies")
	var bb io.Reader
	t.Run("get data of resp body", func(t *testing.T) {
		bb = bytes.NewReader(bodies)
		rc, ok := bb.(io.ReadCloser)
		if !ok && bodies != nil {
			rc = ioutil.NopCloser(bb)
		}
		resp.Body = rc
		b := httputil.ReadBody(resp)
		assert.NotNil(t, b)
		assert.Equal(t, bodies, b)
	})
}

func TestGetRespCookie(t *testing.T) {
	resp := &http.Response{
		Header: make(map[string][]string),
	}
	cookies := []*http.Cookie{
		{
			Name:  "k1",
			Value: "v1",
		},
		{
			Name:  "k2",
			Value: "v2",
		},
	}
	for _, v := range cookies {
		httputil.SetRespCookie(resp, v)
	}
	t.Run("get exist key for cookie", func(t *testing.T) {
		b := httputil.GetRespCookie(resp, "k1")
		assert.NotNil(t, b)
		assert.Equal(t, b, []byte("v1"))
	})
	t.Run("get not exist key for cookie", func(t *testing.T) {
		b := httputil.GetRespCookie(resp, "k3")
		assert.Nil(t, b)
	})
}
