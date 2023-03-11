
// create an API gateway for HTTP API
resource "aws_apigatewayv2_api" "api_gw" {
  name          = "application-gateway"
  protocol_type = "HTTP"
  description   = "Serverless API gateway for HTTP API and AWS Lambda function"

  cors_configuration {
    allow_headers = ["*"]
    allow_methods = [
      "POST",
    ]
    allow_origins = [
      "*" // NOTE: here we should provide a particular domain, but for the sake of simplicity we will use "*"
    ]
    expose_headers = []
    max_age        = 0
  }
}

// create a stage for API GW
resource "aws_apigatewayv2_stage" "api_gw" {
  api_id = aws_apigatewayv2_api.api_gw.id

  name        = "dev"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
  depends_on = [aws_cloudwatch_log_group.api_gw]
}

// create logs for API GW
resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.api_gw.name}"

  retention_in_days = 7
}


// create lambda function to invoke lambda when specific HTTP request is made via API GW
resource "aws_apigatewayv2_integration" "create_order_lambda" {
  api_id = aws_apigatewayv2_api.api_gw.id
  # integration_uri  = module.oders_module.api_lambda_function_arn
  integration_uri  = var.create_order_lambda_arn 
  integration_type = "AWS_PROXY"
}

// specify route that will be used to invoke lambda function
resource "aws_apigatewayv2_route" "create_order_lambda" {
  api_id    = aws_apigatewayv2_api.api_gw.id
  route_key = "POST /order"
  target    = "integrations/${aws_apigatewayv2_integration.create_order_lambda.id}"
}

// provide permission for API GW to invoke lambda function
resource "aws_lambda_permission" "create_order_lambda" {
  statement_id  = "create-order-api-gateway"
  action        = "lambda:InvokeFunction"
  # function_name = module.oders_module.api_lambda_function_name
  function_name = var.create_order_lambda_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api_gw.execution_arn}/*/*"
}



resource "aws_apigatewayv2_integration" "api_payment" {
    api_id = aws_apigatewayv2_api.api_gw.id
    # integration_uri = module.payment_module.api_lambda_function_arn
    integration_uri = var.process_payment_lambda_arn
    integration_type = "AWS_PROXY"
}


resource "aws_apigatewayv2_route" "api_payment" {
  api_id = aws_apigatewayv2_api.api_gw.id
  route_key = "POST /payment"
  target = "integrations/${aws_apigatewayv2_integration.api_payment.id}"
}

resource "aws_lambda_permission" "payment_agw" {
    statement_id = "process-payment-api-gateway"
    action = "lambda:InvokeFunction"
    # function_name = module.payment_module.api_lambda_function_name
    function_name = var.process_payment_lambda_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_apigatewayv2_api.api_gw.execution_arn}/*/*"
}
output "api_url" {
  value = aws_apigatewayv2_stage.api_gw.invoke_url
}
