package LinuxApps

type AppInfo struct {
	Name        string
	Description string
	IconName    string
	ExecName    string
}

func GetApps() []*AppInfo {
	desktopFilesPaths := getEnvPaths()

	// Get all desktop, visible, files
	desktopFiles := getDesktopFiles(DesktopFilesPath)
	for _, path := range desktopFilesPaths {
		desktopFiles = append(desktopFiles, getDesktopFiles(path)...)
	}

	desktopOverrideFiles := getDesktopFiles(DesktopFilesOverridePath)

	appsMap := make(map[string]*AppInfo, 0)
	for _, file := range desktopFiles {
		appInfo, err := decodeDesktopFile(file)
		if err != nil {
			continue
		}
		appsMap[appInfo.ExecName] = appInfo
	}

	for _, file := range desktopOverrideFiles {
		appInfo, err := decodeDesktopFile(file)
		if err != nil {
			continue
		}
		appsMap[appInfo.ExecName] = appInfo
	}

	apps := make([]*AppInfo, 0)
	for _, app := range appsMap {
		apps = append(apps, app)
	}

	return apps
}