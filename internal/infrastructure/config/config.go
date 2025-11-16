package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"veg-store-backend/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
This file handles loading configuration settings from YAML/YML files using Viper.
Logic:
- Determine the mode (dev, prod, etc.) from the MODE environment variable (default to "dev").
- Init configuration from ./config/config.{mode}.yaml or .yml, with a fallback to config.yaml.
- Support environment variable expansion in the format ${VAR} or ${VAR:default}.
- Log loaded configuration values (masking sensitive info like passwords).
- Unmarshal the configuration into a config struct for easy access.

Example YAML structure:
server:
  port: "${SERVER_PORT:8080}"
  api_prefix: "/rest_api"
  api_version: "v1"

database:
  host: "${DB_HOST:localhost}"
  port: 5432

==> config Struct after unmarshalling:
Server:
  Port: "8080"
  ApiPrefix: "/rest_api"
  ApiVersion: "v1"

Data:
  Host: "localhost"
  Port: 5432
*/

// Config - Mapping configuration with yaml structure
type Config struct {
	Mode string

	App struct {
		Name            string `mapstructure:"name"`
		Version         string `mapstructure:"version"`
		EnableDebugLogs bool   `mapstructure:"enable_debug_logs"`
	} `mapstructure:"app"`

	Server struct {
		Port          string `mapstructure:"port"`
		ApiPrefix     string `mapstructure:"api_prefix"`
		ApiVersion    string `mapstructure:"api_version"`
		DefaultLocale string `mapstructure:"default_locale"`
	} `mapstructure:"server"`

	Security struct {
		EnableHTTPS     bool     `mapstructure:"enable_https"`
		KeyStore        string   `mapstructure:"key_store"`
		TrustStore      string   `mapstructure:"trust_store"`
		PublicEndpoints []string `mapstructure:"public_endpoints"`
	} `mapstructure:"security"`

	JWT struct {
		ExpectedIssuer    string   `mapstructure:"expected_issuer"`
		ExpectedAudiences []string `mapstructure:"expected_audiences"`
		AccessDuration    string   `mapstructure:"access_duration"`
		RefreshDuration   string   `mapstructure:"refresh_duration"`
		PrivateKey        string   `mapstructure:"private_key"`
		PublicKey         string   `mapstructure:"public_key"`
	} `mapstructure:"jwt"`

	Cors struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	} `mapstructure:"cors"`

	Swagger struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"swagger"`

	Data struct {
		EnableDataSeeding bool `mapstructure:"enable_data_seeding"`
		Postgres          struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			DBName   string `mapstructure:"db_name"`
			DDLMode  string `mapstructure:"ddl_mode"`
			SSL      struct {
				Mode     string `mapstructure:"mode"`
				RootCert string `mapstructure:"root_cert"`
				Cert     string `mapstructure:"cert"`
				Key      string `mapstructure:"key"`
			} `mapstructure:"ssl"`
		} `mapstructure:"postgres"`
	} `mapstructure:"data"`
}

// Init LoadConfig loads configuration from ./config/config.{mode}.yaml or .yml
func Init(mode string) *Config {
	zap.L().Info(fmt.Sprintf("Init configs for '%s' mode.", mode))

	// Init .env file
	_ = godotenv.Load()

	readConfigWithFallback(mode)

	// Expand ${VAR[:default]} syntax using env values
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		expanded := expandEnvWithDefault(val)
		if val != expanded {
			//component.Logger.Info()("Expanding variable for '%s': '%s' → '%s'\n", key, val, expanded)
			viper.Set(key, expanded)
		}
	}

	zap.L().Info("App info:",
		zap.String("version", viper.GetString("app.name")),
		zap.String("name", viper.GetString("app.version")))

	// Log all loaded values (masking passwords)
	if mode != "prod" && mode != "production" {
		logAppConfig()
	}

	// Unmarshal (is equivalent to Decode) into config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		zap.L().Fatal("failed to unmarshal config", zap.Error(err))
	}
	config.Mode = mode

	return &config
}

// --- Helper: fallback config reader ---
func readConfigWithFallback(mode string) {
	// setup Viper with fallback order: config.{mode}.yaml → config.{mode}.yml → config.yaml
	viper.SetConfigName(fmt.Sprintf("config.%s", mode))
	viper.SetConfigType("yaml")

	// Set config path to .../.../config
	configPath := util.GetConfigPathFromGoMod("config")
	zap.L().Info(fmt.Sprintf("config path: %s", configPath))
	viper.AddConfigPath(configPath)

	// Try .yaml first
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	// Try .yml
	ymlPath := strings.TrimSuffix(viper.ConfigFileUsed(), filepath.Ext(viper.ConfigFileUsed())) + ".yml"
	viper.SetConfigFile(ymlPath)
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	// Try config.yaml (default)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	zap.L().Fatal("no valid config file found")
}

// --- Helper: supports ${VAR} and ${VAR:default} ---
func expandEnvWithDefault(input string) string {
	// Regex to match ${VAR} or ${VAR:default}
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(?::([^}]*))?}`)

	return re.ReplaceAllStringFunc(input, func(s string) string {
		// Extract variable name and default value
		matches := re.FindStringSubmatch(s)

		// Ensure we have at least the variable name
		if len(matches) < 2 {
			return s
		}

		// matches[1] is the variable name, matches[2] is the default value (if any)
		key := matches[1]
		defaultValue := ""

		// If default value is provided
		if len(matches) == 3 {
			defaultValue = matches[2]
		}

		if value, ok := os.LookupEnv(key); ok && value != "" {
			return value
		}
		return defaultValue
	})
}

func logAppConfig() {
	configFile := viper.ConfigFileUsed()
	var fields []zap.Field

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		if strings.Contains(strings.ToLower(key), "password") {
			val = "********"
		}
		fields = append(fields, zap.Any(key, val))
	}

	zap.L().Info("Application configuration loaded",
		zap.String("config_file", configFile),
		zap.Any("configs", fields),
	)
}
