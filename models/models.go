package models

type CpuInfo struct {
	LoadPercentage float64 `json:"load"`
	Temperature    float64 `json:"temperature"`
}

type MemInfo struct {
	LoadPercentage float64 `json:"load"`
	Total          uint64  `json:"total"` // bytes
	Used           uint64  `json:"used"`  // bytes
}
