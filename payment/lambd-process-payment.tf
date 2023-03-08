
resource "null_resource" "process_payment_binary" {
    provisioner "local-exec" {
        interpreter = [
          "PowerShell","-Command"
        ]
        command = "$env:GOOS=\"linux\" ;$env:GOARCH=\"amd64\"; $env:CGO_ENABLED=0; $env:GOFLAGS=\"-trimpath\"; go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ./${local.src_path}"
    }
}

data "archive_file" "process_payment_archive"{
    depends_on = [null_resource.process_payment_binary]
    type = "zip"
    source_file = local.binary_path
    output_path = local.archive_path
}


resource "aws_lambda_function" "process_payment" {
  function_name = local.function_name
  description = "Service to process a payment and complete the order"
  role = aws_iam_role.process_payment.arn
  handler = local.binary_name
  memory_size = 128
  filename = local.archive_path
  source_code_hash = data.archive_file.process_payment_archive.output_base64sha256
  runtime = "go1.x"
  environment {
    variables = {
      QUEUE_URL = var.process_payment_queue_url
    }
  }
  
}
resource "aws_cloudwatch_log_group" "payment_log" {
  name              = "/aws/lambda/${aws_lambda_function.process_payment.function_name}"
  retention_in_days = 7
  
}