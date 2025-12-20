package main

func setupMock() stores {
	// test station store
	ss := inMemoryStationStore{
		"testStation": station{
			Name: "testStation",
			Platforms: map[int]bool{
				1: true,
				2: true,
				3: true,
			},
			Neighbors: []string{
				"anotherTestStation",
			},
		},

		"anotherTestStation": station{
			Name: "anotherTestStation",
			Platforms: map[int]bool{
				1: true,
				2: true,
			},
			Neighbors: []string{
				"testStation",
			},
		},
	}

	testNodeIdOne := "03d3a549-4a2b-44c5-ac81-4548d89a5340"
	testNodeIdTwo := "e3c91099-5c65-4b9c-8646-0acb39dfba6d"
	rs := inMemoryRouteStore{
		"8d7e2ad1-2c16-44f7-9c55-bc181d83900b": route{
			Id: "8d7e2ad1-2c16-44f7-9c55-bc181d83900b",
			StartNode: routeStationNode{
				Id:           testNodeIdOne,
				StationName:  "testStation",
				PreviousNode: nil,
				NextNode:     &testNodeIdTwo,
			},
			EndNode: routeStationNode{
				Id:           testNodeIdTwo,
				StationName:  "anotherTestStation",
				PreviousNode: &testNodeIdOne,
				NextNode:     nil,
			},
		},
	}
	return stores{
		ss: ss,
		rs: rs,
	}
}
