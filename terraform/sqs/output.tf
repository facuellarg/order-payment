output "sqs_complete_order_url" {
    value = aws_sqs_queue.complete_order.url
}
output "sqs_complete_order_arn" {
    value = aws_sqs_queue.complete_order.arn
}
output "sqs_create_order_url" {
    value = aws_sqs_queue.create_order.url
}
output "sqs_create_order_arn" {
    value = aws_sqs_queue.create_order.arn
}