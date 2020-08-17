package fast_http

import (
	"fmt"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) CreateProduct(ctx *fasthttp.RequestCtx) {
	product, err := parseProduct(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.InsertProduct(ctx, product)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx, fmt.Sprintf("product [%s] created success", product.PublicID), fasthttp.StatusOK)
}

func (s *FastHttpServer) DeleteProduct(ctx *fasthttp.RequestCtx) {
	query, isHard, err := parseDynamicQuery(ctx, "hard")
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	if isHard {
		err = s.repo.DeleteHardProduct(ctx, query)
	} else {
		err = s.repo.DeleteSoftProduct(ctx, query)
	}

	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}
}

func (s *FastHttpServer) UpdateProduct(ctx *fasthttp.RequestCtx) {
	product, err := parseProduct(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.UpdateProduct(ctx, product)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx, fmt.Sprintf("product [%s] update success", product.PublicID), fasthttp.StatusOK)
}

func (s *FastHttpServer) ProductsList(ctx *fasthttp.RequestCtx) {
	query, err := parsePaginationAndSort(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	var list []*models.Product
	if query.Assoc {
		list, err = s.repo.GetProducts(ctx, query.Pagination, query.Sort, struct{}{})
	} else {
		list, err = s.repo.GetProducts(ctx, query.Pagination, query.Sort)
	}

	messagePrint(ctx, list, fasthttp.StatusOK)
}

func (s *FastHttpServer) SingleProduct(ctx *fasthttp.RequestCtx) {
	query, assoc, err := parseDynamicQuery(ctx, "assoc")
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	var prod *models.Product
	if assoc {
		prod, err = s.repo.GetProductByNameOrID(ctx, query, struct{}{})
	} else {
		prod, err = s.repo.GetProductByNameOrID(ctx, query)
	}

	messagePrint(ctx, prod, fasthttp.StatusOK)
}

func (s *FastHttpServer) AddProductCategory(ctx *fasthttp.RequestCtx) {

}
