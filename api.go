package main

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"path"
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
	BasePath string
}

type Api struct {
	Schema
	Config
	Client *http.Client
}

func main() {
	config := Config{
		BasePath: "https://gitlab.com/api/v3/",
	}

	api := New(config)
	m := make(map[string]interface{})
	resp := api.Exec("GetIssuesByProject", m);
	log.Printf("%i", resp.StatusCode)

	var respBody interface{}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	log.Print(resp.ContentLength);

	decoder.Decode(&respBody)

	log.Printf("%v", respBody)
}

/*
	create new API client for gitlab api
 */
func New(config Config) (*Api) {
	var schema Schema;
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
func (api *Api) Exec(commandName string, parameters map[string]interface{}) (*http.Response) {
	command := api.offset(commandName)
	url := api.url(command.Uri, parameters)

	req, err := http.NewRequest(command.Method, url, nil)

	if err != nil {
		log.Println("Bad request", err.Error())
	}

	resp, err := api.Client.Do(req)

	if err != nil {
		log.Println("Bad Response", err.Error())
	}

	return resp
}

/*
	Generate url for request
 */
func (api *Api) url(uri string, parameters map[string]interface{}) string {
	chunks, file := path.Split(uri);
	log.Printf("%v", file)
	log.Printf("%+v", chunks)

	/*for chunk := range chunks {
		chunk
	}

	regexp.ReplaceAll();*/

	return api.Config.BasePath+uri+"?private_token=tFcoFiGM1DStbHayyKmc";
}

func (api *Api) offset(commandName string) (Command) {
	command, ok := api.Schema.Operations[commandName]

	if !ok {
		panic("Command not found")
	}

	return command
}