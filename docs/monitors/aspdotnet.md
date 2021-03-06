<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# aspdotnet

(Windows Only) This monitor reports metrics about requests, errors, sessions,
worker processes for ASP.NET applications.

## Windows Performance Counters
The underlying source for these metrics are Windows Performance Counters.
Most of the performance counters that we query in this monitor are actually Gauges
that represent rates per second and percentages.

This monitor reports the instantaneous values for these Windows Performance Counters.
This means that in between a collection interval, spikes could occur on the
Performance Counters.  The best way to mitigate this limitation is to increase
the reporting interval on this monitor to collect more frequently.

Sample YAML configuration:

```yaml
monitors:
 - type: aspdotnet
```


Monitor Type: `aspdotnet`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/aspdotnet)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `counterRefreshInterval` | no | `int64` | (Windows Only) Number of seconds that wildcards in counter paths should be expanded and how often to refresh counters from configuration. (**default:** `60s`) |
| `printValid` | no | `bool` | (Windows Only) Print out the configurations that match available performance counters.  This is used for debugging. (**default:** `false`) |




## Metrics

The following table lists the metrics available for this monitor. Metrics that are not marked as Custom are standard metrics and are monitored by default.

| Name | Type | Custom | Description |
| ---  | ---  | ---    | ---         |
| `asp_net.application_restarts` | gauge | X | Count of ASP.NET application restarts. |
| `asp_net.applications_running` | gauge | X | Number of running ASP.NET applications. |
| `asp_net.requests_current` | gauge | X | Current number of ASP.NET requests. |
| `asp_net.requests_queue` | gauge | X | Number of queued ASP.NET requests. |
| `asp_net.requests_rejected` | gauge | X | Count of rejected ASP.NET requests. |
| `asp_net.worker_process_restarts` | gauge | X | Count of ASP.NET worker process restarts. |
| `asp_net.worker_processes_running` | gauge | X | Number of running ASP.NET worker processes. |
| `asp_net_applications.errors_during_execution` | gauge | X | Count of errors encountered by ASP.NET application durring execution. |
| `asp_net_applications.errors_total_sec` | gauge | X | Error rate per second for the given ASP.NET application. |
| `asp_net_applications.errors_unhandled_during_execution_sec` | gauge | X | Unhandled error rate per second countered while an ASP.NET application is running. |
| `asp_net_applications.pipeline_instance_count` | gauge | X | Number of instances in the ASP.NET application pipeline. |
| `asp_net_applications.requests_failed` | gauge | X | Count of failed requests in the ASP.NET application |
| `asp_net_applications.requests_sec` | gauge | X | Rate of requests in the ASP.NET application per second. |
| `asp_net_applications.session_sql_server_connections_total` | gauge | X | Number of connections to microsoft sql server by an ASP.NET application. |
| `asp_net_applications.sessions_active` | gauge | X | Number of active sessions in the ASP.NET application. |


To specify custom metrics you want to monitor, add a `metricsToInclude` filter
to the agent configuration, as shown in the code snippet below. The snippet
lists all available custom metrics. You can copy and paste the snippet into
your configuration file, then delete any custom metrics that you do not want
sent.

Note that some of the custom metrics require you to set a flag as well as add
them to the list. Check the monitor configuration file to see if a flag is
required for gathering additional metrics.

```yaml

metricsToInclude:
  - metricNames:
    - asp_net.application_restarts
    - asp_net.applications_running
    - asp_net.requests_current
    - asp_net.requests_queue
    - asp_net.requests_rejected
    - asp_net.worker_process_restarts
    - asp_net.worker_processes_running
    - asp_net_applications.errors_during_execution
    - asp_net_applications.errors_total_sec
    - asp_net_applications.errors_unhandled_during_execution_sec
    - asp_net_applications.pipeline_instance_count
    - asp_net_applications.requests_failed
    - asp_net_applications.requests_sec
    - asp_net_applications.session_sql_server_connections_total
    - asp_net_applications.sessions_active
    monitorType: aspdotnet
```




