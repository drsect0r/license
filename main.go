package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"os"
	"path"
	"sync"
)

// pathExists returns true if the path exists
func pathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	args := os.Args[1:]
	wg := sync.WaitGroup{}
	var mainErr error

	// Check existence of license data directory
	// and start making it if it is not present.
	// If we cannot find the home directory, silently ignore
	// for now and handle in the specific command function called.
	if home, err := homedir.Dir(); err == nil && len(args) >= 1 && !(args[0] == "update" || args[0] == "bootstrap") && !pathExists(path.Join(home, base.LicenseDirectory, base.DataDirectory)) {
		wg.Add(1)
		go func() {
			base.Bootstrap([]string{"--quiet"})
			wg.Done()
		}()
	}

	if len(args) < 1 {
		mainErr = base.Help()
	} else {
		command := args[0]

		switch command {
		// Help information
		case "--help":
			fallthrough
		case "help":
			mainErr = base.Help()

		// Version information
		case "--version":
			fallthrough
		case "version":
			mainErr = base.Version()

		// Update to latest remote licenses
		case "update":
			fallthrough
		case "bootstrap":
			mainErr = base.Bootstrap(args[1:])

		// List remote licenses
		case "ls-remote":
			fallthrough
		case "list-remote":
			mainErr = base.ListRemote()

		// List local licenses
		case "ls":
			fallthrough
		case "list":
			wg.Wait()
			mainErr = base.ListLocal()

		default:
			wg.Wait()
			mainErr = base.Generate(args)
		}
	}

	wg.Wait()

	exit(mainErr)
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
