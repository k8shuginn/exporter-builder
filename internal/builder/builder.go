// Package builder 설정을 기반으로 exporter 프로젝트를 생성하는 핵심 로직을 제공합니다.
package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/k8shuginn/exporter-builder/internal/config"
)

// GenerateExporter cfg를 받아 exporter 프로젝트 디렉토리와 소스 파일을 생성합니다.
// 생성 실패 시 error를 반환합니다.
func GenerateExporter(cfg *config.Config) error {
	cfg.Collectors = removeDuplicateCollector(cfg.Collectors)

	if err := createDirectory(cfg.Name, cfg.Collectors); err != nil {
		return err
	}
	if err := generateMain(cfg); err != nil {
		return err
	}
	if err := generateCollector(cfg); err != nil {
		return err
	}
	if err := generateGoMod(cfg); err != nil {
		return err
	}

	return nil
}

// removeDuplicateCollector collectors 슬라이스에서 중복된 항목을 제거하여 반환합니다.
func removeDuplicateCollector(collectors []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, collector := range collectors {
		if _, ok := seen[collector]; !ok {
			seen[collector] = struct{}{}
			result = append(result, collector)
		}
	}

	return result
}

// createDirectory dir 이름의 루트 디렉토리와 collectors 각각의 하위 디렉토리를 생성합니다.
// 생성 실패 시 error를 반환합니다.
func createDirectory(dir string, collectors []string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed create directory: %s", err)
	}

	for _, collector := range collectors {
		path := filepath.Join(dir, "collector", collector)
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("failed create collector directory: %s", err)
		}
	}

	return nil
}

// generateMain cfg를 받아 exporter의 main.go 파일을 생성합니다.
// 파일 생성 또는 템플릿 실행 실패 시 error를 반환합니다.
func generateMain(cfg *config.Config) error {
	out, err := os.Create(filepath.Join(cfg.Name, mainTemplate.Name()))
	if err != nil {
		return fmt.Errorf("failed create main file: %s", err)
	}
	defer out.Close()

	return mainTemplate.Execute(out, cfg)
}

// generateCollector cfg를 받아 각 collector의 소스 파일을 생성합니다.
// 파일 생성 또는 템플릿 실행 실패 시 error를 반환합니다.
func generateCollector(cfg *config.Config) error {
	for _, collector := range cfg.Collectors {
		out, err := os.Create(filepath.Join(cfg.Name, "collector", collector, collectorTemplate.Name()))
		if err != nil {
			return fmt.Errorf("failed create collector file: %s", err)
		}
		defer out.Close()

		if err := collectorTemplate.Execute(out, CollectorValues{Name: collector}); err != nil {
			return fmt.Errorf("failed execute collector template: %s", err)
		}
	}

	return nil
}

// generateGoMod cfg를 받아 go.mod 파일을 생성하고 go mod tidy 및 go mod vendor를 실행합니다.
// 파일 생성, 템플릿 실행, 또는 명령 실행 실패 시 error를 반환합니다.
func generateGoMod(cfg *config.Config) error {
	out, err := os.Create(filepath.Join(cfg.Name, goModTemplate.Name()))
	if err != nil {
		return fmt.Errorf("failed create go.mod file: %s", err)
	}
	defer out.Close()

	if err = goModTemplate.Execute(out, cfg); err != nil {
		return fmt.Errorf("failed execute go.mod template: %s", err)
	}

	if err := runCommand(cfg.Name, "go", "mod", "tidy"); err != nil {
		return err
	}

	return runCommand(cfg.Name, "go", "mod", "vendor")
}

// runCommand dir 디렉토리에서 args 명령을 실행합니다.
// 실행 실패 시 명령 출력을 포함한 error를 반환합니다.
func runCommand(dir string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	if result, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed execute %s: %s", args[0]+" "+args[1], result)
	}
	return nil
}
