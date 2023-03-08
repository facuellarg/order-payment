
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }

  required_version = ">= 1.3.7"
}


module "oders-module" {
  source = "./order"
  depends_on = [
    aws_sqs_queue.complete_order
  ]
  
}

module "payment-module" {
  source = "./payment"
  depends_on = [
    aws_sqs_queue.complete_order
  ]
  process_payment_queue_url = aws_sqs_queue.complete_order.url
  process_payment_queue_arn = aws_sqs_queue.complete_order.arn
}

provider "aws" {
  region = "us-east-1"
  # profile = "tutorial-terraform-profile"

  default_tags {
    tags = {
      app = "tutorial-terraform"
    }
  }
}