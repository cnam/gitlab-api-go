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
	link, err := url.Parse("https://gitlab.com/api/v3/");

	if (err != nil) {
		log.Panicf("Invalid url %s", err.Error())
	}

	config := &gitlabapi.Config{
		BasePath: link,
		PrivateToken: "jc6-QyBGSsF-ySyEMgLn",
	}

	api := gitlabapi.NewApi(config)

	m := make(map[string]string)
	var issues []Issue
	m["project_id"] = "83866";
	api.Exec("GetIssuesByProject", m, &issues);

	log.Printf("%+v", issues)

	var issue Issue
	m["project_id"] = "83866";
	m["issue_id"] = "144751";

	api.Exec("GetIssue", m, &issue);
	log.Printf("%+v", issue)
}