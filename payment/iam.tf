data "aws_iam_policy_document" "assum_lambda_payment_role"{
    statement {
      actions=["sts:AssumeRole"]
        principals {
          type = "Service"
          identifiers = ["lambda.amazonaws.com"]
        }
    }
}


resource "aws_iam_role" "process_payment" {
    name = "ProcessPaymentRole"
    description = "Role for payment process"
    assume_role_policy = data.aws_iam_policy_document.assum_lambda_payment_role.json
}


data "aws_iam_policy_document" "allow_lambda_logging"{
    statement {
      effect = "Allow"
      actions = [
        "logs:CreateLogStream",
        "logs:PutLogEvents",
      ]
      resources = [
        "arn:aws:logs:*:*:*",
      ]
    }
}

resource "aws_iam_policy" "payment_loggin_policy" {
    name = "AllowPaymentLogginPolicy"
    description = "Policy for allow payment log"
    policy = data.aws_iam_policy_document.allow_lambda_logging.json
}



resource "aws_iam_role_policy_attachment" "lambda_loggin_policy_attachment" {
  role = aws_iam_role.process_payment.id
  policy_arn = aws_iam_policy.payment_loggin_policy.arn
}


//dynamo
data "aws_iam_policy_document" "allow_dynamodb_table_operations"{
  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:Scan"
    ]
    resources = [
      aws_dynamodb_table.payment.arn,
    ]
  }
}


resource "aws_iam_policy" "dynamodb_payment_policy" {
  name = "PaymentDynamoPolicy"
  description = "Allow payment write read and update in dynamo db"
  policy = data.aws_iam_policy_document.allow_dynamodb_table_operations.json
}


resource "aws_iam_role_policy_attachment" "payment_dynamodb_policy_attachment"{
  role = aws_iam_role.process_payment.id
  policy_arn = aws_iam_policy.dynamodb_payment_policy.arn
  depends_on = [
    aws_iam_role.process_payment
  ]
}



data "aws_iam_policy_document" "allow_sqs_operation"{
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage"
    ]
    resources = [
      var.process_payment_queue_arn,
    ]
  }
}

resource "aws_iam_policy" "sqs_payment_policy" {
  name = "PaymentSQSPolicy"
  description = "Allow payment send message in sqs queue"
  policy = data.aws_iam_policy_document.allow_sqs_operation.json
}

resource "aws_iam_role_policy_attachment" "payment_process_payment_policy_attachment"{
  role = aws_iam_role.process_payment.id
  policy_arn = aws_iam_policy.sqs_payment_policy.arn
  depends_on = [
    aws_iam_role.process_payment
  ]
}