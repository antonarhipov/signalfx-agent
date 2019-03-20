package client

import (
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
)

func prepareGaugeHelper(metricName string, dims map[string]string, metricValue *int64) *datapoint.Datapoint {
	if metricValue == nil {
		return nil
	}
	return sfxclient.Gauge(metricName, dims, *metricValue)
}

func prepareCumulativeHelper(metricName string, dims map[string]string, metricValue *int64) *datapoint.Datapoint {
	if metricValue == nil {
		return nil
	}
	return sfxclient.Cumulative(metricName, dims, *metricValue)
}
