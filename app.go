package LinuxApps

type AppInfo struct {
	Name        string
	Description string
	IconName    string
	ExecName    string
}

func GetApps() []*AppInfo {
	desktopFiles := getDesktopFiles(DesktopFilesPath)
	desktopOverrideFiles := getDesktopFiles(DesktopFilesPath)

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