package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const tokenPath = "./token.json"

type config struct {
	Timeout time.Duration `yaml:"timeout"`
	Creds   creds         `yaml:"creds"`
	Account account       `yaml:"account"`
}

type account struct {
	Spec string `yaml:"spec"`
	ID   int    `yaml:"id"`
}

type creds struct {
	Name     string    `yaml:"name"`
	App      string    `yaml:"appId"`
	Version  string    `yaml:"appVersion"`
	DeviceID uuid.UUID `yaml:"deviceID"`
	Password string    `yaml:"password"`
	ClientID string    `yaml:"client-id"`
	Secret   uuid.UUID `yaml:"secret"`
}

var c client

type client struct {
	ctx          context.Context
	md, api      *tradovate.WS
	spec         string
	id           int
	chartChannel chan *tradovate.Chart
}

func (c *client) shutdown() {
	for _, v := range []struct {
		name   string
		socket *tradovate.WS
	}{
		{"api", c.api},
		{"market data", c.md},
	} {
		if v.socket != nil {
			if err := v.socket.Close(); err != nil {
				fmt.Printf("failed proper shutdown of %s socket: %s\n", v.name, err)
			}
		}
	}
}

func (c *client) readTokenCache() *tradovate.Token {
	f, err := os.Open(tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		fmt.Println("failed reading token cache:", err)
		return nil
	}

	var t tradovate.Token
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		fmt.Println("token is invalid JSON, skipping:", err)
		return nil
	}

	return &t
}

func writeToken(r *tradovate.REST) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	t, err := r.Token(ctx)
	if err != nil || t == nil || t.Expired() {
		return // if token wasnt fetched in test, dont bother
	}

	buf, err := json.Marshal(t)
	if err != nil {
		fmt.Println("failed writing token cache:", err)
		return
	}

	os.WriteFile(tokenPath, buf, 0700)
}

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") != "1" {
		return
	}

	buf, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		fmt.Println("Must set $CONFIG to point to config file, but got err:", err)
		os.Exit(1)
	}

	var cfg config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rest := tradovate.NewREST(tradovate.RESTStage, &http.Client{Timeout: time.Second * 5}, &tradovate.Creds{
		Name:       cfg.Creds.Name,
		Password:   cfg.Creds.Password,
		AppID:      cfg.Creds.App,
		AppVersion: cfg.Creds.Version,
		ClientID:   cfg.Creds.ClientID,
		DeviceID:   cfg.Creds.DeviceID,
		Secret:     cfg.Creds.Secret,
	})
	defer writeToken(rest)

	exit := 1
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer func() {
		cancel()
		os.Exit(exit)
	}()

	c = client{
		ctx:          ctx,
		spec:         cfg.Account.Spec,
		id:           cfg.Account.ID,
		chartChannel: make(chan *tradovate.Chart, 1),
	}

	if t := c.readTokenCache(); t != nil {
		rest.SetToken(t)
	}

	defer c.shutdown()

	c.api, err = tradovate.NewSocket(ctx, tradovate.WSSSandboxURL, nil, rest,
		tradovate.WithErrHandler(func(err error) { fmt.Println("caught error in api websocket err handler:", err) }),
	)
	if err != nil {
		fmt.Printf("failed creating socket connection %s for test: %v\n", tradovate.WSSSandboxURL, err)
		return
	}

	c.md, err = tradovate.NewSocket(ctx, tradovate.WSSMarketDataSandboxURL, nil, rest,
		tradovate.WithErrHandler(func(err error) {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("caught error in marketdata websocket err handler:", err)
		}),
		tradovate.WithChartHandler(func(chart *tradovate.Chart) {
			fmt.Println(chart)
			select {
			case <-ctx.Done():
			case c.chartChannel <- chart:
			}
		}),
	)
	if err != nil {
		fmt.Printf("failed creating socket connection %s for test: %v\n", tradovate.WSSMarketDataSandboxURL, err)
		return
	}

	exit = m.Run()
}
