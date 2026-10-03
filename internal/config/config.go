package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ConfigEnv struct {
	// RabbitMQDial     string `mapstructure:"RABBITMQ_DIAL"`
	RabbitMQUser     string `mapstructure:"RABBITMQ_USER"`
	RabbitMQDial     string `mapstructure:"RABBITMQ_DIAL"`
	RabbitMQPassword string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQHost     string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort     string `mapstructure:"RABBITMQ_PORT"`
	AppPort          string `mapstructure:"APP_PORT"`
}

var keys = []string{
	// "RABBITMQ_DIAL", put at vps host
	"RABBITMQ_USER",
	"RABBITMQ_DIAL",
	"RABBITMQ_PASSWORD",
	"RABBITMQ_HOST",
	"RABBITMQ_PORT",
	"APP_PORT",
}

func LoadConfig(path string) (cfg ConfigEnv, err error) {
	v := viper.New()

	v.SetConfigFile(filepath.Join(path, "app.env")) // note app.env opsional, klo gak ada cek ke CI ENV
	v.SetConfigType("env")
	v.AutomaticEnv()

	//note: bedakan env var(di CI runner gitlab or github) dan app.env as file
	for _, k := range keys {
		log.Println("keys app env: ", k)
		if err := v.BindEnv(k); err != nil {
			return cfg, err
		}
	}

	if err = v.ReadInConfig(); err != nil {
		var pathErr *fs.PathError
		if !errors.As(err, &pathErr) || !errors.Is(err, fs.ErrNotExist) {
			// path ke app.env tidak error dan fs exist
			// jadi return config aja
			log.Println("cuss sini..")
			return cfg, err
		}
		// path error dan fs not exist, jdi baca ke env variable
		err = nil

	}

	if err = v.Unmarshal(&cfg); err != nil {
		return
	}

	err = overrideFromFile(v, "RABBITMQ_PASSWORD", &cfg.RabbitMQPassword)
	if err != nil {
		return
	}
	err = overrideFromFile(v, "RABBITMQ_DIAL", &cfg.RabbitMQDial)
	if err != nil {
		return
	}

	return

}

func overrideFromFile(v *viper.Viper, key string, target *string) error {
	p := v.GetString(key + "_FILE") // baca ke docker-compose
	if p == "" {
		return nil
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("read %s_FILE: %w", key, err)
	}
	*target = strings.TrimSpace(string(b))
	return nil
}

// get password from docker secret env
