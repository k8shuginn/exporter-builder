// Package builder 설정을 기반으로 exporter 프로젝트를 생성하는 핵심 로직을 제공합니다.
package builder

import (
	_ "embed"
	"text/template"
)

// CollectorValues collector 템플릿 렌더링에 사용되는 값 구조체입니다.
type CollectorValues struct {
	Name string
}

var (
	//go:embed template/go.mod.tmpl
	goModBytes    []byte
	goModTemplate = parseTemplate("go.mod", goModBytes)

	//go:embed template/main.go.tmpl
	mainBytes    []byte
	mainTemplate = parseTemplate("main.go", mainBytes)

	//go:embed template/collector.go.tmpl
	collectorBytes    []byte
	collectorTemplate = parseTemplate("collector.go", collectorBytes)
)

// parseTemplate name과 bytes를 받아 파싱된 Template 포인터를 반환합니다.
// 파싱 실패 시 패닉을 발생시킵니다.
func parseTemplate(name string, bytes []byte) *template.Template {
	return template.Must(template.New(name).Parse(string(bytes)))
}
