terraform {
  required_version = ">= 0.15"
  required_providers { aws = { version = ">= 3.36.0" } }
  backend "s3" {
    bucket = "data-engineer-terraform-states"
    region = "us-east-2"
    key = "lambda/app/lambda-import-file.state"
  }
}

data "aws_caller_identity" "current" {}

provider "aws" {
 region     = "us-east-2"
}