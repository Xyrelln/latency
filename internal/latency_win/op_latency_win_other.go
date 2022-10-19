//go:build !windows
// +build !windows

package latencywin

// InputConf ...
type InputConf struct {
	Type   string `json:"type"` // mouse or keyboard
	IsAuto bool   `json:"isAuto"`
	KeyTap string `json:"keyTap,omitempty"`
	// MousePos [2]int `json:"mouse_pos,omitempty"`
}

// Config ...
type Config struct {
	InputConf          InputConf `json:"inputCconf,omitempty"`
	ImageDiffThreshold int       `json:"imageDiff_threshold"`
	Frames             int       `json:"frames,omitempty"`
	StartKey           string    `json:"startKey"`
	// OffsetMs           int       `json:"offset_ms"`
}

// OpLatencyWindowsManager ...
type OpLatencyWindowsManager struct {
}
