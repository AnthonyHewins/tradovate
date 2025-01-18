package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
	"gopkg.in/yaml.v3"
)

type config struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

var client *tradovate.Client

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") != "1" {
		return
	}

	buf, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		fmt.Println("Must set $CONFIG to point to config file, but got err:", err)
		os.Exit(1)
	}

	var c config
	if err := yaml.Unmarshal(buf, c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = tradovate.NewClient(tradovate.SandboxURL, &http.Client{Timeout: time.Second * 10}, &tradovate.Creds{
		Name:       "",
		Password:   "",
		AppID:      "",
		AppVersion: "",
		CID:        "",
		DeviceID:   [16]byte{},
		Sec:        [16]byte{},
	})

	os.Exit(m.Run())
}
