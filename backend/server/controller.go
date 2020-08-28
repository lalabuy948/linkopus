package server

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type (
	linkBody struct {
		Link string `json:"link"`
	}

	badResponse struct {
		ErrorMsg string `json:"error_msg"`
	}
)

func (s *server) redirectHandler(ctx *fasthttp.RequestCtx) {
	linkHash := string(ctx.Path())[1:]

	linkMap, err := s.container.FacadeService.HandleLinkExtraction(linkHash)
	if err != nil || linkMap == nil {
		responseErrorJson(ctx, err)
		return
	}

	ctx.Redirect(linkMap.Link, fasthttp.StatusTemporaryRedirect)
	go s.container.CommandService.SaveLinkView(linkMap.Link)

	return
}

func (s *server) linkHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	if ctx.IsGet() {

		linkParam := string(ctx.QueryArgs().Peek("link"))

		linkMap, err := s.container.QueryService.QueryLinkMap(linkParam, "")
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		data, err := json.Marshal(linkMap)
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(data)
		return
	}

	if ctx.IsPost() {

		var linkBody linkBody
		if err := json.Unmarshal(ctx.PostBody(), &linkBody); err != nil {
			responseErrorJson(ctx, err)
			return
		}

		linkHash, err := s.container.FacadeService.HandleLinkMapInsert(linkBody.Link)
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		data, err := json.Marshal(linkHash)
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
		ctx.Response.SetBody(data)
		return
	}

}

func (s *server) statsHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	if ctx.IsGet() {
		linkParam := string(ctx.QueryArgs().Peek("link"))

		linkViews, err := s.container.QueryService.QueryLinkViews(linkParam, "")
		if err != nil && err != mongo.ErrNoDocuments {
			responseErrorJson(ctx, err)
			return
		}

		data, err := json.Marshal(linkViews)
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(data)
		return
	}
}

func (s *server) topHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	if ctx.IsGet() {
		linkViews, err := s.container.FacadeService.HandleTodayTopLinksViewsExtraction()
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		data, err := json.Marshal(linkViews)
		if err != nil {
			responseErrorJson(ctx, err)
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(data)
		return
	}

}

func responseErrorJson(ctx *fasthttp.RequestCtx, err error) {
	errStruct := badResponse{ErrorMsg: err.Error()}
	errMsg, _ := json.Marshal(errStruct)

	ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.Response.SetBody(errMsg)
	return
}
