package manager

const (
	// ErrInitNSUpdate error code for problems while setting up NSUpdate
	ErrInitNSUpdate = iota
)

const (
	// BasePath defines the base path for storing and reading application specific files
	BasePath = "/data"

	// SANDMAN_DNS_TTL environment variable identifier for the time-to-live to be applied
	SANDMAN_DNS_TTL = "SANDMAN_DNS_TTL"

	// SANDMAN_DNS_REMOVAL_DELAY environment variable identifier for the removal delay time to be applied
	SANDMAN_DNS_REMOVAL_DELAY = "SANDMAN_DNS_REMOVAL_DELAY"
)
