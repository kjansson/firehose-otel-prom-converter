variable "name_prefix" {
  description = "Prefix to apply to all resources"
  type        = string
}

variable "create_amp_workspace" {
  description = "Create an AMP workspace"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = any
  default     = {}
}

variable "cloudwatch_metrics_stream_include_filter" {
  description = "Include filter for the CloudWatch metric stream"
  type        = any
  default     = {}
}

variable "cloudwatch_metrics_stream_exclude_filter" {
  description = "Include filter for the CloudWatch metric stream"
  type        = any
  default     = {}
}

variable "firehose_buffer_time_seconds" {
  description = "The buffering interval for the Firehose delivery stream"
  type        = number
  default     = 60
}

variable "create_firehose_log_group" {
  description = "Create a CloudWatch log group for the Firehose delivery stream"
  type        = bool
  default     = true
}

variable "create_lambda_log_group" {
  description = "Create a CloudWatch log group for the Lambda function"
  type        = bool
  default     = true
}

variable "firehose_log_group_retention" {
  description = "Retention period for the Firehose log group"
  type        = number
  default     = 7
}

variable "lambda_log_group_retention" {
  description = "Retention period for the Lambda log group"
  type        = number
  default     = 7
}

variable "firehose_log_group_name" {
  description = "Name of the Firehose log group"
  type        = string
  default    =  ""
}

variable "lambda_log_group_name" {
  description = "Name of the Lambda log group"
  type        = string
  default    =  ""
}

variable "create_s3_bucket" {
  description = "Create an S3 bucket for the Firehose delivery stream"
  type        = bool
  default     = true
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket for the Firehose delivery stream"
  type        = string
    default     = ""
}

variable "s3_bucket_retention_days" {
  description = "Retention period for the S3 bucket"
  type        = number
  default     = 7
}

variable "prometheus_endpoint" {
  description = "The endpoint of the Prometheus workspace"
  type        = string
  default = ""
}

variable "prometheus_region" {
  description = "The region of the Prometheus server"
  type        = string
}

variable "dimension_filter" {
  description = "The dimension filter for the Lambda function"
  type        = string
  default     = ""
}

variable "lambda_image_uri" {
  description = "The URI of the Lambda function image"
  type        = string
}

variable "lambda_image_repository_arn" {
  description = "The ARN of the Lambda function image repository"
  type        = string
}

variable "lambda_log_level" {
  description = "The log level for the Lambda function"
  type        = string
  default     = "INFO"
}