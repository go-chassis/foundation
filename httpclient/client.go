package httpclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"compress/gzip"
	"context"
	"github.com/go-chassis/foundation/string"
	"io"
	"os"
)

//DefaultURLClientOption is a struct object which has default client option
var DefaultURLClientOption = URLClientOption{
	Compressed:            true,
	HandshakeTimeout:      30 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	RequestTimeout:        60 * time.Second,
	ConnsPerHost:          5,
}

//SignRequest sign a http request so that it can talk to API server
//this is global implementation, if you do not set SignRequest in URLClientOption
//client will use this function
var SignRequest func(*http.Request) error

//URLClientOption is a struct which provides options for client
type URLClientOption struct {
	SSLEnabled            bool
	TLSConfig             *tls.Config
	Compressed            bool
	HandshakeTimeout      time.Duration
	ResponseHeaderTimeout time.Duration
	RequestTimeout        time.Duration
	ConnsPerHost          int
	SignRequest           func(*http.Request) error
}
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

//URLClient is a struct used for storing details of a client
type URLClient struct {
	*http.Client
	TLS     *tls.Config
	options URLClientOption
}

func (client *URLClient) HTTPDoWithContext(ctx context.Context, method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	if strings.HasPrefix(rawURL, "https") {
		if transport, ok := client.Client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = client.TLS
		}
	}

	if headers == nil {
		headers = make(http.Header)
	}

	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = []string{"*/*"}
	}
	if _, ok := headers["Accept-Encoding"]; !ok && client.options.Compressed {
		headers["Accept-Encoding"] = []string{"deflate, gzip"}
	}

	req, err := http.NewRequest(method, rawURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create request failed: %s", err.Error()))
	}
	req = req.WithContext(ctx)
	req.Header = headers
	//sign a request, first use function in client options
	//if there is not, use global function
	if client.options.SignRequest != nil {
		if err = client.options.SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	} else if SignRequest != nil {
		if err = SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	}
	resp, err = client.Client.Do(req)
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
		fmt.Println("<")
		fmt.Println("--- END ---")
	}
	return resp, nil
}

func (client *URLClient) HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	return client.HTTPDoWithContext(context.Background(), method, rawURL, headers, body)
}

func setOptionDefaultValue(o *URLClientOption) URLClientOption {
	if o == nil {
		return DefaultURLClientOption
	}

	option := *o
	if option.RequestTimeout <= 0 {
		option.RequestTimeout = DefaultURLClientOption.RequestTimeout
	}
	if option.HandshakeTimeout <= 0 {
		option.HandshakeTimeout = DefaultURLClientOption.HandshakeTimeout
	}
	if option.ResponseHeaderTimeout <= 0 {
		option.ResponseHeaderTimeout = DefaultURLClientOption.ResponseHeaderTimeout
	}
	if option.ConnsPerHost <= 0 {
		option.ConnsPerHost = DefaultURLClientOption.ConnsPerHost
	}
	return option
}

//GetURLClient is a function which which sets client option
func GetURLClient(o *URLClientOption) (client *URLClient, err error) {
	option := setOptionDefaultValue(o)

	if !option.SSLEnabled {
		client = &URLClient{
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

	client = &URLClient{
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
