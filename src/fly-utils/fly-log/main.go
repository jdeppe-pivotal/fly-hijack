package main

import (
	"encoding/json"
	"flag"
	"fly-utils/flyrc"
	"fmt"
	"github.com/donovanhide/eventsource"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Build struct {
	ApiUrl       string `json:"api_url"`
	EndTime      int    `json:"end_time"`
	Id           int    `json:"id"`
	JobName      string `json:"job_name"`
	Name         string `json:"name"`
	PipelineName string `json:"pipeline_name"`
	StartTime    int    `json:"start_time"`
	Status       string `json:"status"`
	TeamName     string `json:"team_name"`
	Url          string `json:"url"`
}

type EventData struct {
	Data struct {
		Origin struct {
			Source string `json:"source",omitempty`
		} `json:"origin",omitempty`
		Payload    string `json:"payload",omitempty`
		Time       int    `json:"time",omitempty`
		ExitStatus int    `json:"exit_status",omitempty`
		Status     string `json:"status",omitempty`
	} `json:"data",omitempty`
	Event string `json:"event",omitempty`
}

func main() {
	flag.Parse()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse URL: %s", err))
	}

	rc, err := flyrc.NewFlyrc()
	if err != nil {
		log.Fatal(err)
	}

	token, err := rc.GetBearerToken(u)
	if err != nil {
		log.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:  "skymarshal_auth",
		Value: fmt.Sprintf("Bearer %s", token),
	}

	build, err := getBuild(u, cookie)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get build: %s", err))
	}

	eventUrl, err := url.Parse(fmt.Sprintf("%s://%s/api/v1/builds/%d/events", u.Scheme, u.Host, build.Id))
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("GET", eventUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	request.AddCookie(cookie)
	var client = &http.Client{}

	stream, err := eventsource.SubscribeWith("", client, request)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to subscribe to event source: %s", err))
	}

	showTask := false
	for x := range stream.Events {
		if x.Event() == "end" {
			break
		}
		event := &EventData{}
		err := json.Unmarshal([]byte(x.Data()), event)
		if err != nil {
			log.Printf("%s - %s - '%s'", err, x.Event(), x.Data())
		}

		if event.Event == "start-task" {
			showTask = true
			continue
		}

		if event.Event == "finish-task" {
			showTask = false
			continue
		}

		if showTask {
			fmt.Printf("%s", event.Data.Payload)
			//fmt.Printf("%v", x.Data())
		}
	}
}

func getBuild(u *url.URL, cookie *http.Cookie) (*Build, error) {
	apiUrl := u
	apiUrl.Path = fmt.Sprintf("/api/v1%s", apiUrl.Path)

	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Accept", "*/*")

	var client = &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error calling endpoint for build: %s", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading body from endpoint for build: %s", err))
	}

	var build Build
	err = json.Unmarshal(body, &build)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to unmarshall getBuild response: %s %s", body, err))
	}

	return &build, nil
}
