//go:build !windows
// +build !windows

package latencywin

// InputConf ...
type InputConf struct {
	Type   string `json:"type"` // mouse or keyboard
	IsAuto bool   `json:"isAuto"`
	KeyTap string `json:"keyTap"`
}

// Config ...
type Config struct {
	InputConf InputConf `json:"inputConf,omitempty"`
	Frames    int       `json:"frames,omitempty"`
	StartKey  string    `json:"startKey"`
}

// OpLatencyWindowsManager ...
type OpLatencyWindowsManager struct {
}

func NewOpLatencyWindowsManager() *OpLatencyWindowsManager {
	return &OpLatencyWindowsManager{}
}
