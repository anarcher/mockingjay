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
	"time"
)

var (
	proxyURL string
)

func init() {
	flag.StringVar(&proxyURL, "proxyURL", "http://localhost:8080", "")
}

func main() {

	flag.Parse()

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

	sess := session.Must(session.NewSession(cfg))
	svc := cloudwatch.New(sess)
	params := &cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String("MetricName"),
		Namespace:  aws.String("Namespace"),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("AutoScalingGroupName"),
				Value: aws.String("TEST"),
			},
		},
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
