package fast_http

import (
	"fmt"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) CreateCategory(ctx *fasthttp.RequestCtx) {
	category, err := parseCategory(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.InsertCategory(ctx, category)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx, fmt.Sprintf("category [%s] created success", category.PublicID), fasthttp.StatusOK)
}

func (s *FastHttpServer) DeleteCategory(ctx *fasthttp.RequestCtx) {
	query, isHard, err := parseDynamicQuery(ctx, "hard")
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	if isHard {
		err = s.repo.DeleteHardCategory(ctx, query)
	} else {
		err = s.repo.DeleteSoftCategory(ctx, query)
	}

	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}
}

func (s *FastHttpServer) UpdateCategory(ctx *fasthttp.RequestCtx) {
	category, err := parseCategory(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.repo.UpdateCategory(ctx, category)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	messagePrint(ctx, fmt.Sprintf("category [%s] update success", category.PublicID), fasthttp.StatusOK)
}

func (s *FastHttpServer) CategoriesList(ctx *fasthttp.RequestCtx) {
	query, err := parsePaginationAndSort(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	var list []*models.Category
	if query.Assoc {
		list, err = s.repo.GetCategories(ctx, query.Pagination, query.Sort, struct{}{})
	} else {
		list, err = s.repo.GetCategories(ctx, query.Pagination, query.Sort)
	}

	messagePrint(ctx, list, fasthttp.StatusOK)
}

func (s *FastHttpServer) SingleCategory(ctx *fasthttp.RequestCtx) {
	query, assoc, err := parseDynamicQuery(ctx, "assoc")
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	var cat *models.Category
	if assoc {
		cat, err = s.repo.GetCategoryByNameOrID(ctx, query, struct{}{})
	} else {
		cat, err = s.repo.GetCategoryByNameOrID(ctx, query)
	}

	messagePrint(ctx, cat, fasthttp.StatusOK)
}
