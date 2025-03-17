package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	remote "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite"
	"github.com/prometheus/prometheus/prompb"

	//"github.com/prometheus/prometheus/prompb"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
)

var errInvalidOTLPFormatStart = errors.New("unable to decode data length from message")

func unmarshalMetrics(record []byte) (pmetric.Metrics, error) {
	md := pmetric.NewMetrics()
	dataLen, pos := len(record), 0
	for pos < dataLen {
		n, nLen := proto.DecodeVarint(record[pos:])
		if nLen == 0 && n == 0 {
			return md, errInvalidOTLPFormatStart
		}
		req := pmetricotlp.NewExportRequest()
		pos += nLen
		err := req.UnmarshalProto(record[pos : pos+int(n)])
		pos += int(n)
		if err != nil {
			return pmetric.Metrics{}, fmt.Errorf("unable to unmarshal input at %d: %w", pos, err)
		}
		for i := 0; i < req.Metrics().ResourceMetrics().Len(); i++ {
			rm := req.Metrics().ResourceMetrics().At(i)
			for j := 0; j < rm.ScopeMetrics().Len(); j++ {
				sm := rm.ScopeMetrics().At(j)
				sm.Scope().SetName("github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver")
				sm.Scope().SetVersion("latest")
			}
		}
		req.Metrics().ResourceMetrics().MoveAndAppendTo(md.ResourceMetrics())
	}

	return md, nil
}

func sendRequest(ts []prompb.TimeSeries) (*http.Response, error) {

	r := &prompb.WriteRequest{
		Timeseries: ts,
	}
	tsProto, err := r.Marshal()
	if err != nil {
		panic(err)
	}

	encoded := snappy.Encode(nil, tsProto)
	body := bytes.NewReader(encoded)

	// Create an HTTP request from the body content and set necessary parameters.
	req, err := http.NewRequest("POST", os.Getenv("PROMETHEUS_REMOTE_WRITE_URL"), body)
	if err != nil {
		panic(err)
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	roleArn := os.Getenv("AWS_AMP_ROLE_ARN")

	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var awsCredentials *credentials.Credentials
	if roleArn != "" {
		awsCredentials = stscreds.NewCredentials(sess, roleArn, func(p *stscreds.AssumeRoleProvider) {
			p.RoleSessionName = "aws-sigv4-proxy-" + host
		})
	} else {
		awsCredentials = sess.Config.Credentials
	}

	signer := v4.NewSigner(awsCredentials)

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	_, err = signer.Sign(req, body, "aps", os.Getenv("PROMETHEUS_REGION"), time.Now())

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Println("Request to AMP failed with status: ", resp.StatusCode)
	}

	if err != nil {
		panic(err)
	}
	return resp, err
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, firehoseEvent events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {

	dimensionFilterEnv := os.Getenv("DIMENSION_FILTER")
	useDimensionFilter := dimensionFilterEnv != ""

	dimensionFilterExpression := regexp.MustCompile(dimensionFilterEnv)

	var response events.KinesisFirehoseResponse
	for _, record := range firehoseEvent.Records {

		metrics, err := unmarshalMetrics(record.Data)
		if err != nil {
			panic(err)
		}

		promTs, err := remote.FromMetrics(metrics, remote.Settings{
			SendMetadata:        false,
			DisableTargetInfo:   false,
			ExportCreatedMetric: false,
			AddMetricSuffixes:   false,
		})
		if err != nil {
			panic(err)
		}

		tsArray := []prompb.TimeSeries{}

		for _, ts := range promTs {

			labels := ts.GetLabels()

			dims := labels[0].Value // Extraxt the dimensions label
			newLables := labels[1:] // Remove the dimensions label

			var dimensions map[string]string
			err := json.Unmarshal([]byte(dims), &dimensions)
			if err != nil {
				panic(err)
			}

			valid := false
			for k, v := range dimensions { // Create new labels from dimensions
				newLables = append(newLables, prompb.Label{Name: k, Value: v})
				if useDimensionFilter && dimensionFilterExpression.MatchString(k) {
					valid = true
					continue
				}
			}
			if !valid {
				continue
			}

			ts.Labels = newLables // Replace labels
			tsArray = append(tsArray, *ts)

			fmt.Println(ts.String())

		}

		// _, err = sendRequest(tsArray)
		// if err != nil {
		// 	panic(err)
		// }

		var transformedRecord events.KinesisFirehoseResponseRecord
		transformedRecord.RecordID = record.RecordID
		transformedRecord.Result = events.KinesisFirehoseTransformedStateOk
		transformedRecord.Data = record.Data

		response.Records = append(response.Records, transformedRecord)
	}
	return response, nil
}
