package app_context

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/sftp"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/product-service/config"
	"github.com/usagifm/product-service/utils/logger"
)

type appContext struct {
	db               *sqlx.DB
	requestValidator *validator.Validate
	redisClient      *redis.Client
	SFTPClient       *sftp.Client
}

var appCtx appContext

func Init(ctx context.Context) error {
	db, err := initDb(ctx, config.Get().Postgres)
	if err != nil {
		return err
	}

	redisClient, err := initRedis(ctx, config.Get().Redis)
	if err != nil {
		return err
	}

	appCtx = appContext{
		db:               db,
		redisClient:      redisClient,
		requestValidator: validator.New(),
	}

	return nil
}

func initDb(ctx context.Context, cfg config.Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.ConnURI)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to connect the database err:%v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.GetLogger(ctx).Errorf("failed to ping the database err:%v", err)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxPoolSize)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)
	db.SetConnMaxLifetime(cfg.MaxLifeTime)
	return db, nil
}

func initRedis(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		logger.GetLogger(ctx).Error("init redis fail: ", err)
		return nil, err
	}
	//redisClient.AddHook(nrredis.NewHook(opts))
	return redisClient, nil
}

func RequestValidator() *validator.Validate {
	return appCtx.requestValidator
}

func GetDB() *sqlx.DB {
	return appCtx.db
}

func GetRedisClient() *redis.Client {
	return appCtx.redisClient
}

func GetSFTPClient() *sftp.Client {
	return appCtx.SFTPClient
}
