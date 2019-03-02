package LinuxApps

import (
	"io/ioutil"
	"strings"
)

const DesktopFilesPath = "/usr/share/applications/"
const DesktopFilesOverridePath = "~/.local/share/applications/"

func getDesktopFiles(path string) []string {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0)
	for _, info := range fileInfo {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".desktop") {
			files = append(files, path + info.Name())
		}
	}

	return files
}

