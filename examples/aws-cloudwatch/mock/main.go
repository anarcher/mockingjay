package main

import (
	"github.com/anarcher/mockingjay/pkg/log"

	"github.com/asdine/storm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"net/http"
	"os"
)

type Handler struct {
	mux        map[string]func(http.ResponseWriter, *http.Request)
	db         *storm.DB
	cloudwatch *cloudwatch.CloudWatch
}

func main() {
	logger := log.Logger

	db, err := storm.Open("metrics.db")
	defer db.Close()

	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	awsConfig := aws.NewConfig()
	awsConfig.WithRegion("us-east-1")
	awsSessionOptions := session.Options{
		Config:  *awsConfig,
		Profile: "development",
	}
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	cloudwatch := cloudwatch.New(awsSession)

	handler := &Handler{
		mux:        make(map[string]func(http.ResponseWriter, *http.Request)),
		db:         db,
		cloudwatch: cloudwatch,
	}
	handler.mux["GetMetricStatistics"] = handler.GetMetricStatistics
	handler.mux["SetDesiredCapacity"] = handler.SetDesiredCapacity

	server := http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Log("err", err)
	}

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logger := log.Logger
	action := r.FormValue("Action")
	logger.Log("action", action)
	if h, ok := h.mux[action]; ok {
		h(w, r)
		return
	} else {
		http.Error(w, "Action not found", 500)
	}
}
