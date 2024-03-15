package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// DBPath returns the appropriate database location for the platform.
func DBPath(opts map[string]interface{}, platform string) (string, error) {
	// if `--db` was provided, use that path
	if opts["--db"] != nil {
		return opts["--db"].(string), nil
	}

	// otherwise, find the data path that's appropriate for the platform
	path := os.Getenv("XDG_DATA_HOME")

	switch platform {
	case "aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios",
		"linux", "netbsd", "openbsd", "plan9", "solaris":

		if path == "" {
			path = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}

	case "windows":
		path = os.Getenv("APPDATA")

	default:
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	return filepath.Join(path, "llm.db"), nil
}
