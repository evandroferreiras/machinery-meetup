package githubApi 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var gitHubURL = "https://api.github.com/search/repositories?q=language:%s&sort=stars&order=desc&page=%d"

type repoItems struct {
	FullName        string `json:"full_name,omitempty"`
	Name            string `json:"name,omitempty"`
	HTMLURL         string `json:"html_url,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	Forks           int    `json:"forks,omitempty"`
}

type repo struct {
	Items []repoItems
}

// GetRepositoriesByLanguageAndPage ...
func GetRepositoriesByLanguageAndPage(language string, page int) ([]string, error) {
	uri := fmt.Sprintf(gitHubURL, url.QueryEscape(language), page)
	fmt.Println(uri)
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	var dat repo
	if err := json.Unmarshal(data, &dat); err != nil {
		return nil, err
	}

	result := make([]string, 0)
	if len(dat.Items) > 0 {
		for _, value := range dat.Items {
			result = append(result, value.Name)
		}
	}

	return result, nil

}
