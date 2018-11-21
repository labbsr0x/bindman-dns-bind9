package nsupdate

// NSUpdate holds the information necessary to successfully run nsupdate requests
type NSUpdate struct {
	Server      string
	Port        string
	KeyFilePath string
}

// New constructs a new NSUpdate instance from environment variables
func New(server, port, keyFilePath string) (*NSUpdate, error) {
	return &NSUpdate{
		Server:      server,
		Port:        port,
		KeyFilePath: keyFilePath,
	}, nil
}

// RemoveRR removes a Resource Record
func (nsu *NSUpdate) RemoveRR(name string) (success bool, err error) {
	success = true
	err = nil
	// TODO
	return
}

// AddRR adds a Resource Record
func (nsu *NSUpdate) AddRR(name string, ipaddr string, ttl int) (success bool, err error) {
	success = true
	err = nil
	// TODO
	return
}
