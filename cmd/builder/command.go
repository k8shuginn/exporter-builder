// Copyright 2024 k8shuginn exporter_builder
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main exporter-builder CLI의 진입점 패키지입니다.
package main

import (
	"fmt"
	"log"

	"github.com/k8shuginn/exporter-builder/internal/builder"
	"github.com/k8shuginn/exporter-builder/internal/config"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
)

const (
	// ExampleMessage CLI 사용 예시 메시지입니다.
	ExampleMessage = "builder --config-file config.yaml"
)

// Command builder CLI 커맨드를 생성하여 반환합니다.
func Command() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "builder",
		Short:   "builder is a command line tool to generate exporter",
		Example: ExampleMessage,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := initConfig(configPath)
			if err != nil {
				return err
			}
			if err := builder.GenerateExporter(cfg); err != nil {
				return err
			}

			log.Println("build exporter success")
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "config-file", "./config.yaml", "config file path")

	return cmd
}

// initConfig configPath 경로의 YAML 파일을 읽어 Config를 반환합니다.
// 파일 로드 또는 언마샬 실패 시 error를 반환합니다.
func initConfig(configPath string) (*config.Config, error) {
	k := koanf.New(".")
	cfg := config.NewConfig()

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return cfg, nil
}
