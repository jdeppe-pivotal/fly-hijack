package flyrc

import (
	"os"
	"errors"
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"
	"fmt"
	"net/url"
)

type Flyrc struct {
	Targets map[string]flyrcTarget `yaml:"targets"`
}

type flyrcTarget struct {
	Api   string   `yaml:"api"`
	Token ApiToken `yaml:"token"`
}

type ApiToken struct {
	Type  string `yaml:"api"`
	Value string `yaml:"value"`
}

func NewFlyrc() (*Flyrc, error) {

	home := os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("Unable to determine users' HOME - is $HOME set?")
	}

	flyrcData, err := ioutil.ReadFile(path.Join(home, ".flyrc"))
	if err != nil {
		return nil, err
	}

	f := &Flyrc{}
	err = yaml.Unmarshal(flyrcData, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (this *Flyrc) GetTarget(u *url.URL) (string, error) {
	hostPart := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	return this.getTarget(hostPart)
}

func (this *Flyrc) GetBearerToken(u *url.URL) (string, error) {
	schemeHost := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	for _, v := range this.Targets {
		if v.Api == schemeHost {
			return v.Token.Value, nil
		}
	}
	return "", errors.New("Unable to match URL with target in ~/.flyrc")
}

func (this *Flyrc) getTarget(schemeHost string) (string, error) {
	for k, v := range this.Targets {
		if v.Api == schemeHost {
			return k, nil
		}
	}

	return "", errors.New("Unable to match URL with target in ~/.flyrc")
}
