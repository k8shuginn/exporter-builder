# Introduction
Prometheus provides a variety of exporters to enable monitoring of metrics from different systems and services.
However, as many developers contribute, the structure of each exporter can change, making it challenging to develop one's own new exporter based on existing ones.
To address this issue, the exporter_builder tool was developed. Exporter_builder is a tool that simplifies the process of developing your own exporter by easily creating the basic structure of a Prometheus exporter with the necessary settings.
Using this tool, you can develop an exporter that effectively collects metrics and delivers them to the Prometheus server without complex settings. Let's take a closer look at how you can efficiently develop a custom exporter using exporter_builder.
Utilize exporter_builder to effectively leverage your Prometheus environment and expand the scope of your monitoring.

# Exporter architecture
A Prometheus exporter is composed of multiple collectors, each performing a specific role.
Collectors gather data from systems or applications and convert it into metrics that the Prometheus server can understand.
Each collector is designed to suit its unique characteristics, which may affect the types of data gathered and the form of the metrics.
Therefore, it is crucial to clearly understand the purpose and characteristics of each collector when defining them.
For example, a collector measuring the CPU usage of a system must collect and convert performance-related data into metrics, necessitating precise metric definitions and appropriate setting of collection intervals.
Through this process, Prometheus can effectively monitor data and accurately assess the state of systems or applications.

![exporter](https://github.com/k8shuginn/exporter-builder/assets/79127050/7977708b-b14f-4890-9bbf-83ef62872a86)

# Technology Stack
Exporter_builder is a tool developed to facilitate the easy construction of Prometheus exporters, using a specific technology stack.
The core components of this tool are as follows:
- Go: Exporter_builder is developed in the Go language.
- Cobra: Cobra is a command-line interface (CLI) library developed in Go. Exporter_builder uses Cobra to allow users to easily manipulate and configure the tool through commands.
- Prometheus Go metric library: This library is a metric collection library for Prometheus systems written in Go. Exporter_builder uses this library to support Go-developed exporters in communicating smoothly with the Prometheus server and accurately collecting and delivering metrics.

# Installation
To install exporter_builder, run the following command:
```bash
go install github.com/k8shuginn/exporter-builder/cmd/builder@latest
```

To verify that the builder is installed correctly, run the following command:
```bash
# Check if the builder is installed
ls -al $(go env GOPATH)/bin | grep builder
```

If you encounter a command not found error, set the environment variable as follows:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

# How to use
This section introduces two main methods of configuring the structure of a Prometheus exporter and generating basic Prometheus metrics using exporter_builder: vector and instance methods. These two methods are distinguished by how they handle and provide metrics when receiving requests from the Prometheus server.
- Instance method: This method generates new metrics each time a request for metrics is received from the Prometheus server. This is useful when metrics need to be kept up to date, and real-time updates are essential. For example, it is used when measuring and providing the latest server status or transaction information for each request.
- Vector method: In this method, metrics are generated in advance and registered in a vector. When a request for metrics is received from the Prometheus server, the metrics registered in the vector are returned. This approach is advantageous for shortening request processing times because metrics can be collected and processed in advance. For example, you can continuously monitor the current state of a system, store the results in a vector, and provide them immediately upon request.

First, create a config.yaml file as follows to generate an exporter using exporter_builder:
```yaml
name: my-exporter
module: github.com/my-project/my-exporter
collectors:
  - vector
  - instance
```

After preparing the config.yaml file, the process of building and running the actual exporter using the exporter_builder tool is as follows. This process generates a structured exporter:
```bash
builder --config ./config.yaml
```

This command reads the settings defined in the config.yaml file and configures a new exporter according to those settings.

![tree](https://github.com/k8shuginn/exporter-builder/assets/79127050/f5960d80-4c9e-4493-9959-0f62c9e40924)

main.go: This file is the main execution file of the exporter, initializing and running a server that communicates with the Prometheus server and provides metrics.
collector folder: The files in this folder contain logic for collecting metrics. The collector.go file implements each metric collection method.

# How to develop
## Instance method
The instance method is a method of generating metrics in real-time in a Prometheus exporter and delivering them to the Prometheus server. This method is particularly useful when metrics need to be immediately generated and returned in response to Prometheus requests, reflecting the most up-to-date metric data.

my-exporter/collector/instance/collector.go
```go
package instance

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	val1 = "value1"
	val2 = "value2"

	sampleDesc = prometheus.NewDesc(
		"sample_metric",
		"sample metric",
		[]string{"key1", "key2"}, nil)
)

func SetFlags() {
	flag.StringVar(&val1, "val1", val1, "value1")
	flag.StringVar(&val2, "val2", val2, "value2")
}

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sampleDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(sampleDesc, prometheus.GaugeValue, 1, val1, val2)
}
```

The SetFlags function is a function that sets flags. When flags are set, the flags are automatically registered in main.go.

## Vector method
The vector method is a method of registering predefined metrics in a vector in a Prometheus exporter and delivering them to the Prometheus server when requests are received. This method is advantageous for improving processing speed by updating metric values periodically and providing them immediately upon request.

my-exporter/collector/vector/collector.go

```go
package vector

import (
	"github.com/prometheus/client_golang/prometheus"
)

func SetFlags() {

}

type Collector struct {
	promCollectCount *prometheus.CounterVec
}

func NewCollector() *Collector {
	return &Collector{
		promCollectCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "prometheus_collect_count",
				Help: "Number of times the Prometheus collector is called",
			},
			[]string{"method"},
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.promCollectCount.With(prometheus.Labels{"method": "Describe"}).Inc()
	c.promCollectCount.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.promCollectCount.With(prometheus.Labels{"method": "Collect"}).Inc()
	c.promCollectCount.Collect(ch)
}

```

# exporter_builder flags
When running an exporter created with exporter_builder, you can adjust its operation by setting various flags. These flags allow you to configure the metrics collection endpoint, network address, profiling, and log level according to user needs. Below are the functions and default values of each flag:
### web.telemetry-path
- Function: Sets the HTTP endpoint where the exporter can collect metrics.
- Default value: /metrics
- Example use: Use this flag to specify the URL path where the Prometheus server will collect metrics.
### web.listen-address
- Function: Sets the network address and port used by the exporter to serve metrics.
- Default value: :9090
- Example use: Change the network address and port to prevent port conflicts with other services or to adjust according to security policies.
### profiling
- Function: Sets whether to enable profiling during the exporter's operation.
- Default value: false
- Example use: If performance analysis is needed, set this flag to true to collect profiling data.
### log.level
- Function: Sets the detail level of the logs. This flag allows you to adjust the log levels such as debug, info, and warning.
- Default value: info
- Example use: Increase or decrease the level of detail in the logs to aid in troubleshooting or to manage the volume of logs.


