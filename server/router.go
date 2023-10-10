package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usagifm/product-service/dependency"
	"github.com/usagifm/product-service/handler"
)

func Router(r *chi.Mux, deps *dependency.APIDependencies) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Group(func(api chi.Router) {
		api.Route("/api/v1", func(productRouter chi.Router) {

			productRouter.Get("/product", handler.GetProductListHandler(&deps.Services.Product))
			productRouter.Post("/product", handler.CreateProductHandler(&deps.Services.Product))
			productRouter.Get("/product/{product_id}", handler.GetProductByIdHandler(&deps.Services.Product))

		})
	})

}
