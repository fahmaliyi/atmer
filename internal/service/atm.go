package service

type Machine struct {
	IP   string
	Name string
}

type PingResult struct {
	IP     string
	Name   string
	Status string
}
