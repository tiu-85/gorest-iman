package providers

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/config"

	"tiu-85/gorest-iman/pkg/common/infra/values"
)

type Provider interface {
	Populate(interface{}) error
	PopulateByKey(string, interface{}) error
}

func NewProviderByArgs(args values.Args) (Provider, error) {
	return NewProviderByFile(args.ConfigFilename)
}

func NewProviderByOptions(options ...config.YAMLOption) (Provider, error) {
	dotenvFile, _ := godotenv.Read(".env.dist", ".env")
	configProvider, err := config.NewYAML(append(
		[]config.YAMLOption{
			config.Expand(EnvExpandFunc(dotenvFile)),
		},
		options...,
	)...)
	if err != nil {
		return nil, err
	}

	return &provider{configProvider}, nil
}

func EnvExpandFunc(extraEnv map[string]string) func(key string) (val string, ok bool) {
	var returnDefault bool
	val, ok := os.LookupEnv("USE_DEFAULT_ENV")
	if ok && strings.ToLower(val) == "true" {
		returnDefault = true
	}
	return func(key string) (val string, ok bool) {
		val, ok = os.LookupEnv(key)
		if !ok {
			val, ok = extraEnv[key]
			if !ok && returnDefault {
				val = "OS_ENV_" + key + "_UNDEFINED"
				ok = true
			}
		}
		return val, ok
	}
}

func NewProviderByFile(filename string) (Provider, error) {
	return NewProviderByOptions(config.File(filename))
}

type provider struct {
	provider config.Provider
}

func (p *provider) Populate(target interface{}) error {
	return p.PopulateByKey("", target)
}

func (p *provider) PopulateByKey(key string, target interface{}) error {
	return p.provider.Get(key).Populate(target)
}
