package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	listStation   = "ls"
	newStation    = "ns"
	deleteStation = "ds"
	newRoute      = "nr"
	help          = "h"
	quit          = "q"
)

const url = "http://localhost:8080"

type (
	postStationRequest struct {
		Name      string   `json:"name"`
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}

	getStationResponse struct {
		Name      string   `json:"name"`
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}
)

func showHelp() {
	printCommand(listStation, "list created stations")
	printCommand(newStation, "create new station")
	printCommand(deleteStation, "delete station")
	printCommand(newRoute, "create new route between stations")
	printCommand(help, "delete new station")
	printCommand(quit, "create new route between stations")
}

func main() {
	fmt.Println("Started: ")
	showHelp()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		err := scanner.Err()
		if mustNotBeErr(err) {
			continue
		}
		s := scanner.Text()
		switch s {
		case help:
			showHelp()
		case quit:
			fmt.Println("Quitting...")
			return
		case listStation:
			handleListStations()
		case newStation:
			handleNewStation(*scanner)
		case deleteStation:
			handleDeleteStation(*scanner)
		default:
			fmt.Println("Unrecognised command:", s)
		}
	}
}

func mustNotBeErr(err error) bool {
	if err != nil {
		fmt.Println("Error:", err)
		return true
	}
	return false
}

func getPath(path string) string { return url + path }

func handleDeleteStation(s bufio.Scanner) {
	fmt.Println("Name to delete")
	s.Scan()
	if mustNotBeErr(s.Err()) {
		return
	}
	name := s.Text()
	client := http.DefaultClient

	req, err := http.NewRequest("DELETE", getPath("/stations/")+name, nil)

	resp, err := client.Do(req)
	if mustNotBeErr(err) {
		return
	}
	fmt.Println("Status code:", resp.StatusCode)
}

func handleNewStation(scanner bufio.Scanner) {
	fmt.Println("Name:")
	scanner.Scan()
	if mustNotBeErr(scanner.Err()) {
		return
	}
	name := scanner.Text()
	fmt.Println("Num of platforms:")

	scanner.Scan()
	if mustNotBeErr(scanner.Err()) {
		return
	}
	numPlatformsInput := scanner.Text()
	numPlatforms, err := strconv.Atoi(numPlatformsInput)
	if mustNotBeErr(err) {
		return
	}
	platforms := []int{}
	for i := range numPlatforms {
		platforms = append(platforms, i)
	}

	fmt.Println("Neighbors")
	scanner.Scan()
	if mustNotBeErr(scanner.Err()) {
		return
	}
	// neighbor1 neighbor2
	neighborsInput := scanner.Text()
	neighbors := strings.Split(neighborsInput, ":")
	st := postStationRequest{
		Name:      name,
		Platforms: platforms,
		Neighbors: neighbors,
	}

	json, err := json.Marshal(st)
	if mustNotBeErr(err) {
		return
	}

	resp, err := http.Post(getPath("/stations"), "application/json", bytes.NewReader(json))
	if mustNotBeErr(err) {
		return
	}
	fmt.Println("Response:", resp.StatusCode)
}

func handleListStations() {
	resp, err := http.Get(getPath("/stations"))
	if mustNotBeErr(err) {
		return
	}
	fmt.Println("Status code:", resp.StatusCode)

	var stations []getStationResponse
	err = json.NewDecoder(resp.Body).Decode(&stations)
	if mustNotBeErr(err) {
		return
	}
	fmt.Println("Stations:")
	fmt.Println("Name | Station | Platform")
	for _, v := range stations {
		fmt.Println(v.Name, v.Neighbors, v.Platforms)
	}
}

func printCommand(cmd string, desc string) {
	fmt.Printf("%s - %s\n", cmd, desc)
}
