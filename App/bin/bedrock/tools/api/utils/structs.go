package utils

type Backup struct {
	Name         string `json:"Name"`
	Size         int64  `json:"Size"`
	CreationDate string `json:"CreationDate"`
}

type ServerInfo struct {
	CPUUsage     float64 `json:"CPUUsage"`
	FreeDisk     int64   `json:"FreeDisk"`
	FreeMemory   int64   `json:"FreeMemory"`
	MemoryUsage  float64 `json:"MemoryUsage"`
	PortsUsed    string  `json:"PortsUsed"`
	ServerIP     string  `json:"ServerIP"`
	ServerSize   int64   `json:"ServerSize"`
	ServerStatus string  `json:"ServerStatus"`
	TotalDisk    int64   `json:"TotalDisk"`
	TotalMemory  int64   `json:"TotalMemory"`
	UsedDisk     struct {
		Percentage float64 `json:"Percentage"`
		Size       int64   `json:"Size"`
	} `json:"UsedDisk"`
}
