### Http Client

As every one knows go default http client can not be used in production,
[you must implement your own http client](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779).

http-client give you a production level client.

```go
	options := &httpclient.URLClientOption{
		SSLEnabled: opt.EnableSSL,
		TLSConfig:  opt.TLSConfig,
		Compressed: opt.Compressed,
		Verbose:    opt.Verbose,
	}
	c, err = httpclient.GetURLClient(options)
	if err != nil {
		return err
	}
	resp, err := uc.HTTPDo("GET", "https://fakeURL", nil, nil)
```