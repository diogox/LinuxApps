package LinuxApps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Starts an app, or focuses an existing instance of that same app.
func StartAppOrFocusExisting(app *AppInfo) error {
	return StartAppOrFocusExistingByCommand(app.ExecName)
}

// Takes a command that starts an app and ensures that a new instance of that app is not created if there's one already running.
// It also gives focus to that running instance.
func StartAppOrFocusExistingByCommand(command string) error {
	executable := strings.Split(command, " ")

	// Give focus to existing instance, if it exists
	execNames := getExecutables(executable...)

	isFocusGiven := false
	for _, execName := range execNames {
		err := resolveFocusForExecutable(execName)

		if err == nil {
			isFocusGiven = true
		}
	}

	// If it worked, don't create new instance
	if isFocusGiven { return nil }

	// Run new instance
	return StartAppByCommand(command)
}

// Start app. Doesn't give focus to existing instances.
func StartApp(app *AppInfo) error{

	return StartAppByCommand(app.ExecName)
}

// Starts app through it's executable command. Doesn't give focus to existing instances.
func StartAppByCommand(command string) error{
	executable := strings.Split(command, " ")

	cmd := exec.Command(executable[0], executable[1:]...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	_ = cmd.Process.Release()

	return nil
}

// Sometimes, a command can have multiple executables eg. `ls | grep query`.
// In this case, we'd get the executables `ls` and `grep.`
func getExecutables(command ...string) []string {
	executables := make([]string, 0)

	for _, part := range command {
		_, err := getExecPath(part)

		// If there's no error, then it has an executable path, therefore it's an executable
		if err == nil {
			executables = append(executables, part)
		}
	}

	return executables
}

// Checks if the executable already has a running instance
func hasRunningInstanceWithExactName(execName string) bool {
	findInstanceCmd := exec.Command("ps", "-C", execName)

	// Save output
	var output bytes.Buffer
	findInstanceCmd.Stdout = &output

	// Run "ps -C %execName"
	err := findInstanceCmd.Run()
	if err != nil {
		return false
	}

	return strings.Contains(output.String(), execName)
}

// Checks if the executable already has a running instance
func hasRunningInstanceWithNameThatIncludes(execName string) bool {
	runningProcessedCmd := exec.Command("bash", "-c", fmt.Sprintf("ps -A | grep %s", execName))

	// Run "ps -A | grep %execName." & get output
	res, err := runningProcessedCmd.Output()
	if err != nil {
		return false
	}

	output := string(res)
	return strings.Contains(output, execName)
}

// Gets the path for an executable that has it's path in the env variables.
func getExecPath(execName string) (string, error) {
	findExecPathCmd := exec.Command("which", execName)

	// Save output
	var output bytes.Buffer
	findExecPathCmd.Stdout = &output

	// Run "ps -C %execName"
	err := findExecPathCmd.Run()
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// Gives focus to the window of a proccess with the given name.
func giveFocusToProccessWithName(execName string) error {
	cmd := exec.Command("wmctrl", "-a", execName)
	return cmd.Run()
}

// It tries to give focus to the app, by trying all possible executable name variations until it succeeds.
func resolveFocusForExecutable(exec string) error {
	// If it's a path, get only the name
	// Need to do this because, sometimes, an executable's proccess uses only the name (because it's in the env variables)
	// but the '.desktop' file's command runs the whole path.
	execName := getExecNameIfPath(exec)

	if hasRunningInstanceWithExactName(execName) {

		return giveFocusToProccessWithName(execName)
	} else if execPath, err := getExecPath(execName); err == nil && hasRunningInstanceWithExactName(execPath) {

		return giveFocusToProccessWithName(execPath)
	} else if hasRunningInstanceWithNameThatIncludes(execName + ".") {

		// Runs if there's a match for 'execName.*', where * can be any type of extension (eg. 'sh', 'py', etc.)
		// The focus-giving command works without the extension.
		return giveFocusToProccessWithName(execName)
	}

	return nil
}

// Returns the name of an executable if it's given a path.
func getExecNameIfPath(exec string) string {

	// Check if path
	isExecPath := strings.Contains(exec, string(os.PathSeparator))
	if isExecPath {

		// Get name (last part)
		execParts := strings.Split(exec, string(os.PathSeparator))
		return execParts[len(execParts) - 1]
	}

	return exec
}