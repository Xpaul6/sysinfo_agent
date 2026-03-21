package models

// General response model
type SysInfo struct {
	CPU  CpuInfo    `json:"cpu"`
	Mem  MemInfo    `json:"mem"`
	Disks []DiskInfo `json:"disk"`
}

// Specific hardware models
type CpuInfo struct {
	LoadPercentage float64 `json:"load"`
	Temperature    float64 `json:"temperature"`
}

type MemInfo struct {
	LoadPercentage float64 `json:"load"`
	Total          uint64  `json:"total"` // bytes
	Used           uint64  `json:"used"`  // bytes
}

type DiskInfo struct {
	MountPoint string `json:"mountpoint"`
	Total      uint64 `json:"total"` // bytes
	Used       uint64 `json:"used"`  // bytes
}
