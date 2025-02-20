package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type config struct {
	Timeout time.Duration `yaml:"timeout"`
	Creds   creds         `yaml:"creds"`
	Account account       `yaml:"account"`
}

type account struct {
	Spec string `yaml:"spec"`
	ID   string `yaml:"id"`
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
	ws           *tradovate.WS
	spec         string
	id           int
	chartChannel chan *tradovate.Chart
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

	exit := 1
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer func() {
		cancel()
		os.Exit(exit)
	}()

	chartChannel := make(chan *tradovate.Chart, 50)
	ws, err := tradovate.NewSocket(ctx, tradovate.WSSSandboxURL, nil, rest,
		tradovate.WithErrHandler(func(err error) { fmt.Println("caught error in err handler:", err) }),
		tradovate.WithChartHandler(func(c *tradovate.Chart) {
			chartChannel <- c
		}),
	)
	if err != nil {
		fmt.Printf("failed creating socket connection for test: %v\n", err)
		return
	}
	defer ws.Close()

	c = client{
		ctx:          ctx,
		ws:           ws,
		spec:         c.spec,
		id:           c.id,
		chartChannel: chartChannel,
	}
	exit = m.Run()
}
