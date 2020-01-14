package machinery

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"sync"
	"time"
)

var CONFIG_PATH = "machinery/config.yml"
var CONSUMER_TAG = "worker"
var UNLIMITED_CONCURRENCY_TASKS = 10

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

	//Register signaturesToRegister
	signaturesToRegister := map[string]interface{}{
		"getRepositoriesByLanguageAndPage": GetRepositoriesByLanguageAndPage,
		"saveConsolidatedResults":          SaveConsolidatedResults,
	}

	err = server.RegisterTasks(signaturesToRegister)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *server) StartWorkers(errorsChan chan error) {
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
				Type:  "int",
				Value: 1,
			},
		},
	}
	_, _ = s.server.SendTask(&signature)
}

func (s *server) GenerateConsolidatedReport(language string) {
	var githubTasksSignatures = make([]*tasks.Signature, 0)
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
		githubTasksSignatures = append(githubTasksSignatures, &ta)
	}

	group, _ := tasks.NewGroup(githubTasksSignatures...)

	var saveConsolidatedResultsSignature = tasks.Signature{
		Name: "saveConsolidatedResults",
	}
	
	chord, _ := tasks.NewChord(group, &saveConsolidatedResultsSignature)
	chordAsyncResult, err := s.server.SendChord(chord, UNLIMITED_CONCURRENCY_TASKS)
	if err != nil {
		fmt.Println("Could not send chord: %s", err.Error())
	}

	results, err := chordAsyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		// getting result of a chord failed
		// do something with the error
	}
	for _, result := range results {
		fmt.Println(result.Interface())
	}

}
