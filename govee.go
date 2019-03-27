package govee

import (
	"fmt"
	"time"

	"github.com/blang/semver"
)

// Version represents a semantic version number.
type Version struct {
	semver    semver.Version
	githash   string
	gitbranch string
	gituser   string
	os        string
	arch      string
	compiler  string
	release   string
	timestamp time.Time
	warnings  []string
	err       error
}

// VersionConfig represents the version coniguration.
type VersionConfig struct {
	VersionString string // semver string representation.
	GitHash       string
	GitBranch     string
	GitUser       string
	OS            string
	Arch          string
	Compiler      string
	Release       string
	TStamp        string
}

// NewVersion creates a new version object from a VersionConfig.
func NewVersion(c *VersionConfig) (Version, error) {
	var err error
	v := Version{}
	v.githash = c.GitHash
	v.gitbranch = c.GitBranch
	v.gituser = c.GitUser
	v.os = c.OS
	v.arch = c.Arch
	v.compiler = c.Compiler
	v.release = c.Release

	v.semver, err = semver.Make(c.VersionString)
	if err != nil {
		return Version{}, err
	}

	v.timestamp, err = time.Parse(time.UnixDate, c.TStamp)
	if err != nil {
		return Version{}, err
	}

	if len(v.semver.Pre) > 0 {
		warning := fmt.Sprintf(
			"This version is tagged as a pre-release \"%+v\". Please don't use in production.",
			v.semver.Pre,
		)
		v.warnings = append(v.warnings, warning)
	}

	if v.release != "production" && v.release != "prod" {
		warning := fmt.Sprintf(
			"This version is tagged as release \"%s\". Please don't use in production.",
			v.release,
		)
		v.warnings = append(v.warnings, warning)
	}
	return v, nil
}

// Implement the Stringer interface.
func (v Version) String() string {
	return v.semver.String()
}

// Semver returns the complete semantic version number as a string.
func (v Version) Semver() string {
	return v.semver.String()
}

// Major returns the major version number.
func (v Version) Major() int {
	return int(v.semver.Major)
}

// Minor return the minor version number.
func (v Version) Minor() int {
	return int(v.semver.Minor)
}

// Patch returns the patch version number.
func (v Version) Patch() int {
	return int(v.semver.Patch)
}

// Pre returns the pre-release version information.
func (v Version) Pre() string {
	return fmt.Sprintf("%v", v.semver.Pre[0])
}

// Warnings returns the version warnings.
func (v Version) Warnings() []string {
	return v.warnings
}

// Err returns the version error.
func (v Version) Err() error {
	return v.err
}

// GitHash returns the git hash.
func (v Version) GitHash() string {
	return v.githash
}

// GitBranch returns the git branch.
func (v Version) GitBranch() string {
	return v.gitbranch
}

// GitUser returns the git user.
func (v Version) GitUser() string {
	return v.gituser
}

// OS returns the operating system.
func (v Version) OS() string {
	return v.os
}

// Arch returns the architecture.
func (v Version) Arch() string {
	return v.arch
}

// Release returns the release information.
func (v Version) Release() string {
	return v.release
}

// TStamp returns the timestamp,
func (v Version) TStamp() string {
	return v.timestamp.Format(time.RFC3339)
}

// Compiler returns the compiler version.
func (v Version) Compiler() string {
	return v.compiler
}
