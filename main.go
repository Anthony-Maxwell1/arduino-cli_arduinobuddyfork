package arduinocli

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/internal/cli"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
)

// CommandResult is a simple JSON-serializable struct
type CommandResult struct {
	Level   string `json:"level"` // "info", "warn", "error"
	Message string `json:"message"`
}

// RunArduinoCommand runs an Arduino CLI command and sends results to a callback.
// callback is called once per message as JSON string.
func RunArduinoCommand(args []string, callback func(string)) error {
	ctx := context.Background()
	srv := commands.NewArduinoCoreServer()

	// capture warnings
	configFile := "" // you could pass this in if needed
	openReq := &rpc.ConfigurationOpenRequest{SettingsFormat: "yaml"}
	if configData, err := os.ReadFile(configFile); err == nil {
		openReq.EncodedSettings = string(configData)
	}
	if _, err := srv.ConfigurationOpen(ctx, openReq); err != nil {
		return err
	}

	arduinoCmd := cli.NewCommand(srv)

	// override stdout/stderr with our callback
	writer := &callbackWriter{callback: callback}

	arduinoCmd.SetOut(writer)
	arduinoCmd.SetErr(writer)

	// execute command
	arduinoCmd.SetArgs(args)
	if err := arduinoCmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}

// callbackWriter implements io.Writer
type callbackWriter struct {
	callback func(string)
}

func (w *callbackWriter) Write(p []byte) (n int, err error) {
	// wrap each message in JSON
	msg := CommandResult{
		Level:   "info",
		Message: string(p),
	}
	b, _ := json.Marshal(msg)
	w.callback(string(b))
	return len(p), nil
}

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

    err := RunArduinoCommand(args, callback)
    return outputBuilder.String(), err
}