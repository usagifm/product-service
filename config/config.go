package config

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/usagifm/product-service/utils/logger"

	"github.com/spf13/viper"
)

/*
	All config should be required.
	Optional only allowed if zero value of the type is expected being the default value.
	time.Duration units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”. as in time.ParseDuration().
*/

type (
	Postgres struct {
		ConnURI            string        `mapstructure:"PG_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"PG_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"PG_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"`
	}
	Redis struct {
		Host           string `mapstructure:"REDIS_HOST" validate:"required"`
		RootDocument   string `mapstructure:"REDIS_ROOT_DOCUMENT"`
		Password       string `mapstructure:"REDIS_PASSWORD"`
		InvalidateTime int    `mapstructure:"REDIS_INVALIDATE_TIME"`
	}

	Config struct {
		Postgres    Postgres      `mapstructure:",squash"`
		Redis       Redis         `mapstructure:",squash"`
		APITimeout  time.Duration `mapstructure:"API_TIMEOUT" validate:"required"`
		BindAddress int           `mapstructure:"BIND_ADDRESS" validate:"required"`
	}
)

var cfg *Config

func Init(ctx context.Context) error {
	cfg = &Config{}

	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.GetLogger(ctx).Errorf("failed to read config:%v", err)
		return err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		logger.GetLogger(ctx).Errorf("failed to bind config:%v", err)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.GetLogger(ctx).Errorf("invalid config:%v", err)
		}
		logger.GetLogger(ctx).Errorf("failed to load config")
		return err
	}

	logger.GetLogger(ctx).Infof("Config loaded: %+v", cfg)
	return nil
}

func Get() *Config {
	return cfg
}
