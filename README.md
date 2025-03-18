# firehose-otel-prom-converter

A Lambda processor for AWS Firehose that reads a Cloudwatch metric steam in OTEL format, converts and processes into Prometheus remote write format, and writes to Amazon Managed Prometheus.

## Usage

- Create a Amazon Managed Grafana workspace
- Create a Cloudwatch metric stream to AWS Firehose
- Deploy Lambda and configure as a processor in Firehose stream

## Configuration

The lambda takes the following environment variables as configuration;

- PROMETHEUS_REMOTE_WRITE_URL: The URL to the remote write endpoint of your AMP workspace
- DIMENSION_FILTER - (optional) A regular expression of Cloudwatch dimensions to include metrics for, more information in the dimensions section below

## Dimensions

### Filtering
Since Cloudwatch treats every unique set of dimensions as a metric, a lot of metrics with irrelevant dimension sets will be exported in the metric stream. To reduce the number of timeseries pushed into Prometheus, the environment variable DIMENSION_FILTER can be used.  

For example, EC2 metrics will be exported for the dimension ImageId (AMI), with no reference in the dimension of InstanceId. Timeseries summed up by these kind of dimensions are rarely used and the number of timeseries grows exponentially when number of dimensions increase.  
By filtering on the dimensions that should procude timeseries one can reduce them to the useful ones. For EC2 metrics, a filter of "^InstanceId$" could be a good choice.

### Processing
When exported from Cloudwatch, metrics contain a set of labels that are converted into Prometheus labels. They also contain the dimensions for the Cloudwatch metric nested into one single label, which makes them difficult to use in queries. For this reason, the "dimension label" is removed, and the nested contents are turned into regular labels.  

```
# Original metric labels
{"Dimensions": {"InstanceId: "i-1234567890"}, "Label1": "Value1", "Label2": "Value2"}
# Processed metric labels
{"InstanceId: "i-1234567890", "Label1": "Value1", "Label2": "Value2"}
```
