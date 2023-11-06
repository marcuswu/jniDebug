package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

func adbCommand(command string) (string, error) {
	args := strings.Split(command, " ")
	cmd := exec.Command("adb", args...)
	res, err := cmd.Output()
	return string(res), err
}

func commandWithDevice(device string, command string) string {
	dev := ""
	if len(device) > 0 {
		dev = fmt.Sprintf("-s %s ", device)
	}
	return fmt.Sprintf("%s%s", dev, command)
}

func Push(device string, source string, dest string) error {
	_, err := adbCommand(commandWithDevice(device, fmt.Sprintf("push %s %s", source, dest)))
	return err
}

func Forward(device string, localPort string, destPort string) error {
	_, err := adbCommand(commandWithDevice(device, fmt.Sprintf("forward tcp:%s tcp:%s", localPort, destPort)))
	return err
}

func ShellCommand(device string, runAs string, command string) (string, error) {
	assume := ""
	if len(runAs) > 0 {
		assume = fmt.Sprintf("run-as %s ", runAs)
	}
	cmd := fmt.Sprintf("%s%s", assume, command)
	return adbCommand(commandWithDevice(device, cmd))
}
