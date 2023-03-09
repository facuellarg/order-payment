// build the binary for the lambda function in a specified path
resource "null_resource" "function_binary_complete_order" {
  provisioner "local-exec" {
    interpreter = ["PowerShell", "-Command"]

    command = "$env:GOOS=\"linux\" ;$env:GOARCH=\"amd64\"; $env:CGO_ENABLED=0; $env:GOFLAGS=\"-trimpath\"; go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path_complete_order} ./${local.src_path_complete_order}"
  }
}

// zip the binary, as we can upload only zip files to AWS lambda
data "archive_file" "complete_order_archive" {
  depends_on = [null_resource.function_binary_complete_order]

  type        = "zip"
  source_file = local.binary_path_complete_order
  output_path = local.archive_path_complete_order
}

// create the lambda function from zip file
resource "aws_lambda_function" "complete_order" {
  function_name = local.function_name_complete_order
  description   = "Listen sqs to complet an order"
  role          = aws_iam_role.lambda.arn
  handler       = local.binary_name_complete_order
  memory_size   = 128

  filename         = local.archive_path_complete_order
  source_code_hash = data.archive_file.complete_order_archive.output_base64sha256

  // skip timeout
  runtime = "go1.x"
  // skip tags

  environment {
    variables = {
      QUEUE_URL = var.process_payment_queue_url
    }
  }
  // skip environment variables
}

resource "aws_cloudwatch_log_group" "log_group_complete_order" {
  name              = "/aws/lambda/${aws_lambda_function.complete_order.function_name}"
  retention_in_days = 7
}


resource "aws_lambda_event_source_mapping" "complete_order_sqs" {
  event_source_arn = var.process_payment_queue_arn
  function_name    = aws_lambda_function.complete_order.arn
  # starting_position = "LATEST"
}