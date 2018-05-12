package github

import (
	"net/http"
	"porter/api/src/config"
	"io/ioutil"
	"encoding/json"
)

type FileResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
	DownloadUrl string `json:"download_url"`
}

func FetchFolder(repositoryURI string, path string) ([]FileResponse, error) {
	resp, err := http.Get(repositoryURI + "/contents/"+ path +"?access_token=" +  config.GetValue("githubToken"))
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var fileResponses []FileResponse
	json.Unmarshal(body, &fileResponses)

	return fileResponses, nil
}