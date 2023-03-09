
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

module "sqs" {
  source = "./sqs"
  
}


module "oders-module" {
  source = "./order"
  depends_on = [
    module.sqs
  ]
  process_payment_queue_url = module.sqs.sqs_complete_order_url
  process_payment_queue_arn = module.sqs.sqs_complete_order_arn
  create_order_queue_url = module.sqs.sqs_create_order_url
  create_order_queue_arn = module.sqs.sqs_create_order_arn
  
}

module "payment-module" {
  source = "./payment"
  depends_on = [
    module.sqs
  ]
  process_payment_queue_url = module.sqs.sqs_complete_order_url
  process_payment_queue_arn = module.sqs.sqs_complete_order_arn
  create_order_queue_url = module.sqs.sqs_create_order_url
  create_order_queue_arn = module.sqs.sqs_create_order_arn
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