package models

type CpuInfo struct {
	LoadPercentage float64
	Temperature    float64
}

type Meminfo struct {
	LoadPercentage float64
	Total          uint64 // bytes
	Used           uint64 // bytes
}
