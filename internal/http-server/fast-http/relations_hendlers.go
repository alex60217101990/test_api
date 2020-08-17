package fast_http

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) AddCategoryProductRelation(ctx *fasthttp.RequestCtx) {
	rel, err := parseRelation(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.AddRelationCategory(ctx, rel.ProductID, rel.CategoryID)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx,
		fmt.Sprintf("category [%s] relation for product [%s] added success",
			rel.CategoryID, rel.ProductID),
		fasthttp.StatusOK)
}

func (s *FastHttpServer) DeleteCategoryProductRelation(ctx *fasthttp.RequestCtx) {
	rel, err := parseRelation(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.DelRelationCategory(ctx, rel.ProductID, rel.CategoryID)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx,
		fmt.Sprintf("category [%s] relation for product [%s] remove success",
			rel.CategoryID, rel.ProductID),
		fasthttp.StatusOK)
}
