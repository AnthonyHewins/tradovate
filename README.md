# Tradovate

Golang tradovate client

```shell
go get https://github.com/AnthonyHewins/tradovate
```

- [Tradovate](#tradovate)
	- [Usage](#usage)

## Usage

The current usage for this package is highly tilted toward the websocket API. The REST client has only the
functionality to fetch a token. This is on purpose for now, because the websocket offers
better performance and covers more use cases. Contributing to the REST client is easy if there's
already implementations available for the socket

For all use cases, you need the rest client:

REST client:

```go
client := tradovate.NewREST(tradovate.SandboxURL, &http.Client{Timeout: time.Second * 10}, &tradovate.Creds{
	Name:       "",
	Password:   "",
	AppID:      "",
	AppVersion: "",
	CID:        "",
	DeviceID:   [16]byte{},
	Sec:        [16]byte{},
})
```

Socket client requires a rest client. Tradovate requires you to pick the right URL if you want market data or if you
just plan on interacting with the API. Your connections determines which one

```go
client := tradovate.NewREST(tradovate.SandboxURL, &http.Client{Timeout: time.Second * 10}, &tradovate.Creds{
	Name:       "",
	Password:   "",
	AppID:      "",
	AppVersion: "",
	CID:        "",
	DeviceID:   uuid.UUID{},
	Sec:        uuid.UUID{},
})

// you don't need these options, but you can use them if you want to
opts := &websocket.DialOpts{
	// part of the websocket library this code uses
}

// API may change to remove the need for the REST client as an argument
s, err := tradovate.NewSocket(ctx, tradovate.WSSSandboxURL, opts, client,
	tradovate.WithToken(&tradovate.Token{}), // skip auth by passing token directly
	tradovate.WithTimeout(time.Second*3), // time out websocket requests you make
	tradovate.WithErrHandler(func(e error) {}), // pass connection-related errors here
	tradovate.WithPingRetries(3), // retry ping failures
	tradovate.WithEntityHandler(func(*EntityMsg) {}), // when an entity in your account is updated, send update here
	tradovate.WithChartHandler(x func(*Chart) {}), // when subbed to marked data, send that chart here
)
```