package config

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"

	commonConf "github.com/terui-ryota/offer-item/internal/common/config"
	"golang.org/x/net/publicsuffix"

	libtime "github.com/terui-ryota/offer-item/pkg/time"
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
	Rakuten          *RakutenConfig               `yaml:"rakuten"`
	HttpClient       HttpClient                   `yaml:"http_client"`
}

type ValidationConfig struct {
	MaxInputAssigneeListNum int `yaml:"max_input_assignee_list_num"`
}

type RakutenConfig struct {
	ApplicationID []string            `yaml:"application_id"`
	RateLimit     int                 `yaml:"rate_limit"`
	RakutenIchiba RakutenIchibaConfig `yaml:"ichiba"`
}

type RakutenIchibaConfig struct {
	Format        string `yaml:"format"`
	PartnerID     string `yaml:"partner_id"`
	ClickIDPrefix string `yaml:"click_id_prefix"`
}

type HttpClient struct {
	MaxIdleConn           int              `yaml:"max_idle_conn"`
	MaxIdleConnsPerHost   int              `yaml:"max_idle_conns_per_host"`
	IdleConnTimeout       libtime.Duration `yaml:"idle_conn_timeout"`
	DialTimeout           libtime.Duration `yaml:"dial_timeout"`
	DialKeepAlive         libtime.Duration `yaml:"dial_keep_alive"`
	HttpClientTimeout     libtime.Duration `yaml:"http_client_timeout"`
	InsecureSkipVerify    bool             `yaml:"insecure_skip_verify"`
	TLSHandshakeTimeout   libtime.Duration `yaml:"tls_handshake_timeout"`
	ResponseHeaderTimeout libtime.Duration `yaml:"response_header_timeout"`
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

func LoadHttpClient(config *GRPCConfig) (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	client := &http.Client{
		Timeout: config.HttpClient.HttpClientTimeout.Duration,
		//Transport: &https.MetricsTransport{Transport: makeTransport(config)},
		Jar: jar, // 認証系で時間がかかったりするため
	}

	return client, nil
}
