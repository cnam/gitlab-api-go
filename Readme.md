# Simple sdk for gitlab api

###Install

> go get github.com/cnam/gitlab-api-go/gitlabapi

```go
    import(
        "github.com/cnam/gitlab-api-go/gitlabapi"
    )
```

### Example usage

```go
func main (){
    var issues []Issue
	var issue Issue

	api := gitlabapi.NewApi("https://gitlab.com/api/v3", "qwerty")

	p := make(map[string]string)

	p["project_id"] = "1"

	command := api.NewCommand("GetIssuesByProject", p, &issues)
	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issues)

	p["project_id"] = "1"
	p["issue_id"] = "1"

	command = api.NewCommand("GetIssue", p, &issue)

	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issue)

	p["project_id"] = "1"
	p["title"] = "new issue"

	command = api.NewCommand("CreateIssue", p, &issue)

	log.Printf("Request url %v", command.Request.URL)

	command.Execute()

	log.Printf("%+v", issue)
}
```