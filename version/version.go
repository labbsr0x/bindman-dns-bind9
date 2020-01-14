package version

import "fmt"

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildTime = "unknown-buildtime"
)

const versionF = `Bindman-DNS Bind9
  Version: %s
  GitCommit: %s
  BuildTime: %s
`

func FormattedMessage() string {
	return fmt.Sprintf(versionF, Version, GitCommit, BuildTime)
}
