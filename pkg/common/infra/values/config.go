package values

import (
	"time"
)

type DbConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	Schema          string        `yaml:"schema"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxRetries      int           `yaml:"max_retries"`
	MinIdleConns    int           `yaml:"min_idle_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

type Config struct {
	Env     string `yaml:"env"`
	AppName string `yaml:"app"`

	GatewayPort          string `yaml:"gateway_port"`
	PostCrudServicePort  string `yaml:"post_crud_service_port"`
	PostFetchServicePort string `yaml:"post_fetch_service_port"`

	PostCrudServiceUrl  string `yaml:"post_crud_service_url"`
	PostFetchServiceUrl string `yaml:"post_fetch_service_url"`

	DBs            map[string]*DbConfig `yaml:"dbs"`
	ExternalApiUrl string               `yaml:"external_api_url"`
}
