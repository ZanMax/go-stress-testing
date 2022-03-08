package main

import "fmt"

var googleAgents = []string{}
var userAgents = []string{}

func getRandomGoogleAgent() string {
	googleAgents, err := readLines("google-agents.txt")
	if err != nil {
		fmt.Println("Fail load agents from file", err.Error())
	}
	agent := randChoice(googleAgents)
	return agent
}

func getRandomAgent() string {
	userAgents, err := readLines("user-agents.txt")
	if err != nil {
		fmt.Println("Fail load agents from file", err.Error())
	}
	agent := randChoice(userAgents)
	return agent
}
