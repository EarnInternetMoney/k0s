package info

import (
	"encoding/json"

	"k0s.io/conntroll/pkg/hub"
	"k0s.io/conntroll/pkg/version"
)

var _ hub.Info = (*Info)(nil)

type Meta struct {
	OS       string `json:"os"`
	Pwd      string `json:"pwd"`
	Arch     string `json:"arch"`
	Distro   string `json:"distro,omitempty"`
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}

type Info struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
	Htpasswd map[string]string `json:"htpasswd,omitempty"`

	Meta            `json:"meta"`
	version.Version `json:"version"`
}

func Decode(data []byte) (hub.Info, error) {
	v := &Info{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (info *Info) GetOS() string {
	return info.OS
}

func (info *Info) GetPwd() string {
	return info.Pwd
}

func (info *Info) GetArch() string {
	return info.Arch
}

func (info *Info) GetDistro() string {
	return info.Distro
}

func (info *Info) GetUsername() string {
	return info.Username
}

func (info *Info) GetHostname() string {
	return info.Hostname
}

func (info *Info) GetID() string {
	return info.ID
}

func (info *Info) GetName() string {
	return info.Name
}

func (info *Info) GetTags() []string {
	return info.Tags
}

func (info *Info) GetHtpasswd() map[string]string {
	return info.Htpasswd
}
