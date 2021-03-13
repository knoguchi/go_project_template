// Package version executes and returns the version string
// for the currently running process.
package version

import (
	"fmt"
)

const (
	Major = 0          // Major version component of the current release
	Minor = 0          // Minor version component of the current release
	Patch = 1          // Patch version component of the current release
	Meta  = "unstable" // Version metadata to append to the version string
)

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}()

// WithMeta holds the textual version string including the metadata.
var WithMeta = func() string {
	v := Version
	if Meta != "" {
		v += "-" + Meta
	}
	return v
}()

func WithCommit(gitCommit, gitDate string) string {
	vsn := WithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	if (Meta != "stable") && (gitDate != "") {
		vsn += "-" + gitDate
	}
	return vsn
}
