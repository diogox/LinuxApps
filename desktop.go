package LinuxApps

import (
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"path"
	"strings"
)

const DesktopFilesPath = "/usr/share/applications/"
var DesktopFilesOverridePath = "~/.local/share/applications/"

func init() {
	// Expand `DesktopFilesOverridePath`
	overridePath, err := homedir.Expand(DesktopFilesOverridePath)
	if err != nil {
		panic(err)
	}

	DesktopFilesOverridePath = overridePath
}

func getDesktopFiles(dirPath string) []string {
	fileInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	files := make([]string, 0)
	for _, info := range fileInfo {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".desktop") {
			files = append(files, path.Join(dirPath, info.Name()))
		}
	}

	return files
}

