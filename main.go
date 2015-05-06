package main

import (
	//"github.com/cnam/gitlab-api-go/src/gitlabapi"
	"./src/gitlabapi"
	"log"
	"net/url"
)

func main() {
	link, err := url.Parse("https://gitlab.com/api/v3/");

	if (err != nil) {
		log.Panicf("Invalid url %s", err.Error())
	}

	config := &gitlabapi.Config{
		BasePath: link,
		PrivateToken: "jc6-QyBGSsF-ySyEMgLn",
	}

	api := gitlabapi.New(config)


	m := make(map[string]string)
	m["project_id"] = "83866";
	resp := api.Exec("GetIssuesByProject", m);

	log.Printf("%v", resp)

	m["project_id"] = "83866";
	m["issue_id"] = "144751";

	resp = api.Exec("GetIssue", m);
	log.Printf("%v", resp)
}