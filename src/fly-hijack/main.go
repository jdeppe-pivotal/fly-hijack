package main

import (
	"flag"
	"log"
	"strings"
	"fmt"
	"syscall"
	"os"
)

const (
	PIPELINE = 4
	JOB      = 6
	BUILD    = 8
)

func main() {
	var instance string

	flag.StringVar(&instance, "t", "", "The concourse instance name")
	flag.Parse()

	if instance == "" {
		log.Fatal("Please provide an instance with '-t'")
	}

	// Ex: http://concourse.gemfire.pivotalci.info:8080/pipelines/gemfire-9.0.0/jobs/OperationsTest/builds/4
	parts := strings.Split(flag.Arg(0), "/")

	args := []string{
		"fly",
		"-t", instance,
		"hijack",
		"-j", fmt.Sprintf("%s/%s", parts[PIPELINE], parts[JOB]),
	}

        if len(parts) > BUILD {
		args = append(args, "-b", parts[BUILD])
        }

	err := syscall.Exec("/usr/local/bin/fly", args, os.Environ())
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to launch fly: %s", err))
	}
}
