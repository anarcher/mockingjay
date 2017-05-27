package handler

import (
	"github.com/anarcher/mockingjay/pkg/config"
	"github.com/anarcher/mockingjay/pkg/log"

	"github.com/elazarl/goproxy"
	kitlog "github.com/go-kit/kit/log"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReqCondition struct {
	cfg    config.Filter
	logger kitlog.Logger
}

func (rq *ReqCondition) HandleReq(_req *http.Request, ctx *goproxy.ProxyCtx) (result bool) {
	rq.logger.Log("req", fmt.Sprintf("%v", _req.URL))

	req := new(http.Request)
	*req = *_req

	{
		buf, err := ioutil.ReadAll(_req.Body)
		if err != nil {
			rq.logger.Log("err", err)
		}

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = rdr1
		defer func() {
			_req.Body = rdr2
		}()
	}

	if err := req.ParseForm(); err != nil {
		rq.logger.Log("err", err)
		return
	}

	result = true

	for k, v := range rq.cfg.Form {
		if req.Form.Get(k) != v {
			result = false
			break
		}
	}

	return
}

func (rq *ReqCondition) HandleResp(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
	return true
}

func NewReqCondition(cfg config.Filter) *ReqCondition {
	rc := &ReqCondition{
		cfg:    cfg,
		logger: kitlog.With(log.Logger, "m", "ReqCondition"),
	}

	return rc
}
