package adb

import (
	"fmt"
)

func CopyLLDB(device string, pkg string, lldbPath string) error {
	err := Push(device, lldbPath, "/data/local/tmp")
	if err != nil {
		return err
	}
	_, err = ShellCommand(device, pkg, fmt.Sprintf("cp /data/local/tmp/lldb-server /data/data/%s/", pkg))
	return err
}

func SetWaitForDebugger(device string, wait bool) error {
	_, err := ShellCommand(device, "", fmt.Sprintf("setprop debug.debuggerd.wait_for_debugger %t", wait))
	return err
}

func StartApp(device string, pkg string, activity string) error {
	_, err := ShellCommand(device, "", fmt.Sprintf("am start -n %s/%s", pkg, activity))
	return err
}

func GetAppPid(device string, pkg string) (string, error) {
	return ShellCommand(device, "", fmt.Sprintf("ps -A | grep %s | awk '{ print $2 }'", pkg))
}

func StartLLDB(device string, pkg string, port string) error {
	_, err := ShellCommand(device, pkg, fmt.Sprintf("/data/data/%s/lldb-server --server --listen \"*:%s\"", pkg, port))
	return err
}
