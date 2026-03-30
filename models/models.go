package models

// General response model
type SysInfo struct {
	CPU   CpuInfo    `json:"cpu"`
	Mem   MemInfo    `json:"mem"`
	Disks []DiskInfo `json:"disks"`
	Net   []NetInfo    `json:"net"`
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

type NetInfo struct {
	Name  string  `json:"name"`
	RMbps float64 `json:"rmbps"`
	SMbps float64 `json:"smbps"`
}
