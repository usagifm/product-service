package dependency

import (
	"context"

	"github.com/usagifm/product-service/app_context"
	"github.com/usagifm/product-service/config"
	sqlxRepo "github.com/usagifm/product-service/repository/sqlx"
	"github.com/usagifm/product-service/service"
	"github.com/usagifm/product-service/utils/logger"
)

type APIRepositories struct {
	Product sqlxRepo.ProductRepository
}

type APIServices struct {
	Product service.ProductService
}

type APIDependencies struct {
	Repositories *APIRepositories
	Services     *APIServices
}

func instantiateAPIRepositories(ctx context.Context) *APIRepositories {

	rProduct, err := sqlxRepo.InitProductRepositoryRepository(app_context.GetDB(), app_context.GetRedisClient(), &config.Get().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatalf("InitProductRepository err:%v", err)
		return nil
	}

	apiRepositories := APIRepositories{
		Product: rProduct,
	}

	return &apiRepositories
}

func instantiateAPIServices(ctx context.Context, r *APIRepositories) *APIServices {
	var dep APIServices

	dep.Product = *service.NewProductService(r.Product)

	return &dep
}

func InstantiateAPIDependencies(ctx context.Context) *APIDependencies {
	repositories := instantiateAPIRepositories(ctx)
	services := instantiateAPIServices(ctx, repositories)

	return &APIDependencies{
		Repositories: repositories,
		Services:     services,
	}
}
