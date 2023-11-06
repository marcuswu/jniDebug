package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/marcuswu/jnidebug/adb"
	"github.com/marcuswu/jnidebug/vscode"
)

func checkLLDB(path string) bool {
	if !strings.HasSuffix(path, "lldb-server") {
		return false
	}
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

func printUsage(err string) {
	fmt.Printf(`
	%s
	Set up VSCode debugging for a native (JNI) Android project.
	Usage: jniDebug -lldb /path/to/lldb -package your.package -activity main.activity -debug /path/to/binary -vscode /path/to/launch.json [OPTIONS]

	jniDebug requires adb to operate. Ensure it is installed and in your path.

	Options:
	-config The name to use for the VSCode run configuration
	-device The adb device to connect to (check adb devices -l)
	-port   The port number to listen on
	-wait   If present, pause app execution and wait for the debugger
	`, err)
}

func main() {
	// adb must be in $PATH
	// Handle args, etc
	lldbPtr := flag.String("lldb", "", "Path to target architecture lldb")
	packagePtr := flag.String("package", "", "The name of the package for the app")
	activityPtr := flag.String("activity", "", "The name of the activity to run")
	debugPtr := flag.String("debug", "", "The location of the native library to debug")
	launchPtr := flag.String("vscode", "", "The location of the vscode launch.json")

	configPtr := flag.String("config", "Go Mobile Debugging", "The name to use for the VSCode run configuration")
	devicePtr := flag.String("device", "", "The adb device id to connect to (check adb devices -l)")
	portPtr := flag.String("port", "23456", "The port number to listen on")
	waitPtr := flag.Bool("wait", false, "Pause execution and wait for the debugger")

	flag.Parse()

	if !checkLLDB(*lldbPtr) {
		printUsage("No LLDB was found")
		os.Exit(1)
	}

	if len(*packagePtr) < 1 {
		printUsage("No package was provided")
		os.Exit(1)
	}

	if len(*activityPtr) < 1 {
		printUsage("No activity was provided")
		os.Exit(1)
	}

	if len(*debugPtr) < 1 {
		printUsage("No debug target was provided")
		os.Exit(1)
	}

	if len(*launchPtr) < 1 {
		printUsage("The location of launch.json was not provided")
		os.Exit(1)
	}

	// push lldb-server
	err := adb.CopyLLDB(*devicePtr, *packagePtr, *lldbPtr)
	if err != nil {
		fmt.Printf("Failed to push LLDB to the device: %s\n", err.Error())
		os.Exit(1)
	}
	// setprop wait_for_debugger
	adb.SetWaitForDebugger(*devicePtr, *waitPtr)
	if err != nil {
		fmt.Printf("Failed to set device to wait for debugger: %s\n", err.Error())
		os.Exit(2)
	}
	// forward debugging port
	adb.Forward(*devicePtr, *portPtr, *portPtr)
	if err != nil {
		fmt.Printf("Failed to forward the device port: %s\n", err.Error())
		os.Exit(3)
	}
	// start the app
	adb.StartApp(*devicePtr, *packagePtr, *activityPtr)
	if err != nil {
		fmt.Printf("Failed to start the app: %s\n", err.Error())
		os.Exit(4)
	}
	// find the app's pid
	pid, err := adb.GetAppPid(*devicePtr, *packagePtr)
	if err != nil {
		fmt.Printf("Failed to get the app's PID: %s\n", err.Error())
		os.Exit(5)
	}
	// run lldb on target device
	adb.StartLLDB(*devicePtr, *packagePtr, *portPtr)
	if err != nil {
		fmt.Printf("Failed to start LLDB on the device: %s\n", err.Error())
		os.Exit(6)
	}
	// write vscode config
	config := vscode.GenerateVscodeConfig(*configPtr, *devicePtr, *portPtr, pid, *debugPtr)
	vscode.AlterVscodeConfig(*launchPtr, config, "", "")
	if err != nil {
		fmt.Printf("Failed to wrote VSCode launch.json: %s\n", err.Error())
		os.Exit(7)
	}

	fmt.Printf(`Success! Now run the "%s" run configuration in VSCode to begin debugging!`, *configPtr)
}
