package arduinobuddycli

import (
    "github.com/arduino/arduino-cli"
)

// RunSimple runs an Arduino CLI command (comma-separated args) and returns the output.
func RunSimple(argsCSV string) (string, error) {
    return arduinocli.RunCommand(argsCSV)
}