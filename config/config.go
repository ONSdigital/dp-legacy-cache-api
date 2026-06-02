package config

import (
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

const (
	CacheTimesCollection = "CacheTimesCollection"
	port29100            = ":29100"
	localhost8082        = "http://localhost:8082"
	localhost27017       = "localhost:27017"
	databaseName         = "cache"
	collectionName       = "cachetimes"
)

type MongoConfig = mongodb.MongoDriverConfig

// Config represents service configuration for dp-legacy-cache-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	IsPublishing               bool          `envconfig:"IS_PUBLISHING"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	MongoConfig
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   port29100,
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		IsPublishing:               false,
		ZebedeeURL:                 localhost8082,
		MongoConfig: MongoConfig{
			ClusterEndpoint:               localhost27017,
			Username:                      "",
			Password:                      "",
			Database:                      databaseName,
			Collections:                   map[string]string{CacheTimesCollection: collectionName},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
	}

	return cfg, envconfig.Process("", cfg)
}
