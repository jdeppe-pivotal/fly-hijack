package main

import (
	"flag"
	"log"
	"strings"
	"fmt"
	"syscall"
	"os"
	"net/url"
	"fly-utils/flyrc"
	"os/exec"
)

const (
	PIPELINE = 4
	JOB = 6
	BUILD = 8
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
		rc, err := flyrc.NewFlyrc()
		if err != nil {
			log.Fatal(err)
		}
		instance, err = rc.GetTarget(u)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(parts) < JOB + 1 {
		log.Fatal("URL appears too short - make sure it at least contains a 'jobs' value")
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

	flyPath, err := exec.LookPath("fly")
	if err != nil {
		log.Fatal("Unable to find 'fly' executable in your PATH")
	}

	err = syscall.Exec(flyPath, args, os.Environ())
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to launch fly: %s", err))
	}
}

