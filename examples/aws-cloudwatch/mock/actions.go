package main

import (
	"github.com/anarcher/mockingjay/pkg/log"
	"github.com/anarcher/mockingjay/pkg/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	"net/http"
)

func (h *Handler) GetMetricStatistics(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger

	logger.Log("req.Form", fmt.Sprintf("%v", r.Form))

	metricName := r.Form.Get("MetricName")
	namespace := r.Form.Get("Namespace")
	period := r.Form.Get("Period")
	startTime := r.Form.Get("StartTime")
	endTime := r.Form.Get("EndTime")
	dims := membersToMap(r, "Dimensions")

	if metricName == "metric-forward" {
		err := fwCloudwatchGetMetricStatistics(h.cloudwatch, w, r)
		logger.Log("forward", metricName, "err", err)
		return
	}

	logger.Log("metricName", metricName,
		"namespace", namespace,
		"startTime", startTime, "endTime", endTime,
		"period", period,
		"dims", fmt.Sprintf("%+v", dims))

	m1, err := MetricStartEndTimeMatcher(startTime, endTime)
	if err != nil {
		logger.Log("err", err)
		http.Error(w, err.Error(), 500)
		return
	}
	logger.Log("m1", fmt.Sprintf("%+v", m1))

	m2 := MetricDimMatcher(dims)
	logger.Log("m2", fmt.Sprintf("%+v", m2))

	m3 := MetricNameMatcher(metricName, namespace)
	logger.Log("m3", fmt.Sprintf("%v", m3))

	query := h.db.Select(m2, m3).Limit(1).OrderBy("CreatedAt").Reverse()

	cnt, err := query.Count(&Metric{})
	if err != nil {
		logger.Log("err", err)
		http.Error(w, err.Error(), 500)
		return
	}
	logger.Log("cnt", cnt)

	var datapoints []*cloudwatch.Datapoint
	if cnt > 0 {
		var metrics []*Metric
		if err := query.Find(&metrics); err != nil {
			logger.Log("err", err)
			http.Error(w, err.Error(), 500)
			return
		}

		for _, metric := range metrics {
			logger.Log("Value", metric.Value, "CreatedAt", fmt.Sprintf("%+v", metric.CreatedAt))
			datapoint := &cloudwatch.Datapoint{
				Minimum:   aws.Float64(metric.Value),
				Timestamp: aws.Time(metric.CreatedAt),
			}
			datapoints = append(datapoints, datapoint)
		}
	}

	output := &cloudwatch.GetMetricStatisticsOutput{
		Datapoints: datapoints,
		Label:      aws.String("test"),
	}
	xmlRes, err := xml.Response("GetMetricStatistics", output, "")
	if err != nil {
		logger.Log("err", err)
		http.Error(w, "XML Error", 500)
		return
	}

	fmt.Fprintf(w, xmlRes)
}

func (h *Handler) SetDesiredCapacity(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger

	if err := r.ParseForm(); err != nil {
		logger.Log("err", err)
		http.Error(w, "ParseForm Error", 500)
		return
	}

	namespace := "AWS/AutoScaling"
	asgName := r.FormValue("AutoScalingGroupName")
	dc := r.FormValue("DesiredCapacity")

	logger.Log("asgName", asgName, "dc", dc)

	metric, err := NewMetricString(namespace, asgName, dc)
	if err != nil {
		logger.Log("err", err)
		http.Error(w, "dc error:", 500)
		return
	}
	metric.AutoScalingGroupName = asgName

	if err := h.db.Save(metric); err != nil {
		logger.Log("err", err)
		http.Error(w, "db error:", 500)
		return
	}

	output := &autoscaling.SetDesiredCapacityOutput{}

	xmlr, err := xml.Response("SetDesiredCapacity", output, "")
	if err != nil {
		logger.Log("err", err)
		http.Error(w, "XMLError", 500)
		return
	}

	fmt.Fprintf(w, xmlr)
}
