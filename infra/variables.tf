variable "path_source_code" {
  default = "cmd"
}

variable "function_name" {
  default = "reprocessing_import_invoice"
}

variable "handler" {
  default = "main"
}

variable "runtime" {
  default = "go1.x"
}

variable "AWS_REGION" {
  default = "us-east-1"
}

variable "ELASTIC_IP_ASSOCIATE" {
  default = "18.206.166.162"
}
