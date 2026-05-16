package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/JordanllHarper/trainsgo/shared"
)

type (
	opFunc  func()
	handler struct {
		helpText string
		f        opFunc
	}
	handlers map[string]handler
)

func getTrains(c http.Client, baseUrl string) []shared.Train {
	res, err := c.Get(fmt.Sprintf("%s/trains", baseUrl))
	if err != nil {
		panic(err)
	}
	var trains []shared.Train
	if err := json.NewDecoder(res.Body).Decode(&trains); err != nil {
		panic(err)
	}
	return trains
}

func main() {
	baseUrl := "http://127.0.0.1:8080"
	c := http.Client{}
	s := bufio.NewScanner(os.Stdin)
	handlers := map[string](handler){
		"gt": handler{"[g]et [t]rains", func() {
			trains := getTrains(c, baseUrl)
			for _, v := range trains {
				fmt.Println(v)
			}
		}},
	}
	fmt.Println("Started")
	printHelp(handlers)
	for {
		text := scanAndReadText(s)
		handler, exists := handlers[text]
		if !exists {
			fmt.Println("Invalid command, please type one of the following options")
			printHelp(handlers)
			continue
		}
		handler.f()
	}
}

func scanAndReadText(s *bufio.Scanner) string {
	s.Scan()
	return s.Text()
}

func printHelp(h handlers) {
	for k, v := range h {
		fmt.Println(k, "=>", v.helpText)
	}
}
