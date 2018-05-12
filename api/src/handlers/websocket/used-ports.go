package websocket

import (
	"porter/api/src/handlers/websocket/message"
	"fmt"
	"porter/api/src/lib/github"
	"porter/api/src/lib/redis"
	"bytes"
	"porter/api/src/lib/dockerfile"
	"strings"
	"time"
	"sync"
)

const redisCacheKey = "PORTER:REPOSITORIES"

func exposedPortLookup(repository github.RepositoryResponse,  msgCh chan message.Message) (github.Repository)  {
	ports := make([]string, 0)

	message.SendMessage(fmt.Sprintf( "[%s] - fetching Dockerfile", repository.Name),"green", msgCh, false)
	body, err := github.FetchFileContent(repository.Name, "Dockerfile")
	if err != nil || body == "" {
		message.SendMessage(fmt.Sprintf( "[%s] - couldn't get Dockerfile", repository.Name),"red", msgCh, false)
	} else {
		buf := bytes.NewBuffer([]byte(body))
		message.SendMessage(fmt.Sprintf( "[%s] - parsing Dockerfile into an AST", repository.Name),"green", msgCh, false)
		port, _ := dockerfile.ExtractExposedPort(buf)

		if port != "" {
			message.SendMessage(fmt.Sprintf( "[%s] - found the port %d in the Dockerfile", port, repository.Name),"green", msgCh, false)
			ports = append(ports, port)
		}

		message.SendMessage(fmt.Sprintf( "[%s] - did not found any port in the Dockerfile", repository.Name),"red", msgCh, false)
	}

	if len(ports) == 0 {
		message.SendMessage(fmt.Sprintf( "[%s] - searching ports in k8 config files", repository.Name),"green", msgCh, false)
		k8Ports := github.FetchK8Port(repository)

		portsLength := len(k8Ports)

		if portsLength == 0 {
			message.SendMessage(fmt.Sprintf( "[%s] - did not find any port in k8 config files", repository.Name),"red", msgCh, false)
		} else {
			message.SendMessage(fmt.Sprintf( "[%s] - found ports in k8 config files %s", repository.Name, strings.Join(k8Ports, ",")),"green", msgCh, false)
		}

		for i:= 0; i<portsLength; i++  {
			ports = append(ports,k8Ports[i])
		}
	}

	return github.Repository{
		ExposedPorts: ports,
		Name: repository.Name,
		Url: repository.Url,
		Html_url: repository.Html_url,
	}
}

func usedPortsHandler(ch chan github.Repository, msgCh chan message.Message) {
	message.SendMessage("Fetching repositories in parallel","green" , msgCh, false)

	var wg sync.WaitGroup
	repositories := github.FetchAllRepositories()
	repositoriesPortCh := make(chan github.Repository)
	repositoriesLength := len(repositories)

	if repositoriesLength == 0{
		message.SendMessage("Aborting session, could not find repositories", "red", msgCh, true)
		close(ch)
	}

	message.SendMessage(fmt.Sprintf( "Done fetching repositories, %d total", len(repositories)),"green" , msgCh, false)

	wg.Add(repositoriesLength)

	for i := 0; i < repositoriesLength ; i++ {
		go func(i int) {
			if pushToChannelFromCache(repositories[i], repositoriesPortCh){
				message.SendMessage(fmt.Sprintf( "[%s] - fetched from cache", repositories[i].Name),"yellow" , msgCh, false)
			} else {
				fetchAndPushToChannel(repositories[i], repositoriesPortCh, msgCh)
				message.SendMessage(fmt.Sprintf( "[%s] - fetching ports from github.com", repositories[i].Name),"green" , msgCh, false)
			}
		}(i)
	}

	for i := range repositoriesPortCh {
		if len(i.ExposedPorts) > 0 {
			ch <- i
		}
		wg.Done()
	}

	wg.Wait()
	close(ch)
}

func pushToChannelFromCache(repositoryResponse github.RepositoryResponse, ch chan github.Repository) bool {
	redisClient := redis.GetClient()
	value := redisClient.Get(redisKey(repositoryResponse.Name)).Val()
	if value == "" {
		return false
	}

	var repository github.Repository
	repository.UnMarshalBinary([]byte(value))

	ch <- repository
	return true
}

func fetchAndPushToChannel(repositoryResponse github.RepositoryResponse, repositoriesPortsCh chan github.Repository, messagesCh chan message.Message) {
	redisClient := redis.GetClient()
	repository := exposedPortLookup(repositoryResponse, messagesCh)
	redisClient.Set(redisKey(repository.Name), repository , time.Hour)
	repositoriesPortsCh <- repository
}

func redisKey(repositoryName string) string {
	return redisCacheKey + ":" + repositoryName
}