package govee

import "testing"

func TestNewVersion(t *testing.T) {
	expect := "1.2.3"

	vconf := VersionConfig{
		VersionString: "1.2.3",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "linux",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "prod",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}

	if v.Semver() != expect {
		t.Errorf("Expected %s, got %s", expect, v.Semver())
	}
}

func TestNewPreVersion(t *testing.T) {
	expect := "2-ga1b2c3d"

	vconf := VersionConfig{
		VersionString: "1.2.3-2-ga1b2c3d",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "darwin",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "prod",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}
	if v.Pre() != expect {
		t.Errorf("Expected %s, got %s", expect, v.Pre())
	}
}

func TestMajorVersion(t *testing.T) {
	expect := 1

	vconf := VersionConfig{
		VersionString: "1.2.3-2-ga1b2c3d",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "linux",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "production",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}
	if v.Major() != expect {
		t.Errorf("Expected %d, got %d", expect, v.Major())
	}
}

func TestMinorVesion(t *testing.T) {
	expect := 2

	vconf := VersionConfig{
		VersionString: "1.2.3-2-ga1b2c3d",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "linux",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "prod",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}
	if v.Minor() != expect {
		t.Errorf("Expected %d, got %d", expect, v.Patch())
	}
}

func TestPatchVersion(t *testing.T) {
	expect := 3

	vconf := VersionConfig{
		VersionString: "1.2.3-2-ga1b2c3d",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "linux",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "prod",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}
	if v.Patch() != expect {
		t.Errorf("Expected %d, got %d", expect, v.Patch())
	}
}

func TestVersionWarnings(t *testing.T) {
	expect := "This version is tagged as a pre-release \"[2-ga1b2c3d]\". Please don't use in production."
	expectCount := 2

	vconf := VersionConfig{
		VersionString: "1.2.3-2-ga1b2c3d",
		GitHash:       "1234567890abcdef",
		GitBranch:     "testing",
		GitUser:       "Jane Doe",
		OS:            "linux",
		Arch:          "amd64",
		Compiler:      "go1.11.1",
		Release:       "test",
		TStamp:        "Thu Feb 14 15:04:05 SAST 2019",
	}

	v, err := NewVersion(&vconf)
	if err != nil {
		t.Error(err)
	}
	warnings := v.Warnings()

	if len(warnings) != expectCount {
		t.Errorf("Expected %d warnings, got %d", expectCount, len(warnings))
	}
	if warnings[0] != expect {
		t.Errorf("Expected %s, got %s", expect, warnings[0])
	}
}
