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
	Client *http.Client
}

/*
	create new API client for gitlab api
 */
func New(config *Config) (*Api) {
	var schema *Schema;
	fileContent, err := ioutil.ReadFile("clients/command.json");

	if err != nil {
		log.Println("Error read", err.Error())
	}

	err = json.Unmarshal(fileContent, &schema)

	if (err != nil) {
		log.Println("Error parse", err.Error())
	}

	api := &Api{
		Config:config,
		Schema:schema,
		Client:&http.Client{},
	}

	return api
}

/*
	Exec new command
 */
func (api *Api) Exec(commandName string, parameters map[string]string) interface{} {
	command := api.offset(commandName)
	url := api.url(command.Uri, parameters)

	log.Println(url)

	req, err := http.NewRequest(command.Method, url, nil)

	if err != nil {
		log.Println("Bad request", err.Error())
	}

	resp, err := api.Client.Do(req)

	if err != nil {
		log.Println("Bad Response", err.Error())
	}

	return api.parseResponse(resp)
}

func (api *Api) parseResponse(resp *http.Response) interface{} {

	var respBody interface{}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&respBody)

	if resp.StatusCode != http.StatusOK ||
	   resp.StatusCode != http.StatusCreated ||
	   resp.StatusCode != http.StatusNoContent{
		log.Panicf("Bad response code %s", resp.StatusCode)
	}

	return respBody
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