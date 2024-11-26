package http

import (
	"errors"
	"net/http"

	"github.com/rancher/apiserver/pkg/apierror"

	"github.com/cloudweav/cloudweav/pkg/util"
)

type CloudweavServerHandler interface {
	Do(ctx *Ctx) (interface{}, error)
}

type cloudweavServerHandler struct {
	httpHandler CloudweavServerHandler
}

func NewHandler(httpHandler CloudweavServerHandler) http.Handler {
	return &cloudweavServerHandler{
		httpHandler: httpHandler,
	}
}

func (handler *cloudweavServerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := newDefaultCloudweavServerCtx(rw, req)
	resp, err := handler.httpHandler.Do(ctx)
	if err != nil {
		status := http.StatusInternalServerError
		var e *apierror.APIError
		if errors.As(err, &e) {
			status = e.Code.Status
		}
		rw.WriteHeader(status)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	if resp == nil {
		rw.WriteHeader(ctx.StatusCode())
		return
	}

	util.ResponseOKWithBody(rw, resp)
}
