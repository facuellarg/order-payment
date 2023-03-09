// build the binary for the lambda function in a specified path
resource "null_resource" "function_binary" {
  provisioner "local-exec" {
    interpreter = ["PowerShell", "-Command"]

    command = "$env:GOOS=\"linux\" ;$env:GOARCH=\"amd64\"; $env:CGO_ENABLED=0; $env:GOFLAGS=\"-trimpath\"; go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ./${local.src_path_make_order}"
  }
}

// zip the binary, as we can upload only zip files to AWS lambda
data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "function" {
  function_name = "create-order"
  description   = "function that is trigger by api, create an order"
  role          = aws_iam_role.lambda.arn
  handler       = local.binary_name
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  // skip timeout
  runtime = "go1.x"
  // skip tags
  environment {
    variables = {
      QUEUE_URL = var.create_order_queue_url
    }
  }

  // skip environment variables
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}