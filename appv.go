// Package appv lets you set your application's version info from the application's Git
// repository at compile time using -ldflags.
package appv

type Versioner interface {
	// Semver returns the complete semantic version number as a string.
	Semver() string

	// Major returns the major version number.
	Major() int

	// Minor returns the minor version number.
	Minor() int

	// Patch returns the patch version number.
	Patch() int

	// Pre returns the pre-release version information as a string.
	Pre() string

	// Warnings returns the version warnings as []string.
	Warnings() []string

	// VersionError returns the version error.
	VError() error
}
