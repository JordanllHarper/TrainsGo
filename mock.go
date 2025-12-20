package main

func setupMock() dependencies {
	// Define a simple 2 station connection
	// Creates a route back and fourth between them

	stationOne := station{
		Name: "stationOne",
		Platforms: map[int]bool{
			1: true,
			2: true,
			3: true,
		},
		Neighbors: []string{
			"stationTwo",
		},
	}
	stationTwo := station{
		Name: "stationTwo",
		Platforms: map[int]bool{
			1: true,
			2: true,
		},
		Neighbors: []string{
			"stationOne",
		},
	}

	ss := inMemoryStationStore{
		"stationOne": stationOne,
		"stationTwo": stationTwo,
	}

	rs := inMemoryRouteStore{
		"8d7e2ad1-2c16-44f7-9c55-bc181d83900b": route{
			RouteId: "8d7e2ad1-2c16-44f7-9c55-bc181d83900b",
			Route: map[int]station{
				1: stationOne,
				2: stationTwo,
			},
		},

		"a4b6776b-5f38-456b-9d8c-7abd7a0f10d4": route{
			RouteId: "a4b6776b-5f38-456b-9d8c-7abd7a0f10d4",
			Route: map[int]station{
				1: stationTwo,
				2: stationOne,
			},
		},
	}
	// NOTE: Mocking this might include setting up defined routes already :D
	rb := inMemoryRouteBuilder{}
	return dependencies{
		ss: ss,
		rs: rs,
		rb: rb,
	}
}
