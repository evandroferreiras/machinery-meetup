package machinery

import (
	githubApi "github.com/evandroferreiras/machinery-meetup/machinery/github_api"
)

// GitHubResponse ...
type GitHubResponse struct {
	Language     string   `json:"language,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}


// GetRepositoriesByLanguageAndPage ...
func GetRepositoriesByLanguageAndPage(language string, page int) ([]string, error) {
	repositories, err := githubApi.GetRepositoriesByLanguageAndPage(language, page)
	if err != nil {
		return nil, err
	}
	return repositories, err
}

// SaveConsolidatedResults ...
func SaveConsolidatedResults(args ... []string) ([]string, error) {
	consolidatedResults := make([]string, 0)
	for _, r := range args {
		consolidatedResults = append(consolidatedResults, r...)
	}
	return consolidatedResults, nil
}
