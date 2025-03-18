
resource "aws_cloudwatch_metric_stream" "this" {
  name          = format("%s-metrics-stream", var.name_prefix)
  role_arn      = aws_iam_role.metric_stream.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.this.arn
  output_format = "opentelemetry1.0"

  dynamic "include_filter" {
    for_each = var.cloudwatch_metrics_stream_include_filter
    content {
      namespace    = include_filter.value.namespace
      metric_names = try(include_filter.value.metric_names, [])
    }
  }
  dynamic "exclude_filter" {
    for_each = var.cloudwatch_metrics_stream_exclude_filter
    content {
      namespace    = exclude_filter.value.namespace
      metric_names = try(exclude_filter.value.metric_names, [])
    }
  }
}

