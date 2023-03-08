resource "aws_apigatewayv2_api" "payment_gateway_api" {
    name = "payment_agw"
    protocol_type = "HTTP"
    description = "Gate way for payment, I should use the same for both serices"
    
    cors_configuration {
      allow_headers = ["*"]
      allow_methods  = [
        "POST",
      ]
      allow_origins = ["*"]
      expose_headers = []
      max_age = 0
    }
}

resource "aws_apigatewayv2_stage" "payment_stage" {
  api_id = aws_apigatewayv2_api.payment_gateway_api.id
  name = "dev"
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
  depends_on = [ aws_cloudwatch_log_group.api_gw ]
}
  
  
  
resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.payment_gateway_api.name}"
  retention_in_days = 7
}


resource "aws_apigatewayv2_integration" "api_payment" {
    api_id = aws_apigatewayv2_api.payment_gateway_api.id
    integration_uri = aws_lambda_function.process_payment.arn
    integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "api_payment" {
  api_id = aws_apigatewayv2_api.payment_gateway_api.id
  route_key = "POST /"
  target = "integrations/${aws_apigatewayv2_integration.api_payment.id}"
}

resource "aws_lambda_permission" "payment_agw" {
    statement_id = "process-payment-api-gateway"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.process_payment.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_apigatewayv2_api.payment_gateway_api.execution_arn}/*/*"
}
output "api_url" {
  value = aws_apigatewayv2_api.payment_gateway_api.api_endpoint
}