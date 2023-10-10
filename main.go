package main

import (
	"context"

	"github.com/usagifm/product-service/app_context"
	"github.com/usagifm/product-service/config"
	"github.com/usagifm/product-service/server"
	"github.com/usagifm/product-service/utils/logger"
)

func main() {
	ctx := context.Background()
	logger.Init(ctx)

	if err := config.Init(ctx); err != nil {
		panic(err)
	}

	if err := app_context.Init(ctx); err != nil {
		panic(err)
	}

	server.StartAPIServer(ctx)
}
