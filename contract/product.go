package contract

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type ProductIdRequest struct {
	ProductId int `validate:"required"`
}

type CreateProductRequest struct {
	Name        string `validate:"required" json:"name"`
	Price       int    `validate:"required" json:"price"`
	Description string `validate:"required" json:"description"`
	Quantity    int    `validate:"required" json:"quantity"`
}

type GetProductListRequest struct {
	Limit    int ` validate:"required"`
	Offset   int ` validate:"required"`
	SortBy   string
	SortType string
	Search   string
}

// BUILD AND VALIDATE

func BuildAndValidateProductIdRequest(r *http.Request) (ProductIdRequest, error) {
	var payload ProductIdRequest

	productId, _ := strconv.Atoi(chi.URLParam(r, "product_id"))

	payload.ProductId = productId
	validate := validator.New()
	errValidate := validate.Struct(payload)

	if errValidate != nil {
		return ProductIdRequest{}, errValidate
	}

	return payload, nil
}

func BuildAndValidateCreateProductRequest(r *http.Request) (CreateProductRequest, error) {
	var payload CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return CreateProductRequest{}, err
	}

	validate := validator.New()
	errValidate := validate.Struct(payload)
	if errValidate != nil {
		return CreateProductRequest{}, errValidate
	}

	return payload, nil
}

func BuildAndValidateGetProductListRequest(r *http.Request) (GetProductListRequest, error) {
	var payload GetProductListRequest

	sortBy := r.URL.Query().Get("sort_by")
	sortType := strings.ToUpper(r.URL.Query().Get("sort_type"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	payload = GetProductListRequest{
		SortBy:   sortBy,
		SortType: sortType,
		Limit:    limit,
		Offset:   offset,
		Search:   search,
	}

	validate := validator.New()
	errValidate := validate.Struct(payload)
	if errValidate != nil {
		return GetProductListRequest{}, errValidate
	}

	return payload, nil
}
