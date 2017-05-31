package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	svc := autoscaling.New(sess)
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String("TEST"),
		DesiredCapacity:      aws.Int64(1000),
	}
	resp, err := svc.SetDesiredCapacity(input)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}

	fmt.Printf("%s", resp)
	fmt.Println()
}
