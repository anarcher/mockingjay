package cloudwatch

import (
	"github.com/anarcher/mockingjay/pkg/xml"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"

	"fmt"
	"net/http"
	"strconv"
	"time"
)

type GetMetricStatistics struct {
	cloudwatchAPI cloudwatchiface.CloudWatchAPI
}

func NewGetMetricStatistics(cloudwatchAPI cloudwatchiface.CloudWatchAPI) *GetMetricStatistics {
	f := &GetMetricStatistics{
		cloudwatchAPI: cloudwatchAPI,
	}

	return f
}

func (f *GetMetricStatistics) Forward(w http.ResponseWriter, r *http.Request) error {

	startTime, _ := time.Parse(time.RFC3339, r.FormValue("StartTime"))
	endTime, _ := time.Parse(time.RFC3339, r.FormValue("EndTime"))
	period, _ := strconv.Atoi(r.FormValue("Period"))

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(r.FormValue("Namespace")),
		MetricName: aws.String(r.FormValue("MetricName")),
		Period:     aws.Int64(int64(period)),
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Unit:       aws.String(r.FormValue("Unit")),
		Statistics: formValueToStringSlice(r, "Statistics", 5),
		Dimensions: formValueToDimensions(r),
	}

	if input.Unit == nil || *input.Unit == "" {
		input.Unit = aws.String("None")
	}

	req, output := f.cloudwatchAPI.GetMetricStatisticsRequest(input)

	if err := req.Send(); err != nil {
		return err
	}

	xmlr, err := xml.Response("GetMetricStatistics", output, "")
	if err != nil {
		return err
	}

	fmt.Fprintf(w, xmlr)

	return nil
}

func formValueToStringSlice(r *http.Request, name string, n int) []*string {
	var ret []*string

	for i := 1; i <= n; i++ {
		v := r.FormValue(fmt.Sprintf("%s.member.%d", name, i))
		if v != "" {
			ret = append(ret, &v)
		}
	}
	return ret
}

func formValueToDimensions(r *http.Request) []*cloudwatch.Dimension {
	var ret []*cloudwatch.Dimension

	for i := 1; i <= 10; i++ {
		n := r.FormValue(fmt.Sprintf("Dimensions.member.%d.Name"))
		v := r.FormValue(fmt.Sprintf("Dimensions.member.%d.Value"))

		if n == "" {
			continue
		}
		d := &cloudwatch.Dimension{
			Name:  aws.String(n),
			Value: aws.String(v),
		}
		ret = append(ret, d)

	}

	return ret
}
