# Introduction
Prometheus는 다양한 exporter를 제공하여 서로 다른 시스템과 서비스의 메트릭을 모니터링할 수 있도록 지원합니다. 그러나 많은 개발자들이 참여하면서 exporter 각각의 구조가 달라지게 되어 기존의 exporter를 참고하여 자신만의 새로운 exporter를 개발하는데 있어 어려움을 겪을 수 있습니다. 이러한 문제를 해결하기 위해 exporter_builder 도구가 개발되었습니다.
exporter_builder는 Prometheus의 exporter를 개발할 때 필요한 설정으로 쉽게 exporter의 기본 구조를 생성해주는 도구로, 사용자가 자신만의 exporter를 개발하는 과정을 간소화합니다. 이 도구를 사용하면, 복잡한 설정 없이도 효과적으로 메트릭을 수집하고 Prometheus 서버에 전달할 수 있는 exporter를 만들 수 있습니다. exporter_builder를 사용하여 어떻게 효율적으로 사용자 정의 exporter를 개발할 수 있는지 자세히 알아보겠습니다. exporter_builder를 이용하여 자신만의 exporter로 Prometheus 환경을 보다 효과적으로 활용하고, 모니터링의 범위를 확장할 수 있는 기회를 얻어 보세요.

# Exporter architecture
Prometheus exporter는 여러 개의 collector로 구성되어 있으며, 이들 각각은 특정한 역할을 수행합니다. Collector는 시스템이나 응용 프로그램으로부터 데이터를 수집하고, 이를 Prometheus 서버가 이해할 수 있는 메트릭 형태로 변환하여 전달합니다. 각 collector는 고유의 특성에 맞게 설계되어 있으며, 이에 따라 수집하는 데이터의 종류나 메트릭의 형태가 다를 수 있습니다.
따라서 collector를 정의할 때는 해당 collector의 목적과 특성을 명확히 이해하고, 이에 맞게 설정을 조정하는 것이 중요합니다. 예를 들어, 시스템의 CPU 사용률을 측정하는 collector는 시스템의 성능 관련 데이터를 수집하고 이를 메트릭으로 변환해야 하므로, 관련된 메트릭을 정확하게 정의하고 수집 주기를 적절히 설정하는 것이 필수적입니다. 이러한 과정을 통해 Prometheus는 효과적으로 데이터를 모니터링하고, 시스템이나 응용 프로그램의 상태를 정확하게 파악할 수 있습니다.

![exporter](https://github.com/k8shuginn/exporter-builder/assets/79127050/7977708b-b14f-4890-9bbf-83ef62872a86)

# Technology Stack
exporter_builder는 Prometheus exporter를 쉽게 구축할 수 있도록 도와주는 도구로, 특정 기술 스택을 사용하여 개발되었습니다. 이 도구의 핵심 구성 요소는 다음과 같습니다:
- Go: exporter_builder는 Go 언어로 개발되었습니다.
- Cobra: Cobra는 Go 언어로 개발된 커맨드 라인 인터페이스(CLI) 라이브러리입니다. exporter_builder에서는 Cobra를 사용하여 사용자가 명령어를 통해 도구를 쉽게 조작하고 설정할 수 있도록 합니다.
- Prometheus Go metric library: 이 라이브러리는 Go 언어로 작성된 Prometheus 시스템을 위한 메트릭 수집 라이브러리입니다. exporter_builder는 이 라이브러리를 사용하여 Go 언어로 개발된 exporter가 Prometheus 서버와 원활하게 통신하고, 메트릭을 정확하게 수집 및 전달할 수 있도록 지원합니다.

# Installation
exporter_builder를 설치하려면 다음 명령어를 실행하세요.
```bash
go install github.com/k8shuginn/exporter-builder/cmd/builder@latest
```

정상적으로 설치되었는지 확인하려면 다음 명령어를 실행하세요.
```bash
# Check if the builder is installed
ls -al $(go env GOPATH)/bin | grep builder
```

만약 command not found 에러가 발생하면, 환경 변수를 다음과 같이 설정합니다.
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

# How to use
exporter_builder를 사용하여 Prometheus exporter의 구조를 설정하고 기본적으로 Prometheus 메트릭을 생성하는 두 가지 주요 방법을 소개하겠습니다: vector 방식과 instance 방식. 이 두 방법은 Prometheus 서버로부터 요청을 받았을 때 메트릭을 어떻게 처리하고 제공하는지에 따라 구분됩니다.
- instance 방식 : 이 방식은 Prometheus 서버로부터 메트릭 요청이 들어올 때마다 새로운 메트릭을 생성하여 반환합니다. 이는 메트릭이 최신 상태로 유지되어야 하는 경우에 유용하며, 메트릭의 실시간 업데이트가 중요한 환경에서 적합합니다. 예를 들어, 매 요청마다 최신의 서버 상태나 트랜잭션 정보를 측정하고 제공해야 할 때 사용됩니다.
- vector 방식 : 이 방식에서는 메트릭을 미리 생성하고 vector에 등록합니다. Prometheus 서버로부터 메트릭 요청이 들어올 때, 사전에 vector에 등록된 메트릭들을 반환합니다. 이 접근 방식은 메트릭 수집과 처리를 미리 할 수 있기 때문에 요청 처리 시간을 단축시키는 데 유리합니다. 예를 들어, 시스템의 현재 상태를 지속적으로 모니터링하며 그 결과를 vector에 저장하고, 요청 시 즉시 이를 제공할 수 있습니다.

먼저 exporter_builder를 사용하여 exporter를 생성하기 위해 다음과 같이 config.yaml 파일을 작성합니다.
```yaml
name: my-exporter
module: github.com/my-project/my-exporter
collectors:
  - vector
  - instance
```

config.yaml 파일을 준비한 후, exporter_builder 도구를 사용하여 실제로 exporter를 빌드하고 실행하는 과정은 다음과 같습니다. 이 과정을 통해 구조화된 exporter가 생성됩니다.
```bash
builder --config ./config.yaml
```
이 명령어는 config.yaml 파일에 정의된 설정을 읽고, 해당 설정에 맞추어 새로운 exporter를 구성합니다.

![tree](https://github.com/k8shuginn/exporter-builder/assets/79127050/f5960d80-4c9e-4493-9959-0f62c9e40924)

main.go: 이 파일은 exporter의 주 실행 파일로, Prometheus 서버와 통신하고 메트릭을 제공하는 서버를 초기화하고 실행합니다.
collector 폴더: 이 폴더 내의 파일들은 메트릭을 수집하는 로직을 담고 있습니다. collector.go 파일은 각 메트릭 수집 방법을 구현합니다.

# How to develop
## instance 방식
Instance 방식은 Prometheus exporter에서 실시간으로 메트릭을 생성하고 Prometheus 서버에 전달하는 방법입니다. 이 방식은 Prometheus의 요청에 따라 메트릭을 즉시 생성하여 결과를 반환하므로, 메트릭 데이터가 최신 상태를 반영할 필요가 있는 경우에 특히 유용합니다.
이 방식은 메트릭 데이터의 실시간성이 중요한 경우, 예를 들어 서버의 현재 부하 상태나 트랜잭션 수와 같은 동적 정보를 모니터링할 때 효과적입니다. 따라서, 실시간 데이터 반영이 중요한 모니터링 시스템에서는 Instance 방식을 적극적으로 활용하는 것이 좋습니다.
instance방식은 exporter_builder에서 제공하는 sample metric와 같이 metric을 생성하는 코드를 추가하여 사용합니다.

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
SetFlags 함수는 flag를 설정하는 함수로, flag를 설정하면 main.go에서 flag를 자동으로 등록해줍니다.

## vector 방식
Vector 방식은 Prometheus exporter에서 미리 정의된 메트릭을 vector에 등록하고, Prometheus 서버의 요청이 들어올 때 등록된 메트릭을 전달하는 방법입니다. 이 방식은 주기적으로 메트릭 값을 업데이트하여 저장하고, 요청 시 즉시 이 값을 제공함으로써 처리 속도를 향상시킬 수 있습니다. Vector 방식은 시스템의 성능을 최적화하고자 할 때 유용하며, 특히 메트릭 데이터의 빠른 응답 시간이 필요한 경우에 적합합니다. 이를 통해 Prometheus 모니터링 시스템의 전반적인 효율성을 향상시킬 수 있습니다. vector 방식은 sample metric보다는 아래와 같이 코드를 수정하여 사용합니다.

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

# exporter_builder를 사용한 Exporter 실행 시 설정 가능한 플래그
exporter_builder로 생성된 exporter를 실행할 때, 다양한 플래그를 설정하여 exporter의 동작을 조정할 수 있습니다. 이러한 플래그를 통해 메트릭 수집의 엔드포인트, 네트워크 주소, 프로파일링 및 로그 레벨을 사용자의 요구에 맞게 구성할 수 있습니다. 아래는 각 플래그의 기능과 기본값을 설명합니다.
### web.telemetry-path
- 기능: Exporter가 메트릭을 수집할 수 있는 HTTP 엔드포인트를 설정합니다.
- 기본값: /metrics
- 사용 예: 이 플래그를 사용하여 Prometheus 서버가 메트릭을 수집할 URL 경로를 지정할 수 있습니다.
### web.listen-address
- 기능: Exporter가 메트릭을 제공할 때 사용할 네트워크 주소와 포트를 설정합니다.
- 기본값: :9090
- 사용 예: 네트워크 주소와 포트를 변경하여 다른 서비스와의 포트 충돌을 방지하거나 보안 정책에 맞게 조정할 수 있습니다.
### profiling
- 기능: Exporter 실행 중 프로파일링 기능을 활성화할지 여부를 설정합니다.
- 기본값: false
- 사용 예: 성능 분석이 필요한 경우, 이 플래그를 true로 설정하여 프로파일링 데이터를 수집할 수 있습니다.
### log.level
- 기능: 로그의 상세 수준을 설정합니다. 이 플래그를 통해 디버그, 정보, 경고 등의 로그 레벨을 조정할 수 있습니다.
- 기본값: info
- 사용 예: 로그의 상세 수준을 높이거나 낮춰서 문제 해결을 돕거나 로그의 양을 조절할 수 있습니다.
