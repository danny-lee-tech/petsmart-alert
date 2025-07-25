package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danny-lee-tech/petsmart-alert/internal/config"
	"github.com/danny-lee-tech/petsmart-alert/internal/history"
	"github.com/danny-lee-tech/petsmart-alert/internal/notifier"
	"github.com/danny-lee-tech/petsmart-alert/internal/rakuten"
	"gopkg.in/yaml.v2"
)

var DefaultConfigLocation = "configs/config.yml"

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal("Error getting configurations:", err)
		panic(err)
	}

	notif := &notifier.Notifier{
		PushBullet: &notifier.PushBulleter{
			APIKey: config.PushBulletApiKey,
			Tag:    "petsmart-deal-alert",
			Title:  "New Petsmart Alert",
		},
	}

	cashback, err := rakuten.RetrieveCashback("petsmart")
	if err != nil {
		log.Fatal("Error retrieving cashback:", err)
		panic(err)
	}

	history := history.Init("rakuten", 1)
	exists, err := history.CheckIfExists(strconv.Itoa(cashback))
	if err != nil {
		log.Fatal("Error retrieving history:", err)
		panic(err)
	}

	if !exists {
		notif.Notify(fmt.Sprintf("New Rakuten Cashback: %d%%", cashback))
		history.RecordItemIfNotExist(strconv.Itoa(cashback))
	}
}

func getConfig() (config.Config, error) {
	configLocation := getConfigLocation()
	cfgBytes, err := os.ReadFile(configLocation)
	if err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func getConfigLocation() string {
	configLocation := os.Getenv("CONFIG_LOCATION")
	if configLocation != "" {
		return configLocation
	}
	return DefaultConfigLocation
}

func validateConfig(cfg *config.Config) error {
	if cfg.PushBulletApiKey == "" {
		return errors.New("missing required field: pushbullet_api_key")
	}

	return nil
}
