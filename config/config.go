package config

import (
	"time"

	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

const (
	DatasetsCollection         = "DatasetsCollection"
	ContactsCollection         = "ContactsCollection"
	EditionsCollection         = "EditionsCollection"
	InstanceCollection         = "InstanceCollection"
	DimensionOptionsCollection = "DimensionOptionsCollection"
	InstanceLockCollection     = "InstanceLockCollection"
)

// Config represents service configuration for dp-legacy-cache-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	MongoConfig
}

type MongoConfig struct {
	mongodriver.MongoDriverConfig

	// CodeListAPIURL string `envconfig:"CODE_LIST_API_URL"`
	// DatasetAPIURL  string `envconfig:"DATASET_API_URL"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:29100",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MongoConfig: MongoConfig{
			MongoDriverConfig: mongodriver.MongoDriverConfig{
				ClusterEndpoint:               "localhost:27017",
				Username:                      "",
				Password:                      "",
				Database:                      "datasets",
				Collections:                   map[string]string{DatasetsCollection: "datasets", ContactsCollection: "contacts", EditionsCollection: "editions", InstanceCollection: "instances", DimensionOptionsCollection: "dimension.options", InstanceLockCollection: "instances_locks"},
				ReplicaSet:                    "",
				IsStrongReadConcernEnabled:    false,
				IsWriteConcernMajorityEnabled: true,
				ConnectTimeout:                5 * time.Second,
				QueryTimeout:                  15 * time.Second,
				TLSConnectionConfig: mongodriver.TLSConnectionConfig{
					IsSSL: false,
				},
			},
		},
	}

	return cfg, envconfig.Process("", cfg)
}
