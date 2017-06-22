package main

import (
	"github.com/anarcher/mockingjay/pkg/xml"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"fmt"
	"net/http"
)

func (h *Handler) GetMetricStatistics(w http.ResponseWriter, r *http.Request, logger log.Logger) {

	logger.Log("req.Form", fmt.Sprintf("%v", r.Form))

	metricName := r.Form.Get("MetricName")
	namespace := r.Form.Get("Namespace")
	period := r.Form.Get("Period")
	startTime := r.Form.Get("StartTime")
	endTime := r.Form.Get("EndTime")
	dims := membersToMap(r, "Dimensions")

	if metricName == "metric-forward" {
		err := fwCloudwatchGetMetricStatistics(h.cloudwatch, w, r)
		level.Info(logger).Log("forward", metricName, "err", err)
		return
	}

	level.Debug(logger).Log("metricName", metricName,
		"namespace", namespace,
		"startTime", startTime, "endTime", endTime,
		"period", period,
		"dims", fmt.Sprintf("%+v", dims))

	m1, err := MetricStartEndTimeMatcher(startTime, endTime)
	if err != nil {
		level.Error(logger).Log("err", err)
		http.Error(w, err.Error(), 500)
		return
	}
	level.Debug(logger).Log("m1", fmt.Sprintf("%+v", m1))

	m2 := MetricDimMatcher(dims)
	level.Debug(logger).Log("m2", fmt.Sprintf("%+v", m2))

	m3 := MetricNameMatcher(metricName, namespace)
	level.Debug(logger).Log("m3", fmt.Sprintf("%v", m3))

	query := h.db.Select(m2, m3).Limit(1).OrderBy("CreatedAt").Reverse()

	cnt, err := query.Count(&Metric{})
	if err != nil {
		logger.Log("err", err)
		http.Error(w, err.Error(), 500)
		return
	}
	level.Info(logger).Log("metricName", metricName, "namespace", namespace, "dims", fmt.Sprintf("%v", dims), "cnt", cnt)

	var datapoints []*cloudwatch.Datapoint
	if cnt > 0 {
		var metrics []*Metric
		if err := query.Find(&metrics); err != nil {
			level.Error(logger).Log("err", err)
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
		level.Error(logger).Log("err", err)
		http.Error(w, "XML Error", 500)
		return
	}

	fmt.Fprintf(w, xmlRes)
}

func (h *Handler) SetDesiredCapacity(w http.ResponseWriter, r *http.Request, logger log.Logger) {

	namespace := "AWS/AutoScaling"
	asgName := r.FormValue("AutoScalingGroupName")
	dc := r.FormValue("DesiredCapacity")

	level.Info(logger).Log("asgName", asgName, "dc", dc)

	metric, err := NewASGInServiceInstancesMetric(namespace, asgName, dc)
	if err != nil {
		level.Error(logger).Log("err", err)
		http.Error(w, "dc error:", 500)
		return
	}

	if err := h.db.Save(metric); err != nil {
		level.Error(logger).Log("err", err)
		http.Error(w, "db error:", 500)
		return
	}

	output := &autoscaling.SetDesiredCapacityOutput{}

	xmlr, err := xml.Response("SetDesiredCapacity", output, "")
	if err != nil {
		level.Error(logger).Log("err", err)
		http.Error(w, "XMLError", 500)
		return
	}

	fmt.Fprintf(w, xmlr)
}
