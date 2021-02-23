package main

// Config is the root configuration object
// represented as a map of EventConfigs
type Config map[string]EventConfig

// EventConfig represents an event target configuration
type EventConfig struct {
	Type                   string            `yaml:"type"`
	MatterMostWebsocketURL string            `yaml:"mm_ws_url"`
	MatterMostURL          string            `yaml:"mm_url"`
	MatterMostToken        string            `yaml:"mm_token"`
	MatterMostSendAuth     bool              `yaml:"mm_send_auth"`
	MatterMostEvents       []string          `yaml:"mm_events"`
	URL                    string            `yaml:"url"`
	Headers                map[string]string `yaml:"headers"`
}
