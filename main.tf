
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
  source = "./terraform/sqs"
  
}


module "oders_module" {
  source = "./order/terraform"
  depends_on = [
    module.sqs
  ]
  process_payment_queue_url = module.sqs.sqs_complete_order_url
  process_payment_queue_arn = module.sqs.sqs_complete_order_arn
  create_order_queue_url = module.sqs.sqs_create_order_url
  create_order_queue_arn = module.sqs.sqs_create_order_arn
  
}

module "payment_module" {
  source = "./payment/terraform"
  depends_on = [
    module.sqs
  ]
  process_payment_queue_url = module.sqs.sqs_complete_order_url
  process_payment_queue_arn = module.sqs.sqs_complete_order_arn
  create_order_queue_url = module.sqs.sqs_create_order_url
  create_order_queue_arn = module.sqs.sqs_create_order_arn
}


module "api_gw"{
  source = "./terraform/api_gw"
  depends_on = [
    module.oders_module,
    module.payment_module,
  ]
  create_order_lambda_arn = module.oders_module.api_lambda_function_arn
  create_order_lambda_name = module.oders_module.api_lambda_function_name
  process_payment_lambda_arn = module.payment_module.api_lambda_function_arn
  process_payment_lambda_name = module.payment_module.api_lambda_function_name
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

output "api_url" {
  value = module.api_gw.api_url
}