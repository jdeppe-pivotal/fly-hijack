package main

import (
	"flag"
	"log"
	"strings"
	"fmt"
	"syscall"
	"os"
	"errors"
	"path"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"net/url"
)

const (
	PIPELINE = 3
	JOB      = 5
	BUILD    = 7
)

func main() {
	var instance string

	flag.StringVar(&instance, "t", "", "The concourse instance name")
	flag.Parse()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse URL: %s", err))
	}

	// Ex: /teams/main/pipelines/gemfire-9.0.0/jobs/OperationsTest/builds/4
	parts := strings.Split(u.Path, "/")

	// We didn't give -t option
	if instance == "" {
		hostPart := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

		instance, err = getTarget(hostPart)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	args := []string{
		"fly",
		"-t", instance,
		"hijack",
		"-j", fmt.Sprintf("%s/%s", parts[PIPELINE], parts[JOB]),
	}

	if len(parts) > BUILD {
		args = append(args, "-b", parts[BUILD])
	}

	err = syscall.Exec("/usr/local/bin/fly", args, os.Environ())
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to launch fly: %s", err))
	}
}

type Flyrc struct {
	Targets map[string]flyrcTarget `yaml:"targets"`
}

type flyrcTarget struct {
	Api string `yaml:"api"`
}

func getTarget(host string) (string, error) {
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
