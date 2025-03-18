resource "aws_cloudwatch_log_group" "firehose" {
  count             = var.create_firehose_log_group ? 1 : 0
  name              = format("%s-firehose", var.name_prefix)
  retention_in_days = 7
}

resource "aws_cloudwatch_log_stream" "firehose" {
  name           = format("%s-firehose-stream", var.name_prefix)
  log_group_name = var.create_firehose_log_group ? aws_cloudwatch_log_group.firehose[0].name : var.create_firehose_log_group
}

resource "aws_kinesis_firehose_delivery_stream" "this" {
  name        = format("%s-metrics-stream", var.name_prefix) // TODO fix name
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose.arn
    bucket_arn = aws_s3_bucket.this[0].arn

    processing_configuration {
      enabled = "true"
      processors {
        type = "Lambda"
        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = aws_lambda_function.this.arn
        }
      }
    }
    cloudwatch_logging_options {
      enabled         = true
      log_group_name  = var.create_firehose_log_group ? aws_cloudwatch_log_group.firehose[0].name : var.create_firehose_log_group
      log_stream_name = aws_cloudwatch_log_stream.firehose.name
    }
    buffering_interval = var.firehose_buffer_time_seconds
  }
}
