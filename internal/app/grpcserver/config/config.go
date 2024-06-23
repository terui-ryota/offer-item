package config

import (
	"fmt"
	"os"

	//"github.com/ca-media-nantes/libgo/v2/configs/definitions"
	//"github.com/ca-media-nantes/libgo/v2/tracing"
	commonConf "github.com/terui-ryota/offer-item/internal/common/config"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

type GRPCConfig struct {
	GrpcPort uint16 `yaml:"grpc_port"`
	//MonitorPort uint16      `yaml:"monitor_port"`
	Logger *zap.Config `yaml:"logger"`
	//Tracing     *tracing.Config      `yaml:"tracing"`
	Database *commonConf.Database `yaml:"database"`
	// Databases        *commonConf.Databases        `yaml:"database"` // TODO:
	ExternalContexts *commonConf.ExternalContexts `yaml:"external_contexts"`
	Validation       *ValidationConfig            `yaml:"validation"`
}

type ValidationConfig struct {
	MaxInputAssigneeListNum int `yaml:"max_input_assignee_list_num"`
}

func LoadConfig() *GRPCConfig {
	defaultFile := "config.yaml"
	confFile := ""
	env := os.Getenv("ENV")

	switch env {
	case "dev", "stg", "prd":
		confFile = fmt.Sprintf("%s_%s", env, defaultFile)
	case "", "local":
		confFile = defaultFile
	}

	confContent, err := os.ReadFile(fmt.Sprintf("./configs/grpcserver/%s", confFile))
	if err != nil {
		panic(err)
	}

	// expand environment variables
	expandedConfContent := os.ExpandEnv(string(confContent))

	var cfg GRPCConfig
	if err := yaml.Unmarshal([]byte(expandedConfContent), &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
