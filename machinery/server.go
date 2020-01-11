package machinery

import (
	"fmt"
	"sync"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var configPath = "machinery/config.yml"
var consumerTag = "worker"

type server struct {
	server *machinery.Server
}

var instance *server
var once sync.Once

//GetServer gets server singleton instance
func GetServer() *server {
	once.Do(func() {
		instance = &server{}
		instance.server, _ = startServer()
	})
	return instance
}

func loadConfig() (*config.Config, error) {
	return config.NewFromYaml(configPath, true)
}

func startServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Create server instance
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	//Register tasks
	tasks := map[string]interface{}{
		"getTopGitHubRepoByLanguage": GetTopGitHubRepoByLanguage,
		"printAllResults":            PrintAllResults,
	}

	err = server.RegisterTasks(tasks)
	if err != nil {
		return nil, err
	}

	
	return server, nil
}

func (s *server) StartWorkers(errorsChan chan error){
	worker := s.server.NewWorker(consumerTag, 10)
	worker.LaunchAsync(errorsChan)
}

func (s *server) SendGitHubTask(language string) {
	var signature = tasks.Signature{
		Name: "getTopGitHubRepoByLanguage",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: fmt.Sprintf("%v", language),
			},
		},
	}
	s.server.SendTask(&signature)
}
