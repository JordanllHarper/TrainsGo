package main

type (
	station struct {
		// unique identifier of the station as well as descriptive name
		Name string `json:"name"`
		// set of platform numbers
		Platforms map[int]struct{} `json:"platforms"`

		Neighbors []string `json:"neighbors"`
	}

	// A path to take from a start node to an end node
	route struct {
		RouteId string `json:"routeId"`
		// map of order to station
		Route map[int]station `json:"route"`
	}
)
