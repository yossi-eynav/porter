package github

import (
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"porter/api/src/config"
	"sync"
	"fmt"
)

type RepositoryResponse struct {
	Url string `json:"url"`
	Name string `json:"name"`
	Html_url string `json:"html_url"`
}

const repositoryPages = 20

type Repository struct {
	Name         string
	ExposedPorts []string
	Url 		string
	Html_url string
}


func (rp Repository) MarshalBinary() ([]byte, error) {
	return json.Marshal(&rp)
}

func (rp *Repository) UnMarshalBinary(data []byte) (error) {
	return json.Unmarshal(data, &rp)
}

func FetchRepository(page int,ch chan []RepositoryResponse){
	var repositoryResponse []RepositoryResponse
	resp,_ := http.Get("https://api.github.com/orgs/fiverr/repos?access_token="+ config.GetValue("githubToken") +"&limit=100&sort=pushed&page=" + strconv.Itoa(page))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &repositoryResponse)
	ch <- repositoryResponse
}

func FetchAllRepositories() ([]RepositoryResponse) {
	repositories := make([]RepositoryResponse, 0)
	repositoriesResponseCh := make(chan []RepositoryResponse)
	var wg sync.WaitGroup

	for i:= 1; i <= repositoryPages ; i++  {
		go FetchRepository(i, repositoriesResponseCh)
	}

	go func(ch chan []RepositoryResponse) {
		wg.Add(repositoryPages)
		wg.Wait()
		close(ch)
	}(repositoriesResponseCh)

	for repositoryResponses := range repositoriesResponseCh {
		for i :=0 ; i < len(repositoryResponses); i++  {
			repositories = append(repositories, repositoryResponses[i])
		}
		wg.Done()
	}

	return repositories
}
