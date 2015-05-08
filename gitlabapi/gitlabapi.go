package gitlabapi

import (
	"github.com/cnam/apibuilder"
	"net/url"
)

func NewApi (host, token string) (*apibuilder.Api) {

	link, _ := url.Parse(host);

	return  apibuilder.NewApi(&apibuilder.Config{
		link,
		token,
		"clients/",
		"index.json",
	})
}
