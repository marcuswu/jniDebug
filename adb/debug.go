package adb

import (
	"fmt"
	"strings"
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
	pid, err := ShellCommand(device, "", fmt.Sprintf("ps -A | grep %s | awk '{ print $2 }'", pkg))
	return strings.TrimSpace(pid), err
}

func StopLLDB(device string, pkg string) error {
	// kill existing lldb prior to starting
	pid, err := ShellCommand(device, "", "ps -A | grep lldb-server | head -1 | awk '{ print $2 }'")
	pid = strings.TrimSpace(pid)
	if err == nil && len(pid) >= 1 {
		_, err = ShellCommand(device, pkg, fmt.Sprintf("kill -9 %s", pid))
	}
	return err
}

func StartLLDB(device string, pkg string, port string) error {
	_, err := ShellCommand(device, pkg, fmt.Sprintf("/data/data/%s/lldb-server platform --server --listen \"*:%s\"", pkg, port))
	return err
}
