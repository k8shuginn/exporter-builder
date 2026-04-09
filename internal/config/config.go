// Package config exporter 생성에 필요한 설정 구조체와 기본값을 제공합니다.
package config

const (
	// DefaultName 설정 파일에 name이 없을 때 사용하는 기본 exporter 이름입니다.
	DefaultName = "DiyExporter"
	// DefaultModule 설정 파일에 module이 없을 때 사용하는 기본 Go 모듈 경로입니다.
	DefaultModule = "DiyModule"
	// DefaultGoVersion 설정 파일에 go_version이 없을 때 사용하는 기본 Go 버전입니다.
	DefaultGoVersion = "1.26"
)

// Config exporter 생성에 필요한 설정 정보를 담는 구조체입니다.
type Config struct {
	Name       string   `mapstructure:"name"`
	Module     string   `mapstructure:"module"`
	GoVersion  string   `mapstructure:"go_version"`
	Collectors []string `mapstructure:"collectors"`
}

// NewConfig 기본값이 설정된 Config 포인터를 반환합니다.
func NewConfig() *Config {
	return &Config{
		Name:       DefaultName,
		Module:     DefaultModule,
		GoVersion:  DefaultGoVersion,
		Collectors: make([]string, 0),
	}
}
