package handler

import (
	"github.com/anarcher/mockingjay/pkg/config"
	"github.com/anarcher/mockingjay/pkg/log"

	"github.com/elazarl/goproxy"
	kitlog "github.com/go-kit/kit/log"

	"fmt"
	"net/http"
)

type ReqHandler struct {
	cfg    config.Proxy
	logger kitlog.Logger
}

// An implementation of goproxy.ReqHandler
func (rh *ReqHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//rh.logger.Log("request", fmt.Sprintf("%+v", req.URL))
	if rh.cfg.Target == "" {
		return req, nil // pass via proxy
	}

	targetURL := rh.cfg.Target
	rh.logger.Log("target.URL", targetURL)
	resp, err := DoRequest(req, targetURL)
	if err != nil {
		rh.logger.Log("err", err)
		return nil, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusBadGateway, fmt.Sprintf("err: %s", err))
	}

	return nil, resp
}

func NewReqHandler(cfg config.Proxy) *ReqHandler {
	rq := &ReqHandler{
		cfg:    cfg,
		logger: kitlog.With(log.Logger, "m", "ReqHandler"),
	}

	return rq
}
