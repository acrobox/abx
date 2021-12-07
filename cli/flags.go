package cli

import "strconv"

// flags represents the command flag parameters.
type flags struct {
	addr    string // acrobox.io service addr
	auth    string // acrobox.io api token
	host    string // machine hostname
	port    string // machine ssh port without the colon
	verbose bool
	init    flagsInit
	cancel  flagsCancel
	renew   flagsRenew
	destroy flagsDestroy
	status  flagsStatus
	restore flagsRestore
}

// flagsInit represents the flags for initializing a new machine.
type flagsInit struct {
	Product     string `json:"product"`      // acrobox.io product
	Name        string `json:"name"`         // default "acrobox"
	Region      string `json:"region"`       // default "nyc1"
	Size        string `json:"size"`         // default "s-1vcp-1gb-intel"
	DataSize    int    `json:"data_size"`    // default 1GB
	PublicKey   string `json:"public_key"`   // ssh authorized key
	AccessToken string `json:"access_token"` // $DIGITALOCEAN_ACCESS_TOKEN
	fullYear    bool
	force       bool
}

// flagsCancel represents the flags for cancelling service.
type flagsCancel struct {
	force bool
}

// flagsRenew represents the flags for cancelling service.
type flagsRenew struct {
	force bool
}

// flagsDestroy represents the flags for destroying a machine.
type flagsDestroy struct {
	AccessToken string `json:"access_token"` // $DIGITALOCEAN_ACCESS_TOKEN
	force       bool
}

// flagsStatus represents the flags for machine status.
type flagsStatus struct {
	format string
}

// flagsRestore represents the flags for restoring a machine.
type flagsRestore struct {
	force bool
}

// flagInt represents an integer flag.
type flagInt struct{}

// Parse returns value as an integer in base 10 and bit size 0.
//
// Parse implements the FlagKind interface.
func (f flagInt) Parse(value string) interface{} {
	i, _ := strconv.Atoi(value)
	return i
}

// HasArg implements the FlagKind interface.
func (f flagInt) HasArg() bool {
	return true
}
