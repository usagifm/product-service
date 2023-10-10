package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/usagifm/product-service/contract"
	"github.com/usagifm/product-service/model"
	"github.com/usagifm/product-service/utils/logger"
)

type ProductService interface {
	GetProductById(ctx context.Context, getProductByIdReq contract.ProductIdRequest) (*model.Product, error)
	CreateProduct(ctx context.Context, reqContract contract.CreateProductRequest) error
	GetProductList(ctx context.Context, reqContract contract.GetProductListRequest) (model.ProductList, error)
}

func GetProductByIdHandler(productSvc ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, errValidate := contract.BuildAndValidateProductIdRequest(r)
		if errValidate != nil {
			logger.GetLogger(ctx).Errorf("Get product by id validation err:%v\n", errValidate)
		}
		product, err := productSvc.GetProductById(ctx, req)
		if err != nil {
			logger.GetLogger(ctx).Errorf("Get product by id err:%v\n", err)

			if errors.Is(err, sql.ErrNoRows) {
				contract.ErrorResponse(ctx, w, "Get product failed", "Product not found", "Error", http.StatusNotFound)
				return
			}

			contract.ErrorResponse(ctx, w, "Get product failed", "Internal Server Error", "Error", http.StatusInternalServerError)
			return
		}

		contract.SuccessResponse(ctx, w, "Get product success", product, http.StatusOK)
	}
}

func CreateProductHandler(productSvc ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, errValidate := contract.BuildAndValidateCreateProductRequest(r)
		if errValidate != nil {
			logger.GetLogger(ctx).Errorf("Create product validation err:%v\n", errValidate)
			contract.ErrorResponse(ctx, w, "Failed to create product", "Bad Request", "Error", http.StatusBadRequest)
			return
		}

		err := productSvc.CreateProduct(ctx, req)
		if err != nil {
			logger.GetLogger(ctx).Errorf("Create product err:%v\n", err)
			contract.ErrorResponse(ctx, w, "Failed to create product", err.Error(), "Error", http.StatusInternalServerError)
			return
		}

		contract.SuccessResponse(ctx, w, "Create product success", nil, http.StatusOK)

	}
}

func GetProductListHandler(productSvc ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, errValidate := contract.BuildAndValidateGetProductListRequest(r)
		if errValidate != nil {
			logger.GetLogger(ctx).Errorf("Get product list validation err:%v\n", errValidate)
			contract.ErrorResponse(ctx, w, "Failed to get product list", "Bad Request", "Error", http.StatusBadRequest)
			return
		}

		productList, err := productSvc.GetProductList(ctx, req)
		if err != nil {
			logger.GetLogger(ctx).Errorf("Get product list err:%v\n", err)

			contract.ErrorResponse(ctx, w, "Get product list failed", "Internal Server Error", "Error", http.StatusInternalServerError)
			return
		}

		contract.SuccessResponse(ctx, w, "Get product list success", productList, http.StatusOK)

	}
}
