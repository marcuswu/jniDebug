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
/*
AlterVscodeConfig is heavily based on insert_commands_into_vscode_config from
https://android.googlesource.com/platform/development/+/master/scripts/gdbclient.py#365
which is licensed under the Apache 2.0 license:

Copyright (C) 2015 The Android Open Source Project

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
