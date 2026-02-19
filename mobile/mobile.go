package arduinobuddycli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/internal/cli"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"go.bug.st/serial"
)

// BindablePortCallback is the Kotlin/Java-facing interface
type BindablePortCallback interface {
	Call(action string, dataJSON string) (resultJSON string, err error)
}

// BindablePortCallback is the Kotlin/Java-facing interface
type BindableExecutableCallback interface {
	Run(command string) (result string)
}

// serialAdapter implements serial.PortCallbackInterface
type serialAdapter struct {
	cb BindablePortCallback
}

type executableAdapter struct {
	cb BindableExecutableCallback
}

func (s *executableAdapter) Run(command string) (result string) {
	return // TODO
}

func (s *serialAdapter) Call(action string, data map[string]interface{}) (map[string]interface{}, error) {
	jsonInput, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	jsonOutput, err := s.cb.Call(action, string(jsonInput))
	if err != nil {
		return nil, err
	}

	var out map[string]interface{}
	err = json.Unmarshal([]byte(jsonOutput), &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SetProvider binds the Kotlin/Java callback to the internal provider
func SetProvider(cb BindablePortCallback) {
	serial.SetProvider(&serialAdapter{cb})
}

var (
	coreServer rpc.ArduinoCoreServiceServer
	ctx        context.Context
	instanceID *rpc.Instance
)

type callbackWriter struct {
	fn func(string)
}

func (w *callbackWriter) Write(p []byte) (int, error) {
	w.fn(string(p))
	return len(p), nil
}

// Init sets paths for Arduino CLI and initializes the core server.
func Init(baseDir string) error {
	arduinoData := filepath.Join(baseDir, "arduino")
	configFile := filepath.Join(arduinoData, "arduino-cli.yaml")

	if err := os.MkdirAll(arduinoData, 0o755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	ctx = context.Background()
	coreServer = commands.NewArduinoCoreServer()

	// ---- Create minimal config if missing ----
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := fmt.Sprintf(`
board_manager:
  additional_urls: []
directories:
  data: "%s"
`, arduinoData)

		if err := os.WriteFile(configFile, []byte(defaultConfig), 0o644); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	// ---- Ensure index files exist ----
	packageIndex := filepath.Join(arduinoData, "package_index.json")
	if _, err := os.Stat(packageIndex); os.IsNotExist(err) {
		os.WriteFile(packageIndex, []byte(`{"packages":[]}`), 0o644)
	}

	libraryIndex := filepath.Join(arduinoData, "library_index.json")
	if _, err := os.Stat(libraryIndex); os.IsNotExist(err) {
		os.WriteFile(libraryIndex, []byte(`{"libraries":[]}`), 0o644)
	}

	// ---- Temp dir ----
	tempDir := filepath.Join(arduinoData, "tmp")
	os.MkdirAll(tempDir, 0o755)

	os.Setenv("TMPDIR", tempDir)
	os.Setenv("ARDUINO_DATA_DIR", arduinoData)
	os.Setenv("ARDUINO_CONFIG_FILE", configFile)

	// ---- Open CLI config ----
	openReq := &rpc.ConfigurationOpenRequest{
		SettingsFormat: "yaml",
	}

	if data, err := os.ReadFile(configFile); err == nil {
		openReq.EncodedSettings = string(data)
	}

	_, err := coreServer.ConfigurationOpen(ctx, openReq)
	if err != nil {
		return fmt.Errorf("failed to open configuration: %w", err)
	}

	createResp, err := coreServer.Create(ctx, &rpc.CreateRequest{})
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}
	instanceID = createResp.Instance

	return nil
}

// LogOutput is a struct that can be returned by gomobile bind
// to pass logs back to the native environment.
type LogOutput struct {
	logs []string
}

// GetLogs returns the logs as a single, newline-separated string
// which is a gomobile-compatible return type.
func (l *LogOutput) GetLogs() string {
	return strings.Join(l.logs, "\n")
}

func BoardList() (string, error) {
	if instanceID == nil {
		return "", fmt.Errorf("instance not initialized")
	}

	req := &rpc.BoardListRequest{
		Instance: instanceID,
		Timeout:  2000,
	}

	resp, err := coreServer.BoardList(ctx, req)

	if err != nil {
		return "", err
	}

	if len(resp.Ports) == 0 {
		return "No boards found", nil
	}

	// Convert RPC to JSON or readable text
	b, _ := json.MarshalIndent(resp, "", "  ")
	return string(b), nil
}

// RunSimple executes the command with the given arguments.
// It returns a *LogOutput struct (compatible with gomobile) and an error.
func RunSimple(argsCSV string) (*LogOutput, error) { // TODO - FIX LOG CAPTURING
	if coreServer == nil {
		return nil, fmt.Errorf("core server not initialized; call Init() first")
	}

	// Assuming a simple CSV split is sufficient for this context
	args := strings.Split(argsCSV, ",")

	// Create an instance of the struct to collect logs
	output := &LogOutput{
		logs: []string{},
	}

	writer := &callbackWriter{
		fn: func(line string) {
			output.logs = append(output.logs, line)
		},
	}

	arduinoCmd := cli.NewCommand(coreServer)
	arduinoCmd.SetOut(writer)
	arduinoCmd.SetErr(writer)
	arduinoCmd.SetArgs(args)

	err := arduinoCmd.ExecuteContext(ctx)

	// Return the populated struct
	return output, err
}

// RunNative executes a native function by name with a request object.
// req must be the correct request type for the function.
func RunNative(fnName string) (string, error) {
	if coreServer == nil {
		return "", fmt.Errorf("core server not initialized; call Init() first")
	}

	var resp any
	var err error

	switch fnName {
	case "BoardList":
		resp, err = BoardList()
	default:
		return "", fmt.Errorf("unknown native function: %s", fnName)
	}

	if err != nil {
		return "", err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(data), nil
}
