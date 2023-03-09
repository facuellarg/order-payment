// build the binary for the lambda function in a specified path
resource "null_resource" "function_binary_order_created" {
  provisioner "local-exec" {
    interpreter = ["PowerShell", "-Command"]

    command = "$env:GOOS=\"linux\" ;$env:GOARCH=\"amd64\"; $env:CGO_ENABLED=0; $env:GOFLAGS=\"-trimpath\"; go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path_order_created} ./${local.src_path_order_created}"
  }
}

// zip the binary, as we can upload only zip files to AWS lambda
data "archive_file" "order_created_archive" {
  depends_on = [null_resource.function_binary_order_created]

  type        = "zip"
  source_file = local.binary_path_order_created
  output_path = local.archive_path_order_created
}

// create the lambda function from zip file
resource "aws_lambda_function" "order_created" {
  function_name = local.function_name_order_created
  description   = "Listen sqs when a order is created"
  role          = aws_iam_role.process_payment.arn
  handler       = local.binary_name_order_created
  memory_size   = 128

  filename         = local.archive_path_order_created
  source_code_hash = data.archive_file.order_created_archive.output_base64sha256

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

resource "aws_cloudwatch_log_group" "log_group_order_created" {
  name              = "/aws/lambda/${aws_lambda_function.order_created.function_name}"
  retention_in_days = 7
}


resource "aws_lambda_event_source_mapping" "complete_order_sqs" {
  event_source_arn = var.create_order_queue_arn
  function_name    = aws_lambda_function.order_created.arn
  # starting_position = "LATEST"
}