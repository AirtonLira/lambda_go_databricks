
data "archive_file" "init" {
  type        = "zip"
  source_dir  = "${abspath("../")}/cmd/"
  output_path = "${abspath("./")}/archive.zip"
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${var.function_name}"

  retention_in_days = 7
}


resource "aws_iam_role" "lambda_exec_role" {
  name        = "lambda_exec"
  path        = "/"
  description = "Allows Lambda Function to call AWS services on your behalf."

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "function_logging_policy" {
  name   = "function-logging-policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        Action : [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Effect : "Allow",
        Resource : "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "function_logging_policy_attachment" {
  role = aws_iam_role.lambda_exec_role.id
  policy_arn = aws_iam_policy.function_logging_policy.arn
}

resource "aws_lambda_function" "trigger" {
  role             = "${aws_iam_role.lambda_exec_role.arn}"
  handler          = "${var.handler}"
  runtime          = "${var.runtime}"
  filename         = "./archive.zip"
  function_name    = "${var.function_name}"
  source_code_hash = "${data.archive_file.init.output_base64sha256}"
}

output "aws_lambda_function" {
  value = aws_lambda_function.trigger.function_name
}




resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.trigger.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.invoice_files.arn
}


resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.invoice_files.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.trigger.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}