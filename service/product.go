package service

import (
	"context"

	"github.com/usagifm/product-service/contract"
	"github.com/usagifm/product-service/model"
	"github.com/usagifm/product-service/utils/logger"
)

type ProductRepository interface {
	GetProductById(ctx context.Context, productId int) (model.Product, error)
	GetProductList(ctx context.Context, searchByName string, sortBy string, sortType string, limit int, offset int) (model.ProductList, error)
	CreateProduct(ctx context.Context, product model.Product) error
}

type ProductService struct {
	rProduct ProductRepository
}

func NewProductService(rProduct ProductRepository) *ProductService {

	return &ProductService{
		rProduct: rProduct,
	}

}

func (s *ProductService) GetProductById(ctx context.Context, getProductIdContract contract.ProductIdRequest) (*model.Product, error) {
	log := logger.GetLogger(ctx).WithField("Action : ", "get product by id")
	ctx = logger.WithLogger(ctx, log)

	product, err := s.rProduct.GetProductById(ctx, getProductIdContract.ProductId)
	if err != nil {
		return nil, err
	}

	return &product, nil

}

func (s *ProductService) GetProductList(ctx context.Context, reqContract contract.GetProductListRequest) (model.ProductList, error) {
	log := logger.GetLogger(ctx).WithField("Action : ", "get product list")
	ctx = logger.WithLogger(ctx, log)

	productList, err := s.rProduct.GetProductList(ctx, reqContract.Search, reqContract.SortBy, reqContract.SortType, reqContract.Limit, reqContract.Offset)
	if err != nil {
		return model.ProductList{}, err
	}

	return productList, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, reqContract contract.CreateProductRequest) error {
	log := logger.GetLogger(ctx).WithField("Action : ", "create product")
	ctx = logger.WithLogger(ctx, log)

	product := model.Product{
		Name:        reqContract.Name,
		Price:       reqContract.Price,
		Description: reqContract.Description,
		Quantity:    reqContract.Quantity,
	}

	err := s.rProduct.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
