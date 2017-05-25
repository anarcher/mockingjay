package xml

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func TestResponse(t *testing.T) {

	output := cloudwatch.GetMetricStatisticsOutput{
		Datapoints: []*cloudwatch.Datapoint{
			&cloudwatch.Datapoint{
				Average: aws.Float64(1),
				Unit:    aws.String("None"),
			},
			&cloudwatch.Datapoint{
				Average: aws.Float64(2),
			},
		},
		Label: aws.String("label"),
	}

	xml, err := Response("GetMetricStatistics", output, "")
	if err != nil {
		t.Error(err)
	}

	expected := "<GetMetricStatisticsResponse><GetMetricStatisticsResult><Datapoints><member><Average>1</Average><Unit>None</Unit></member><member><Average>2</Average></member></Datapoints><Label>label</Label></GetMetricStatisticsResult></GetMetricStatisticsResponse>"
	if xml != expected {
		t.Error("have: %v want: %v", xml, expected)
	}

}
