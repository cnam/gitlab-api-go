package main

import (
	"github.com/cnam/gitlab-api-go/src/gitlabapi"
	"log"
	"encoding/json"
	"net/url"
)

func main() {
	link, err := url.Parse("https://gitlab.com/api/v3/");

	if (err != nil) {
		log.Panicf("Invalid url %s", err.Error())
	}

	config := &gitlabapi.Config{
		BasePath: link,
		PrivateToken: "qwerty",
	}

	api := gitlabapi.New(config)
	m := make(map[string]string)
	m["project_id"] = "83866";
	resp := api.Exec("GetIssuesByProject", m);
	log.Printf("%i", resp.StatusCode)

	var respBody interface{}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	log.Print(resp.ContentLength);

	decoder.Decode(&respBody)

	log.Printf("%v", respBody)


	m["project_id"] = "83866";
	m["issue_id"] = "144751";

	resp = api.Exec("GetIssue", m);
	log.Printf("%i", resp.StatusCode)

	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)

	log.Print(resp.ContentLength);

	decoder.Decode(&respBody)

	log.Printf("%v", respBody)
}