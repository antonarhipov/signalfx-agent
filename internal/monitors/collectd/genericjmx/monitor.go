// +build linux

package genericjmx

import (
	"github.com/signalfx/signalfx-agent/internal/monitors"
	yaml "gopkg.in/yaml.v2"
)

const monitorType = "collectd/genericjmx"

// Monitor is the main type that represents the monitor
type Monitor struct {
	*JMXMonitorCore
}

func init() {
	err := yaml.Unmarshal([]byte(defaultMBeanYAML), &DefaultMBeans)
	if err != nil {
		panic("YAML for GenericJMX MBeans is invalid: " + err.Error())
	}

	monitors.Register(monitorType, func() interface{} {
		return &Monitor{
			NewJMXMonitorCore(DefaultMBeans, "java"),
		}
	}, &Config{})
}
