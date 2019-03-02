package LinuxApps

import "sync"

type AppInfo struct {
	Name        string
	Description string
	IconName    string
	ExecName    string
}

func GetApps() []*AppInfo {
	desktopFiles := getDesktopFiles(DesktopFilesPath)
	desktopOverrideFiles := getDesktopFiles(DesktopFilesPath)

	apps := make([]*AppInfo, 0)

	var wg sync.WaitGroup
	for _, file := range desktopFiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			appInfo, err := decodeDesktopFile(file)
			if err != nil {
				panic(err)
			}
			apps = append(apps, appInfo)
		}()
	}
	wg.Wait()

	for _, file := range desktopOverrideFiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			appInfo, err := decodeDesktopFile(file)
			if err != nil {
				panic(err)
			}
			apps = append(apps, appInfo)

			for _, app := range apps {
				if app.ExecName == appInfo.ExecName {
					return
				}
				apps = append(apps, app)
			}
		}()
	}

	return apps
}