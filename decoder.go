package LinuxApps

import (
	"errors"
	"fmt"
	jj "github.com/cloudfoundry-attic/jibber_jabber"
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

	// Get info compatible with the System-Language
	userLanguage, err := jj.DetectLanguage()
	if err == nil {
		// Get name
		nameKey := entry.Key(fmt.Sprintf("Name[%s]", userLanguage))
		if nameKey != nil {
			name = nameKey.Value()
		}
		// Get description
		descriptionKey := entry.Key(fmt.Sprintf("Comment[%s]", userLanguage))
		if descriptionKey != nil {
			description = descriptionKey.Value()
		}
	}

	return &AppInfo{
		Name:        name,
		Description: description,
		ExecName: execName,
		IconName: iconName,
	}, nil
}
