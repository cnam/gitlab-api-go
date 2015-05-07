package gitlabapi

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"strings"
	"net/url"
)


type Parameter struct {
	Location string `json:"location"`
	Require bool `json:"required"`
	Type string `json:"type"`
}

type Command struct {
	Uri string `json:"uri"`
	Method string `json:"httpMethod"`
	Parameters map[string]Parameter `json:"parameters"`
	ResponseClass string `json:"responseClass"`
}

type Schema struct {
	Name string `json:"name"`
	ApiVersion string `json:"apiVersion"`
	Description string `json:"description"`
	Operations map[string]Command `json:"operations"`
}

type Config struct {
	BasePath *url.URL
	PrivateToken string
}

type Api struct {
	*Schema
	*Config
	*http.Client
}

/*
	create new API client for gitlab api
 */
func NewApi(config *Config) (*Api) {
	var schema *Schema;
	fileContent, err := ioutil.ReadFile("clients/command.json");

	if err != nil {
		log.Println("Error read", err.Error())
	}

	err = json.Unmarshal(fileContent, &schema)

	if (err != nil) {
		log.Println("Error parse", err.Error())
	}

	return &Api{
		schema,
		config,
		&http.Client{},
	}
}

/*
	Exec new command
 */
func (api *Api) Exec(commandName string, parameters map[string]string, mapping interface{}) {
	command := api.offset(commandName)
	url := api.url(command.Uri, parameters)

	req, err := http.NewRequest(command.Method, url, nil)

	if err != nil {
		log.Panicf("Bad request", err.Error())
	}

	resp, err := api.Client.Do(req)

	if err != nil {
		log.Panicf("Bad Response", err.Error())
	}

	api.parseResponse(&command, resp, mapping)
}

/*
	Parse response from gitlab
 */
func (api *Api) parseResponse(command *Command, resp *http.Response, mapping interface{}) {

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&mapping)

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Bad response code", resp.StatusCode)
	}
}

/*
	Generate url for request
 */
func (api *Api) url(uri string, parameters map[string]string) string {
	chunks := strings.Split(uri, "/");

	for index, chunk := range chunks {
		if strings.HasPrefix(chunk,"{") {
			chunk = strings.TrimSuffix(strings.TrimPrefix(chunk, "{"),"}")

			chunkValue, ok := parameters[chunk]

			delete(parameters, chunk)

			if !ok {
				log.Panicf("Parameter %s require", chunk)
			}

			chunks[index] = chunkValue
		}
	}

	uri = strings.Join(chunks, "/");

	baseUrl := &url.URL{
		Scheme: api.Config.BasePath.Scheme,
		Host: api.Config.BasePath.Host,
		Path: strings.TrimSuffix(api.Config.BasePath.Path,"/")+"/"+uri,
	}

	queryParams := url.Values{}

	for name, params := range parameters  {
		queryParams.Add(name, params)
	}

	queryParams.Add("private_token", api.Config.PrivateToken)
	baseUrl.RawQuery = queryParams.Encode()

	return baseUrl.String();
}

/*
	Gets command by name
 */
func (api *Api) offset(commandName string) (Command) {
	command, ok := api.Schema.Operations[commandName]

	if !ok {
		log.Panicf("Command not %s found", commandName)
	}

	return command
}