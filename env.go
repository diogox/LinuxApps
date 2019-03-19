package LinuxApps

import (
	"os"
	"path"
	"strings"
)

const XDG_DATA_DIRS = "XDG_DATA_DIRS"

func getEnvPaths() []string {
	paths := make([]string, 0)
	allPaths := os.Getenv(XDG_DATA_DIRS)

	for _, p := range strings.Split(allPaths, ":") {
		appPath := path.Join(p, "applications")
		paths = append(paths, appPath)
	}
	return paths
}
