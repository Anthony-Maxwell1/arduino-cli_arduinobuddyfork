package serial

import (
	"errors"
	"time"
)

// -------------------------
// Internal provider interface
// -------------------------

type PortCallbackInterface interface {
	Call(action string, data map[string]interface{}) (map[string]interface{}, error)
}

var ProviderInterface PortCallbackInterface

func SetProvider(p PortCallbackInterface) {
	ProviderInterface = p
}

// -------------------------
// Port interface
// -------------------------

type Port interface {
	SetMode(mode *Mode) error
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Drain() error
	ResetInputBuffer() error
	ResetOutputBuffer() error
	SetDTR(dtr bool) error
	SetRTS(rts bool) error
	GetModemStatusBits() (*ModemStatusBits, error)
	SetReadTimeout(t time.Duration) error
	Close() error
	Break(duration time.Duration) error
}

var NoTimeout time.Duration = -1

type ModemStatusBits struct {
	CTS bool
	DSR bool
	RI  bool
	DCD bool
}

type ModemOutputBits struct {
	RTS bool
	DTR bool
}

type Mode struct {
	BaudRate          int
	DataBits          int
	Parity            Parity
	StopBits          StopBits
	InitialStatusBits *ModemOutputBits
}

type Parity int

const (
	NoParity Parity = iota
	OddParity
	EvenParity
	MarkParity
	SpaceParity
)

type StopBits int

const (
	OneStopBit StopBits = iota
	OnePointFiveStopBits
	TwoStopBits
)

// -------------------------
// Open / GetPortsList
// -------------------------

func Open(portName string, mode *Mode) (Port, error) {
	if ProviderInterface == nil {
		return nil, &PortError{code: PortNotFound}
	}

	res, err := ProviderInterface.Call("open", map[string]interface{}{
		"portName": portName,
		"mode":     mode,
	})
	if err != nil {
		return nil, err
	}

	pRaw, ok := res["port"]
	if !ok {
		return nil, &PortError{code: InvalidSerialPort}
	}

	portID, ok := pRaw.(string)
	if !ok {
		return nil, errors.New("invalid port returned from provider")
	}

	return &providerPort{name: portID, mode: mode}, nil
}

func GetPortsList() ([]string, error) {
	if ProviderInterface == nil {
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}

	res, err := ProviderInterface.Call("list", nil)
	if err != nil {
		return nil, err
	}

	portsRaw, ok := res["ports"]
	if !ok {
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}

	switch v := portsRaw.(type) {
	case []interface{}:
		ports := make([]string, len(v))
		for i, x := range v {
			ports[i], _ = x.(string)
		}
		return ports, nil
	case string:
		return []string{v}, nil
	default:
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}
}

// -------------------------
// providerPort implementation
// -------------------------

type providerPort struct {
	name    string
	mode    *Mode
	timeout time.Duration
}

func (p *providerPort) SetMode(mode *Mode) error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("setMode", map[string]interface{}{
		"port": p.name,
		"mode": mode,
	})
	if err == nil {
		p.mode = mode
	}
	return err
}

func (p *providerPort) Read(buf []byte) (int, error) {
	if ProviderInterface == nil {
		return 0, &PortError{code: PortClosed}
	}
	res, err := ProviderInterface.Call("read", map[string]interface{}{
		"port": p.name,
		"size": len(buf),
	})
	if err != nil {
		return 0, err
	}
	dataRaw, ok := res["data"].([]byte)
	if !ok {
		return 0, errors.New("invalid data returned from provider")
	}
	n := copy(buf, dataRaw)
	return n, nil
}

func (p *providerPort) Write(buf []byte) (int, error) {
	if ProviderInterface == nil {
		return 0, &PortError{code: PortClosed}
	}
	res, err := ProviderInterface.Call("write", map[string]interface{}{
		"port": p.name,
		"data": buf,
	})
	if err != nil {
		return 0, err
	}
	written, ok := res["written"].(int)
	if !ok {
		written = len(buf)
	}
	return written, nil
}

func (p *providerPort) Drain() error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("drain", map[string]interface{}{"port": p.name})
	return err
}

func (p *providerPort) ResetInputBuffer() error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("resetInputBuffer", map[string]interface{}{"port": p.name})
	return err
}

func (p *providerPort) ResetOutputBuffer() error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("resetOutputBuffer", map[string]interface{}{"port": p.name})
	return err
}

func (p *providerPort) SetDTR(dtr bool) error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("setDTR", map[string]interface{}{"port": p.name, "dtr": dtr})
	return err
}

func (p *providerPort) SetRTS(rts bool) error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("setRTS", map[string]interface{}{"port": p.name, "rts": rts})
	return err
}

func (p *providerPort) GetModemStatusBits() (*ModemStatusBits, error) {
	if ProviderInterface == nil {
		return nil, &PortError{code: PortClosed}
	}
	res, err := ProviderInterface.Call("getModemStatusBits", map[string]interface{}{"port": p.name})
	if err != nil {
		return nil, err
	}
	bits := &ModemStatusBits{}
	if val, ok := res["CTS"].(bool); ok {
		bits.CTS = val
	}
	if val, ok := res["DSR"].(bool); ok {
		bits.DSR = val
	}
	if val, ok := res["RI"].(bool); ok {
		bits.RI = val
	}
	if val, ok := res["DCD"].(bool); ok {
		bits.DCD = val
	}
	return bits, nil
}

func (p *providerPort) SetReadTimeout(t time.Duration) error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	p.timeout = t
	_, err := ProviderInterface.Call("setReadTimeout", map[string]interface{}{"port": p.name, "timeout": t})
	return err
}

func (p *providerPort) Close() error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("close", map[string]interface{}{"port": p.name})
	return err
}

func (p *providerPort) Break(duration time.Duration) error {
	if ProviderInterface == nil {
		return &PortError{code: PortClosed}
	}
	_, err := ProviderInterface.Call("break", map[string]interface{}{"port": p.name, "duration": duration})
	return err
}

// -------------------------
// PortError
// -------------------------

type PortError struct {
	code     PortErrorCode
	causedBy error
}

type PortErrorCode int

const (
	PortBusy PortErrorCode = iota
	PortNotFound
	InvalidSerialPort
	PermissionDenied
	InvalidSpeed
	InvalidDataBits
	InvalidParity
	InvalidStopBits
	InvalidTimeoutValue
	ErrorEnumeratingPorts
	PortClosed
	FunctionNotImplemented
)

func (e PortError) Error() string {
	if e.causedBy != nil {
		return e.EncodedErrorString() + ": " + e.causedBy.Error()
	}
	return e.EncodedErrorString()
}

func (e PortError) EncodedErrorString() string {
	switch e.code {
	case PortBusy:
		return "Serial port busy"
	case PortNotFound:
		return "Serial port not found"
	case InvalidSerialPort:
		return "Invalid serial port"
	case PermissionDenied:
		return "Permission denied"
	case InvalidSpeed:
		return "Port speed invalid or not supported"
	case InvalidDataBits:
		return "Port data bits invalid or not supported"
	case InvalidParity:
		return "Port parity invalid or not supported"
	case InvalidStopBits:
		return "Port stop bits invalid or not supported"
	case InvalidTimeoutValue:
		return "Timeout value invalid or not supported"
	case ErrorEnumeratingPorts:
		return "Could not enumerate serial ports"
	case PortClosed:
		return "Port has been closed"
	case FunctionNotImplemented:
		return "Function not implemented"
	default:
		return "Other error"
	}
}

func (e PortError) Code() PortErrorCode {
	return e.code
}