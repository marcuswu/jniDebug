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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	Usage: jniDebug -device emulator-0000 -lldb /path/to/lldb -package your.package -activity main.activity -debug /path/to/binary -vscode /path/to/launch.json [OPTIONS]

	jniDebug requires adb to operate. Ensure it is installed and in your path.

	Options:
	-config  The name to use for the VSCode run configuration
	-device  The adb device to connect to (check adb devices -l)
	-port    The port number to listen on
	-wait    If present, pause app execution and wait for the debugger
	-verbose If present, output verbose logging
	`, err)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// adb must be in $PATH
	// Handle args, etc
	devicePtr := flag.String("device", "", "The adb device id to connect to (check adb devices -l)")
	lldbPtr := flag.String("lldb", "", "Path to target architecture lldb")
	packagePtr := flag.String("package", "", "The name of the package for the app")
	activityPtr := flag.String("activity", "", "The name of the activity to run")
	debugPtr := flag.String("debug", "", "The location of the native library to debug")
	launchPtr := flag.String("vscode", "", "The location of the vscode launch.json")

	configPtr := flag.String("config", "Go Mobile Debugging", "The name to use for the VSCode run configuration")
	portPtr := flag.String("port", "23456", "The port number to listen on")
	waitPtr := flag.Bool("wait", false, "Pause execution and wait for the debugger")
	verbosePtr := flag.Bool("verbose", false, "Use verbose logging")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verbosePtr {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if !checkLLDB(*lldbPtr) {
		printUsage("No LLDB was found")
		os.Exit(1)
	}

	if len(*devicePtr) < 1 {
		printUsage("No device was provided")
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

	// Stop LLDB so we don't get errors when we try to copy
	adb.StopLLDB(*devicePtr, *packagePtr)

	// push lldb-server
	err := adb.CopyLLDB(*devicePtr, *packagePtr, *lldbPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to push LLDB to the device")
		os.Exit(1)
	}
	// setprop wait_for_debugger
	adb.SetWaitForDebugger(*devicePtr, *waitPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set device to wait for debugger")
		os.Exit(2)
	}
	// forward debugging port
	adb.Forward(*devicePtr, *portPtr, *portPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to forward the device port")
		os.Exit(3)
	}
	// start the app
	adb.StartApp(*devicePtr, *packagePtr, *activityPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start the app")
		os.Exit(4)
	}
	// find the app's pid
	pid, err := adb.GetAppPid(*devicePtr, *packagePtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get the app's PID")
		os.Exit(5)
	}
	// write vscode config
	config := vscode.GenerateVscodeConfig(*configPtr, *devicePtr, *portPtr, pid, *debugPtr)
	launchConfig, err := os.ReadFile(*launchPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read VSCode launch.json")
		os.Exit(6)
	}
	newConfig, err := vscode.AlterVscodeConfig(string(launchConfig), config, "// #lldbclient-generated-begin", "// #lldbclient-generated-end")
	if err != nil {
		log.Error().Err(err).Msg("Failed to alter VSCode launch.json")
		os.Exit(7)
	}

	err = os.WriteFile(*launchPtr, []byte(newConfig), 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write VSCode launch.json")
		os.Exit(8)
	}
	// run lldb on target device
	log.Info().Msgf(`Success! Now run the "%s" run configuration in VSCode to begin debugging!`, *configPtr)
	adb.StartLLDB(*devicePtr, *packagePtr, *portPtr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start LLDB on the device")
		os.Exit(9)
	}
}
