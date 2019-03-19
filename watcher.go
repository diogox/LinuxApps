package LinuxApps

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
	"strings"
)

func NewAppWatcher(onChange func(*AppInfo) error, onRemove func() error) *AppWatcher {
	return &AppWatcher{
		OnChange: onChange,
		OnRemove: onRemove,
	}
}

type AppWatcher struct {
	OnChange    func(*AppInfo) error
	OnRemove    func() error
}

func (aw *AppWatcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	// Watch app path
	err = watcher.Add(DesktopFilesPath)
	if err != nil {
		panic(err)
	}

	// Watch override file, only if it exists
	if _, err := os.Stat(DesktopFilesOverridePath); !os.IsNotExist(err) {
		err = watcher.Add(DesktopFilesOverridePath)
		if err != nil {
			panic(err)
		}
	}

	done := make(chan bool)

	go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if !strings.HasSuffix(event.Name, ".desktop") {
						continue
					}

					filePathComponents := strings.Split(event.Name, "/")
					fileName := filePathComponents[len(filePathComponents) - 1]

					if strings.Contains(event.Name, DesktopFilesOverridePath) {
						// Check that the file is not overriding any other
						if _, err := os.Stat(DesktopFilesPath + fileName); os.IsNotExist(err) {
							aw.OnRemove()
							continue
						}
					}

					// Check for original (non-overridden) file
					app, err := decodeDesktopFile(DesktopFilesPath + fileName)
					if err != nil {
						panic(err)
					}

					err = aw.OnChange(app)
					if err != nil {
						panic(err)
					}
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					if !strings.HasSuffix(event.Name, ".desktop") {
						continue
					}

					app, err := decodeDesktopFile(event.Name)
					if err != nil {
						continue
					}

					// Ignore if being overridden
					if !strings.HasPrefix(event.Name, DesktopFilesOverridePath) {
						// Check if it's been overridden
						filePathComponents := strings.Split(event.Name, string(os.PathSeparator))
						fileName := filePathComponents[len(filePathComponents) - 1]

						if _, err = os.Stat(path.Join(DesktopFilesOverridePath, fileName)); !os.IsNotExist(err) {
							// File exists in override folder
							continue
						}
					}

					err = aw.OnChange(app)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}()

	<-done

	return nil
}
