
resource "aws_sqs_queue" "complete_order" {
    name = "complete_order_queue"
    delay_seconds = 0
    max_message_size = 262144
    message_retention_seconds = 345600
    receive_wait_time_seconds = 10
    visibility_timeout_seconds = 300
}