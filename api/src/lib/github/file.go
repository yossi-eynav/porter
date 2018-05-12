package github

import (
	"net/http"
	"porter/api/src/config"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
)
type ContentFileResponse struct {
	Content string `json:"content"`
}

func (cf *ContentFileResponse) parseContent() (string) {
	content, err :=base64.StdEncoding.DecodeString(cf.Content)
	if err != nil {
		return ""
	}

	return  string(content)
}

func FetchFileContent(repositoryURI string, path string) (string, error) {
	var fileResponse ContentFileResponse
	resp, err := http.Get(repositoryURI + "/contents/"+ path +"?access_token=" + config.GetValue("githubToken"))
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(body, &fileResponse)

	return fileResponse.parseContent(), nil
}