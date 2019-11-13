package httpclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-chassis/foundation/string"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//SignRequest sign a http request so that it can talk to API server
//this is global implementation, if you do not set SignRequest in Options
//client will use this function
var SignRequest func(*http.Request) error

type gzipBodyReader struct {
	*gzip.Reader
	Body io.ReadCloser
}

func (w *gzipBodyReader) Close() error {
	w.Reader.Close()
	return w.Body.Close()
}

func NewGZipBodyReader(body io.ReadCloser) (io.ReadCloser, error) {
	reader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	return &gzipBodyReader{reader, body}, nil
}

//Requests is a restful client
type Requests struct {
	*http.Client
	TLS     *tls.Config
	options Options
}

func (r *Requests) Get(ctx context.Context, url string, headers http.Header) (resp *http.Response, err error) {
	return r.Do(ctx, "GET", url, headers, nil)
}
func (r *Requests) Post(ctx context.Context, url string, headers http.Header, body []byte) (resp *http.Response, err error) {
	return r.Do(ctx, "POST", url, headers, body)
}
func (r *Requests) Put(ctx context.Context, url string, headers http.Header, body []byte) (resp *http.Response, err error) {
	return r.Do(ctx, "PUT", url, headers, body)
}
func (r *Requests) Delete(ctx context.Context, url string, headers http.Header) (resp *http.Response, err error) {
	return r.Do(ctx, "DELETE", url, headers, nil)
}
func (r *Requests) Do(ctx context.Context, method string, url string, headers http.Header, body []byte) (resp *http.Response, err error) {
	if strings.HasPrefix(url, "https") {
		if transport, ok := r.Client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = r.TLS
		}
	}
	if headers == nil {
		headers = make(http.Header)
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = []string{"*/*"}
	}
	if _, ok := headers["Accept-Encoding"]; !ok && r.options.Compressed {
		headers["Accept-Encoding"] = []string{"deflate, gzip"}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create request failed: %s", err.Error()))
	}
	req = req.WithContext(ctx)
	req.Header = headers
	//sign a request, first use function in r options
	//if there is not, use global function
	if r.options.SignRequest != nil {
		if err = r.options.SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	} else if SignRequest != nil {
		if err = SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	}
	resp, err = r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := NewGZipBodyReader(resp.Body)
		if err != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
			return nil, err
		}
		resp.Body = reader
	}

	if os.Getenv("HTTP_DEBUG") == "1" {
		fmt.Println("--- BEGIN ---")
		fmt.Printf("> %s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
		for key, header := range req.Header {
			for _, value := range header {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		fmt.Println(">")
		fmt.Println(stringutil.Bytes2str(body))
		fmt.Printf("< %s %s\n", resp.Proto, resp.Status)
		for key, header := range resp.Header {
			for _, value := range header {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("< %s \n", bodyBytes)
		fmt.Println("--- END ---")
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return resp, nil
}

func setOptionDefaultValue(o *Options) Options {
	if o == nil {
		return DefaultOptions
	}

	option := *o
	if option.RequestTimeout <= 0 {
		option.RequestTimeout = DefaultOptions.RequestTimeout
	}
	if option.HandshakeTimeout <= 0 {
		option.HandshakeTimeout = DefaultOptions.HandshakeTimeout
	}
	if option.ResponseHeaderTimeout <= 0 {
		option.ResponseHeaderTimeout = DefaultOptions.ResponseHeaderTimeout
	}
	if option.ConnsPerHost <= 0 {
		option.ConnsPerHost = DefaultOptions.ConnsPerHost
	}
	return option
}

//New is a function which which sets client option
func New(o *Options) (client *Requests, err error) {
	option := setOptionDefaultValue(o)
	if !option.SSLEnabled {
		client = &Requests{
			Client: &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost:   option.ConnsPerHost,
					TLSHandshakeTimeout:   option.HandshakeTimeout,
					ResponseHeaderTimeout: option.ResponseHeaderTimeout,
					DisableCompression:    !option.Compressed,
				},
			},
			options: option,
		}

		return
	}

	client = &Requests{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   option.HandshakeTimeout,
				ResponseHeaderTimeout: option.ResponseHeaderTimeout,
				DisableCompression:    !option.Compressed,
			},
			Timeout: option.RequestTimeout,
		},
		TLS:     option.TLSConfig,
		options: option,
	}
	return
}
