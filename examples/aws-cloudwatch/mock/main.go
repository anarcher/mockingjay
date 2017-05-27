package main

import (
	"github.com/anarcher/mockingjay/pkg/log"
	"github.com/anarcher/mockingjay/pkg/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	"net/http"
)

type handler struct{}

var (
	mux map[string]func(http.ResponseWriter, *http.Request)
)

func main() {
	logger := log.Logger

	server := http.Server{
		Addr:    ":8081",
		Handler: &handler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["GetMetricStatistics"] = GetMetricStatistics

	if err := server.ListenAndServe(); err != nil {
		logger.Log("err", err)
	}

}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger
	if err := r.ParseForm(); err != nil {
		logger.Log("err", err)
		http.Error(w, "ParseForm Error", 500)
	}
	action := r.FormValue("Action")
	logger.Log("action", action)
	if h, ok := mux[action]; ok {
		h(w, r)
		return
	} else {
		http.Error(w, "Action not found", 500)
	}
}

func GetMetricStatistics(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger

	if err := r.ParseForm(); err != nil {
		logger.Log("err", err)
		http.Error(w, "ParseForm Error", 500)
		return
	}

	logger.Log("req.Form", fmt.Sprintf("%v", r.Form))

	output := &cloudwatch.GetMetricStatisticsOutput{
		Label: aws.String("test"),
	}
	xmlRes, err := xml.Response("GetMetricStatistics", output, "")
	if err != nil {
		logger.Log("err", err)
		http.Error(w, "XML Error", 500)
	}

	fmt.Fprintf(w, xmlRes)
}
