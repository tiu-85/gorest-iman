package adapters

import (
	"tiu-85/gorest-iman/pkg/common/infra/providers"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

var cfg values.Config

func NewDefaultConfig(provider providers.Provider) (*values.Config, error) {
	err := provider.Populate(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
