package handler

import (
	"github.com/anarcher/mockingjay/pkg/config"

	"github.com/elazarl/goproxy"

	"net/http"
)

type ProxyHandler struct {
	cfg config.Proxy
	rh  *ReqHandler
	rc  *ReqCondition
}

func NewProxyHandler(cfg config.Proxy) (*ProxyHandler, error) {
	ph := &ProxyHandler{
		cfg: cfg,
	}
	ph.rh = NewReqHandler(cfg)
	ph.rc = NewReqCondition(cfg.Filter)

	return ph, nil
}

// An implementation of goproxy.ReqCondition interface
func (ph *ProxyHandler) HandleReq(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	return ph.rc.HandleReq(req, ctx)
}

func (ph *ProxyHandler) HandleResp(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
	return ph.rc.HandleResp(resp, ctx)
}

// An implementation of goproxy.ReqHandler
func (ph *ProxyHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return ph.rh.Handle(req, ctx)
}
