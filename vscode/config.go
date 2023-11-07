package vscode

import (
	"fmt"
	"strings"
)

func GenerateVscodeConfig(configName string, device string, port string, pid string, debugFile string) []string {
	return []string{
		`{`,
		`    "name": "` + configName + `",`,
		`    "type": "lldb",`,
		`    "request": "custom",`,
		`    "initCommands": ["platform select remote-android", "file ` + debugFile + `"],`,
		`    "processCreateCommands": ["platform connect connect://` + device + `:` + port + `", "attach ` + pid + `"]`,
		`}`,
	}
}

// This code is a port of https://android.googlesource.com/platform/development/+/master/scripts/gdbclient.py#365
func AlterVscodeConfig(launchConfig string, newLines []string, startMarker string, endMarker string) (string, error) {
	foundBegin := false
	beginLine := -1

	lines := strings.Split(launchConfig, "\n")
	output := make([]string, 0, len(lines))

	for lineNum, line := range lines {
		if beginLine >= 0 {
			if strings.TrimSpace(line) == endMarker {
				beginLine = -1
			} else {
				continue
			}
		}
		output = append(output, line)
		if strings.TrimSpace(line) == startMarker {
			foundBegin = true
			beginLine = lineNum
			markerIndent := line[:strings.Index(line, startMarker)]
			// Can't append these all at once due to needing to keep indentation
			for _, newLine := range newLines {
				output = append(output, markerIndent+newLine)
			}
		}
	}

	if !foundBegin {
		return launchConfig, fmt.Errorf("did not find begin marker line %s in the VSCode launch file", startMarker)
	}

	if beginLine != -1 {
		return launchConfig, fmt.Errorf("unterminated begin marker at line %d in the VSCode launch file. Add end marker line to file: '%s'", beginLine+1, endMarker)
	}

	return strings.Join(output, "\n"), nil
}
