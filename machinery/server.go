package machinery

import (
	"fmt"
	"sync"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var CONFIG_PATH = "machinery/config.yml"
var CONSUMER_TAG = "worker"
var UNLIMITED_CONCURRENCY_TASKS = 0

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
	return config.NewFromYaml(CONFIG_PATH, true)
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

	//Register signatures
	signatures := map[string]interface{}{
		"getRepositoriesByLanguageAndPage": GetRepositoriesByLanguageAndPage,
		"saveConsolidatedResults":  SaveConsolidatedResults,
	}

	err = server.RegisterTasks(signatures)
	if err != nil {
		return nil, err
	}

	
	return server, nil
}

func (s *server) StartWorkers(errorsChan chan error){
	worker := s.server.NewWorker(CONSUMER_TAG, 10)
	worker.LaunchAsync(errorsChan)
}

func (s *server) GenerateReport(language string) {
	var signature = tasks.Signature{
		Name: "getRepositoriesByLanguageAndPage",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: fmt.Sprintf("%v", language),
			},
			{
				Type: "int",
				Value: 1,
			},
		},
	}
	_, _ = s.server.SendTask(&signature)
}

func (s *server) GenerateConsolidatedReport(language string) {

	var signatures = make([]*tasks.Signature, 0)
	for i := 1; i <= 10; i++ {
		var ta = tasks.Signature{
			Name: "getRepositoriesByLanguageAndPage",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: fmt.Sprintf("%v", language),
				},
				{
					Type:  "int",
					Value: i,
				},
			},
		}
		signatures = append(signatures, &ta)
	}
	group, _ := tasks.NewGroup(signatures ...)
	_, err := s.server.SendGroup(group, UNLIMITED_CONCURRENCY_TASKS)
	if err != nil {

	}

}