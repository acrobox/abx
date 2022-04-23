package cli

import "time"

type acroboxdStatus struct {
	SystemBootedAt            time.Time     `json:"system_booted_at"`
	AcroboxStartedAt          time.Time     `json:"acrobox_started_at"`
	LicenseUpdatedAt          time.Time     `json:"license_updated_at"`
	NextBackupAt              time.Time     `json:"next_backup_at"`
	LastBackupAt              time.Time     `json:"last_backup_at"`
	LastBackupTook            time.Duration `json:"last_backup_took"`
	NextUpdateAt              time.Time     `json:"next_update_at"`
	LastUpdateAt              time.Time     `json:"last_update_at"`
	LastUpdateTook            time.Duration `json:"last_update_took"`
	ProxyActiveRequests       uint64        `json:"proxy_active_requests"`
	ProxyActiveRequestsPeak   uint64        `json:"proxy_active_requests_peak"`
	ProxyActiveRequestsPeakAt time.Time     `json:"proxy_active_requests_peak_at"`
	ProxyRequestsPerSec       uint64        `json:"proxy_requests_per_sec"`
	ProxyRequestsPerSecPeak   uint64        `json:"proxy_requests_per_sec_peak"`
	ProxyRequestsPerSecPeakAt time.Time     `json:"proxy_requests_per_sec_peak_at"`
	ProxyLatencyP99           []float64     `json:"proxy_latency_p99"`
	MemStatsGoRoutines        int           `json:"memstats_go_routines"`  // count
	MemStatsMemoryTotal       uint64        `json:"memstats_memory_total"` // bytes
	MemStatsMemoryAlloc       uint64        `json:"memstats_memory_alloc"` // bytes
	MemStatsMemoryCount       uint64        `json:"memstats_memory_count"` // count
}
