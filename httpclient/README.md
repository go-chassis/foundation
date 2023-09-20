### Http Client

As every one knows go default http client can not be used in production,
[you must implement your own http client](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779).

httpclient is a production ready client.

```go
	options := &httpclient.URLClientOption{
		TLSConfig:  opt.TLSConfig,
		Compressed: opt.Compressed,
	}
	c, err = httpclient.New(options)
	if err != nil {
		return err
	}
	resp, err := uc.Do("GET", "https://fakeURL", nil, nil)
```

Debug Mode
```shell script
export HTTP_DEBUG=1
```
