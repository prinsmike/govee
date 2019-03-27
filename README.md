# govee

Go semantic versioning library that generates version information from a git repository.

## Introduction

This package provides the ability to embed semantic version ([SemVer](https://semver.org/)) information in your Go application at compile time. The intent is to use the Go compiler's -ldflags option to pull information from the application's Git repository, and embed it in the application binary at compile time. Use the following instructions to set up an example application that demonstrates the **govee** functionality.

## Installing the govee package

**govee** can be installed using `go get`:

```
go get -u github.com/prinsmike/govee
```

You can also use Go modules to manage this package as a dependency in your project, but this is left as an exercise for the reader.

## Create the main entry point for your application

Create a new file and name it `main.go`, and paste the following code in it:

**main.go:**
```golang
package main

func main() {
	a := app{}
	a.name = "MyApp"
	a.setVersion()
	a.opts()
	a.printVersion()
}
```

## Create the app source file

Create a file called `app.go`. This file will contain a central application data structure, and any methods for our application's functionality. This architecture allows us to keep the top-level scope of our application clean. If we have objects that need to be accessed from multiple places, such as a database connection, we'll add this object to the app data structure and utilise it from methods attached to the same struct.

**app.go:**
```golang
package main

import (
	"flag"

	"github.com/prinsmike/govee"
)

// app represents the application environment.
type app struct {
	name        string
	showVersion bool
	version     govee.Version
}

// Parse command line options.
func (a *app) opts() {
	flag.BoolVar(&a.showVersion, "version", false, "Display the version info and exit.")
	flag.BoolVar(&a.showVersion, "V", false, "Display the version info and exit (shorthand).")
	flag.Parse()
}
```

## Create the version source file

Create a file called `version.go` and add the following code to it:

**version.go:**
```golang
package main

import (
        "fmt"
        "log"
        "os"
        "runtime"

        "github.com/prinsmike/govee"
)

// These variable should be set at compile time using -ldflags.
var VString, GitHash, GitBranch, GitUser, OS, Arch, Release, TStamp string

// setVersion will set our application's version information at runtime.
func (a *app) setVersion() {
        var err error
        vconf := govee.VersionConfig{
                VersionString: VString,
                GitHash:       GitHash,
                GitBranch:     GitBranch,
                GitUser:       GitUser,
                OS:            OS,
                Arch:          Arch,
                Compiler:      runtime.Version(),
                Release:       Release,
                TStamp:        TStamp,
        }

        a.version, err = govee.NewVersion(&vconf)
        if err != nil {
                log.Printf("Could not set version: [Error: %s] [VersionConfig: %#v]", err.Error(), vconf)
        }
}

// showVersion will print the application's version information (if a.showVersion is true) to standard
// output, and then terminate the application with an exit code of 3.
func (a *app) printVersion() {
        if a.showVersion {
                fmt.Printf("%s version: %s\n", a.name, a.version)
                fmt.Printf("Git hash: %s\n", a.version.GitHash())
                fmt.Printf("Git branch: %s\n", a.version.GitBranch())
                fmt.Printf("Git user: %s\n", a.version.GitUser())
                fmt.Printf("OS: %s\n", a.version.OS())
                fmt.Printf("Arch: %s\n", a.version.Arch())
                fmt.Printf("Compiler: %s\n", a.version.Compiler())
                fmt.Printf("Release: %s\n", a.version.Release())
                fmt.Printf("Timestamp: %s\n", a.version.TStamp())
				if s.version.Err() != nil {
					fmt.Printf("Error: %s\n", s.version.Err().Error())
				}
				if len(s.version.Warnings()) > 0 {
					for _, warning := range s.version.Warnings() {
						fmt.Printf("\t- %s\n", warning)
					}
				}
                os.Exit(3)
        }
}
```

## Create a build script

Next we'll create shell script to facilitate the compilation of our application.

**build.sh:**
```bash
#!/usr/bin/env bash

V="0.1.0" # The version number of this build script.
APP_NAME="myapp"
SCRIPT=$0
RELEASE=""
OS=""
ARCH=""
EXT=""

main() {
        # Parse the command line options.
        while getopts ":hVr:o:a" opt; do
                case ${opt} in
                        h )
                                usage
                                exit 2
                                ;;
                        V )
                                version
                                exit 3
                                ;;
                        r )
                                REL=$OPTARG
                                ;;
                        o )
                                OS=$OPTARG
                                ;;
                        a )
                                ARCH=$OPTARG
                                ;;
                        \? )
                                echo "Invalid option: -$OPTARG" 1>&2
                                echo
                                usage
                                exit 1
                                ;;
                        : )
                                echo "Invalid option: -$OPTARG requires an argument" 1>&2
                                echo
                                usage
                                exit 1
                                ;;
                esac
        done
        shift $((OPTIND -1))

        run
}

# Provide usage information for our script.
usage() {
        echo "Usage: $SCRIPT"
        echo "      Build my App."
        echo "  -r"
        echo "      The release tag for this binary (test | testing | prod | production)."
        echo "      (default: test)"
        echo "  -o"
        echo "      The target operating system for this binary."
        echo "      See https://github.com/golang/go/blob/master/src/go/build/syslist.go"
        echo "      for a full list of available operating systems (const goosList)."
        echo "      (default: linux)"
        echo "  -a"
        echo "      The target architecture for this binary."
        echo "      See https://github.com/golang/go/blob/master/src/go/build/syslist.go"
        echo "      for a full list of available architectures (const goarchList)."
        echo "      (default: amd64)"
}

# Display the version information for this build script.
version() {
        echo "$SCRIPT version $V"
}

# Initilize our variables and run the build.
run() {
        if [[ -z "$REL" ]]; then
                REL="test"
        fi
        if [[ -z "$OS" ]]; then
                OS="linux"
        fi
        if [[ -z "$ARCH" ]]; then
                ARCH="amd64"
        fi
        if [ $OS = "windows" ]; then
                $EXT=".exe"
        fi

        build
}

# Build the application.
build() {
        githash=$(git rev-parse HEAD)
        gitdescribe=$(git describe --tags --match semver/* | cut -f2 -d"/")
        gitbranch=$(git rev-parse --abbrev-ref HEAD)
        gituser=$(git config user.name)
        tstamp=$(date)

        ld_hash='"main.GitHash='${githash}'"'
        ld_desc='"main.VString='${gitdescribe}'"'
        ld_branch='"main.GitBranch='${gitbranch}'"'
        ld_user='"main.GitUser='${gituser}'"'
        ld_os='"main.OS='$OS'"'
        ld_arch='"main.Arch='${ARCH}'"'
        ld_rel='"main.Release='${REL}'"'
        ld_tstamp='"main.TStamp='${tstamp}'"'

        echo
        echo "-ldflags:"
        echo $ld_hash
        echo $ld_desc
        echo $ld_branch
        echo $ld_user
        echo $ld_os
        echo $ld_arch
        echo $ld_rel
        echo $ld_tstamp

        CGO_ENABLED="0" GOOS="${OS}" GOARCH="${ARCH}" \
                go build \
                -ldflags "-X ${ld_desc} -X ${ld_hash} -X ${ld_branch} -X ${ld_user}
                -X ${ld_os} -X ${ld_arch} -X ${ld_rel} -X ${ld_tstamp}" \
                -tags "$REL" -o "${APP_NAME}-${OS}_${ARCH}${EXT}"
}

# Call the main() function and pass all arguments to it.
main $@
```

## Create a version tag for the application

At this point the code above should be part of a git repository, and all files should be committed. We'll pretend that the project has had previous releases, and tag this one as semver/1.2.3. Enter the following command to tag the repository with our new release:

```
git tag -a semver/1.2.3 -m "v1.2.3"
```

## Build it!

Now that our repository has been tagged with the appropriate semantic version number, let's build our application and ensure the version information is successfully compiled with our binary.

Execute the build script:

```
./build.sh
```

This will build a test version of the application. To build a production version, ensure that no commits have been made to the branch since the git tag was created, and execute the build script with the `-r prod` flag.

```
./build.sh -r prod
```

You can see usage information for the build script by asking it for help: `./build.sh --help`

Once the script executes successfully, you should see something like this:

```
$ ./build.sh

-ldflags:
"main.GitHash=de6e4f2f6afbe97e9e96c12efef66fb4b65a0a3d"
"main.VString=1.2.3"
"main.GitBranch=master"
"main.GitUser=Jane Doe"
"mainl.OS=linux"
"main.Arch=amd64"
"main.Release=test"
"main.TStamp=Thu Feb 21 12:32:52 SAST 2019"
```

You'll end up with a binary executable: `myapp-linux_amd64`. The filename might be different depending on which operating system and architecture you've built for.

Once the application has been successfully built, execute it and ask it for its version info:

```
./myapp-linux_amd64 --version
```

You should see something like this:

```
MyApp version: 1.2.3
Git hash: de6e4f2f6afbe97e9e96c12efef66fb4b65a0a3d
Git branch: master
Git user: Jane Doe
OS: linux
Arch: amd64
Compiler: go1.11.5
Release: test
Timestamp: 2019-02-21T13:01:43+02:00
Warnings:
	- This version is tagged as release "test". Please don't use in production.
```

## Troubleshooting

### -bash: ./build.sh: Permission denied

Make sure the build script is executable: `chmod +x ./build.sh`

### fatal: not a git repository (or any of the parent directories): .git

The current directory is not a git repository. Use `git init` to initialize your repository, and then commit your changes and tag the repository with the version number.

### fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree.

Add your files to the repository and create a commit.

```
git add . # Pay attention to the dot.
git commit -a -m "First commit."
```

### fatal: No names found, cannot describe anything.

This error occurs when there are no semver tags. Create a git tag for your version: `git tag -a v1.2.3 -m "1.2.3"`.
