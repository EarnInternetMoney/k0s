package config

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/btwiuse/pretty"
	"gopkg.in/yaml.v2"
	"k0s.io/conntroll/pkg"
	"k0s.io/conntroll/pkg/client"
	"k0s.io/conntroll/pkg/version"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	if value == "" {
		return nil
	}
	tags := strings.Split(value, ",")
	*i = append(*i, tags...)
	return nil
}

type config struct {
	Redir            string                     `json:"-" yaml:"redir"`
	Socks            string                     `json:"-" yaml:"socks"`
	Socks5ToHTTP     string                     `json:"-" yaml:"socks5tohttp"`
	Verbose          bool                       `json:"-" yaml:"verbose"`
	Insecure         bool                       `json:"-" yaml:"insecure"`
	CacheCredentials bool                       `json:"-" yaml:"cache_credentials"`
	Credentials      map[string]client.KeyStore `json:"-" yaml:"credentials"`
	ConfigLocation   string                     `json:"-" yaml:"-"`
	Hub              string                     `json:"-" yaml:"hub"`

	uri *url.URL `json:"-"` // where server scheme, host, port, addr are defined

	Version pkg.Version `json:"version" yaml:"-"`
}

func (c *config) GetConfigLocation() string {
	return c.ConfigLocation
}

func (c *config) GetCacheCredentials() bool {
	return c.CacheCredentials
}

func (c *config) GetCredentials() map[string]client.KeyStore {
	return c.Credentials
}

func (c *config) GetSocks() string {
	return c.Socks
}

func (c *config) GetSocks5ToHTTP() string {
	return c.Socks5ToHTTP
}

func (c *config) GetRedir() string {
	return c.Redir
}

func (c *config) GetVersion() pkg.Version {
	return c.Version
}

func (c *config) GetVerbose() bool {
	return c.Verbose
}

func (c *config) GetInsecure() bool {
	return c.Insecure
}

func (c *config) GetPort() string {
	if c.uri.Port() == "" {
		switch c.uri.Scheme {
		case "http":
			return "80"
		case "https":
			return "443"
		}
	}
	return c.uri.Port()
}

func (c *config) GetAddr() string {
	var (
		scheme = c.GetScheme()
		host   = c.GetHost()
		port   = c.GetPort()
	)
	// omit port if already on standard port
	switch {
	case scheme == "http" && port == "80":
		return host
	case scheme == "https" && port == "443":
		return host
	default:
		return fmt.Sprintf("%s:%s", host, port)
	}
}

func (c *config) GetSchemeWS() string {
	switch c.GetScheme() {
	case "https":
		return "wss"
	default:
		return "ws"
	}
}

func (c *config) GetScheme() string {
	if c.uri.Scheme == "http" && c.uri.Hostname() == "" && c.uri.Port() == "443" {
		return "https"
	}
	return c.uri.Scheme
}

type Opt func(c *config)

func SetCacheCredentials(cc bool) Opt {
	return func(c *config) {
		c.CacheCredentials = cc
	}
}

func SetHub(h string) Opt {
	return func(c *config) {
		c.Hub = h
	}
}

func SetRedir(h string) Opt {
	return func(c *config) {
		c.Redir = h
	}
}

func SetSocks(h string) Opt {
	return func(c *config) {
		c.Socks = h
	}
}

func SetSocks5ToHTTP(h string) Opt {
	return func(c *config) {
		c.Socks5ToHTTP = h
	}
}

func SetInsecure(h bool) Opt {
	return func(c *config) {
		c.Insecure = h
	}
}

func SetURI() Opt {
	return func(c *config) {
		var hubapi = c.Hub
		// prepend host 127.0.0.1 to port without an explicit host :8000
		if strings.HasPrefix(hubapi, ":") {
			hubapi = "127.0.0.1" + hubapi
		}
		// default to http
		if !(strings.HasPrefix(hubapi, "http://") || strings.HasPrefix(hubapi, "https://")) {
			hubapi = "http://" + hubapi
		}

		uri, err := url.Parse(hubapi)
		if err != nil {
			log.Fatalln(err)
		}
		c.uri = uri
	}
}

func SetVerbose(v bool) Opt {
	return func(c *config) {
		c.Verbose = v
	}
}

func (c *config) GetHost() string {
	host := c.uri.Hostname()
	if host == "" {
		return "127.0.0.1"
	}
	return host
}

func isExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func probeConfigFile() string {
	var (
		globalConfig = "/etc/conntroll/client.yaml"
		userConfig   = os.ExpandEnv("${HOME}/.conntroll/client.yaml")
		localConfig  = "client.yaml"
	)
	for _, conf := range []string{
		localConfig,
		userConfig,
		globalConfig,
	} {
		if isExist(conf) {
			return conf
		}
	}
	return ""
}

func loadConfigFile(file string) *config {
	c := &config{
		Hub:            pkg.DEFAULT_HUB_ADDRESS,
		Version:        version.GetVersion(),
		ConfigLocation: file,
	}
	if file == "" {
		return c
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
		return c
	}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(c)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
	return c
}

func Parse(args []string) client.Config {
	var (
		fset = flag.NewFlagSet("client", flag.ExitOnError)

		opts = []Opt{}

		hubapi       *string = fset.String("hub", pkg.DEFAULT_HUB_ADDRESS, "Hub address.")
		redir        *string = fset.String("redir", pkg.REDIR_PROXY_PORT, "Redir port.")
		socks        *string = fset.String("socks", pkg.SOCKS5_PROXY_PORT, "Socks port.")
		socks5tohttp *string = fset.String("socks5tohttp", pkg.HTTP_PROXY_PORT, "Http port.")
		verbose      *bool   = fset.Bool("verbose", false, "Verbose log.")
		version      *bool   = fset.Bool("version", false, "Show agent/hub version info.")
		insecure     *bool   = fset.Bool("insecure", false, "Allow insecure server connections when using SSL.")
		// cc           *bool   = fset.Bool("cc", false, "Cache credentials.")
		c *string = fset.String("c", probeConfigFile(), "Config file location.")
	)

	err := fset.Parse(args)
	if err != nil {
		log.Fatalln(err)
	}

	fset.Visit(func(f *flag.Flag) {
		/*
			if f.Name == "cc" {
				opts = append(opts, SetCacheCredentials(*cc))
			}
		*/
		if f.Name == "hub" {
			opts = append(opts, SetHub(*hubapi))
		}
		if f.Name == "redir" {
			opts = append(opts, SetRedir(*redir))
		}
		if f.Name == "socks" {
			opts = append(opts, SetSocks(*socks))
		}
		if f.Name == "socks5tohttp" {
			opts = append(opts, SetSocks5ToHTTP(*socks5tohttp))
		}
		if f.Name == "verbose" {
			opts = append(opts, SetVerbose(*verbose))
		}
		if f.Name == "insecure" {
			opts = append(opts, SetInsecure(*insecure))
		}
	})

	//  The 1st positional argument is used if you leave out the -hub part.
	if len(fset.Args()) != 0 {
		opts = append(opts, SetHub(fset.Args()[0]))
	}

	opts = append(opts, SetURI())

	baseConfig := loadConfigFile(*c)

	for _, opt := range opts {
		opt(baseConfig)
	}

	if *version {
		printClientVersion(baseConfig)
		printHubVersion(baseConfig)
		os.Exit(0)
	}

	return baseConfig
}

type clientVersion struct {
	Client pkg.Version
}

type hubVersion struct {
	Hub pkg.Version
}

func printClientVersion(c client.Config) {
	av := &clientVersion{c.GetVersion()}
	fmt.Println(pretty.YAMLString(av))
}

func printHubVersion(c client.Config) {
	var (
		ub = &url.URL{
			Scheme: c.GetScheme(),
			Host:   c.GetAddr(),
			Path:   "/version",
		}
		u = ub.String()
	)
	resp, err := http.Get(u)
	if err != nil {
		log.Fatalln(err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// log.Println(string(buf))
	v, err := version.Decode(buf)
	if err != nil {
		log.Fatalln(err)
	}

	hv := &hubVersion{v}
	fmt.Print(pretty.YAMLString(hv))
}

func (c *config) String() string {
	return pretty.JsonString(c)
}

/*
func Decode(data []byte) (agent.Info, error) {
	v := &config{
		Info: info.EmptyInfo(),
	}
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, err
}
*/