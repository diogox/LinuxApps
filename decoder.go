package LinuxApps

import (
	"errors"
	"gopkg.in/ini.v1"
)

func decodeDesktopFile(filepath string) (*AppInfo, error) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, err
	}

	entry := cfg.Section("Desktop Entry")

	isNoDisplay := entry.Key("NoDisplay").Value()
	if isNoDisplay == "true" {
		return nil, errors.New("app not to be displayed")
	}

	name := entry.Key("Name").Value()
	description := entry.Key("Comment").Value()
	execName := entry.Key("Exec").Value()
	iconName := entry.Key("Icon").Value()

	return &AppInfo{
		Name:        name,
		Description: description,
		ExecName: execName,
		IconName: iconName,
	}, nil
}
