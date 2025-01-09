// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package lrutracedecisioncache // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/tracedecisioncache/lrutracedecisioncache"

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/tracedecisioncache/lrutracedecisioncache/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		metadata.Type,
		createDefaultConfig,
		createExtension,
		metadata.ExtensionStability,
	)
}

type Config struct {
	Size int `mapstructure:"size"`
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createExtension(_ context.Context, _ extension.Settings, cfg component.Config) (extension.Extension, error) {
	c, ok := cfg.(*Config)
	if !ok {
		// TODO
		return nil, nil
	}

	return NewLRUDecisionCache(c.Size)
}
