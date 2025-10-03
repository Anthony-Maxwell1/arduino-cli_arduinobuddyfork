package arduinobuddycli

import (
    "github.com/Anthony-Maxwell1/arduino-cli_arduinobuddyfork/arduinocli"
)

// RunSimple runs an Arduino CLI command (comma-separated args) and returns the output.
func RunSimple(argsCSV string) (string, error) {
    return RunCommand(argsCSV)
}