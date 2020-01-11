package main

import (
	"fmt"
	"github.com/evandroferreiras/machinery-meetup/route"
	"github.com/evandroferreiras/machinery-meetup/machinery"
)

func main() {
	fmt.Println(" ___  ___           _     _                       ")
	fmt.Println(" |  \\/  |          | |   (_)                      ")
	fmt.Println(" | .  . | __ _  ___| |__  _ _ __   ___ _ __ _   _ ")
	fmt.Println(" | |\\/| |/ _` |/ __| '_ \\| | '_ \\ / _ \\ '__| | | |")
	fmt.Println(" | |  | | (_| | (__| | | | | | | |  __/ |  | |_| |")
	fmt.Println(" \\_|  |_/\\__,_|\\___|_| |_|_|_| |_|\\___|_|   \\__, |")
	fmt.Println("                                             __/ |")
	fmt.Println("                                            |___/ ")

	errorsChan := make(chan error)
	var server = machinery.GetServer()

	server.StartWorkers(errorsChan)


	router := route.Init()
	router.Logger.Fatal(router.Start(":1323"))
	if err := <- errorsChan; err != nil {
		panic(err)
	} 	
}
