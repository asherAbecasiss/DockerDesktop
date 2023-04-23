package types

type Memory struct {
	MemTotal     float64 `json:"total"`
	MemFree      float64 `json:"free"`
	MemAvailable float64 `json:"avilable"`
	MemPercent   int     `json:"mempercent"`
}

type Ip struct {
	LocalIp string `json:"hostip"`
}

type Cpu struct {
	CpuUsage float64 `json:"cpuUsage"`
	Busy     float64 `json:"busy"`
	Total    float64 `json:"total"`
}

type ProcessList struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	CpuPercent float64 `json:"CpuPercent"`
}

type Agents struct {
	Alist []Agent `json:"ip"`
}

type Agent struct {
	Ip   string `json:"ip"`
	Host string `json:"host"`
}
