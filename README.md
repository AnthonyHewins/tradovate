# Tradovate

Golang tradovate client

```shell
go get https://github.com/AnthonyHewins/tradovate
```

- [Tradovate](#tradovate)
  - [Usage](#usage)

## Usage

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

Socket client:

```go

```