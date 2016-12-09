package flyrc

import (
	"os"
	"errors"
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"
)

type Flyrc struct {
	Targets map[string]flyrcTarget `yaml:"targets"`
}

type flyrcTarget struct {
	Api string `yaml:"api"`
}

func GetTarget(host string) (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("Unable to determine users' HOME - is $HOME set?")
	}

	flyrcData, err := ioutil.ReadFile(path.Join(home, ".flyrc"))
	if err != nil {
		return "", err
	}

	f := Flyrc{}
	err = yaml.Unmarshal(flyrcData, &f)
	if err != nil {
		return "", err
	}


	for k, v := range f.Targets {
		if v.Api == host {
			return k, nil
		}
	}

	return "", errors.New("Unable to match URL with target in ~/.flyrc")
}
