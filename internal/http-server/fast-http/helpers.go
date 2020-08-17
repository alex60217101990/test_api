package fast_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	http_server "github.com/alex60217101990/test_api/internal/http-server"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/valyala/fasthttp"
)

func messagePrint(ctx *fasthttp.RequestCtx, data interface{}, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	encoder := json.NewEncoder(ctx)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func errorPrint(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	if err1 := json.NewEncoder(ctx).Encode(map[string]string{
		"error": err.Error(),
	}); err1 != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func parsePaginationAndSort(ctx *fasthttp.RequestCtx) (query *models.ListRequest, err error) {
	query = &models.ListRequest{
		Pagination: &models.Pagination{
			From: string(ctx.QueryArgs().Peek("from")),
		},
		Sort: &models.SortedBy{
			FieldName: string(ctx.QueryArgs().Peek("field")),
		},
	}

	to := string(ctx.QueryArgs().Peek("to"))
	if len(string(to)) > 0 {
		i, err := strconv.ParseInt(string(to), 10, 16)
		if err != nil {
			return query, fmt.Errorf("empty or invalid 'to' query parameter")
		}
		query.Pagination.To = uint16(i)
	}

	desc := string(ctx.QueryArgs().Peek("desc"))
	if len(string(desc)) > 0 {
		b, err := strconv.ParseBool(string(desc))
		if err != nil {
			return query, fmt.Errorf("empty or invalid 'desc' query parameter")
		}
		query.Sort.Desc = b
	}

	assoc := string(ctx.QueryArgs().Peek("assoc"))
	if len(string(assoc)) > 0 {
		a, err := strconv.ParseBool(string(assoc))
		if err != nil {
			return query, fmt.Errorf("empty or invalid 'assoc' query parameter")
		}
		query.Assoc = a
	}
	return query, err
}

func parseDynamicQuery(ctx *fasthttp.RequestCtx, paramName string) (query string, dyn bool, err error) {
	publicID := ctx.QueryArgs().Peek("id")
	if len(publicID) == 0 {
		publicID = ctx.QueryArgs().Peek("name")
		if len(publicID) == 0 {
			return query, dyn, fmt.Errorf("empty 'name' or 'id' query parameter")
		}
	}

	d := string(ctx.QueryArgs().Peek(paramName)) // hard
	if len(string(d)) > 0 {
		dyn, err = strconv.ParseBool(string(d))
		if err != nil {
			return query, dyn, fmt.Errorf("empty or invalid '%s' query parameter", paramName)
		}
	}
	return string(publicID), dyn, err
}

func parseRelation(ctx *fasthttp.RequestCtx) (rel *models.RelationRequest, err error) {
	rel = &models.RelationRequest{}
	err = parsePostBody(ctx, rel)
	return rel, err
}

func parseCategory(ctx *fasthttp.RequestCtx) (cat *models.Category, err error) {
	cat = &models.Category{}
	err = parsePostBody(ctx, cat)
	return cat, err
}

func parseProduct(ctx *fasthttp.RequestCtx) (prod *models.Product, err error) {
	prod = &models.Product{}
	err = parsePostBody(ctx, prod)
	return prod, err
}

func parseCreeds(ctx *fasthttp.RequestCtx) (creeds *models.Credentials, err error) {
	creeds = &models.Credentials{}
	if ctx.IsPost() {
		err = parsePostBody(ctx, creeds)
		if err != nil {
			creeds, err = parseCreedsFromForm(ctx)
		}
	}
	err = http_server.ValidateCreeds(creeds)
	return creeds, err
}

func parseCreedsFromForm(ctx *fasthttp.RequestCtx) (*models.Credentials, error) {
	creeds := models.Credentials{
		Username: string(ctx.FormValue("username")),
		Email:    string(ctx.FormValue("email")),
		Password: string(ctx.FormValue("password")),
	}
	if len(creeds.Password) > 0 {
		return &creeds, nil
	}

	return nil, errors.New("empty credentials data")
}

func parsePostBody(ctx *fasthttp.RequestCtx, data interface{}) error {
	// Get the JSON body and decode into credentials
	return json.Unmarshal(ctx.PostBody(), data)
}
