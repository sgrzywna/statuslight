package models

import (
	"fmt"
	"strings"
)

const (
	// On turns light on.
	On = "on"
	// Off turns light off.
	Off = "off"
	// SeqRunning represents state of the running sequence.
	SeqRunning = "running"
	// SeqStopped represents state of the stopped sequence.
	SeqStopped = "stopped"
	// SeqPaused represents state of the paused sequence.
	SeqPaused = "paused"
)

// Sequence represents light control sequence.
type Sequence struct {
	Name  string         `json:"name"`
	Steps []SequenceStep `json:"steps"`
}

// SequenceStep represents single step from the light control sequence.
type SequenceStep struct {
	Light    Light `json:"light"`
	Duration int   `json:"duration"`
}

// Light represents command to control light.
type Light struct {
	Color      *string `json:"color"`
	Brightness *int    `json:"brightness"`
	Switch     *string `json:"switch"`
}

// SetColor sets color name.
func (l *Light) SetColor(color string) {
	l.Color = new(string)
	*l.Color = color
}

// SetBrightness sets light brightness.
func (l *Light) SetBrightness(brightness int) {
	l.Brightness = new(int)
	*l.Brightness = brightness
}

// SetSwitch sets light state.
func (l *Light) SetSwitch(state bool) {
	l.Switch = new(string)
	if state {
		*l.Switch = On
	} else {
		*l.Switch = Off
	}
}

// Clear sets all attributes to their zero values.
func (l *Light) Clear() {
	l.Color = nil
	l.Brightness = nil
	l.Switch = nil
}

// String implements string representation for the Light structure.
func (l *Light) String() string {
	var items []string
	if l.Color != nil {
		items = append(items, fmt.Sprintf("color:%s", *l.Color))
	} else {
		items = append(items, "color:nil")
	}
	if l.Brightness != nil {
		items = append(items, fmt.Sprintf("brightness:%d", *l.Brightness))
	} else {
		items = append(items, "brightness:nil")
	}
	if l.Switch != nil {
		items = append(items, fmt.Sprintf("switch:%s", *l.Switch))
	} else {
		items = append(items, "switch:nil")
	}
	return strings.Join(items, ",")
}

// SequenceState represents sequence state.
type SequenceState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}
