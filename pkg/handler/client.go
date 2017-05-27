package handler

import (
	"net/http"
	"net/url"
)

func DoRequest(req *http.Request, requrl string) (*http.Response, error) {
	/*
		reqURL, err := url.Parse(requrl)
		if err != nil {
			return nil, err
		}

		req.RequestURI = ""
		req.URL = reqURL
		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	*/

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay
	if req.ContentLength == 0 {
		outreq.Body = nil // Issue 16036: nil Body for http.Transport retries
	}

	outreq.RequestURI = ""

	reqURL, err := url.Parse(requrl)
	if err != nil {
		return nil, err
	}
	outreq.URL = reqURL

	c := &http.Client{}
	resp, err := c.Do(outreq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
