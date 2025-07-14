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

type ServiceRecord struct {
	Location       string `json:"location"`
	WANIP          string `json:"wan_ip"`
	LANIP          string `json:"lan_ip"`
	ConnectionType string `json:"connection_type"`
	Bandwidth      any    `json:"bandwidth"`
	LineType       string `json:"line_type"`
	ServiceNumber  any    `json:"service_number"`
	AccountNumber  any    `json:"account_number"`
}
