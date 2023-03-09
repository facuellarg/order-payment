
output "api_lambda_function_name" {
    value = aws_lambda_function.process_payment.function_name
}
output "api_lambda_function_arn" {
    value = aws_lambda_function.process_payment.arn
}