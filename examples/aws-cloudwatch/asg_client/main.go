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
	asgName  string
	desired  int
	debug    bool
)

func init() {
	flag.StringVar(&proxyURL, "proxy", "http://localhost:8080", "")
	flag.StringVar(&asgName, "a", "", "ASG Name")
	flag.IntVar(&desired, "d", 1, "DesiredCapacity")
	flag.BoolVar(&debug, "debug", false, "")
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if asgName == "" {
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
	if debug {
		cfg = cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	}
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(_proxyURL),
		},
	}

	sess := session.Must(session.NewSession(cfg))
	svc := autoscaling.New(sess)
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(int64(desired)),
	}
	resp, err := svc.SetDesiredCapacity(input)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}

	fmt.Printf("%s", resp)
	fmt.Println()
}
