package main

import (
	//"github.com/cnam/gitlab-api-go/src/gitlabapi"
	"./src/gitlabapi"
	"log"
	"net/url"
)

type Issue struct {
	ProjectId int `json:"project_id"`
	Title  string `json:"title"`
}

func main() {
	var issues []Issue
	var issue Issue

	link, _ := url.Parse("https://gitlab.com/api/v3/");

	config := &gitlabapi.Config{
		BasePath: link,
		PrivateToken: "qwerty",
	}

	api := gitlabapi.NewApi(config)
	p := make(map[string]string)

	p["project_id"] = "83866";

	command := api.NewCommand("GetIssuesByProject", p, &issues)
	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issues)

	p["project_id"] = "83866";
	p["issue_id"] = "144751";

	command = api.NewCommand("GetIssue", p, &issue)

	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issue)
}