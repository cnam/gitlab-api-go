package gitlabapi

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"strings"
	"net/url"
	"bytes"
)


type Parameter struct {
	Location string `json:"location"`
	Required bool `json:"required"`
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

type ApiCommand struct {
	Name string
	Parameters map[string]string
	MapTo interface{}
	*Api
	Command
	*http.Request
	*http.Response
}

//
//	create new API client for gitlab api
//
func NewApi(config *Config) (*Api) {
	var schema *Schema

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

//
//Create new command for execute
//
func (api *Api) NewCommand(name string, parameters map[string]string, mapping interface{}) (*ApiCommand) {

	command := api.offset(name)
	requestUrl, parameters  := api.url(&command, parameters)

	data := &url.Values{}

	for name, parameter := range parameters  {
		data.Set(name, parameter)
	}

	req, err := http.NewRequest(command.Method, requestUrl, bytes.NewBufferString(data.Encode()))

	if err != nil {
		log.Panicf("Bad request", err.Error())
	}

	return &ApiCommand{
		name,
		parameters,
		mapping,
		api,
		command,
		req,
		&http.Response{},
	}
}

//
//  execute created command
//
func (command *ApiCommand) Execute() {
	resp, err := command.Api.Client.Do(command.Request)

	if err != nil {
		log.Panicf("Bad Response", err.Error())
	}

	command.parseResponse(resp)
}



// parse response after execute command
func (command *ApiCommand) parseResponse(resp *http.Response) {

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&command.MapTo)
	command.Response = resp

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Panicf("Bad response code", resp.StatusCode)
	}
}


//	Generate url for request 
func (api *Api) url(command *Command, parameters map[string]string) (string, map[string]string) {

	uri := command.Uri
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

	for name, parameter := range command.Parameters {
		if (parameter.Location == "query") {
			parameter, ok := parameters[name]
			if (ok) {
				queryParams.Add(name, parameter)
				delete(parameters, name)
			}
		}
	}

	queryParams.Add("private_token", api.Config.PrivateToken)
	baseUrl.RawQuery = queryParams.Encode()

	return baseUrl.String(), parameters;
}

//
//	Gets command by name
//
func (api *Api) offset(commandName string) (Command) {
	command, ok := api.Schema.Operations[commandName]

	if !ok {
		log.Panicf("Command not %s found", commandName)
	}

	return command
}
