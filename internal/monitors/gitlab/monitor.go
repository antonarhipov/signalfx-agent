package gitlab

import (
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
	log "github.com/sirupsen/logrus"
)

const monitorType = "gitlab"

var logger = log.WithFields(log.Fields{"monitorType": monitorType})

func init() {
	monitors.Register(monitorType, func() interface{} { return &Monitor{} }, &Config{})
}

// Config for this monitor
type Config struct {
	prometheusexporter.Config `yaml:",inline"`
}

// Monitor for GitLab
type Monitor struct {
	prometheusexporter.Monitor
}

// Configure the embedded prometheus-exporter monitor
func (m *Monitor) Configure(conf *Config) error {
	return nil
}
