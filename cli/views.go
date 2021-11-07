package cli

import "time"

type acroboxdStatus struct {
	SystemBootedAt   time.Time     `json:"system_booted_at"`
	AcroboxStartedAt time.Time     `json:"acrobox_started_at"`
	LicenseUpdatedAt time.Time     `json:"license_updated_at"`
	NextBackupAt     time.Time     `json:"next_backup_at"`
	LastBackupAt     time.Time     `json:"last_backup_at"`
	LastBackupTook   time.Duration `json:"last_backup_took"`
	NextUpdateAt     time.Time     `json:"next_update_at"`
	LastUpdateAt     time.Time     `json:"last_update_at"`
	LastUpdateTook   time.Duration `json:"last_update_took"`
}
