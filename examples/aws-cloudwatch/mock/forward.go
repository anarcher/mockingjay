package main

import (
	forward "github.com/anarcher/mockingjay/pkg/forward/cloudwatch"

	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"net/http"
)

func fwCloudwatchGetMetricStatistics(cw *cloudwatch.CloudWatch, w http.ResponseWriter, r *http.Request) error {
	GetMetricStatistics := forward.NewGetMetricStatistics(cw)
	return GetMetricStatistics.Forward(w, r)
}
