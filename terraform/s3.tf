resource "random_string" "this" {
  count   = var.create_s3_bucket ? 1 : 0
  length  = 6
  special = false
  upper   = false
}

resource "aws_s3_bucket" "this" {
  count         = var.create_s3_bucket ? 1 : 0
  bucket        = var.create_s3_bucket ? format("%s-%s", var.name_prefix, random_string.this[0].result) : var.s3_bucket_name
  force_destroy = true
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  bucket = aws_s3_bucket.this[0].id

  rule {
    id     = "retention-rule"
    status = "Enabled"
    filter {
        object_size_greater_than = 1
    }
    expiration {
      days = var.s3_bucket_retention_days
    }
  }
}
