package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	proxyURL   string
	metricName string
	namespace  string
	dims       string
)

func init() {
	flag.StringVar(&proxyURL, "proxy-url", "http://localhost:8080", "")
	flag.StringVar(&metricName, "m", "", "metric name")
	flag.StringVar(&namespace, "n", "", "namespace")
	flag.StringVar(&dims, "d", "", "Dimensions (k=v,k=v)")

}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if metricName == "" || namespace == "" {
		Usage()
		os.Exit(1)
	}

	_proxyURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := aws.NewConfig()
	cfg = cfg.WithRegion("us-east-1")
	cfg.DisableSSL = aws.Bool(true)
	cfg = cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(_proxyURL),
		},
	}

	dimensions := DimensionsFromFlag()

	sess := session.Must(session.NewSession(cfg))
	svc := cloudwatch.New(sess)
	params := &cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(metricName),
		Namespace:  aws.String(namespace),
		Dimensions: dimensions,
		Period:     aws.Int64(60),
		StartTime:  aws.Time(time.Now().Add(-10 * time.Second)),
		EndTime:    aws.Time(time.Now()),
		Statistics: []*string{aws.String("Minimum")},
	}
	resp, err := svc.GetMetricStatistics(params)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}

	fmt.Printf("%s", resp)
	fmt.Println()
}

func DimensionsFromFlag() []*cloudwatch.Dimension {
	var dimensions []*cloudwatch.Dimension

	for _, kv := range strings.Split(dims, ",") {
		kv := strings.Split(kv, "=")
		if len(kv) <= 1 {
			continue
		}
		d := &cloudwatch.Dimension{
			Name:  aws.String(kv[0]),
			Value: aws.String(kv[1]),
		}
		dimensions = append(dimensions, d)
	}
	return dimensions
}
