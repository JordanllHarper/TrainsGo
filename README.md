# Rail Management System in Go

## Feature list

- [x] CRUD stations.
- [ ] CRUD routes  - paths between stations.
- [ ] CRUD journeys - timetabled instances where a train will pick up and drop off passengers along a particular route.
- [ ] Update the status of an ongoing journey - where it is on the route, estimated arrival time to each


## Stations

- Named entity on a rail network
- Contain references to their neighbors.
- Platforms and their availability

## Routes

- An ordered list of stations from a start station to an end station.
- Path from start to end

## Journey

- A timetabled train fulfills a particular route.
- Contains the route, a timetable of expected and estimated times, as well as a status - at station, traveling, waiting, etc.

## Service

- A repeating occurrence of journeys along a route - automatic scheduling
