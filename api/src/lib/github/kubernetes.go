package github

import (
	"regexp"
	"fmt"
	"strconv"
	"strings"
)

var portRegex = regexp.MustCompile("servicePort: ([0-9]{4})")

func FetchK8Port(repository RepositoryResponse) ([] string)  {
	files, err := FetchFolder(repository.Url, "kube/core")
	if err != nil{
		return []string{}
	}

	ports := make([]string,0)
	for i:=0; i< len(files) ; i++  {
		if strings.Contains(files[i].Name, "deployment"){
			content, err := FetchFileContent(repository.Url, files[i].Path)

			if err != nil {
				continue
			}

			results := portRegex.FindAllString(string(content), -1)
			for i := 0; i< len(results)  ; i++  {
				port := strings.Split(results[i], " ")[1]
				if _, err := strconv.Atoi(port); err != nil {
					fmt.Println("Can't parse " + port)
				} else {
					ports = append(ports, port)
				}
			}
		}
	}

	return ports
}
