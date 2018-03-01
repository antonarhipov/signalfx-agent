// Package genericjmx coordinates the various monitors that rely on the
// GenericJMX Collectd plugin to pull JMX metrics.  All of the GenericJMX
// monitors share the same instance of the coreInstance struct that can be
// gotten by the Instance func in this package.  This ultimately means that all
// GenericJMX config will be written to one file to make it simpler to control
// dependencies.
package genericjmx

import (
	"sync"

	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors/collectd"
)

//go:generate collectd-template-to-go genericjmx.tmpl

type connection struct {
	ServiceName     string
	MBeansToCollect []string
}

// Config has configuration that is specific to GenericJMX. This config should
// be used by a monitors that use the generic JMX collectd plugin.
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true"`

	Host string  `yaml:"host" validate:"required"`
	Port uint16  `yaml:"port" validate:"required"`
	Name *string `yaml:"name"`

	// This is how the service type is identified in the SignalFx UI so that
	// you can get built-in content for it.  For custom JMX integrations, it
	// can be set to whatever you like and metrics will get the dimension
	// `sf_hostHasService` set to this value.
	ServiceName string `yaml:"serviceName"`
	// The JMX connection string.  This is rendered as a Go template and has
	// access to the other values in this config.
	ServiceURL     string `yaml:"serviceURL" default:"service:jmx:rmi:///jndi/rmi://{{.Host}}:{{.Port}}/jmxrmi"`
	InstancePrefix string `yaml:"instancePrefix"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	// A list of the MBeans defined in `mBeanDefinitions` to actually collect.
	// If not provided, then all defined MBeans will be collected.
	MBeansToCollect []string `yaml:"mBeansToCollect"`
	// Specifies how to map JMX MBean values to metrics.  If using a specific
	// service monitor such as cassandra, kafka, or activemq, they come
	// pre-loaded with a set of mappings, and any that you add in this option
	// will be merged with those.  See
	// [collectd GenericJMX](https://collectd.org/documentation/manpages/collectd-java.5.shtml#genericjmx_plugin)
	// for more details.
	MBeanDefinitions MBeanMap `yaml:"mBeanDefinitions"`
}

// JMXMonitorCore should be embedded by all monitors that use the collectd
// GenericJMX plugin.  It has most of the logic they will need.  The individual
// monitors mainly just need to provide their set of default mBean definitions.
type JMXMonitorCore struct {
	collectd.MonitorCore

	defaultMBeans      MBeanMap
	defaultServiceName string
	lock               sync.Mutex
}

// NewJMXMonitorCore makes a new JMX core as well as the underlying MonitorCore
func NewJMXMonitorCore(defaultMBeans MBeanMap, defaultServiceName string) *JMXMonitorCore {
	mc := &JMXMonitorCore{
		MonitorCore:        *collectd.NewMonitorCore(CollectdTemplate),
		defaultMBeans:      defaultMBeans,
		defaultServiceName: defaultServiceName,
	}

	mc.MonitorCore.UsesGenericJMX = true
	return mc
}

// Configure configures and runs the plugin in collectd
func (m *JMXMonitorCore) Configure(conf *Config) error {
	conf.MBeanDefinitions = m.defaultMBeans.MergeWith(conf.MBeanDefinitions)
	if conf.MBeansToCollect == nil {
		conf.MBeansToCollect = conf.MBeanDefinitions.MBeanNames()
	}
	if conf.ServiceName == "" {
		conf.ServiceName = m.defaultServiceName
	}

	return m.SetConfigurationAndRun(conf)
}