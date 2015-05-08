package main

import (
	"github.com/cnam/apibuilder"
	"log"
	"net/url"
)

type Issue struct {
	IssueId int `json:"id"`
	ProjectId int `json:"project_id"`
	Title  string `json:"title"`
}

func main() {
	var issues []Issue
	var issue Issue

	link, _ := url.Parse("https://gitlab.com/api/v3/");

	api := apibuilder.NewApi(&apibuilder.Config{
		link,
		"qwerty",
		"clients/",
		"index.json",
	})

	p := make(map[string]string)

	p["project_id"] = "83866"

	command := api.NewCommand("GetIssuesByProject", p, &issues)
	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issues)

	p["project_id"] = "83866"
	p["issue_id"] = "144751"

	command = api.NewCommand("GetIssue", p, &issue)

	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issue)

	p["project_id"] = "83866"
	p["title"] = "new issue"

	command = api.NewCommand("CreateIssue", p, &issue)

	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issue)

}