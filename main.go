package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const (
	region         = "us-east-1"
	namespace      = "AWS/Billing"
	metricName     = "EstimatedCharges"
	dimensionName  = "Currency"
	dimensionValue = "USD"
)

func getBilling() float64 {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal(err)
	}

	cw := cloudwatch.New(sess)

	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(21600),
		StartTime:  aws.Time(time.Now().Add(time.Duration(21600) * time.Second * -1)),
		EndTime:    aws.Time(time.Now()),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
		Unit: aws.String(cloudwatch.StandardUnitNone),
	}

	res, err := cw.GetMetricStatistics(params)
	if err != nil {
		log.Fatal(err)
	}

	return float64(*res.Datapoints[0].Maximum)
}

func slackNotify(billing float64) {
	json := fmt.Sprintf(
		`
{
  "text": "%sの AWS 請求金額",
  "attachments": [
    {
      "title": "Total",
      "text": "%s %.2f（ ¥ %.2f ）",
      "color": "#2eb886"
    }
  ]
}
	`, currentMonth(), dimensionValue, billing, billing*110.54)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("SLACK_INCOMING_WEBHOOK_URL"),
		bytes.NewBuffer([]byte(json)),
	)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
}

func currentMonth() string {
	nowUTC := time.Now().UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)
	return nowJST.Format("2006年01月")
}

func Run() {
	slackNotify(getBilling())
}

func main() {
	// Run local
	//Run()

	// Run lambda
	lambda.Start(Run)
}
