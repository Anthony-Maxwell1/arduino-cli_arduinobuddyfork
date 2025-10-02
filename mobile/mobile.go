package arduinocli

import (
    "github.com/Anthony-Maxwell1/android-cli_arduinobuddyfork/arduinocli"
    "strings"
)

// RunCommand runs an Arduino CLI command and returns all output as a single string.
// argsCSV = "board,list,all" or "compile,/path/to/sketch"
func RunCommand(argsCSV string) (string, error) {
    args := strings.Split(argsCSV, ",")
    var outputBuilder strings.Builder

    // original callback
    callback := func(line string) {
        outputBuilder.WriteString(line)
        outputBuilder.WriteString("\n")
    }

    err := arduinocli.RunArduinoCommand(args, callback)
    return outputBuilder.String(), err
}