package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	VersionDate = "0001-01-01 00:00:00"
	Version     = "dev"
	Hash        = "COMMIT ID"
)

// DefaultConfig возвращает конфигурацию по умолчанию.
// будет использоваться в случае когда config.yaml не найден
func DefaultConfig() *Config {
	return &Config{
		HttpPort:                 8080,
		BaseURL:                  "http://127.0.0.1:8080",
		NewsScraperAddress:       "127.0.0.1",
		NewsScraperPort:          8081,
		CommentsServiceAddress:   "127.0.0.1",
		CommentsServicePort:      8082,
		ModerationServiceAddress: "127.0.0.1",
		ModerationServicePort:    8083,
	}
}

// конфигурация приложения, подразумевается yaml-формат
type Config struct {
	HttpPort                 int    `yaml:"http_port"`
	BaseURL                  string `yaml:"base_url"`
	NewsScraperAddress       string `yaml:"news_scraper_address"`
	NewsScraperPort          int    `yaml:"news_scraper_port"`
	CommentsServiceAddress   string `yaml:"comments_service_address"`
	CommentsServicePort      int    `yaml:"comments_service_port"`
	ModerationServiceAddress string `yaml:"moderation_service_address"`
	ModerationServicePort    int    `yaml:"moderation_service_port"`
}

func VersionString() string {
	return fmt.Sprintf("Version: %s Commit: %s BuildDate: %s", Version, Hash, VersionDate)
}

func New() (*Config, error) {
	var config *Config

	var configPath string
	var printConfig bool
	var printVersion bool
	flag.StringVar(&configPath, "config", "./config.yaml", "path to a YAML config file")
	flag.BoolVar(&printConfig, "print-config", false, "print loaded config")
	flag.BoolVar(&printVersion, "version", false, "print build version")
	flag.Parse()

	if printVersion {
		fmt.Println(VersionString())
		return nil, nil
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("not found config file at %s, using defaults\n", configPath)
		config = DefaultConfig()
	} else {
		log.Printf("reading config at %s\n", configPath)
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	if printConfig {
		yamlData, _ := yaml.Marshal(&config)
		fmt.Println()
		fmt.Println(string(yamlData))
		return nil, nil
	}

	return config, nil
}
